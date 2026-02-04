package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	webs := []string{
		"https://site1.com",
		"https://site2.com",
		"https://site3.com",
		"https://site4.com",
		"https://site5.com",
		"https://site6.com",
		"https://site7.com",
		"https://site8.com",
		"https://site9.com",
		"https://site10.com",
	}

	// Configuración
	workers := 3
	var wg sync.WaitGroup

	// Canales
	// inChan recibe strings (URLs) para procesar
	inChan := make(chan string, len(webs))
	// outChan recibe bools (true = exito, false = fallo)
	outChan := make(chan bool)

	falloExito := 3

	TotalTimeStart := time.Now()

	// 1. ENVIAR TRABAJOS (Producer)
	// Enviamos las URLs al canal en una goroutine separada para no bloquear
	// la lectura de resultados más abajo.
	go func() {
		for _, url := range webs {
			inChan <- url
		}
		close(inChan) // Cerramos el canal para avisar a los workers que no hay más trabajos
	}()

	// 2. INICIAR WORKERS (Pool de trabajadores)
	// Aquí es donde interactúa la variable "workers". Creamos solo 3 goroutines.
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			// El worker escucha el canal hasta que se cierre (close(inChan))
			for url := range inChan {
				// Simular proceso
				fmt.Printf("Started %d processing: %s\n", workerID, url)
				tiempo := 500 + rand.Intn(1500)
				time.Sleep(time.Duration(tiempo) * time.Millisecond)
				

				// Lógica de fallo/éxito
				f := rand.Intn(5)
				if f == falloExito {
					fmt.Printf("Worker %d failed URL: %s (error timeout)\n", workerID, url)
					outChan <- false // Enviar fallo
				} else {
					fmt.Printf("Worker %d succeeded URL: %s\n", workerID, url)
					outChan <- true // Enviar éxito
				}
			}
		}(i+1)
	}

	// 3. RECOLECTAR RESULTADOS (Consumer del outChan)
	// Esto se hace en el hilo principal para contar de forma segura
	totalURLs := 0
	successfulURLs := 0
	failedURLs := 0

	// Esperamos exactamente tantos resultados como URLs enviamos
	for i := 0; i < len(webs); i++ {
		resultado := <-outChan
		totalURLs++
		if resultado {
			successfulURLs++
		} else {
			failedURLs++
		}
	}

	// Esperar a que los workers terminen oficialmente (limpieza)
	wg.Wait()
	close(outChan)

	TotalTimeEnd := time.Since(TotalTimeStart)
	fmt.Println("--------------------------------")
	fmt.Println("Total URLs processed:", totalURLs)
	fmt.Println("Successful URLs:", successfulURLs)
	fmt.Println("Failed URLs:", failedURLs)
	fmt.Println("Total time:", TotalTimeEnd.Seconds())
}
