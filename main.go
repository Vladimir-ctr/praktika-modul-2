package main

import "fmt"

type User struct {
	ID      string
	Name    string
	Balance float32
}

func (u *User) Deposit(amount float32) {
	u.Balance += amount
	fmt.Printf("Ваш баланс пополнен на: %.2f, текущий остаток: %.2f\n", amount, u.Balance)
}

func (u *User) Withdrawn(amount float32) {
	if u.Balance >= amount {
		u.Balance -= amount
		fmt.Printf("Списано %.2f, Ваш остаток %.2f\n", amount, u.Balance)
	} else {
		fmt.Println("Ошибка: недостаточно средств. Ваш баланс:\n", u.Balance)
	}
}

func main() {
	user1 := &User{"1", "First", 150.35}
	user2 := &User{"2", "Second", 1539.73}
	user3 := &User{"3", "Third", 539.00}

	user1.Deposit(132.35)
	user1.Withdrawn(32.00)
	fmt.Printf("%+v\n", user1)
	fmt.Println()

	user2.Deposit(132.35)
	user2.Withdrawn(32.00)
	fmt.Printf("%+v\n", user2)
	fmt.Println()

	user3.Deposit(132.35)
	user3.Withdrawn(32.00)
	fmt.Printf("%+v\n", user3)
	fmt.Println()

}
