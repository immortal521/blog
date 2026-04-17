// Package repository provides dependency injection wiring for repository layer components.
package repository

import "go.uber.org/fx"

// Module returns an fx.Option that registers repository constructors within
// the "repository" module scope.
//
// The returned option is intended to be composed into a larger fx.App. It groups
// repository providers under a named module to improve graph organization and
// debugging visibility.
//
// Only explicitly listed constructors are registered; omitted providers are not
// part of the dependency graph.
//
// This function is stateless and safe for concurrent use.
func Module() fx.Option {
	return fx.Module(
		"repository",
		fx.Provide(
			NewUserRepo,
			// NewLinkRepo,
			NewPostRepo,
		),
	)
}
