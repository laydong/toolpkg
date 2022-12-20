package serverx

import (
	// "fmt"
	"net"
)

// Worker 是处理一个请求的Goroutine，worker可以复用
type Worker struct {
	Work    func(net.Conn)
	id      int
	pos     uint32
	iServer *iServer
	conn    chan net.Conn
}

func (worker *Worker) run() {
	go func() {
		for {
			select {
			case conn := <-worker.conn:
				worker.Work(conn)
				conn.Close()
				worker.iServer.releaseWorker(worker)
				// fmt.Printf("%d worker release\n", worker.id)
				// fmt.Printf("%d workers available\n", worker.iServer.workersPool.Available())
			}
		}
	}()
}

// SetID set worker id
func (worker *Worker) SetID(id int) {
	worker.id = id
}

// ID return worker id
func (worker *Worker) ID() int {
	return worker.id
}
