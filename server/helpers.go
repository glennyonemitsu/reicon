package server

func CreateServer(basepath string) *Server {
	s := new(Server)
	s.Bootstrap(basepath)
	return s
}
