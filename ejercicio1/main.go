package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

func main() {
	
	var numeroVeces int = 5
	var wg sync.WaitGroup
	for i := 0; i < numeroVeces; i++ {
		wg.Add(1)
		numeroString := strconv.Itoa(i)
		go func(numero string, wg *sync.WaitGroup) {
			parallelGreetings(numero, wg)
		}	(numeroString, &wg)
	}

	wg.Wait()
	fmt.Println("All goroutines completed!")
}

func parallelGreetings(numero string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("Goroutine "+numero+": Hello!")
	base := 100
	randomNumber := rand.Intn(400)
	tiempoEspera := base + randomNumber
	time.Sleep(time.Duration(tiempoEspera) * time.Millisecond)
	fmt.Println("Goroutine "+numero+": Goodbye!")
}
