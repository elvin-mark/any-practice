package astcheck

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// Result contains the outcome of an AST validation.
type Result struct {
	Valid      bool
	ParseErr   error
	Violations []string
}

// CheckFile parses a Go file and inspects it against a list of forbidden identifiers/packages.
// forbidden rules can be:
// - package names: "sort", "slices", "container/heap"
// - qualified names: "sort.Slice", "slices.Sort"
func CheckFile(filePath string, forbidden []string) Result {
	if len(forbidden) == 0 {
		return Result{Valid: true}
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return Result{
			Valid:    false,
			ParseErr: err,
		}
	}

	var violations []string

	// Track imported packages and their aliases
	// alias/pkgName -> full import path
	pkgAliases := make(map[string]string)

	for _, imp := range node.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)
		var localName string
		if imp.Name != nil {
			localName = imp.Name.Name
		} else {
			parts := strings.Split(importPath, "/")
			localName = parts[len(parts)-1]
		}
		pkgAliases[localName] = importPath

		// Check if the import path itself is forbidden
		for _, f := range forbidden {
			if importPath == f || strings.HasSuffix(importPath, "/"+f) {
				pos := fset.Position(imp.Pos())
				violations = append(violations, fmt.Sprintf("line %d: import of forbidden package %q", pos.Line, importPath))
			}
		}
	}

	// Inspect AST nodes for selector expressions (pkg.Function)
	ast.Inspect(node, func(n ast.Node) bool {
		sel, ok := n.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		ident, ok := sel.X.(*ast.Ident)
		if !ok {
			return true
		}

		callee := ident.Name
		member := sel.Sel.Name
		qualified := fmt.Sprintf("%s.%s", callee, member)

		// Check if callee is an alias of a forbidden package or matches directly
		if fullPath, ok := pkgAliases[callee]; ok {
			for _, f := range forbidden {
				if fullPath == f || strings.HasSuffix(fullPath, "/"+f) {
					pos := fset.Position(sel.Pos())
					violations = append(violations, fmt.Sprintf("line %d: call to forbidden package function %s", pos.Line, qualified))
				}
			}
		}

		// Also check qualified name directly against forbidden list
		for _, f := range forbidden {
			if qualified == f || callee == f {
				pos := fset.Position(sel.Pos())
				violations = append(violations, fmt.Sprintf("line %d: use of forbidden symbol %q", pos.Line, qualified))
			}
		}

		return true
	})

	return Result{
		Valid:      len(violations) == 0,
		Violations: violations,
	}
}
