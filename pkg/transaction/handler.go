package transaction

import (
	"database/sql"
	"strings"

	"github.com/avast/retry-go"
	"golang.org/x/net/context"
	"gorm.io/gorm"

	"simple-blog-system/config/db"
)

// KeyTrx type for key context value transaction
type KeyTrx string

// KeyTransaction concrete type for key context value transaction
const KeyTransaction KeyTrx = KeyTrx("postgredb-hr-trx")

type SqlTransaction struct {
	dbx *db.GormDB
}

// Transaction wraps sql transaction within 1 function with given context
func (s SqlTransaction) Transaction(c context.Context, fn func(wrappedCtx context.Context) error, opts ...*sql.TxOptions) error {
	const errorMsg = "could not serialize access due to read/write dependencies among transactions"
	return retry.Do(
		func() error {
			return s.dbx.Transaction(func(tx *gorm.DB) error {
				db := &db.GormDB{
					DB: tx,
				}
				ctx := context.WithValue(c, KeyTransaction, db)
				return fn(ctx)
			}, opts...)
		},
		retry.RetryIf(func(err error) bool {
			return strings.Contains(err.Error(), errorMsg)
		}),
		retry.LastErrorOnly(true),
	)
}

func NewSqlTransaction(db *db.GormDB) ISqlTransaction {
	return SqlTransaction{
		dbx: db,
	}
}
