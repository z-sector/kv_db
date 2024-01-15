package dsem

type SemaphoreChan struct {
	tickets chan struct{}
}

func NewSemaphoreChan(limit int) *SemaphoreChan {
	return &SemaphoreChan{
		tickets: make(chan struct{}, limit),
	}
}

func (s *SemaphoreChan) Acquire() {
	s.tickets <- struct{}{}
}

func (s *SemaphoreChan) Release() {
	<-s.tickets
}

func (s *SemaphoreChan) WithSemaphore(action func()) {
	if action == nil {
		return
	}

	s.Acquire()
	action()
	s.Release()
}
