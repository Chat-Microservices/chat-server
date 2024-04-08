package transaction

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/semho/chat-microservices/chat-server/internal/client/db"
	"github.com/semho/chat-microservices/chat-server/internal/client/db/pg"
)

type manager struct {
	db db.Transactor
}

func NewTransactionManager(db db.Transactor) db.TxManager {
	return &manager{
		db: db,
	}
}

func (m *manager) transaction(ctx context.Context, opts pgx.TxOptions, fn db.Handler) (err error) {
	// если транзакция вложенная, пропускаем инициализацию новой транзакции и выполняем функцию
	tx, ok := ctx.Value(pg.TxKey).(pgx.Tx)
	if ok {
		return fn(ctx)
	}
	// новая транзакция
	tx, err = m.db.BeginTx(ctx, opts)
	if err != nil {
		return errors.Wrap(err, "can't begin transaction")
	}
	// отправляем в контекст
	ctx = pg.MakeContextTx(ctx, tx)

	// откат коммита
	defer func() {
		// востановление
		if r := recover(); r != nil {
			err = errors.Errorf("panic recovered: %v", r)
		}
		// откат транзакции
		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				err = errors.Wrapf(err, "can't rollback transaction: %v", errRollback)
			}
			return
		}
		err = tx.Commit(ctx)
		if err != nil {
			err = errors.Wrap(err, "can't commit transaction")
		}
	}()

	// выполняем код транзакции
	if err = fn(ctx); err != nil {
		return errors.Wrap(err, "failed executing code inside transaction")
	}

	return err
}

func (m *manager) ReadCommitted(ctx context.Context, f db.Handler) (err error) {
	txOpts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	return m.transaction(ctx, txOpts, f)
}
