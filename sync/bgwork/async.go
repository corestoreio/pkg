package bgwork

import "sync"

// Async runs the given functions asynchronous. It returns all errors, if there
// are. The channel gets closed once all functions are done.
func Async(fns ...func() error) <-chan error {
	errChan := make(chan error)
	var wg sync.WaitGroup
	wg.Add(len(fns))
	go func() {
		wg.Wait()
		close(errChan)
	}()
	for i := range fns {
		go func(i int) {
			defer wg.Done()
			errChan <- fns[i]()
		}(i)
	}

	return errChan
}
