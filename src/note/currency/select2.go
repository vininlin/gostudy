package main

import "sync"

/**
设置channel为nil,它就会阻塞，没办法再被
 */
func main()  {
	var wg sync.WaitGroup
	wg.Add(3)

	a, b := make(chan int), make(chan int)

	go func() {
		wg.Done()
		for {

			select {
			case x, ok := <-a:
				if !ok {
					a = nil
					break
				}
				println("a", x)
			case x, ok := <-b:
				if !ok {
					b = nil
					break
				}
				println("b", x)

			}
			if a == nil && b == nil {
				return
			}
		}

	}()

	go func() {
		defer wg.Done()
		defer close(a)

		for i := 0; i < 3; i++ {
			a<-i
		}

	}()

	go func() {
		defer wg.Done()
		defer close(b)

		for i := 0; i < 5; i++ {
			b<-i * 10
		}

	}()

	wg.Wait()
}
