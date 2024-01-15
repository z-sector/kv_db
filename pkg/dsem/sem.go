package dsem

type Semaphore interface {
	Acquire()
	Release()
	WithSemaphore(action func())
}
