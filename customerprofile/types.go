package customerprofile

import (
	"BillingEngine/loanengine"
	"fmt"
	"sync"
)

var totalCustomer = 0
var accountNumberSeries int64 = 1000000000
var mu sync.Mutex

type Customer struct {
	id            string
	firstName     string
	lastName      string
	accountNumber int64
	contactNumber int64
	loanProfile   []loanengine.Loan // change it to map if needed.
	balance       float64
}

func NewCustomer(firstName, lastName string, contactNumber int64) *Customer {
	return &Customer{
		id:            generateCustomerId(),
		firstName:     firstName,
		lastName:      lastName,
		accountNumber: generateAccountNumber(),
		contactNumber: contactNumber,
		balance:       0.00,
		loanProfile:   make([]loanengine.Loan, 0),
	}
}

func (c *Customer) GetLoanProfile(loanId string) loanengine.Loan {
	for _, l := range c.loanProfile {
		if l.GetLoanId() == loanId {
			return l
		}
	}
	return nil
}

func (c *Customer) GetCustomerId() string {
	return c.id
}

func (c *Customer) SetLoanProfile(profile loanengine.Loan) {
	c.loanProfile = append(c.loanProfile, profile)
}

func generateCustomerId() string {
	mu.Lock()
	defer mu.Unlock()
	totalCustomer++
	return fmt.Sprintf("cust-%d", totalCustomer)
}

func generateAccountNumber() int64 {
	mu.Lock()
	defer mu.Unlock()
	accountNumberSeries++
	return accountNumberSeries
}
