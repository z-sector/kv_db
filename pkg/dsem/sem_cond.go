package dsem

import "sync"

type SemaphoreCond struct {
	mutex sync.Mutex
	count int
	cond  *sync.Cond
}

func NewSemaphoreCond(limit int) *SemaphoreCond {
	sem := &SemaphoreCond{count: limit}
	sem.cond = sync.NewCond(&sem.mutex)
	return sem
}

func (s *SemaphoreCond) Acquire() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for s.count <= 0 {
		s.cond.Wait()
	}

	s.count--
}

func (s *SemaphoreCond) Release() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.count++
	s.cond.Signal()
}

func (s *SemaphoreCond) WithSemaphore(action func()) {
	if action == nil {
		return
	}

	s.Acquire()
	action()
	s.Release()
}
