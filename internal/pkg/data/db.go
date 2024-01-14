package data

import (
	"context"
	"harmoni/internal/types/iface"

	"github.com/google/wire"
	"gorm.io/gorm"
)

var ProviderSet = wire.NewSet(
	NewDB,
	wire.Bind(new(iface.Transaction), new(*DB)),
)

var _ iface.Transaction = (*DB)(nil)

// 用来承载事务的上下文
type contextTxKey struct{}

type DB struct {
	db *gorm.DB
}

func NewDB(db *gorm.DB) *DB {
	return &DB{
		db: db,
	}
}
func (d *DB) ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return fn(ctx)
	})
}

// DB 根据此方法来判断当前的 db 是不是使用 事务的 DB
func (d *DB) DB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return d.db.WithContext(ctx)
}
