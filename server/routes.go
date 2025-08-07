package server

import (
	"jaegerin/handlers"
)

func (s *Server) setupRoutes() {

	handlers.HandleTraces(s.mux)

}
