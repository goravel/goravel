package db

import (
	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Expr struct {
	clause.Expr
}

func Raw(expr string, args ...any) any {
	return Expr{gorm.Expr(expr, args...)}
}

func (r Expr) ToSql() (sql string, args []any, err error) {
	return squirrel.Expr(r.SQL, r.Vars...).ToSql()
}
