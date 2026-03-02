package main

import (
	"errors"
	"fmt"
	"sync"
)

type User struct {
	ID      string
	Name    string
	Balance float64
	mu      sync.Mutex
}

func (u *User) Deposit(amount float64) {
	u.mu.Lock()
	u.Balance += amount
	u.mu.Unlock()
}

func (u *User) Withdrawn(amount float64) error {
	u.mu.Lock()
	if u.Balance < amount {
		return errors.New("insufficient funds")
	}
	u.Balance -= amount
	u.mu.Unlock()
	return nil
}

type Transaction struct {
	FromID string
	ToID   string
	Amount float64
}

type PaymentSystem struct {
	Users        map[string]*User
	Transactions []Transaction
	ch           chan Transaction
}

func (ps *PaymentSystem) AddUser(id, name string, balance float64) {
	if ps.Users == nil {
		ps.Users = make(map[string]*User)
	}
	ps.Users[id] = &User{
		ID:      id,
		Name:    name,
		Balance: balance,
	}
}

func (ps *PaymentSystem) AddTransaction(t Transaction) {
	// Добавления транзакций в очередь (добавлие элемента в слайс)
	ps.Transactions = append(ps.Transactions, t)

}

func (ps *PaymentSystem) ProcessingTransactions(tr Transaction) error {
	for _, t := range ps.Transactions {
		sender, senderExists := ps.Users[t.FromID]
		if !senderExists {
			return fmt.Errorf("Sender not found")
		}

		receiver, recExists := ps.Users[t.ToID]
		if !recExists {
			return fmt.Errorf("Receiver not found")
		}

		err := sender.Withdrawn(t.Amount) //списание средств с пользователя
		if err != nil {
			return fmt.Errorf("insufficient funds for User: %s\n", t.FromID)
		}

		receiver.Deposit(t.Amount)

		fmt.Printf("Balance User: %s changed on %.2f. Now balance %.2f\n", t.ToID, t.Amount, receiver.Balance)
	}

	ps.Transactions = nil
	return nil
}

func worker(ps *PaymentSystem, ch <-chan Transaction, wg *sync.WaitGroup) {
	defer wg.Done()

	for tr := range ch {
		err := ps.ProcessingTransactions(tr)
		if err != nil {
			fmt.Printf("Transaction failed: %v", err)
		} else {
			fmt.Println("Function comleted")
		}
	}
}

func main() {
	paymentSystem := &PaymentSystem{
		Users:        make(map[string]*User),
		Transactions: make([]Transaction, 0),
	}

	fmt.Println("Создаю UserID: 1 с балансом 1000")
	fmt.Println("Создаю UserID: 2 с балансом 500")

	paymentSystem.AddUser("1", "Artem", 1000)
	paymentSystem.AddUser("2", "Dmitry", 500)

	fmt.Println("Перевожу с UserID: 1 на UserID: 2 сумму в размере 200")
	fmt.Println("Перевожу с UserID: 2 на UserID: 1 сумму в размере 50")

	paymentSystem.AddTransaction(Transaction{FromID: "1", ToID: "2", Amount: 200})
	paymentSystem.AddTransaction(Transaction{FromID: "2", ToID: "1", Amount: 50})

	var wg sync.WaitGroup
	ch := make(chan Transaction)

	for i := 1; i <= 6; i++ {
		wg.Add(1)
		go worker(paymentSystem, ch, &wg)
	}

	for _, t := range paymentSystem.Transactions {
		ch <- t
	}
	close(ch)
	wg.Wait()

	for id, user := range paymentSystem.Users {
		fmt.Printf("User %s (%s): balance = %.2f\n", id, user.Name, user.Balance)
	}

	fmt.Println("Итого")
	fmt.Printf("У первого пользователя должно получиться 850, а получилось %.2f\n", paymentSystem.Users["1"].Balance)
	fmt.Printf("У второго пользователя должно получиться 650, а получилось %.2f", paymentSystem.Users["2"].Balance)

}
