package golb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

const libDir = "golb/golb/testdata/lib"

func Bundle(src string) error {
	srcNode, err := perseFile(src)
	if err != nil {
		return err
	}
	findUsedElements(srcNode)

	return nil
}

func perseFile(filename string) (*ast.File, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func findUsedElements(file *ast.File) (nodes []ast.Node) {
	nodes = []ast.Node{}

	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.CallExpr:
			// 関数呼び出しを処理
			switch fun := x.Fun.(type) {
			case *ast.SelectorExpr:
				if ident, ok := fun.X.(*ast.Ident); ok {
					if ident.Name == "lib" {
						nodes = append(nodes, x)
					}
				}
			}
			// prettyPrint(x.Fun)
		case *ast.SelectorExpr:
			// 構造体メソッドを処理
		case *ast.Ident:
			// 定数や変数を処理
		case *ast.ImportSpec:
			// インポート宣言を処理
		case *ast.BinaryExpr:
			// 二項演算式を処理
		case *ast.UnaryExpr:
			// 単項演算式を処理
		case *ast.ParenExpr:
			// 括弧式を処理
		case *ast.BasicLit:
			// 基本リテラルを処理
		case *ast.AssignStmt:
			// 代入文を処理
		case *ast.FuncDecl:
			// 関数定義を処理
		case *ast.TypeSpec:
			// 型定義を処理
		case *ast.ValueSpec:
			// 定数や変数の定義を処理
		case *ast.StructType:
			// 構造体定義を処理
		case *ast.Field:
			// 構造体フィールドを処理
		case *ast.FuncType:
			// 関数型を処理
		case *ast.InterfaceType:
			// インターフェース型を処理
		case *ast.MapType:
			// マップ型を処理
		case *ast.ArrayType:
			// 配列型を処理
		case *ast.ChanType:
			// チャネル型を処理
		case *ast.StarExpr:
			// ポインタ型を処理
		case *ast.SliceExpr:
			// スライス型を処理
		case *ast.IndexExpr:
			// インデックス式を処理
		case *ast.KeyValueExpr:
			// キーと値のペアを処理
		case *ast.CompositeLit:
			// 複合リテラルを処理
		case *ast.ReturnStmt:
			// 戻り値を処理
		case *ast.BlockStmt:
			// ブロック文を処理
		case *ast.IfStmt:
			// if文を処理
		case *ast.ForStmt:
			// for文を処理
		case *ast.RangeStmt:
			// range文を処理
		case *ast.SwitchStmt:
			// switch文を処理
		case *ast.CaseClause:
			// case節を処理
		case *ast.TypeSwitchStmt:
			// 型switch文を処理
		case *ast.TypeAssertExpr:
			// 型アサーションを処理
		case *ast.SelectStmt:
			// select文を処理
		case *ast.CommClause:
			// 通信ケースを処理
		case *ast.SendStmt:
			// 送信文を処理
		case *ast.DeclStmt:
			// 宣言文を処理
		case *ast.GenDecl:
			// 宣言を処理
		case *ast.BadDecl:
			// 不正な宣言を処理
		case *ast.FuncLit:
			// 関数リテラルを処理
		case *ast.Comment:
			// コメントを処理
		case *ast.CommentGroup:
			// コメントグループを処理
		case *ast.Ellipsis:
			// 可変長引数を処理
		case *ast.FieldList:
			// フィールドリストを処理
		case *ast.File:
			// ファイルを処理
		case *ast.Package:
			// パッケージを処理
		case *ast.BadExpr:
			// 不正な式を処理
		}
		return true
	})

	return
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

// func Bundle() {
// 	sourceFile := "golb/testdata/src/a/main.go" // 提出用ファイルのパス
// 	libDir := "golb/testdata/lib"               // ライブラリディレクトリのパス

// 	// Goファイルを解析
// 	fs := token.NewFileSet()
// 	node, err := parser.ParseFile(fs, sourceFile, nil, parser.ParseComments)
// 	if err != nil {
// 		log.Fatalf("failed to parse file: %v", err)
// 	}

// 	// 提出ファイルに含めるべき識別子を収集
// 	identifiers := collectIdentifiers(node)

// 	// ライブラリディレクトリを走査
// 	libFiles, err := os.ReadDir(libDir)
// 	if err != nil {
// 		log.Fatalf("failed to read library directory: %v", err)
// 	}

// 	// 必要な識別子をライブラリから収集し、提出ファイルに追加
// 	var combinedCode strings.Builder
// 	combinedCode.WriteString(extractCode(node, identifiers)) // 元の提出ファイルのコードを追加

// 	for _, file := range libFiles {
// 		if strings.HasSuffix(file.Name(), ".go") {
// 			libNode, err := parser.ParseFile(fs, libDir+"/"+file.Name(), nil, 0)
// 			if err != nil {
// 				log.Printf("failed to parse library file: %v", err)
// 				continue
// 			}
// 			combinedCode.WriteString(extractCode(libNode, identifiers)) // ライブラリのコードを追加
// 		}
// 	}

// 	// 結果のコードを出力
// 	log.Println(combinedCode.String())
// }

// // 指定されたファイルノードから識別子を収集
// func collectIdentifiers(node *ast.File) map[string]bool {
// 	identifiers := make(map[string]bool)

// 	// 識別子を収集
// 	ast.Inspect(node, func(n ast.Node) bool {
// 		if ident, ok := n.(*ast.Ident); ok {
// 			identifiers[ident.Name] = true
// 		}
// 		return true
// 	})

// 	return identifiers
// }

// // 必要な識別子を含むコードを抽出
// func extractCode(node *ast.File, identifiers map[string]bool) string {
// 	var code strings.Builder

// 	// 構造体、関数、定数などを抽出
// 	for _, decl := range node.Decls {
// 		switch decl := decl.(type) {
// 		case *ast.GenDecl:
// 			for _, spec := range decl.Specs {
// 				switch spec := spec.(type) {
// 				case *ast.TypeSpec:
// 					if identifiers[spec.Name.Name] {
// 						fmt.Println(spec.Name.Name)
// 						code.WriteString(node.Name.Name + "\n")
// 					}
// 				case *ast.ValueSpec:
// 					for _, name := range spec.Names {
// 						if identifiers[name.Name] {
// 							fmt.Println(name.Name)
// 							code.WriteString(node.Name.Name + "\n")
// 						}
// 					}
// 				}
// 			}
// 		case *ast.FuncDecl:
// 			if identifiers[decl.Name.Name] {
// 				fmt.Println(decl.Name.Name)
// 				code.WriteString(node.Name.Name + "\n")
// 			}
// 		}
// 	}

// 	fmt.Println(code.String())
// 	return code.String()
// }
