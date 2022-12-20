package serverx

import (
	"fmt"
	"sync"
)

type workersPool struct {
	capacity     uint32
	busy         uint32
	workers      []*Worker
	createWorker func() *Worker
	lock         *sync.Mutex
}

func newWorkersPool(cap uint32) *workersPool {
	if cap == 0 {
		panic("cap must greater than 0, a zero size poolx is useless")
	}
	return &workersPool{
		workers:  make([]*Worker, cap),
		capacity: cap,
		lock:     new(sync.Mutex),
	}
}

// Busy 获取当前有多少个worker正在使用
func (wm *workersPool) Busy() uint32 {
	return wm.busy
}

// Available 获取当前还有多少可用的worker
func (wm *workersPool) Available() uint32 {
	return wm.capacity - wm.busy
}

// Capacity 获取workers的容量
func (wm *workersPool) Capacity() uint32 {
	return wm.capacity
}

// preallocWorkers 预先分配n个Workers，n必须比capacity小
func (wm *workersPool) preallocWorkers(n uint32) {
	if n > wm.capacity {
		panic("n must be between 0 and capacity")
	}
	for idx := 0; idx < int(n); idx++ {
		wm.workers[idx] = wm.newWorkerAndRun(uint32(idx))
	}
}

// 获取一个可用的worker，并标记其为已用，若超过capacity，返回error
func (wm *workersPool) acquireWorker() (*Worker, error) {
	wm.lock.Lock()
	defer wm.lock.Unlock()
	if wm.busy < wm.capacity {
		if wm.workers[wm.busy] == nil {
			wm.workers[wm.busy] = wm.newWorkerAndRun(wm.busy)
		}
		worker := wm.workers[wm.busy]
		wm.busy++
		return worker, nil
	}
	return nil, fmt.Errorf("no available worker exists")
}

// 释放一个worker回到可用区，并标记其为可用
func (wm *workersPool) releaseWorker(worker *Worker) {
	wm.lock.Lock()
	defer wm.lock.Unlock()
	if wm.busy == 0 {
		panic("BUG: all workers are free, release what?")
	}
	wm.busy--
	wm.workers[worker.pos] = wm.workers[wm.busy]
	wm.workers[worker.pos].pos = worker.pos
	wm.workers[wm.busy] = worker
	wm.workers[wm.busy].pos = wm.busy
}

func (wm *workersPool) newWorkerAndRun(pos uint32) *Worker {
	worker := wm.createWorker()
	worker.id = int(pos)
	worker.pos = pos
	worker.run()
	return worker
}
