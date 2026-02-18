package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Expresion struct {
	numero1  int
	operador string
	numero2  int
}
type lineaFichero struct {
	worker    int
	resultado int
	expresion Expresion
}

func main() {
	// numero esclavos
	esclavo := 5

	// Cosas de workgroups
	var wg sync.WaitGroup
	var mutexFicheroSaliente sync.Mutex
	canalLecturaFichero := make(chan string)

	//estadisticas
	contadorGeneral := 0
	var mutexEstadistica sync.Mutex
	esclavosSlice := make([]string, esclavo)
	mases := 0
	menos := 0
	multiplicar := 0
	division := 0
	errores := 0

	// Procesar el fichero (leerlo)
	buffer, error := os.ReadFile("expressions.txt")
	if error != nil {
		fmt.Println("ha habido un error de lectura")
		return
	}
	ficheroSaliente, error2 := os.OpenFile("salida.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if error2 != nil {
		fmt.Println("error con el fichero de salida")
		return
	}

	texto := string(buffer)
	troceado := strings.Split(texto, "\n")

	// mandar cada linea a un canal
	go func() {
		defer close(canalLecturaFichero)
		for i := 0; i < len(troceado); i++ {
			canalLecturaFichero <- troceado[i]
			mutexEstadistica.Lock()
			contadorGeneral++
			mutexEstadistica.Unlock()
		}

	}()

	// Recibir la linea del canal y lo que corresponda
	for i := 0; i < esclavo; i++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			contador := 0
			for linea := range canalLecturaFichero {
				contenido := strings.Split(linea, " ")
				if len(contenido) == 3 {
					var expresion Expresion
					numero1, error1 := strconv.Atoi(contenido[0])
					numero2, error2 := strconv.Atoi(contenido[2])
					expresion.operador = contenido[1]
					if error1 != nil || error2 != nil {
						mutexFicheroSaliente.Lock()
						ficheroSaliente.WriteString("Worker " + strconv.Itoa(worker) + " \"" + linea + "\" = Undefined\n")
						mutexFicheroSaliente.Unlock()
						continue

					}
					expresion.numero1 = numero1
					expresion.numero2 = numero2
					mutexFicheroSaliente.Lock()
					switch expresion.operador {
					case "+":
						{
							mutexEstadistica.Lock()
							mases++
							mutexEstadistica.Unlock()
							ficheroSaliente.WriteString("Worker " + strconv.Itoa(worker) + " \"" + linea + "\" = " + strconv.Itoa(numero1+numero2) + "\n")
						}
					case "-":
						{
							mutexEstadistica.Lock()
							menos++
							mutexEstadistica.Unlock()
							ficheroSaliente.WriteString("Worker " + strconv.Itoa(worker) + " \"" + linea + "\" = " + strconv.Itoa(numero1-numero2) + "\n")
						}
					case "*":
						{
							mutexEstadistica.Lock()
							multiplicar++
							mutexEstadistica.Unlock()
							ficheroSaliente.WriteString("Worker " + strconv.Itoa(worker) + " \"" + linea + "\" = " + strconv.Itoa(numero1*numero2) + "\n")
						}
					case "/":
						{
							mutexEstadistica.Lock()
							if numero2 == 0 {
								errores++
								ficheroSaliente.WriteString("Worker " + strconv.Itoa(worker) + " \"" + linea + "\" = undefined\n")
							} else {
								division++
								ficheroSaliente.WriteString("Worker " + strconv.Itoa(worker) + " \"" + linea + "\" = " + strconv.Itoa(numero1/numero2) + "\n")
							}
							mutexEstadistica.Unlock()
						}
					default:
						{
							mutexEstadistica.Lock()
							errores++
							mutexEstadistica.Unlock()
							ficheroSaliente.WriteString("Worker " + strconv.Itoa(worker) + " \"" + linea + "\" = operador no valido\n")
						}
					}
					mutexFicheroSaliente.Unlock()

				} else {
					mutexEstadistica.Lock()
					errores++
					mutexEstadistica.Unlock()
					mutexFicheroSaliente.Lock()
					ficheroSaliente.WriteString("Worker " + strconv.Itoa(worker) + " \"" + linea + "\" = no cumple la estructura numero, operador, numero\n")
					mutexFicheroSaliente.Unlock()
				}
				mutexEstadistica.Lock()
				contador++
				mutexEstadistica.Unlock()
			}
			mutexEstadistica.Lock()
			esclavosSlice[i] = "Thread " + strconv.Itoa(worker) + ": " + strconv.Itoa(contador)
			mutexEstadistica.Unlock()
		}(i + 1)
	}

	wg.Wait()
	// estadisticas
	fmt.Println("Stats:")
	fmt.Println("------")
	fmt.Println()
	fmt.Println(contadorGeneral)
	fmt.Println()
	fmt.Println("Number of expressions per thread:")
	fmt.Println("---------------------------------")
	for i := 0; i < len(esclavosSlice); i++ {
		fmt.Println(esclavosSlice[i])
	}
	fmt.Println()
	fmt.Println("Percentage of operations:")
	fmt.Println("-------------------------")
	fmt.Println("+: " + strconv.Itoa(mases))
	fmt.Println("-: " + strconv.Itoa(menos))
	fmt.Println("*: " + strconv.Itoa(multiplicar))
	fmt.Println("/: " + strconv.Itoa(division))
	fmt.Println("Error: " + strconv.Itoa(errores))
	fmt.Println("End of main thread.")

	//estadisticas

}
