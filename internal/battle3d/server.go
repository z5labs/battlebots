// Copyright (c) 2025 Z5Labs and Contributors
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package battle3d

import "github.com/z5labs/battlebots/pkgs/battlepb"

type Server struct {
	battlepb.UnimplementedBattle3DServer
}

func NewServer() *Server {
	return &Server{}
}
