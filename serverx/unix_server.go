package serverx

import (
	"net"
	"time"
)

// UnixServer struct
type UnixServer struct {
	iServer
}

var _ Server = &UnixServer{}

// NewUnixServer 创建新的TCPServer
// NewUnixServer(x)表示总共有x个workers，并且全部预先创建好；
// NewUnixServer(y)表示总共有x个workers，只预先创建y个。
func NewUnixServer(n ...uint32) *UnixServer {
	tcpServer := &UnixServer{}
	tcpServer.SetWorkersPoolSize(n...)
	return tcpServer
}

// Run with loop
func (s *UnixServer) Run(file string) (err error) {
	addr, err := net.ResolveUnixAddr("unix", file)
	if err != nil {
		return err
	}
	s.listener, err = net.ListenUnix("unix", addr)
	if err != nil {
		return err
	}

	s.iServer.run(time.Second)
	return nil
}
