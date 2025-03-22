package closer

import (
	"log"
	"os"
	"os/signal"
	"sync"
)

var globalCloser = New()

func Add(closeFunc ...closeFunction) { globalCloser.add(closeFunc...) }

func Wait() { globalCloser.wait() }

// CloseAll - wrapper for Closer.closeAll. Calls all closer functions.
func CloseAll() { globalCloser.closeAll() }

type closeFunction func() error

type Closer struct {
	mu    sync.Mutex
	once  sync.Once
	done  chan struct{}
	funcs []closeFunction
}

func New(sig ...os.Signal) *Closer {
	c := &Closer{done: make(chan struct{})}
	if len(sig) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, sig...)
			<-ch
			signal.Stop(ch)
			c.closeAll()
		}()
	}

	return c
}

// closeAll calls all closer functions.
func (c *Closer) closeAll() {
	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		// call all funcs async
		errs := make(chan error, len(funcs))
		for _, f := range funcs {
			go func(f closeFunction) {
				errs <- f()
			}(f)
		}

		for i := 0; i < cap(errs); i++ {
			if err := <-errs; err != nil {
				log.Println("error returned from Closer", err)
			}
		}
	})
}

func (c *Closer) add(closeFunc ...closeFunction) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.funcs = append(c.funcs, closeFunc...)
}

func (c *Closer) wait() { <-c.done }
