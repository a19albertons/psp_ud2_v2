package main

import (
	"fmt"
	"sync"
	"time"
)

func countdown(seconds *int, mu *sync.Mutex) {
    for *seconds > 0 {
        time.Sleep(1 * time.Second)
        *seconds -= 1
		mu.Unlock()
    }
}

func main() {
    count := 5
	var mu sync.Mutex
    go countdown(&count, &mu)
    for count > 0 {
		mu.Lock()
        time.Sleep(500 * time.Millisecond)
        fmt.Println(count)
		
    }
}