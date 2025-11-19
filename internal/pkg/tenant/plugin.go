package tenant

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ctxKey struct{}

func WithID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKey{}, id)
}

func ID(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKey{}).(string); ok {
		return v
	}
	return ""
}

type Plugin struct{}

func (p *Plugin) Name() string { return "tenant" }

func (p *Plugin) Initialize(db *gorm.DB) error {
	db.Callback().Query().Before("gorm:query").Register("tenant_before_query", before)
	db.Callback().Create().Before("gorm:create").Register("tenant_before_create", beforeCreate)
	db.Callback().Update().Before("gorm:update").Register("tenant_before_update", beforeCreate)
	return nil
}

func before(db *gorm.DB) {
	tid := ID(db.Statement.Context)
	if tid == "" {
		return
	}
	db.Statement.AddClause(TenantClause{tenantID: tid})
}

func beforeCreate(db *gorm.DB) {
	tid := ID(db.Statement.Context)
	if tid == "" {
		return
	}
	if db.Statement.Schema != nil {
		if f := db.Statement.Schema.LookUpField("TenantID"); f != nil {
			f.Set(db.Statement.ReflectValue, tid)
		}
	}
}

type TenantClause struct{ tenantID string }

func (t TenantClause) Build(builder clause.Builder) {
	builder.WriteString("tenant_id = ?")
	builder.AddVar(builder, t.tenantID)
}
