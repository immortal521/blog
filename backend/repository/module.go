// Package repository provides dependency injection wiring for repository layer components.
package repository

import "go.uber.org/fx"

// Package repository wires repository-layer dependencies using Uber Fx.
//
// This module groups all repository providers under a single Fx module scope
// to improve dependency graph structure and debugging clarity.
func Module() fx.Option {
	return fx.Module(
		"repository",
		fx.Provide(
			NewUserRepo,
			NewLinkRepo,
			NewPostRepo,
		),
	)
}
