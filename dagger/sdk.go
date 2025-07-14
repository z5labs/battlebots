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

type Sdk struct {
	// +private
	Source *dagger.Directory
}

func (m *BattleBots) Sdk() *Sdk {
	return &Sdk{
		Source: m.Source,
	}
}

func (sdk *Sdk) Ci(ctx context.Context) error {
	ep := pool.New().WithErrors().WithContext(ctx)

	ep.Go(sdk.Go)

	return ep.Wait()
}

func (sdk *Sdk) Go(ctx context.Context) error {
	c := dag.Go(dagger.GoOpts{Version: "latest"}).Container()

	c = dag.Protobuf("v31.1").
		Go("v1.36.6").
		Grpc().
		Protobuf().
		CopyTo(c)

	return dag.Go(dagger.GoOpts{
		From: c,
	}).
		Module(sdk.Source, dagger.GoModuleOpts{
			Path: "sdk/battlebots-go",
		}).
		Library().
		Ci(ctx)
}
