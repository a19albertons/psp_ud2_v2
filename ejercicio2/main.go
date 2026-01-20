package main

import (
	"fmt"
	"strconv"
	"sync"
)

func main() {
	workers := 4
	rango := 100
	resultados := make([]int, workers)
	var wg sync.WaitGroup

	cantidad := rango / workers
	origen := 0
	fin := cantidad
	for i := 0; i < workers; i++ {
		wg.Add(1)
		inicioModificado := origen + 1
		go func(inicioModificado int, fin int, worker int, wg *sync.WaitGroup) {
			resultados[i] = sumaFragmentada(inicioModificado, fin, i, wg)
		}(inicioModificado, fin, i, &wg)
		origen += cantidad
		fin += cantidad
	}
	wg.Wait()
	sumaTotal := 0
	for i := 0; i < workers; i++ {
		sumaTotal += resultados[i]
	}
	fmt.Println("Total Sum :" + strconv.Itoa(sumaTotal))
}

func sumaFragmentada(inicio int, fin int, worker int, wg *sync.WaitGroup) int {
	defer wg.Done()
	devolver := 0
	for i := inicio; i <= fin; i++ {
		devolver += i
	}
	fmt.Println("Worker " + strconv.Itoa(worker) + ": from " + strconv.Itoa(inicio) + " to " +strconv.Itoa(fin) + " = " + strconv.Itoa(devolver)) 
	return devolver
}
