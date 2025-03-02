// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package app

import (
	"context"

	"github.com/z5labs/battlebots/internal/battle2d"

	"github.com/z5labs/humus/grpc"
)

type Config struct {
	grpc.Config `config:",squash"`
}

func Init(ctx context.Context, cfg Config) (*grpc.Api, error) {
	api := grpc.NewApi()

	battle2d.Register(api)

	return api, nil
}
