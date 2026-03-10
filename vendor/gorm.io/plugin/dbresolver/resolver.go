package dbresolver

import (
	"gorm.io/gorm"
)

type resolver struct {
	sources           []gorm.ConnPool
	replicas          []gorm.ConnPool
	policy            Policy
	dbResolver        *DBResolver
	traceResolverMode bool
}

func (r *resolver) resolve(stmt *gorm.Statement, op Operation) (connPool gorm.ConnPool) {
	if op == Read {
		if len(r.replicas) == 1 {
			connPool = r.replicas[0]
		} else {
			connPool = r.policy.Resolve(r.replicas)
		}
		if r.traceResolverMode {
			markStmtResolverMode(stmt, ResolverModeReplica)
		}
	} else if len(r.sources) == 1 {
		connPool = r.sources[0]
		if r.traceResolverMode {
			markStmtResolverMode(stmt, ResolverModeSource)
		}
	} else {
		connPool = r.policy.Resolve(r.sources)
		if r.traceResolverMode {
			markStmtResolverMode(stmt, ResolverModeSource)
		}
	}

	if stmt.DB.PrepareStmt {
		if preparedStmt, ok := r.dbResolver.prepareStmtStore[connPool]; ok {
			return &gorm.PreparedStmtDB{
				ConnPool: connPool,
				Mux:      preparedStmt.Mux,
				Stmts:    preparedStmt.Stmts,
			}
		}
	}

	return
}

func (r *resolver) call(fc func(connPool gorm.ConnPool) error) error {
	for _, s := range r.sources {
		if err := fc(s); err != nil {
			return err
		}
	}

	for _, re := range r.replicas {
		if err := fc(re); err != nil {
			return err
		}
	}
	return nil
}
