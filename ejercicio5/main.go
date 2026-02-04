package main

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
)

func main() {
	numbers := 20
	var wg1 sync.WaitGroup // Waitgroup for processor stage
	wg1workers := 3
	var wg2 sync.WaitGroup // Waitgroup for validator stage
	wg2workers := 2
	firstChan := make(chan int, numbers)
	secondChan := make(chan int, numbers)
	finalChan := make(chan int, numbers)

	// step 1: Generate stage
	go func() {
		defer close(firstChan)
		for i := 0; i < numbers; i++ {
			j := i + 1
			fmt.Println("Generator produced:", j)
			firstChan <- j
		}

	}()

	// step 2: Generate slaves
	for i := 0; i < wg1workers; i++ {
		wg1.Add(1)
		go func(workerID int) {
			defer wg1.Done()
			for num := range firstChan {
				square := num * num
				fmt.Println("Slave from wg1 " + strconv.Itoa(workerID) + " " + strconv.Itoa(num) + "²" + " = " + strconv.Itoa(square))
				secondChan <- square
			}

		}(i + 1)
	}

	// step 3: Wait for all slaves from wg1 to finish
	wg1.Wait()
	close(secondChan)

	// step 4: Generate validator slaves
	for i := 0; i < wg2workers; i++ {
		wg2.Add(1)
		go func(workerID int) {
			defer wg2.Done()
			for squaredNum := range secondChan {
				if squaredNum%2 == 0 {
					fmt.Printf("Validator %d: %d is even (passed)\n", workerID, squaredNum)
					finalChan <- squaredNum
				} else {
					fmt.Printf("Validator %d: %d is odd (filtered)\n", workerID, squaredNum)
				}
			}
		}(i + 1)
	}

	// step 5: Wait for all slaves from wg2 to finish
	wg2.Wait()
	close(finalChan)

	finalList := []int{}

	// step6 : Collect results
	for finalNum := range finalChan {
		finalList = append(finalList, finalNum)
	}
	sort.Ints(finalList)

	fmt.Println("Final Results:", finalList)

}
