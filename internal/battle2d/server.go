// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package battle2d

import "github.com/z5labs/battlebots/pkgs/battlepb"

type Server struct {
	battlepb.UnimplementedBattle2DServer
}

func NewServer() *Server {
	return &Server{}
}
