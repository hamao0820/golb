package golb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/imports"
)

type Library struct {
	ImportPath   string
	SelectorName string
	Alias        string
}

type Bundler struct {
	libPackage string
	rootDir    string
}

func NewBundler(libPackage, goModDir string) Bundler {
	return Bundler{
		libPackage: libPackage,
		rootDir:    goModDir,
	}
}

func (b Bundler) Bundle(src string) (code string, err error) {
	codes := map[string]string{} // key: file path, value: 展開コード
	usedFuncs := map[string]struct{}{}

	// 再帰的に使われている関数を取得
	var dfs func(string, map[string]struct{})
	dfs = func(file string, visited map[string]struct{}) {
		if _, ok := visited[file]; ok {
			return
		}
		node, err := b.perseFile(file)
		if err != nil {
			return
		}

		for f := range b.getUsedFunction(node) {
			usedFuncs[f] = struct{}{}
		}

		importLibs := b.getImportedLibPackage(node)
		visited[file] = struct{}{}

		for _, lib := range importLibs {
			libDir := b.getDir(lib.ImportPath)
			libFiles, err := b.getFiles(libDir)
			if err != nil {
				panic(err)
			}
			for _, file := range libFiles {
				libPath := path.Join(libDir, file.Name())
				dfs(libPath, visited)
			}
		}
	}

	// 再帰的にファイルを取得
	var dfs2 func(string)
	dfs2 = func(file string) {
		if _, ok := codes[file]; ok {
			return
		}
		node, err := b.perseFile(file)
		if err != nil {
			return
		}

		importLibs := b.getImportedLibPackage(node)
		targetSelectors := map[string]struct{}{}
		targetImports := map[string]struct{}{}
		for _, lib := range importLibs {
			targetSelectors[lib.SelectorName] = struct{}{}
			targetImports[lib.ImportPath] = struct{}{}
		}

		// コメントを削除
		node.Comments = nil

		b.removeSelector(node, targetSelectors)
		b.removeAllImport(node)
		b.removeUnusedFunction(node, usedFuncs)
		code := b.nodeToString(node)
		if file != src {
			code = strings.Join(strings.Split(code, "\n")[1:], "\n") // package行を削除
		}
		codes[file] = code

		for _, lib := range importLibs {
			libDir := b.getDir(lib.ImportPath)
			libFiles, err := b.getFiles(libDir)
			if err != nil {
				panic(err)
			}
			for _, file := range libFiles {
				libPath := path.Join(libDir, file.Name())
				dfs2(libPath)
			}
		}
	}

	dfs(src, map[string]struct{}{})

	dfs2(src)

	sourceCode := codes[src]
	sourceCode += "//" + strings.Repeat("-", 50) + "以下は生成コード" + strings.Repeat("-", 50)
	sourceCode += "\n\n"
	for file, code := range codes {
		if file == src {
			continue
		}
		dirs := strings.Split(file, "/")
		sourceCode += "// " + path.Join(dirs[len(dirs)-2:]...)
		sourceCode += code
		sourceCode += "\n\n"
	}

	formatted, err := imports.Process("", []byte(sourceCode), nil)
	if err != nil {
		return "", err
	}

	formatted, err = format.Source(formatted)
	if err != nil {
		return "", err
	}

	return string(formatted), nil
}

// ASTを取得
func (b Bundler) perseFile(filename string) (*ast.File, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	return node, nil
}

// importされている自作ライブラリのパッケージを取得
func (b Bundler) getImportedLibPackage(file *ast.File) []Library {
	libs := []Library{}

	for _, imp := range file.Imports {
		if !b.isLibrary(imp.Path.Value) {
			continue
		}

		// aliasがある場合はaliasを使う
		if imp.Name != nil {
			libs = append(libs, Library{
				ImportPath:   imp.Path.Value,
				SelectorName: imp.Name.Name,
				Alias:        imp.Name.Name,
			})
			continue
		}

		libs = append(libs, Library{
			ImportPath:   imp.Path.Value,
			SelectorName: path.Base(strings.ReplaceAll(imp.Path.Value, "\"", "")),
			Alias:        "",
		})
	}

	return libs
}

func (b Bundler) isLibrary(value string) bool {
	return strings.Contains(value, b.libPackage)
}

// ライブラリのディレクトリを取得
// "golb/golb/testdata/lib/sample" -> "golb/testdata/lib/sample"
func (b Bundler) getDir(value string) string {
	return path.Join(b.rootDir, strings.TrimPrefix(strings.Trim(value, "\""), b.libPackage))
}

// ディレクトリ内のファイルを取得
func (b Bundler) getFiles(dir string) ([]os.DirEntry, error) {
	f, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// 使用されている関数を取得
func (b Bundler) getUsedFunction(file *ast.File) map[string]struct{} {
	usedFuncs := map[string]struct{}{}

	astutil.Apply(file, func(cursor *astutil.Cursor) bool {
		switch node := cursor.Node().(type) {
		case *ast.CallExpr:
			switch fun := node.Fun.(type) {
			case *ast.Ident:
				usedFuncs[fun.Name] = struct{}{}
			case *ast.SelectorExpr:
				usedFuncs[fun.Sel.Name] = struct{}{}
			}
		}
		return true
	}, nil)

	return usedFuncs
}

// ASTを書き換え
// targetに含まれるselectorを削除する
// vector.X -> X
func (b Bundler) removeSelector(file *ast.File, targets map[string]struct{}) {
	astutil.Apply(file, func(cursor *astutil.Cursor) bool {
		switch node := cursor.Node().(type) {
		case *ast.SelectorExpr:
			if _, ok := targets[node.X.(*ast.Ident).Name]; ok {
				cursor.Replace(node.Sel)
			}
		}
		return true
	}, nil)
}

// ASTを書き換え
// import宣言を削除する
func (b Bundler) removeAllImport(file *ast.File) {
	// import宣言を削除
	astutil.Apply(file, func(cursor *astutil.Cursor) bool {
		switch node := cursor.Node().(type) {
		case *ast.GenDecl:
			if node.Tok == token.IMPORT {
				cursor.Delete()
			}
		}
		return true
	}, nil)
}

// ASTを書き換え
// 使用していない関数宣言を削除する
func (b Bundler) removeUnusedFunction(file *ast.File, usedFuncs map[string]struct{}) {
	// 使用していない関数を削除
	astutil.Apply(file, func(cursor *astutil.Cursor) bool {
		switch node := cursor.Node().(type) {
		case *ast.FuncDecl:
			// main関数は削除しない
			if node.Name.Name == "main" {
				return true
			}

			// methodは削除しない
			// func (v Vector) Add(v2 Vector) Vector {}
			if node.Recv != nil {
				return true
			}

			if _, ok := usedFuncs[node.Name.Name]; !ok {
				cursor.Delete()
			}
		}
		return true
	}, nil)
}

// ノードを文字列として取得する
func (b Bundler) nodeToString(node *ast.File) string {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, token.NewFileSet(), node)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error printing node: %v\n", err)
		return ""
	}
	return buf.String()
}

func prettyPrint(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	err = json.Indent(&buf, data, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(buf.String())
	return nil
}
