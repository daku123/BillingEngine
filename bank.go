package main

import (
	"BillingEngine/customerprofile"
	"BillingEngine/loanengine"
	"fmt"
	"sync"
)

var once sync.Once
var bankInstance *Bank

type Bank struct {
	lock      sync.RWMutex
	customers map[string]*customerprofile.Customer
}

func GetBankInstance() *Bank {
	once.Do(func() {
		bankInstance = &Bank{customers: make(map[string]*customerprofile.Customer)}
	})
	return bankInstance
}

func (b *Bank) getOutstandingBalance(customerId, loanId string) (float64, error) {
	b.lock.RLock()
	defer b.lock.RUnlock()
	if b.customers[customerId] == nil {
		return 0.00, fmt.Errorf("Customer does not exists ")
	}
	loanProfile := b.customers[customerId].GetLoanProfile(loanId)
	if loanProfile == nil || loanProfile.GetLoanStatus() == loanengine.Closed {
		return 0.00, fmt.Errorf("No active loan for customer ")
	}
	return loanProfile.GetOutStanding()
}

func (b *Bank) getDelinquentStatus(customerId, loanId string) (bool, error) {
	b.lock.RLock()
	defer b.lock.RUnlock()
	if b.customers[customerId] == nil {
		return false, fmt.Errorf("Customer does not exist ")
	}
	loanProfile := b.customers[customerId].GetLoanProfile(loanId)
	if loanProfile == nil || loanProfile.GetLoanStatus() == loanengine.Closed {
		return false, fmt.Errorf("No active loan for customer ")
	}
	return loanProfile.IsDelinquent(), nil
}

func (b *Bank) makePayment(customerId, loanId string, payment *loanengine.Payment) error {
	b.lock.RLock()
	defer b.lock.RUnlock()
	if b.customers[customerId] == nil {
		return fmt.Errorf("Customer does not exist ")
	}
	loanProfile := b.customers[customerId].GetLoanProfile(loanId)
	if loanProfile == nil || loanProfile.GetLoanStatus() == loanengine.Closed {
		return fmt.Errorf("No active loan for customer ")
	}
	return loanProfile.MakePayment(payment)
}

func (b *Bank) createCustomer(request *CreateCustomerRequest) *customerprofile.Customer {
	b.lock.Lock()
	defer b.lock.Unlock()
	customer := customerprofile.NewCustomer(request.FirstName, request.LastName, request.ContactNumber)
	b.customers[customer.GetCustomerId()] = customer
	return customer
}
