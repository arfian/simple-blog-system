package transaction

import (
	"context"
	"database/sql"
)

type ISqlTransaction interface {
	Transaction(c context.Context, fn func(wrappedCtx context.Context) error, opts ...*sql.TxOptions) error
}
