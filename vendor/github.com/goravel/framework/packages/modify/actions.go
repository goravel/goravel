package modify

import (
	"go/token"
	"slices"
	"strconv"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/dstutil"

	"github.com/goravel/framework/contracts/packages/modify"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/packages/match"
	"github.com/goravel/framework/support/color"
)

// Add adds an expression to the matched specified function.
func Add(expression string) modify.Action {
	return func(cursor *dstutil.Cursor) {
		expr := MustParseExpr(expression).(dst.Expr)
		stmt := &dst.ExprStmt{
			X: expr,
		}

		node := cursor.Node().(*dst.FuncDecl)
		node.Body.List = append(node.Body.List, stmt)
	}
}

// AddConfig adds a configuration key with the given expression to the config file.
func AddConfig(name, expression string, annotations ...string) modify.Action {
	return func(cursor *dstutil.Cursor) {
		var value *dst.CompositeLit
		switch node := cursor.Node().(type) {
		case *dst.KeyValueExpr:
			value = node.Value.(*dst.CompositeLit)
		case *dst.CallExpr:
			value = node.Args[1].(*dst.CompositeLit)
		}

		key := WrapNewline(&dst.BasicLit{Kind: token.STRING, Value: strconv.Quote(name)})
		newExpr := WrapNewline(&dst.KeyValueExpr{
			Key:   key,
			Value: WrapNewline(MustParseExpr(expression)).(dst.Expr),
		})

		// Add annotations as comments if provided
		if len(annotations) > 0 {
			var comments dst.Decorations
			for _, annotation := range annotations {
				// Ensure the annotation starts with "//" for proper comment formatting
				if !strings.HasPrefix(annotation, "//") {
					annotation = "// " + annotation
				}
				comments = append(comments, annotation)
			}
			newExpr.Decs.Start = comments
		}

		existExprIndex := KeyIndex(value.Elts, key)

		if existExprIndex >= 0 {
			value.Elts[existExprIndex] = newExpr
		} else {
			// add config
			value.Elts = append(value.Elts, newExpr)
		}
	}
}

// AddImport adds an import statement to the file.
func AddImport(path string, name ...string) modify.Action {
	return func(cursor *dstutil.Cursor) {
		node := cursor.Node().(*dst.GenDecl)

		// Check if import already exists
		for _, spec := range node.Specs {
			if importSpec, ok := spec.(*dst.ImportSpec); ok {
				if importSpec.Path.Value == strconv.Quote(path) {
					// Import already exists, no need to add
					return
				}
			}
		}

		// import spec
		im := &dst.ImportSpec{
			Path: &dst.BasicLit{
				Kind:  token.STRING,
				Value: strconv.Quote(path),
			},
		}
		if len(name) > 0 {
			im.Name = &dst.Ident{
				Name: name[0],
			}
		}

		// Insert third-party imports at the top and others at the bottom.
		// When formatting the source code, this helps group and sort imports
		// into stdlib, third-party, and local packages.
		if isThirdParty(path) {
			node.Specs = append([]dst.Spec{WrapNewline(im)}, node.Specs...)
			return
		}
		node.Specs = append(node.Specs, WrapNewline(im))
	}
}

// Register adds a registration to the matched specified array.
func Register(expression string, before ...string) modify.Action {
	return func(cursor *dstutil.Cursor) {
		expr := MustParseExpr(expression).(dst.Expr)
		node := cursor.Node().(*dst.CompositeLit)
		if ExprExists(node.Elts, expr) {
			color.Warningln(errors.PackageRegistrationDuplicate.Args(expression))
			return
		}
		if len(before) > 0 {
			// check if before is "*" and insert registration at the beginning
			if before[0] == "*" {
				node.Elts = slices.Insert(node.Elts, 0, expr)
				return
			}

			// check if beforeExpr is existing and insert registration before it
			beforeExpr := MustParseExpr(before[0]).(dst.Expr)
			if i := ExprIndex(node.Elts, beforeExpr); i >= 0 {
				node.Elts = slices.Insert(node.Elts, i, expr)
				return
			}

			color.Warningln(errors.PackageRegistrationNotFound.Args(before[0]))
		}

		node.Elts = append(node.Elts, expr)
	}
}

// Remove removes an expression from the matched specified function.
func Remove(expression string) modify.Action {
	return func(cursor *dstutil.Cursor) {
		expr := MustParseExpr(expression).(dst.Expr)
		stmt := &dst.ExprStmt{
			X: expr,
		}
		node := cursor.Node().(*dst.FuncDecl)
		node.Body.List = slices.DeleteFunc(node.Body.List, func(ex dst.Stmt) bool {
			return match.EqualNode(stmt).MatchNode(ex)
		})
	}
}

// RemoveConfig removes a configuration key from the config file.
func RemoveConfig(name string) modify.Action {
	return func(cursor *dstutil.Cursor) {
		var value *dst.CompositeLit
		switch node := cursor.Node().(type) {
		case *dst.KeyValueExpr:
			value = node.Value.(*dst.CompositeLit)
		case *dst.CallExpr:
			value = node.Args[1].(*dst.CompositeLit)
		}
		key := WrapNewline(&dst.BasicLit{Kind: token.STRING, Value: strconv.Quote(name)})

		// remove config
		value.Elts = slices.DeleteFunc(value.Elts, func(expr dst.Expr) bool {
			if kv, ok := expr.(*dst.KeyValueExpr); ok {
				return match.EqualNode(key).MatchNode(kv.Key)
			}
			return false
		})
	}
}

// RemoveImport removes an import statement from the file.
func RemoveImport(path string, name ...string) modify.Action {
	return func(cursor *dstutil.Cursor) {
		if IsUsingImport(cursor.Parent().(*dst.File), path, name...) {
			return
		}

		node := cursor.Node().(*dst.GenDecl)
		node.Specs = slices.DeleteFunc(node.Specs, func(spec dst.Spec) bool {
			return match.Import(path, name...).MatchNode(spec)
		})
	}
}

// ReplaceConfig replaces a configuration key with the given expression in the config file.
func ReplaceConfig(name, expression string) modify.Action {
	return func(cursor *dstutil.Cursor) {
		var value *dst.CompositeLit
		switch node := cursor.Node().(type) {
		case *dst.KeyValueExpr:
			value = node.Value.(*dst.CompositeLit)
		case *dst.CallExpr:
			value = node.Args[1].(*dst.CompositeLit)
		}
		key := WrapNewline(&dst.BasicLit{Kind: token.STRING, Value: strconv.Quote(name)})

		// replace config
		if i := KeyIndex(value.Elts, key); i >= 0 {
			value.Elts[i] = WrapNewline(&dst.KeyValueExpr{
				Key:   key,
				Value: WrapNewline(MustParseExpr(expression)).(dst.Expr),
			})
			return
		}
	}
}

// Unregister remove a registration from the matched specified array.
func Unregister(expression string) modify.Action {
	return func(cursor *dstutil.Cursor) {
		expr := MustParseExpr(expression).(dst.Expr)
		node := cursor.Node().(*dst.CompositeLit)
		node.Elts = slices.DeleteFunc(node.Elts, func(ex dst.Expr) bool {
			return match.EqualNode(expr).MatchNode(ex)
		})
	}
}
