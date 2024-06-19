package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

// ASTNode represents a node in the abstract syntax tree.
type ASTNode struct {
	Name     string      `json:"name,omitempty"`
	Type     string      `json:"type"`
	Children []*ASTNode  `json:"children,omitempty"`
	Value    interface{} `json:"value,omitempty"`
	Comments []string    `json:"comments,omitempty"`
}

// marshalAST converts an ast.Node into an ASTNode.
func marshalAST(node ast.Node, visited map[ast.Node]bool) *ASTNode {
	if node == nil {
		return nil
	}

	// Check if the node has been visited before to avoid cycles.
	if visited[node] {
		return nil
	}
	visited[node] = true

	astNode := &ASTNode{Type: fmt.Sprintf("%T", node)}

	// Handle different types of AST nodes.
	// Handle different types of AST nodes.
	switch n := node.(type) {
	case *ast.Ident:
		astNode.Value = n.Name
	case *ast.BasicLit:
		astNode.Value = n.Value
	case *ast.File:
		astNode.Value = n.Name.Name
	case *ast.Ellipsis:
		if n.Elt != nil {
			eltNode := marshalAST(n.Elt, visited)
			if eltNode != nil {
				astNode.Children = append(astNode.Children, eltNode)
			}
		}
	case *ast.GenDecl:
		for _, spec := range n.Specs {
			childNode := marshalAST(spec, visited)
			if childNode != nil {
				astNode.Children = append(astNode.Children, childNode)
			}
		}
	case *ast.FuncDecl:
		astNode.Name = n.Name.Name
		if n.Recv != nil {
			recvNode := marshalAST(n.Recv, visited)
			if recvNode != nil {
				astNode.Children = append(astNode.Children, recvNode)
			}
		}
		if n.Type != nil {
			typeNode := marshalAST(n.Type, visited)
			if typeNode != nil {
				astNode.Children = append(astNode.Children, typeNode)
			}
		}
		if n.Body != nil {
			bodyNode := marshalAST(n.Body, visited)
			if bodyNode != nil {
				astNode.Children = append(astNode.Children, bodyNode)
			}
		}
	case *ast.TypeSpec:
		astNode.Name = n.Name.Name
		typeNode := marshalAST(n.Type, visited)
		if typeNode != nil {
			astNode.Children = append(astNode.Children, typeNode)
		}
	case *ast.ValueSpec:
		for _, name := range n.Names {
			nameNode := marshalAST(name, visited)
			if nameNode != nil {
				astNode.Children = append(astNode.Children, nameNode)
			}
		}
		if n.Type != nil {
			typeNode := marshalAST(n.Type, visited)
			if typeNode != nil {
				astNode.Children = append(astNode.Children, typeNode)
			}
		}
		for _, value := range n.Values {
			valueNode := marshalAST(value, visited)
			if valueNode != nil {
				astNode.Children = append(astNode.Children, valueNode)
			}
		}
	case *ast.AssignStmt:
		for _, lhs := range n.Lhs {
			lhsNode := marshalAST(lhs, visited)
			if lhsNode != nil {
				astNode.Children = append(astNode.Children, lhsNode)
			}
		}
		for _, rhs := range n.Rhs {
			rhsNode := marshalAST(rhs, visited)
			if rhsNode != nil {
				astNode.Children = append(astNode.Children, rhsNode)
			}
		}
	case *ast.ReturnStmt:
		for _, result := range n.Results {
			resultNode := marshalAST(result, visited)
			if resultNode != nil {
				astNode.Children = append(astNode.Children, resultNode)
			}
		}
	case *ast.IfStmt:
		if n.Init != nil {
			initNode := marshalAST(n.Init, visited)
			if initNode != nil {
				astNode.Children = append(astNode.Children, initNode)
			}
		}
		if n.Cond != nil {
			condNode := marshalAST(n.Cond, visited)
			if condNode != nil {
				astNode.Children = append(astNode.Children, condNode)
			}
		}
		if n.Body != nil {
			bodyNode := marshalAST(n.Body, visited)
			if bodyNode != nil {
				astNode.Children = append(astNode.Children, bodyNode)
			}
		}
		if n.Else != nil {
			elseNode := marshalAST(n.Else, visited)
			if elseNode != nil {
				astNode.Children = append(astNode.Children, elseNode)
			}
		}
	case *ast.ForStmt:
		if n.Init != nil {
			initNode := marshalAST(n.Init, visited)
			if initNode != nil {
				astNode.Children = append(astNode.Children, initNode)
			}
		}
		if n.Cond != nil {
			condNode := marshalAST(n.Cond, visited)
			if condNode != nil {
				astNode.Children = append(astNode.Children, condNode)
			}
		}
		if n.Post != nil {
			postNode := marshalAST(n.Post, visited)
			if postNode != nil {
				astNode.Children = append(astNode.Children, postNode)
			}
		}
		if n.Body != nil {
			bodyNode := marshalAST(n.Body, visited)
			if bodyNode != nil {
				astNode.Children = append(astNode.Children, bodyNode)
			}
		}
	case *ast.RangeStmt:
		if n.Key != nil {
			keyNode := marshalAST(n.Key, visited)
			if keyNode != nil {
				astNode.Children = append(astNode.Children, keyNode)
			}
		}
		if n.Value != nil {
			valueNode := marshalAST(n.Value, visited)
			if valueNode != nil {
				astNode.Children = append(astNode.Children, valueNode)
			}
		}
		if n.X != nil {
			xNode := marshalAST(n.X, visited)
			if xNode != nil {
				astNode.Children = append(astNode.Children, xNode)
			}
		}
		if n.Body != nil {
			bodyNode := marshalAST(n.Body, visited)
			if bodyNode != nil {
				astNode.Children = append(astNode.Children, bodyNode)
			}
		}
	case *ast.BlockStmt:
		for _, stmt := range n.List {
			stmtNode := marshalAST(stmt, visited)
			if stmtNode != nil {
				astNode.Children = append(astNode.Children, stmtNode)
			}
		}
	case *ast.ExprStmt:
		if n.X != nil {
			xNode := marshalAST(n.X, visited)
			if xNode != nil {
				astNode.Children = append(astNode.Children, xNode)
			}
		}
	case *ast.CallExpr:
		if n.Fun != nil {
			funNode := marshalAST(n.Fun, visited)
			if funNode != nil {
				astNode.Children = append(astNode.Children, funNode)
			}
		}
		for _, arg := range n.Args {
			argNode := marshalAST(arg, visited)
			if argNode != nil {
				astNode.Children = append(astNode.Children, argNode)
			}
		}
	case *ast.SelectorExpr:
		if n.X != nil {
			xNode := marshalAST(n.X, visited)
			if xNode != nil {
				astNode.Children = append(astNode.Children, xNode)
			}
		}
		if n.Sel != nil {
			selNode := marshalAST(n.Sel, visited)
			if selNode != nil {
				astNode.Children = append(astNode.Children, selNode)
			}
		}

	case *ast.IndexListExpr:
		if n.X != nil {
			xNode := marshalAST(n.X, visited)
			if xNode != nil {
				astNode.Children = append(astNode.Children, xNode)
			}
		}
		for _, index := range n.Indices {
			indexNode := marshalAST(index, visited)
			if indexNode != nil {
				astNode.Children = append(astNode.Children, indexNode)
			}
		}
	case *ast.IndexExpr:
		if n.X != nil {
			xNode := marshalAST(n.X, visited)
			if xNode != nil {
				astNode.Children = append(astNode.Children, xNode)
			}
		}
		if n.Index != nil {
			indexNode := marshalAST(n.Index, visited)
			if indexNode != nil {
				astNode.Children = append(astNode.Children, indexNode)
			}
		}
	case *ast.SliceExpr:
		if n.X != nil {
			xNode := marshalAST(n.X, visited)
			if xNode != nil {
				astNode.Children = append(astNode.Children, xNode)
			}
		}
		if n.Low != nil {
			lowNode := marshalAST(n.Low, visited)
			if lowNode != nil {
				astNode.Children = append(astNode.Children, lowNode)
			}
		}
		if n.High != nil {
			highNode := marshalAST(n.High, visited)
			if highNode != nil {
				astNode.Children = append(astNode.Children, highNode)
			}
		}
		if n.Max != nil {
			maxNode := marshalAST(n.Max, visited)
			if maxNode != nil {
				astNode.Children = append(astNode.Children, maxNode)
			}
		}
	case *ast.StructType:
		if n.Fields != nil {
			fieldsNode := marshalAST(n.Fields, visited)
			if fieldsNode != nil {
				astNode.Children = append(astNode.Children, fieldsNode)
			}
		}
	case *ast.FuncType:
		if n.Params != nil {
			paramsNode := marshalAST(n.Params, visited)
			if paramsNode != nil {
				astNode.Children = append(astNode.Children, paramsNode)
			}
		}
		if n.Results != nil {
			resultsNode := marshalAST(n.Results, visited)
			if resultsNode != nil {
				astNode.Children = append(astNode.Children, resultsNode)
			}
		}
	case *ast.InterfaceType:
		if n.Methods != nil {
			methodsNode := marshalAST(n.Methods, visited)
			if methodsNode != nil {
				astNode.Children = append(astNode.Children, methodsNode)
			}
		}
	case *ast.ArrayType:
		if n.Elt != nil {
			eltNode := marshalAST(n.Elt, visited)
			if eltNode != nil {
				astNode.Children = append(astNode.Children, eltNode)
			}
		}

	case *ast.SelectStmt:
		if n.Body != nil {
			bodyNode := marshalAST(n.Body, visited)
			if bodyNode != nil {
				astNode.Children = append(astNode.Children, bodyNode)
			}
		}
	case *ast.CompositeLit:
		if n.Type != nil {
			typeNode := marshalAST(n.Type, visited)
			if typeNode != nil {
				astNode.Children = append(astNode.Children, typeNode)
			}
		}
		for _, elt := range n.Elts {
			eltNode := marshalAST(elt, visited)
			if eltNode != nil {
				astNode.Children = append(astNode.Children, eltNode)
			}
		}
	case *ast.ParenExpr:
		if n.X != nil {
			xNode := marshalAST(n.X, visited)
			if xNode != nil {
				astNode.Children = append(astNode.Children, xNode)
			}
		}
	case *ast.TypeAssertExpr:
		if n.X != nil {
			xNode := marshalAST(n.X, visited)
			if xNode != nil {
				astNode.Children = append(astNode.Children, xNode)
			}
		}
		if n.Type != nil {
			typeNode := marshalAST(n.Type, visited)
			if typeNode != nil {
				astNode.Children = append(astNode.Children, typeNode)
			}
		}

	case *ast.BadDecl:
		// No specific handling required for BadDecl
	case *ast.BadExpr:
		// No specific handling required for BadExpr
	case *ast.FuncLit:
		if n.Type != nil {
			typeNode := marshalAST(n.Type, visited)
			if typeNode != nil {
				astNode.Children = append(astNode.Children, typeNode)
			}
		}
		if n.Body != nil {
			bodyNode := marshalAST(n.Body, visited)
			if bodyNode != nil {
				astNode.Children = append(astNode.Children, bodyNode)
			}
		}
	case *ast.StarExpr:
		if n.X != nil {
			xNode := marshalAST(n.X, visited)
			if xNode != nil {
				astNode.Children = append(astNode.Children, xNode)
			}
		}
	case *ast.UnaryExpr:
		if n.X != nil {
			xNode := marshalAST(n.X, visited)
			if xNode != nil {
				astNode.Children = append(astNode.Children, xNode)
			}
		}
	case *ast.BinaryExpr:
		if n.X != nil {
			xNode := marshalAST(n.X, visited)
			if xNode != nil {
				astNode.Children = append(astNode.Children, xNode)
			}
		}
		if n.Y != nil {
			yNode := marshalAST(n.Y, visited)
			if yNode != nil {
				astNode.Children = append(astNode.Children, yNode)
			}
		}
	case *ast.KeyValueExpr:
		if n.Key != nil {
			keyNode := marshalAST(n.Key, visited)
			if keyNode != nil {
				astNode.Children = append(astNode.Children, keyNode)
			}
		}
		if n.Value != nil {
			valueNode := marshalAST(n.Value, visited)
			if valueNode != nil {
				astNode.Children = append(astNode.Children, valueNode)
			}
		}
	case *ast.BadStmt:
		// No specific handling required for BadStmt
	case *ast.DeclStmt:
		if n.Decl != nil {
			declNode := marshalAST(n.Decl, visited)
			if declNode != nil {
				astNode.Children = append(astNode.Children, declNode)
			}
		}
	case *ast.EmptyStmt:
		// No specific handling required for EmptyStmt
	case *ast.LabeledStmt:
		if n.Label != nil {
			labelNode := marshalAST(n.Label, visited)
			if labelNode != nil {
				astNode.Children = append(astNode.Children, labelNode)
			}
		}
		if n.Stmt != nil {
			stmtNode := marshalAST(n.Stmt, visited)
			if stmtNode != nil {
				astNode.Children = append(astNode.Children, stmtNode)
			}
		}
	case *ast.SendStmt:
		if n.Chan != nil {
			chanNode := marshalAST(n.Chan, visited)
			if chanNode != nil {
				astNode.Children = append(astNode.Children, chanNode)
			}
		}
		if n.Value != nil {
			valueNode := marshalAST(n.Value, visited)
			if valueNode != nil {
				astNode.Children = append(astNode.Children, valueNode)
			}
		}
	case *ast.IncDecStmt:
		if n.X != nil {
			xNode := marshalAST(n.X, visited)
			if xNode != nil {
				astNode.Children = append(astNode.Children, xNode)
			}
		}
	case *ast.GoStmt:
		if n.Call != nil {
			callNode := marshalAST(n.Call, visited)
			if callNode != nil {
				astNode.Children = append(astNode.Children, callNode)
			}
		}
	case *ast.DeferStmt:
		if n.Call != nil {
			callNode := marshalAST(n.Call, visited)
			if callNode != nil {
				astNode.Children = append(astNode.Children, callNode)
			}
		}
	case *ast.CaseClause:
		for _, expr := range n.List {
			exprNode := marshalAST(expr, visited)
			if exprNode != nil {
				astNode.Children = append(astNode.Children, exprNode)
			}
		}
		for _, stmt := range n.Body {
			stmtNode := marshalAST(stmt, visited)
			if stmtNode != nil {
				astNode.Children = append(astNode.Children, stmtNode)
			}
		}

	case *ast.CommentGroup:
		for _, comment := range n.List {
			commentNode := marshalAST(comment, visited)
			if commentNode != nil {
				astNode.Children = append(astNode.Children, commentNode)
			}
		}
	case *ast.Comment:
		astNode.Comments = append(astNode.Comments, n.Text)

	case *ast.TypeSwitchStmt:
		if n.Init != nil {
			initNode := marshalAST(n.Init, visited)
			if initNode != nil {
				astNode.Children = append(astNode.Children, initNode)
			}
		}
		if n.Assign != nil {
			assignNode := marshalAST(n.Assign, visited)
			if assignNode != nil {
				astNode.Children = append(astNode.Children, assignNode)
			}
		}
		if n.Body != nil {
			bodyNode := marshalAST(n.Body, visited)
			if bodyNode != nil {
				astNode.Children = append(astNode.Children, bodyNode)
			}
		}
	case *ast.CommClause:
		if n.Comm != nil {
			commNode := marshalAST(n.Comm, visited)
			if commNode != nil {
				astNode.Children = append(astNode.Children, commNode)
			}
		}
		for _, stmt := range n.Body {
			stmtNode := marshalAST(stmt, visited)
			if stmtNode != nil {
				astNode.Children = append(astNode.Children, stmtNode)
			}
		}
	case *ast.ImportSpec:
		if n.Name != nil {
			nameNode := marshalAST(n.Name, visited)
			if nameNode != nil {
				astNode.Children = append(astNode.Children, nameNode)
			}
		}
		if n.Path != nil {
			pathNode := marshalAST(n.Path, visited)
			if pathNode != nil {
				astNode.Children = append(astNode.Children, pathNode)
			}
		}
	// case *ast.Package:
	// 	if n.Name != nil {
	// 		nameNode := marshalAST(n.Name, visited)
	// 		if nameNode != nil {
	// 			astNode.Children = append(astNode.Children, nameNode)
	// 		}
	// 	}
	case *ast.Field:
		for _, name := range n.Names {
			nameNode := marshalAST(name, visited)
			if nameNode != nil {
				astNode.Children = append(astNode.Children, nameNode)
			}
		}
		if n.Type != nil {
			typeNode := marshalAST(n.Type, visited)
			if typeNode != nil {
				astNode.Children = append(astNode.Children, typeNode)
			}
		}
	case *ast.FieldList:
		for _, field := range n.List {
			fieldNode := marshalAST(field, visited)
			if fieldNode != nil {
				astNode.Children = append(astNode.Children, fieldNode)
			}
		}
	case *ast.MapType:
		if n.Key != nil {
			keyNode := marshalAST(n.Key, visited)
			if keyNode != nil {
				astNode.Children = append(astNode.Children, keyNode)
			}
		}
		if n.Value != nil {
			valueNode := marshalAST(n.Value, visited)
			if valueNode != nil {
				astNode.Children = append(astNode.Children, valueNode)
			}
		}
	case *ast.ChanType:
		if n.Value != nil {
			valueNode := marshalAST(n.Value, visited)
			if valueNode != nil {
				astNode.Children = append(astNode.Children, valueNode)
			}
		}
	case *ast.BranchStmt:
		if n.Label != nil {
			labelNode := marshalAST(n.Label, visited)
			if labelNode != nil {
				astNode.Children = append(astNode.Children, labelNode)
			}
		}
	case *ast.SwitchStmt:
		if n.Init != nil {
			initNode := marshalAST(n.Init, visited)
			if initNode != nil {
				astNode.Children = append(astNode.Children, initNode)
			}
		}
		if n.Tag != nil {
			tagNode := marshalAST(n.Tag, visited)
			if tagNode != nil {
				astNode.Children = append(astNode.Children, tagNode)
			}
		}
		if n.Body != nil {
			bodyNode := marshalAST(n.Body, visited)
			if bodyNode != nil {
				astNode.Children = append(astNode.Children, bodyNode)
			}
		}

	default:
		// Panic with an error message if an unexpected node type is encountered.
		panic(fmt.Sprintf("unsupported AST node type: %T", node))
	}

	// Traverse child nodes and add them to the current node's children.
	ast.Inspect(node, func(n ast.Node) bool {
		if n != nil {
			childNode := marshalAST(n, visited)
			if childNode != nil {
				astNode.Children = append(astNode.Children, childNode)
			}
		}
		return true
	})

	return astNode
}

// processFile processes a single Go source file and outputs its AST in JSON format.
func processFile(sourceFilePath string) error {
	// Parse the Go source file and generate the AST.
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, sourceFilePath, nil, parser.AllErrors)
	if err != nil {
		return fmt.Errorf("error parsing Go source file %s: %w", sourceFilePath, err)
	}

	// Generate the output file path with a .json extension.
	dir := filepath.Dir(sourceFilePath)
	base := filepath.Base(sourceFilePath)
	ext := filepath.Ext(base)
	baseNameWithoutExt := strings.TrimSuffix(base, ext)
	newBaseName := baseNameWithoutExt + ".json"
	newFilePath := filepath.Join(dir, newBaseName)

	// Create the output file for the JSON representation of the AST.
	outputFile, err := os.Create(newFilePath)
	if err != nil {
		return fmt.Errorf("error creating output file %s: %w", newFilePath, err)
	}
	defer outputFile.Close()

	// Serialize the AST to JSON and write it to the output file.
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonEncoder := json.NewEncoder(outputFile)
	jsonEncoder.SetIndent("", "  ")

	visited := make(map[ast.Node]bool)
	astNode := marshalAST(file, visited)
	err = jsonEncoder.Encode(astNode)
	if err != nil {
		return fmt.Errorf("error serializing AST to JSON for file %s: %w", sourceFilePath, err)
	}

	fmt.Println("AST generated and saved to " + newFilePath)
	return nil
}

// processFolder processes all .go files in the provided folder.
func processFolder(folderPath string) error {
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			return processFile(path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error processing folder %s: %w", folderPath, err)
	}
	return nil
}

func main() {
	// Ensure a Go source file or folder path is provided as a command-line argument.
	if len(os.Args) < 2 {
		fmt.Println("Please provide the path to the Go source file or folder as a command-line argument.")
		os.Exit(1)
	}

	path := os.Args[1]

	// Check if the path is a file or a folder.
	info, err := os.Stat(path)
	if err != nil {
		fmt.Printf("Error accessing the path: %s\n", err)
		os.Exit(1)
	}

	if info.IsDir() {
		// Process all .go files in the folder.
		err = processFolder(path)
		if err != nil {
			fmt.Printf("Error processing folder: %s\n", err)
			os.Exit(1)
		}
	} else {
		// Process the single file.
		err = processFile(path)
		if err != nil {
			fmt.Printf("Error processing file: %s\n", err)
			os.Exit(1)

		}
	}
}
