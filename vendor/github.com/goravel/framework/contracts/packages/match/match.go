package match

import (
	"github.com/dave/dst"
	"github.com/dave/dst/dstutil"
)

type GoNode interface {
	MatchNode(node dst.Node) bool
	MatchCursor(cursor *dstutil.Cursor) bool
}
