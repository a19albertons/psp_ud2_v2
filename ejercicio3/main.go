package main

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"sync"
	"time"
)

func main() {
	trabajos := 5
	var wg sync.WaitGroup
	for i := 0; i < trabajos; i++ {
		wg.Add(1)
		go func (numero string, wg *sync.WaitGroup) {
			descargar(numero, wg)
		} (strconv.Itoa(i), &wg)
	}
	wg.Wait()
	fmt.Println("All downloads completed")
}

func descargar(numero string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Starting download: file"+numero+".zip")
	tiempoEspera:=1
	tiempoEsperaVariable := rand.IntN(2)
	tiempoEsperaReal := tiempoEspera + tiempoEsperaVariable
	tamanoMinimo := 10
	tamanoVariable := rand.IntN(90)
	tamanoReal := tamanoMinimo + tamanoVariable
	time.Sleep(time.Duration(tiempoEsperaReal) * time.Second)
	fmt.Println("Completed download: file"+numero+".zip ("+strconv.Itoa(tamanoReal)+" MB)")
}