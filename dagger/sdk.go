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

	ep.Go(sdk.Go().Ci)

	return ep.Wait()
}

type SdkGo struct {
	// +private
	Library *dagger.GoLibrary
}

func (sdk *Sdk) Go() *SdkGo {
	c := dag.Go(dagger.GoOpts{Version: "latest"}).Container()

	c = dag.Protobuf("v31.1").
		Go("v1.36.6").
		Grpc().
		Protobuf().
		CopyTo(c)

	lib := dag.Go(dagger.GoOpts{
		From: c,
	}).Module(sdk.Source, dagger.GoModuleOpts{
		Path: "sdk/battlebots-go",
	}).
		Library()

	return &SdkGo{
		Library: lib,
	}
}

func (g *SdkGo) Ci(ctx context.Context) error {
	return g.Library.Ci(ctx)
}
