// package main

// import (
// 	"fmt"
// 	"sync"
// )

// func sayHi(wg *sync.WaitGroup) {
// 	defer wg.Done() // Tell the WaitGroup: "I'm done"
// 	fmt.Println("Hi!")
// }

// func main() {
// 	var wg sync.WaitGroup // Create a WaitGroup

// 	wg.Add(1) // We're going to wait for 1 goroutine

// 	go sayHi(&wg) // Start the goroutine

// 	wg.Wait() // Wait until the goroutine calls Done()
// 	fmt.Println("All done!")
// }

// package main

// import (
// 	"fmt"
// 	"sync"
// )

// var balance = 0
// var mu sync.Mutex
// var wg sync.WaitGroup

// func deposit() {
// 	defer wg.Done()
// 	for i := 0; i < 1000; i++ {
// 		mu.Lock()
// 		balance++
// 		mu.Unlock()
// 	}
// }

// func main() {
// 	wg.Add(2)
// 	go deposit()
// 	go deposit()
// 	wg.Wait()

// 	fmt.Println("Final Balance:", balance)
// }

// sync.RWMutex — Read/Write Mutex
// It’s like a smarter Mutex. Instead of just allowing one person (goroutine) in at a time, it:
//     ✅ Allows multiple readers at the same time.
//     🚫 Only one writer can enter, and no readers allowed during writing.

// package main
// import (
// 	"fmt"
// )

// func main() {
// 	ch := make(chan string) // create a channel

// 	go func() {
// 		ch <- "hello" // send value into channel
// 	}()

// 	msg := <-ch // receive value from channel
// 	fmt.Println("Received:", msg)
// }

package main

import (
	"fmt"
)

func main() {
	ch := make(chan int)

	go func() {
		for i := 1; i <= 5; i++ {
			ch <- i
		}
		close(ch) // important: close when done sending
	}()

	for num := range ch {
		fmt.Println("Got:", num)
	}
}
