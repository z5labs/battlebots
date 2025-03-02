// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package battle2d

import (
	"github.com/z5labs/battlebots/pkgs/battlepb"

	"google.golang.org/grpc"
)

type Service struct {
	battlepb.UnimplementedBattle2DServer
}

func Register(sr grpc.ServiceRegistrar) {
	s := &Service{}

	battlepb.RegisterBattle2DServer(sr, s)
}
