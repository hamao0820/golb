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
	src        string
	libPackage string
	rootDir    string
}

func NewBundler(src, libPackage, goModDir string) Bundler {
	return Bundler{
		src:        src,
		libPackage: libPackage,
		rootDir:    goModDir,
	}
}

func (b Bundler) Bundle() (code string, err error) {
	files := map[string]string{} // key: file path, value: 展開コード

	// 再帰的にファイルを取得
	var dfs func(string)
	dfs = func(file string) {
		if _, ok := files[file]; ok {
			return
		}
		node, err := b.perseFile(file)
		if err != nil {
			return
		}

		importLibs := b.getImportedPackage(node)
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
		code := b.nodeToString(node)
		if file != b.src {
			code = strings.Join(strings.Split(code, "\n")[1:], "\n") // package行を削除
		}
		files[file] = code

		for _, lib := range importLibs {
			libDir := b.getDir(lib.ImportPath)
			libFiles, err := b.getFiles(libDir)
			if err != nil {
				panic(err)
			}
			for _, file := range libFiles {
				libPath := path.Join(libDir, file.Name())
				dfs(libPath)
			}
		}
	}

	dfs(b.src)

	sourceCode := files[b.src]
	sourceCode += "/*" + strings.Repeat("-", 50) + "以下は生成コード" + strings.Repeat("-", 50) + "*/"
	sourceCode += "\n\n"
	for file, code := range files {
		if file == b.src {
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
func (b Bundler) getImportedPackage(file *ast.File) []Library {
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
