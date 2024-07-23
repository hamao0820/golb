package golb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"strings"
)

type Library struct {
	ImportPath   string
	SelectorName string
	Alias        string
}

const libDir = "golb/golb/testdata/lib"

func Bundle(src string) error {
	srcNode, err := perseFile(src)
	if err != nil {
		return err
	}
	fmt.Println(getImportedPackage(srcNode))

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
			SelectorName: path.Base(imp.Path.Value),
			Alias:        "",
		})
	}

	return libs
}

func isLibrary(value string) bool {
	return strings.Contains(value, libDir)
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
