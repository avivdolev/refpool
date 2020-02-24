package refpool

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/avivdolev/refpool/example"
)

func Example() {
	// Set a global pool of example.Buffer (see local package example)
	// type Buffer struct {
	// 	bytes.Buffer
	// 	count int64
	// }
	pool := New(func() Element {
		return &example.Buffer{}
	})

	// This Example spawns some workers.
	// All workers will receive a ref counted buffer with data.
	// Each worker waits on a channel of buffers.
	// We use a WaitGroup only to let this example finish.
	numOfWorkers := 5
	channels := make(map[int]chan *example.Buffer)
	wg := &sync.WaitGroup{}
	for i := 0; i < numOfWorkers; i++ {
		channels[i] = make(chan *example.Buffer, 10)
		wg.Add(1)
		go func(id int, in chan *example.Buffer) {
			select {
			case b := <-in:
				fmt.Printf("worker %d got: %s\n", id, b.Bytes())
				pool.Put(b) // safely put the buffer back, even if other goroutines still use it
				wg.Done()
			}
		}(i, channels[i])
	}

	// Some distant routine which gets data input and sends it to workers
	// we want to reuse allocated buffers here
	input := bytes.NewReader([]byte("very important data"))
	go func(r io.Reader) {
		b := pool.Get().(*example.Buffer)
		b.Reset()
		IncElement(b, int64(len(channels)))
		io.Copy(b, r)
		for _, c := range channels {
			c <- b
		}
		wg.Wait()
		fmt.Printf("done, counter should be %d\n", IncElement(b, 0))
	}(input)
	wg.Wait()

	// Unordered output:
	// worker 0 got: very important data
	// worker 1 got: very important data
	// worker 2 got: very important data
	// worker 3 got: very important data
	// worker 4 got: very important data
	// done, counter should be 0
}
