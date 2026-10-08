package main

import (
	"errors"
	"fmt"
)

type Account struct {
	Owner   string
	Balance float64
}

func (a *Account) Deposit(amount float64) {
	a.Balance += amount
}

func (a *Account) Withdraw(amount float64) error {
	if amount > a.Balance {
		return errors.New("Not Enough Money to Withddraw")
	}
	a.Balance -= amount
	return nil
}

func (a Account) Summary() string {
	return fmt.Sprintf("%v has %v rupees", a.Owner, a.Balance)
}

func main() {
	account := Account{"Shankar", 1000}
	fmt.Println(account.Summary())
	account.Deposit(200)
	fmt.Println(account.Summary())
	account.Withdraw(500)
	fmt.Println(account.Summary())
}
