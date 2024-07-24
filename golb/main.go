package golb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"strings"
)

type Library struct {
	ImportPath   string
	SelectorName string
	Alias        string
}

const libPackage = "github.com/hamao0820/ac-library-go"

var goModDir = path.Join("golb", "testdata")

func Bundle(src string) error {
	files := map[string]struct{}{}
	var libs []Library

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

		if file == src {
			libs = importLibs
		}

		for _, lib := range importLibs {
			libDir := getDir(lib.ImportPath)
			libFiles := getFiles(libDir)
			for _, file := range libFiles {
				libPath := path.Join(libDir, file.Name())
				dfs(libPath)
				files[libPath] = struct{}{}
			}
		}
	}

	dfs(src)
	for file := range files {
		fmt.Println(file)
	}
	fmt.Println(libs)

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
