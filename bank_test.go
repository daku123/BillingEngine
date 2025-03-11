package main

import (
	"BillingEngine/customerprofile"
	"BillingEngine/loanengine"
	"testing"
	"time"
)

var b *Bank

func createLoanProfileTest(custId string, loanAmount float64) loanengine.Loan {
	//b := createBankInstance()
	//b.customers["cust-1"]
	loanProfile := loanengine.NewBusinessLoan(loanAmount)
	b.customers[custId].SetLoanProfile(loanProfile)
	return loanProfile
}

func createCustomerTest(firstName, lastName string) *customerprofile.Customer {
	b = createBankInstance()
	return b.createCustomer(&CreateCustomerRequest{LastName: lastName, FirstName: firstName})
}

func createBankInstance() *Bank {
	return GetBankInstance()
}

func TestBank_makePayment1(t *testing.T) {
	b = createBankInstance()

	c := createCustomerTest("ab1", "ab2")
	l := createLoanProfileTest(c.GetCustomerId(), 55000)
	customerId := c.GetCustomerId()
	loanId := l.GetLoanId()

	type args struct {
		customerId string
		loanId     string
		payment    *loanengine.Payment
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy billing test case",
			args: args{loanId: loanId, customerId: customerId},
		},
	}
	for _, tt := range tests {

		for i := 1; i <= 50; i++ {
			week := i * 7 * 24
			tt.args.payment = &loanengine.Payment{Amount: 1210.00, Date: time.Now().Add(time.Duration(week) * time.
				Hour)}
			t.Run(tt.name, func(t *testing.T) {
				if err := b.makePayment(tt.args.customerId, tt.args.loanId, tt.args.payment); (err != nil) != tt.wantErr {
					t.Errorf("makePayment() error = %v, wantErr %v", err, tt.wantErr)
				}
				status, err := b.getDelinquentStatus(tt.args.customerId, tt.args.loanId)
				if status {
					t.Errorf("status is expected fale got true error = %v", err)
				}
			})
		}
		// should fail, loan is closed
		week := 51 * 7 * 24 // 51th week
		tt.args.payment = &loanengine.Payment{Amount: 1210.00, Date: time.Now().Add(time.Duration(week) * time.
			Hour)}
		err := b.makePayment(tt.args.customerId, tt.args.loanId, tt.args.payment)
		if err == nil {
			t.Error("Should fail ")
		}
	}
}

func TestBank_makePayment2(t *testing.T) {
	b = createBankInstance()

	c := createCustomerTest("acf", "ghj")
	customerId := c.GetCustomerId()
	l := createLoanProfileTest(customerId, 55000)
	loanId := l.GetLoanId()

	type args struct {
		customerId string
		loanId     string
		payment    *loanengine.Payment
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "mixed test cases",
			args: args{loanId: loanId, customerId: customerId},
		},
	}
	for _, tt := range tests {

		for i := 5; i <= 60; i++ {
			week := i * 7 * 24
			tt.args.payment = &loanengine.Payment{Amount: 1210.00, Date: time.Now().Add(time.Duration(week) * time.
				Hour)}
			if i < 55 {
				t.Run(tt.name, func(t *testing.T) {
					if err := b.makePayment(tt.args.customerId, tt.args.loanId, tt.args.payment); (err != nil) != tt.wantErr {
						t.Errorf("makePayment() error = %v, wantErr %v", err, tt.wantErr)
					}
				})
				status, err := b.getDelinquentStatus(tt.args.customerId, tt.args.loanId)
				// on 55th iteration payment call will close the loan and this status changed to false.
				if !status && i < 54 {
					t.Errorf("status is expected true got false error = %v", err)
				}
			} else {
				// should fail after 50 repayment as loan is closed, status will be false.
				err := b.makePayment(tt.args.customerId, tt.args.loanId, tt.args.payment)
				if err == nil {
					t.Errorf("should be failing as loan is closed %v ", err)
				}
				// should turn to false as loan is closed
				status, err := b.getDelinquentStatus(tt.args.customerId, tt.args.loanId)
				if status {
					t.Errorf("status is expected false got true error = %v", err)
				}
			}
		}
	}
}
