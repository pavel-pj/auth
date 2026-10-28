package closer

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

//var globalCloser = New()

type Closer struct {
	mu    sync.Mutex
	onse  sync.Once
	done  chan struct{}
	funcs []func() error
}

var globalCloser = New(os.Interrupt, syscall.SIGTERM)

func New(sigs ...os.Signal) *Closer {

	c := &Closer{done: make(chan struct{})}
	if len(sigs) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, sigs...)
			<-ch
			signal.Stop(ch)
			c.CloseAll()
		}()

	}

	return c
}
func Add(f ...func() error) {
	globalCloser.Add(f...)
}

func Wait() {
	globalCloser.Wait()
}

func CloseAll() {
	globalCloser.CloseAll()
}

func (c *Closer) Wait() {
	<-c.done
}

func (c *Closer) Add(f ...func() error) {
	c.mu.Lock()
	c.funcs = append(c.funcs, f...)
	c.mu.Unlock()
}

func (c *Closer) CloseAll() {

	c.onse.Do(func() {
		defer close(c.done)

		var funcs []func() error
		c.mu.Lock()
		for _, f := range c.funcs {
			funcs = append(funcs, f)
		}
		c.funcs = nil
		c.mu.Unlock()

		errs := make(chan error, len(funcs))
		for _, f := range funcs {
			go func(f func() error) {
				errs <- f()
			}(f)
		}

		for x := 0; x < len(funcs); x++ {
			if err := <-errs; err != nil {
				log.Println("Error returned from Closer")
			}
		}

	})

}

/*
var globalCloser = New(os.Interrupt, syscall.SIGTERM)

func Add(f ...func() error) {
	globalCloser.Add(f...)
}

func Wait() {
	globalCloser.Wait()
}

func CloseAll() {
	globalCloser.CloseAll()
}

type Closer struct {
	mu    sync.Mutex
	once  sync.Once
	done  chan struct{}
	funcs []func() error
}

func New(sig ...os.Signal) *Closer {

	c := &Closer{done: make(chan struct{})}

	if len(sig) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, sig...)
			<-ch
			signal.Stop(ch)
			log.Println("ВЫКЛЮЧЕНО ПОМАННЫМ СИНГАЛОМ")
			c.CloseAll()
		}()

	}

	return c

}

func (c *Closer) Add(f ...func() error) {
	c.mu.Lock()
	c.funcs = append(c.funcs, f...)
	c.mu.Unlock()
}

func (c *Closer) Wait() {
	<-c.done
}

func (c *Closer) CloseAll() {

	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()
		//fmt.Println("ЗАКРЫВАЕМ ВСЕ__1")
		errs := make(chan error, len(funcs))
		for _, f := range funcs {
			go func(f func() error) {
				errs <- f()
			}(f)
		}

		for i := 0; i < cap(errs); i++ {
			if err := <-errs; err != nil {
				log.Println("error returned from Closer")
			}
		}

	})

}
*/
