package match

import (
	"reflect"
	"strconv"

	"github.com/dave/dst"
	"github.com/dave/dst/dstutil"

	"github.com/goravel/framework/contracts/packages/match"
)

type (
	// GoNode represents a matcher for Go AST nodes with optional position constraints
	GoNode struct {
		match       func(node dst.Node) bool
		first, last bool
	}
	// GoNodes is a collection of GoNode matchers
	GoNodes []match.GoNode
)

// MatchCursor checks if the cursor's current node matches this GoNode matcher.
// If first or last flags are set, it also verifies the node's position in its parent slice.
// For example, FirstOf(Ident("x")).MatchCursor(cursor) returns true only if cursor points to the first identifier "x" in its parent slice.
func (r GoNode) MatchCursor(cursor *dstutil.Cursor) bool {
	if r.first || r.last {
		if r.MatchNode(cursor.Node()) {
			pr := reflect.Indirect(reflect.ValueOf(cursor.Parent())).FieldByName(cursor.Name())
			if pr.Kind() == reflect.Slice || pr.Kind() == reflect.Array {
				if r.first {
					return cursor.Index() == 0
				}

				if r.last {
					return cursor.Index() == pr.Len()-1
				}
			}
		}

		return false
	}

	return r.MatchNode(cursor.Node())
}

// MatchNode checks if the given node matches this GoNode matcher.
// For example, Ident("x").MatchNode(node) returns true if node is an identifier with the name "x".
func (r GoNode) MatchNode(node dst.Node) bool {
	return r.match(node)
}

// MatchNodes checks if all nodes in the slice match their corresponding matchers.
// Returns true if the GoNodes collection is empty or all nodes match.
// For example, GoNodes{Ident("x"), Ident("y")}.MatchNodes(nodes) returns true if nodes contains exactly two identifiers "x" and "y" in that order.
func (r GoNodes) MatchNodes(nodes []dst.Node) bool {
	if len(r) == 0 {
		return true
	}

	if len(nodes) != len(r) {
		return false
	}

	for i := range nodes {
		if len(r) > i {
			if !r[i].MatchNode(nodes[i]) {
				return false
			}
		}
	}

	return true
}

// AnyOf creates a matcher that succeeds if any of the provided matchers match.
// For example, AnyOf(Ident("foo"), Ident("bar")) matches either an identifier named "foo" or "bar".
func AnyOf(matchers ...match.GoNode) match.GoNode {
	return GoNode{
		match: func(node dst.Node) bool {
			for _, matcher := range matchers {
				if matcher.MatchNode(node) {
					return true
				}
			}

			return false
		},
	}
}

// AnyNode creates a matcher that matches any node.
// For example, CallExpr(Ident("print"), GoNodes{AnyNode()}) matches print(x) where x can be any expression.
func AnyNode() match.GoNode {
	return &GoNode{
		match: func(node dst.Node) bool {
			return true
		},
	}
}

// AnyNodes creates an empty GoNodes collection that matches any sequence of nodes.
// For example, CallExpr(Ident("print"), AnyNodes()) matches print() with any number of arguments.
func AnyNodes() GoNodes {
	return GoNodes{}
}

// ArrayType creates a matcher for array type expressions with element and length matchers.
// For example, ArrayType(Ident("int"), BasicLit("5")) matches [5]int.
func ArrayType(elt, l match.GoNode) match.GoNode {
	return GoNode{
		match: func(n dst.Node) bool {
			if e, ok := n.(*dst.ArrayType); ok {
				return elt.MatchNode(e.Elt) && l.MatchNode(e.Len)
			}

			return false
		},
	}
}

// BasicLit creates a matcher for basic literal expressions with a specific value.
// For example, BasicLit("\"hello\"") matches the string literal "hello".
func BasicLit(value string) match.GoNode {
	return GoNode{
		match: func(n dst.Node) bool {
			if e, ok := n.(*dst.BasicLit); ok {
				return e.Value == value
			}

			return false
		},
	}
}

// CallExpr creates a matcher for function call expressions with specific function and argument matchers.
// For example, CallExpr(Ident("fmt.Println"), GoNodes{BasicLit("\"test\"")}) matches fmt.Println("test").
func CallExpr(fun match.GoNode, args GoNodes) match.GoNode {
	return GoNode{
		match: func(n dst.Node) bool {
			if e, ok := n.(*dst.CallExpr); ok {
				var nodes = make([]dst.Node, len(e.Args))
				for i := range e.Args {
					nodes[i] = e.Args[i]
				}

				return fun.MatchNode(e.Fun) && args.MatchNodes(nodes)
			}

			return false
		},
	}
}

// CompositeLit creates a matcher for composite literal expressions with a type matcher.
// For example, CompositeLit(Ident("Person")) matches Person{} or Person{Name: "John"}.
func CompositeLit(t match.GoNode) match.GoNode {
	return GoNode{
		match: func(n dst.Node) bool {
			if e, ok := n.(*dst.CompositeLit); ok {
				return t.MatchNode(e.Type)
			}

			return false
		},
	}
}

// EqualNode creates a matcher that checks if a node is equal to the provided node.
// For example, if target is &dst.Ident{Name: "x"}, then EqualNode(target) matches only identifier nodes with the exact name "x".
func EqualNode(n dst.Node) match.GoNode {
	return GoNode{
		match: func(node dst.Node) bool {
			return dstNodeEq(n, node)
		},
	}
}

// FirstOf creates a matcher that only matches the first element in a parent slice.
// For example, FirstOf(TypeOf(&dst.ImportSpec{})) matches only the first import statement in the import declarations.
func FirstOf(n match.GoNode) match.GoNode {
	return GoNode{
		first: true,
		match: n.MatchNode,
	}
}

// Func creates a matcher for function declarations with a specific name matcher.
// For example, Func(Ident("main")) matches func main() { ... }.
func Func(name match.GoNode) match.GoNode {
	return GoNode{
		match: func(n dst.Node) bool {
			if e, ok := n.(*dst.FuncDecl); ok {
				return name.MatchNode(e.Name)
			}

			return false
		},
	}
}

// Ident creates a matcher for identifier nodes with a specific name.
// For example, Ident("x") matches the identifier x in expressions like x = 1 or return x.
func Ident(name string) match.GoNode {
	return GoNode{
		match: func(n dst.Node) bool {
			if ident, ok := n.(*dst.Ident); ok {
				return ident.Name == name
			}

			return false
		},
	}
}

// Import creates a matcher for import specifications with a specific path and optional name.
// For example, Import("fmt") matches import "fmt", and Import("fmt", "f") matches import f "fmt".
func Import(path string, name ...string) match.GoNode {
	return GoNode{
		match: func(n dst.Node) bool {
			if im, ok := n.(*dst.ImportSpec); ok {
				if im.Path.Value == strconv.Quote(path) {
					if len(name) > 0 {
						if im.Name != nil {
							return im.Name.Name == name[0]
						}
					}
					return true
				}
			}
			return false
		},
	}
}

// KeyValueExpr creates a matcher for key-value expressions with specific key and value matchers.
// For example, KeyValueExpr(Ident("Name"), BasicLit("\"John\"")) matches Name: "John" in struct literals.
func KeyValueExpr(key, value match.GoNode) match.GoNode {
	return GoNode{
		match: func(n dst.Node) bool {
			if e, ok := n.(*dst.KeyValueExpr); ok {
				return key.MatchNode(e.Key) && value.MatchNode(e.Value)
			}

			return false
		},
	}
}

// LastOf creates a matcher that only matches the last element in a parent slice.
// For example, LastOf(Ident("return")) matches return only if it's the last statement in a block.
func LastOf(n match.GoNode) match.GoNode {
	return GoNode{
		last:  true,
		match: n.MatchNode,
	}
}

// SelectorExpr creates a matcher for selector expressions (e.g., x.sel) with specific matchers.
// For example, SelectorExpr(Ident("fmt"), Ident("Println")) matches fmt.Println.
func SelectorExpr(x, sel match.GoNode) match.GoNode {
	return GoNode{
		match: func(n dst.Node) bool {
			if e, ok := n.(*dst.SelectorExpr); ok {
				return x.MatchNode(e.X) && sel.MatchNode(e.Sel)
			}

			return false
		},
	}
}

// TypeOf creates a matcher that checks if a node is of a specific type T.
// For example, TypeOf(&dst.CallExpr{}) matches any function call like foo(), bar(x), or fmt.Println("hello").
func TypeOf[T any](_ T) match.GoNode {
	return GoNode{
		match: func(node dst.Node) bool {
			_, ok := node.(T)
			return ok
		},
	}
}

// dstNodeEq compares two dst.Node instances for equality.
func dstNodeEq(x, y dst.Node) bool {
	switch x := x.(type) {
	case dst.Expr:
		y, ok := y.(dst.Expr)
		return ok && dstExprEq(x, y)
	case *dst.ImportSpec:
		y, ok := y.(*dst.ImportSpec)
		return ok && dstImportSpecEq(x, y)
	case *dst.ExprStmt:
		y, ok := y.(*dst.ExprStmt)
		return ok && dstExprStmtEq(x, y)
	default:
		panic("unhandled node type, please add it to dstNodeEq")
	}
}

// dstExprEq compares two dst.Expr instances for equality.
func dstExprEq(x, y dst.Expr) bool {
	if x == nil || y == nil {
		return x == y
	}

	switch x := x.(type) {
	case *dst.ArrayType:
		y, ok := y.(*dst.ArrayType)
		return ok && dstArrayTypeEq(x, y)
	case *dst.BasicLit:
		y, ok := y.(*dst.BasicLit)
		return ok && dstBasicLitEq(x, y)
	case *dst.CompositeLit:
		y, ok := y.(*dst.CompositeLit)
		return ok && dstCompositeLitEq(x, y)
	case *dst.Ident:
		y, ok := y.(*dst.Ident)
		return ok && dstIdentEq(x, y)
	case *dst.KeyValueExpr:
		y, ok := y.(*dst.KeyValueExpr)
		return ok && dstKeyValueExprEq(x, y)
	case *dst.MapType:
		y, ok := y.(*dst.MapType)
		return ok && dstMapTypeEq(x, y)
	case *dst.SelectorExpr:
		y, ok := y.(*dst.SelectorExpr)
		return ok && dstSelectorExprEq(x, y)
	case *dst.UnaryExpr:
		y, ok := y.(*dst.UnaryExpr)
		return ok && dstUnaryExprEq(x, y)
	case *dst.CallExpr:
		y, ok := y.(*dst.CallExpr)
		return ok && dstCallExprEq(x, y)
	case *dst.FuncType:
		y, ok := y.(*dst.FuncType)
		return ok && dstFuncTypeEq(x, y)
	default:
		panic("unhandled node type, please add it to dstExprEq")
	}
}

// dstArrayTypeEq compares two dst.ArrayType instances for equality.
func dstArrayTypeEq(x, y *dst.ArrayType) bool {
	if x == nil || y == nil {
		return x == y
	}

	return dstExprEq(x.Elt, y.Elt) && dstExprEq(x.Len, y.Len)
}

// dstBasicLitEq compares two dst.BasicLit instances for equality.
func dstBasicLitEq(x, y *dst.BasicLit) bool {
	if x == nil || y == nil {
		return x == y
	}

	return x.Kind == y.Kind && x.Value == y.Value
}

// dstCompositeLitEq compares two dst.CompositeLit instances for equality.
func dstCompositeLitEq(x, y *dst.CompositeLit) bool {
	if x == nil || y == nil {
		return x == y
	}

	return dstExprEq(x.Type, y.Type) && dstExprSliceEq(x.Elts, y.Elts)
}

// dstCallExprEq compares two dst.CallExpr instances for equality.
// For example, dstCallExprEq(&dst.CallExpr{Fun: &dst.Ident{Name: "print"}}, &dst.CallExpr{Fun: &dst.Ident{Name: "print"}}) returns true.
func dstCallExprEq(x, y *dst.CallExpr) bool {
	if x == nil || y == nil {
		return x == y
	}

	return dstExprEq(x.Fun, y.Fun) && dstExprSliceEq(x.Args, y.Args) && x.Ellipsis == y.Ellipsis
}

// dstExprSliceEq compares two slices of dst.Expr for equality.
func dstExprSliceEq(xs, ys []dst.Expr) bool {
	if len(xs) != len(ys) {
		return false
	}

	for i := range xs {
		if !dstExprEq(xs[i], ys[i]) {
			return false
		}
	}

	return true
}

// dstIdentEq compares two dst.Ident instances for equality.
func dstIdentEq(x, y *dst.Ident) bool {
	if x == nil || y == nil {
		return x == y
	}

	return x.Name == y.Name
}

// dstImportSpecEq compares two dst.ImportSpec instances for equality.
func dstImportSpecEq(x, y *dst.ImportSpec) bool {
	if x == nil || y == nil {
		return x == y
	}

	return x.Path.Value == y.Path.Value && dstIdentEq(x.Name, y.Name)
}

// dstKeyValueExprEq compares two dst.KeyValueExpr instances for equality.
func dstKeyValueExprEq(x, y *dst.KeyValueExpr) bool {
	if x == nil || y == nil {
		return x == y
	}

	return dstExprEq(x.Key, y.Key) && dstExprEq(x.Value, y.Value)
}

// dstMapTypeEq compares two dst.MapType instances for equality.
func dstMapTypeEq(x, y *dst.MapType) bool {
	if x == nil || y == nil {
		return x == y
	}

	return dstExprEq(x.Key, y.Key) && dstExprEq(x.Value, y.Value)
}

// dstSelectorExprEq compares two dst.SelectorExpr instances for equality.
func dstSelectorExprEq(x, y *dst.SelectorExpr) bool {
	if x == nil || y == nil {
		return x == y
	}

	return dstExprEq(x.X, y.X) && dstIdentEq(x.Sel, y.Sel)
}

// dstExprStmtEq compares two dst.ExprStmt instances for equality.
func dstExprStmtEq(x, y *dst.ExprStmt) bool {
	if x == nil || y == nil {
		return x == y
	}

	return dstExprEq(x.X, y.X)
}

// dstFuncTypeEq compares two dst.FuncType instances for equality.
func dstFuncTypeEq(x, y *dst.FuncType) bool {
	if x == nil || y == nil {
		return x == y
	}

	return dstFieldListEq(x.Params, y.Params) && dstFieldListEq(x.Results, y.Results)
}

// dstFieldListEq compares two dst.FieldList instances for equality.
func dstFieldListEq(x, y *dst.FieldList) bool {
	if x == nil || y == nil {
		return x == y
	}

	if len(x.List) != len(y.List) {
		return false
	}

	for i := range x.List {
		if !dstFieldEq(x.List[i], y.List[i]) {
			return false
		}
	}

	return true
}

// dstFieldEq compares two dst.Field instances for equality.
func dstFieldEq(x, y *dst.Field) bool {
	if x == nil || y == nil {
		return x == y
	}

	// Compare names
	if len(x.Names) != len(y.Names) {
		return false
	}
	for i := range x.Names {
		if !dstIdentEq(x.Names[i], y.Names[i]) {
			return false
		}
	}

	// Compare type
	return dstExprEq(x.Type, y.Type)
}

// dstUnaryExprEq compares two dst.UnaryExpr instances for equality.
func dstUnaryExprEq(x, y *dst.UnaryExpr) bool {
	if x == nil || y == nil {
		return x == y
	}

	return x.Op == y.Op && dstExprEq(x.X, y.X)
}
