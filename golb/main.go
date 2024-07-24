package golb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

type Library struct {
	ImportPath   string
	SelectorName string
	Alias        string
}

const libPackage = "github.com/hamao0820/ac-library-go"

var goModDir = path.Join("golb", "testdata")

func Bundle(src string) error {
	files := map[string]string{} // key: file path, value: 展開コード
	// var libs []Library

	// 再帰的にファイルを取得
	var dfs func(string)
	dfs = func(file string) {
		if _, ok := files[file]; ok {
			return
		}
		node, err := perseFile(file)
		if err != nil {
			return
		}

		importLibs := getImportedPackage(node)
		targetSelectors := map[string]struct{}{}
		targetImports := map[string]struct{}{}
		for _, lib := range importLibs {
			targetSelectors[lib.SelectorName] = struct{}{}
			targetImports[lib.ImportPath] = struct{}{}
		}

		removeSelector(node, targetSelectors)
		removeAllImport(node)
		code := nodeToString(node)
		if file != src {
			code = strings.Join(strings.Split(code, "\n")[1:], "\n") // package行を削除
		}
		files[file] = code

		for _, lib := range importLibs {
			libDir := getDir(lib.ImportPath)
			libFiles := getFiles(libDir)
			for _, file := range libFiles {
				libPath := path.Join(libDir, file.Name())
				dfs(libPath)
			}
		}
	}

	dfs(src)

	sourceCode := files[src]
	sourceCode += "//========================================"
	sourceCode += "\n\n"
	for file, code := range files {
		if file == src {
			continue
		}
		dirs := strings.Split(file, "/")
		sourceCode += "// " + path.Join(dirs[len(dirs)-2:]...) + "\n"
		sourceCode += code
		sourceCode += "\n\n"
	}

	fmt.Println(sourceCode)

	return nil
}

// ASTを取得
func perseFile(filename string) (*ast.File, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	return node, nil
}

// importされている自作ライブラリのパッケージを取得
func getImportedPackage(file *ast.File) []Library {
	libs := []Library{}

	for _, imp := range file.Imports {
		if !isLibrary(imp.Path.Value) {
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

func isLibrary(value string) bool {
	return strings.Contains(value, libPackage)
}

// ライブラリのディレクトリを取得
// "golb/golb/testdata/lib/sample" -> "golb/testdata/lib/sample"
func getDir(value string) string {
	return path.Join(goModDir, strings.TrimPrefix(strings.Trim(value, "\""), libPackage))
}

// ディレクトリ内のファイルを取得
func getFiles(dir string) []os.DirEntry {
	f, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	return f
}

// ASTを書き換え
// targetに含まれるselectorを削除する
// vector.X -> X
func removeSelector(file *ast.File, targets map[string]struct{}) {
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
// targetに含まれるimportを削除する
func removeImport(file *ast.File, targets map[string]struct{}) {
	astutil.Apply(file, func(cursor *astutil.Cursor) bool {
		switch node := cursor.Node().(type) {
		case *ast.ImportSpec:
			if _, ok := targets[node.Path.Value]; ok {
				cursor.Delete()
			}
		}
		return true
	}, nil)
}

// ASTを書き換え
// import宣言を削除する
func removeAllImport(file *ast.File) {
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
func nodeToString(node ast.Node) string {
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
