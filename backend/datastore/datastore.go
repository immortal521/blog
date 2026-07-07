package datastore

import (
	"context"
	"fmt"

	"blog-server/config"
	"blog-server/ent"
	"blog-server/logger"
	"blog-server/pkg/txmgr"

	_ "github.com/lib/pq"
)

type contextTxKey struct{}

var _ txmgr.TxManager = (*DataStore)(nil)

// DataStore provides a unified access point to the database layer,
// optionally routing operations through an active transaction.
//
// The zero value is not usable; instances must be constructed via
// NewDataStore. DataStore is not safe for concurrent mutation when a
// transaction is in progress.
type DataStore struct {
	client *ent.Client
}

// NewDSN returns a PostgreSQL DSN string derived from cfg.
//
// It assumes all required fields are present; missing or invalid values
// will surface as connection errors at open time rather than here.
func NewDSN(cfg *config.Config) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
	)
}

// NewDataStore constructs a DataStore backed by an ent client.
//
// In development mode, the client enables debug logging, which may
// significantly impact performance due to verbose SQL output.
//
// If client initialization fails, a non-nil DataStore may still be
// returned alongside the error; callers must check the error before use.
func NewDataStore(cfg *config.Config, log logger.Logger) (*DataStore, error) {
	dsn := NewDSN(cfg)

	client, err := ent.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if cfg.App.IsDev() {
		client = client.Debug()
	}

	return &DataStore{
		client: client,
	}, nil
}

func (ds *DataStore) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, ok := ctx.Value(contextTxKey{}).(*ent.Tx); ok {
		return fn(ctx)
	}

	tx, err := ds.client.Tx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(context.WithValue(ctx, contextTxKey{}, tx)); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Client returns the active ent client.
//
// If a transaction is in progress, it returns the transaction-scoped
// client; otherwise, it returns the base client. The returned client
// should not be cached across transaction boundaries.
func (ds *DataStore) Client(ctx context.Context) *ent.Client {
	if tx, ok := ctx.Value(contextTxKey{}).(*ent.Tx); ok {
		return tx.Client()
	}
	return ds.client
}

// Close releases resources associated with the underlying client.
//
// It does not attempt to resolve or rollback any active transaction.
// Calling Close on an already closed or nil client results in the
// underlying driver's behavior.
func (ds *DataStore) Close() error {
	return ds.client.Close()
}
