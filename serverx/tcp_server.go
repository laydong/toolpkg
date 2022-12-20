package serverx

import (
	"net"
	"time"
)

// TCPServer struct
type TCPServer struct {
	iServer
}

var _ Server = &TCPServer{}

// NewTCPServer 创建新的TCPServer
// NewTcpServer(x)表示总共有x个workers，并且全部预先创建好；
// NewTcpServer(y)表示总共有x个workers，只预先创建y个。
func NewTCPServer(n ...uint32) *TCPServer {
	tcpServer := &TCPServer{}
	tcpServer.SetWorkersPoolSize(n...)
	return tcpServer
}

// Run with loop
func (s *TCPServer) Run(listen string) (err error) {
	addr, err := net.ResolveTCPAddr("tcp", listen)
	if err != nil {
		return err
	}

	s.iServer.listener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	s.iServer.run(time.Second)
	return nil
}
