package minions

import (
	"log"
	"sync"
)

type wChannel[T any] struct {
	ch chan T
}

func newWChannel[T any](size int) wChannel[T] {
	if size > 0 {
		return wChannel[T]{ch: make(chan T, size)}
	}
	return wChannel[T]{ch: make(chan T)}
}

type Minions[in, out any] struct {
	input  wChannel[in]
	output wChannel[out]
	n      int
	fun    func(in) out
}

func ListenHere[in, out any](fun func(in) out, minionsNum, inNum, outNum int) (Minions[in, out], chan in, chan out) {
	input := newWChannel[in](inNum)
	output := newWChannel[out](outNum)
	return Minions[in, out]{
		input:  input,
		output: output,
		n:      minionsNum,
		fun:    fun,
	}, input.ch, output.ch
}

func (m *Minions[in, out]) Go() (minionsDone chan struct{}) {
	wg := &sync.WaitGroup{}
	wg.Add(m.n)
	minionsDone = make(chan struct{})
	go func() {
		for n := 1; n <= m.n; n++ {
			go func(w int) {
				for {
					select {
					case v, ok := <-m.input.ch:
						if !ok {
							log.Printf("Minion %v shurgs: no more work to do\n", w)
							wg.Done()
							return
						}
						m.output.ch <- m.fun(v)
					}
				}
			}(n)
		}
		wg.Wait()
		close(m.output.ch)
		minionsDone <- struct{}{}
		close(minionsDone)
	}()
	return minionsDone
}
