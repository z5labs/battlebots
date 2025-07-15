// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"

	"dagger/battle-bots/internal/dagger"

	"github.com/sourcegraph/conc/pool"
)

type BattleBots struct {
	// +private
	Source *dagger.Directory
}

func New(
	// +defaultPath="."
	source *dagger.Directory,
) *BattleBots {
	return &BattleBots{
		Source: source,
	}
}

func (m *BattleBots) Ci(ctx context.Context) error {
	lib := dag.Go(dagger.GoOpts{
		Version: "latest",
	}).
		Module(m.Source).
		Library()

	ep := pool.New().WithErrors().WithContext(ctx)

	ep.Go(lib.Ci)
	ep.Go(m.Sdk().Ci)

	return ep.Wait()
}
