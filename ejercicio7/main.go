package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
)

type BankAccount struct {
	balance          float64
	mutex            sync.Mutex
	transactionCount int
}

func main() {
	// Create a bank account
	var account BankAccount
	account.balance = 1000

	// wg for goroutines
	var wg sync.WaitGroup
	workers := 10
	subworkers := 100

	//statistics
	successfulDeposits := 0
	successfulWithdrawals := 0
	var succesfulMutex sync.Mutex

	// Simulate concurrent transactions
	for i := 0; i < workers; i++ {
		wg.Add(1)
		// goroutine for each worker logic
		go func(){
			defer wg.Done()
			for i := 0; i < subworkers; i++ {
				addOrRemove := rand.IntN(2)
				amount := 1 + rand.IntN(50)
				if addOrRemove == 0 {
					err := Deposit(&account, float64(amount))
					succesfulMutex.Lock()
					if err == nil {
						successfulDeposits++
					}
					succesfulMutex.Unlock()
				} else {
					err := Withdraw(&account, float64(amount))
					succesfulMutex.Lock()
					if err == nil {
						successfulWithdrawals++
					}
					succesfulMutex.Unlock()
				}

			}
 		}()
	}

	wg.Wait()

	fmt.Println("Simulation Complete!")
	fmt.Printf("Final balance: %.2f\n", account.balance)
	fmt.Printf("Total transactions: %d\n", account.transactionCount)
	fmt.Printf("Successful deposits: %d\n", successfulDeposits)
	fmt.Printf("Successful withdrawals: %d\n", successfulWithdrawals)
	fmt.Printf("Failed Transactions: %d\n", workers*subworkers - account.transactionCount)

}

// This funcion handle deposits, when you add money to the account
func Deposit(account *BankAccount, amount float64) error {
	// enable mutex to avoid race conditions
	account.mutex.Lock()
	defer account.mutex.Unlock()

	// Increase transaction count
	account.transactionCount++

	// validate amount
	if amount >= 0 {
		account.balance += amount
		return nil
	} else {
		account.transactionCount-- // Revert transaction count increase if the amount is invalid
		return fmt.Errorf("Your amount is negative")
	}

}

// This function handle withdrawals, when you remove money from the account
func Withdraw(account *BankAccount, amount float64) error {
	// enable mutex to avoid race conditions
	account.mutex.Lock()
	defer account.mutex.Unlock()

	// Increase transaction count
	account.transactionCount++

	//validate a positive value
	if amount >= 0 {
		if account.balance >= amount {
			account.balance -= amount
			return nil
		} else {
			account.transactionCount-- // Revert transaction count increase if the withdrawal fails
			return fmt.Errorf("Insufficient funds")
		}

	} else {
		account.transactionCount-- // Revert transaction count increase if the amount is invalid
		return fmt.Errorf("Your amount is negative")
	}

}

// this function handle balance consults, when you want to know how much money you have in the account
func GetBalance(account *BankAccount) float64 {
	account.mutex.Lock()
	defer account.mutex.Unlock()

	return account.balance
}


// this functin handle transaction count consults, when you want to know succesful transactions you have made in the account
func GetTransactionCount(account *BankAccount) int  {
	account.mutex.Lock()
	defer account.mutex.Unlock()

	return account.transactionCount
}
