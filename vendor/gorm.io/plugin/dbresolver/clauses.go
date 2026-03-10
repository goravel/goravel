package dbresolver

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Operation specifies dbresolver mode
type Operation string

const (
	writeName = "gorm:db_resolver:write"
	readName  = "gorm:db_resolver:read"
)

// ModifyStatement modify operation mode
func (op Operation) ModifyStatement(stmt *gorm.Statement) {
	var optName string
	if op == Write {
		optName = writeName
		stmt.Settings.Delete(readName)
	} else if op == Read {
		optName = readName
		stmt.Settings.Delete(writeName)
	}

	if optName != "" {
		stmt.Settings.Store(optName, struct{}{})
		if fc := stmt.DB.Callback().Query().Get("gorm:db_resolver"); fc != nil {
			fc(stmt.DB)
		}
	}
}

// Build implements clause.Expression interface
func (op Operation) Build(clause.Builder) {
}

// Use specifies configuration
func Use(str string) clause.Expression {
	return using{Use: str}
}

type using struct {
	Use string
}

const usingName = "gorm:db_resolver:using"

// ModifyStatement modify operation mode
func (u using) ModifyStatement(stmt *gorm.Statement) {
	stmt.Clauses[usingName] = clause.Clause{Expression: u}
	if fc := stmt.DB.Callback().Query().Get("gorm:db_resolver"); fc != nil {
		fc(stmt.DB)
	}
}

// Build implements clause.Expression interface
func (u using) Build(clause.Builder) {
}
