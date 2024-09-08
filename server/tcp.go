package server


type TCPServerDev struct {
	server
}

func (s *TCPServerDev) Start() error {
	return nil
}

func (s *TCPServerDev) Stop() error {
	return nil
}