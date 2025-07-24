// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"
	"dagger/battle-bots/internal/dagger"
)

type Example struct {
	// +private
	Source *dagger.Directory
}

func (m *BattleBots) Example() *Example {
	return &Example{
		Source: m.Source,
	}
}

type RandomMoveBotGo struct {
	// +private
	Application *dagger.GoApplication
}

func (e *Example) RandomMoveBotGo() *RandomMoveBotGo {
	return &RandomMoveBotGo{
		Application: dag.Go(dagger.GoOpts{
			Version: "latest",
		}).
			Module(e.Source, dagger.GoModuleOpts{
				Path: "example/random-move-bot-go",
			}).
			Library().
			Application("."),
	}
}

func (r *RandomMoveBotGo) Ci(ctx context.Context) error {
	return r.Application.Ci(
		ctx,
		"",
		"",
		[]string{},
		"",
		nil,
	)
}

func (r *RandomMoveBotGo) AsService() *dagger.Service {
	return r.Application.AsService()
}
