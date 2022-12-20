package serverx

import (
	"net"
	"time"
)

// Server interface
type Server interface {
	Run(string) error
	SetHandler(func(net.Conn))
	SetRejectHandler(func(net.Conn, error))
	SetAcceptErrorHandler(func(error))
	SetWorkersPoolSize(...uint32)
	CountBusyWorkers() uint32
	CountAvailableWorkers() uint32
	Close()
}

// iServer
type iServer struct {
	close              bool
	listener           net.Listener
	handler            func(net.Conn)
	rejectHandler      func(net.Conn, error)
	acceptErrorHandler func(error)
	workersPool        *workersPool
	workersPoolCap     uint32
	workersPoolPre     uint32
}

// SetHandler 设置连接建立后的处理方法，此方法会分到不同的Worker去执行，不阻塞
// 监听协程。
func (s *iServer) SetHandler(h func(net.Conn)) {
	s.handler = h
}

// SetRejectHandler 设置当连接数过多被拒绝时的处理方法，此方法在监听协程执
// 行，会阻塞监听协程，处理完成之前无法监听新的连接。
func (s *iServer) SetRejectHandler(h func(net.Conn, error)) {
	s.rejectHandler = h
}

// SetAcceptErrorHandler 当监听失败时的处理方法，此方法在监听协程执行，会阻塞监听
// 协程，处理完成之前无法监听新的连接。
func (s *iServer) SetAcceptErrorHandler(h func(error)) {
	s.acceptErrorHandler = h
}

// SetWorkersPoolSize 设置Worker数量
// 如果调用SetWorkersPoolSize(x)表示总共有x个workers，并且全部预先创建好；
// 如果调用SetWorkersPoolSize(y, x)表示总共有x个workers，只预先创建y个。
func (s *iServer) SetWorkersPoolSize(n ...uint32) {
	switch len(n) {
	case 1:
		s.workersPoolPre = n[0]
		s.workersPoolCap = n[0]
	case 2:
		s.workersPoolPre = n[0]
		s.workersPoolCap = n[1]
	default:
		panic("need 1 or 2 parameters")
	}
}

func (s *iServer) CountBusyWorkers() uint32 {
	return s.workersPool.Busy()
}

func (s *iServer) CountAvailableWorkers() uint32 {
	return s.workersPool.Available()
}

// Close 关闭
func (s *iServer) Close() {
	s.close = true
	s.listener.Close()
}

func (s *iServer) releaseWorker(worker *Worker) {
	s.workersPool.releaseWorker(worker)
}

func (s *iServer) run(idle time.Duration) {
	s.prepare()
	defer s.Close()
	for {
		if s.close {
			return
		}
		// Set Accept Timeout
		switch s.listener.(type) {
		case *net.TCPListener:
			l := s.listener.(*net.TCPListener)
			l.SetDeadline(time.Now().Add(idle))
		case *net.UnixListener:
			l := s.listener.(*net.UnixListener)
			l.SetDeadline(time.Now().Add(idle))
		}
		// Accept and error handle
		conn, err := s.listener.Accept()
		if err != nil {
			if s.acceptErrorHandler != nil {
				s.acceptErrorHandler(err)
			}
			continue
		}
		// acquire worker or reject handle
		worker, err := s.workersPool.acquireWorker()
		if err != nil {
			if s.rejectHandler != nil {
				s.rejectHandler(conn, err)
			}
			conn.Close()
			continue
		}
		worker.conn <- conn
	}
}

func (s *iServer) prepare() {
	if s.handler == nil {
		panic("Handler is not set")
	}
	s.workersPool = newWorkersPool(s.workersPoolCap)
	s.workersPool.createWorker = func() *Worker {
		return &Worker{
			Work:    s.handler,
			iServer: s,
			conn:    make(chan net.Conn),
		}
	}
	s.workersPool.preallocWorkers(s.workersPoolPre)
}
