// Package datastore
package datastore

import (
	"context"

	"go.uber.org/fx"
)

// Module returns an fx.Module that registers datastore related dependencies,
// including the database client, transaction manager, and lifecycle hooks.
func Module() fx.Option {
	return fx.Module("datastore",
		fx.Provide(
			NewDataStore,
		), fx.Invoke(registerLifecycle))
}

// registerLifecycle attaches lifecycle hooks for the datastore module.
//
// It ensures that the database connection is properly closed when the
// application stops.
func registerLifecycle(lc fx.Lifecycle, ds *DataStore) {
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return ds.Close()
		},
	})
}
