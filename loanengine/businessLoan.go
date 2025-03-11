package loanengine

import (
	"fmt"
	"sync"
	"time"
)

var totalBusinessLoan = 0
var loanAccountNumberSeries int64 = 1000000000
var mu sync.Mutex

type LoanStatus string

var (
	Open           LoanStatus = "OPEN"
	Closed         LoanStatus = "CLOSED"
	loanTimePeriod            = 50
	interest       float64    = 10
)

type BusinessLoan struct {
	loanId               string
	loanAccountNumber    int64
	totalLoanAmount      float64
	outstandingAmount    float64
	upcomingInstallments []*Installment
	delayedInstallments  []*Installment
	isDelinquent         bool
	loanStatus           LoanStatus
	loanCreationDate     time.Time
	lastRepaymentAt      time.Time
	lock                 sync.RWMutex
}

func NewBusinessLoan(loanAmount float64) Loan {
	b := &BusinessLoan{
		loanId:              generateLoanId(),
		loanAccountNumber:   generateLoanAccountNumber(),
		totalLoanAmount:     loanAmount,
		outstandingAmount:   loanAmount,
		delayedInstallments: make([]*Installment, 0),
		isDelinquent:        false,
		loanStatus:          Open,
		loanCreationDate:    time.Now(),
		lastRepaymentAt:     time.Now(),
	}
	b.upcomingInstallments = b.generateInstallments(loanAmount)
	return b
}

func (b *BusinessLoan) generateInstallments(principal float64) []*Installment {
	payableAmount := principal + (principal * interest / 100)
	b.outstandingAmount = payableAmount
	emiAmount := payableAmount / float64(loanTimePeriod)

	emis := make([]*Installment, loanTimePeriod+1)

	for i := 1; i <= loanTimePeriod; i++ {
		dueDate := i * 24 * 7
		emis[i] = &Installment{installmentAmount: emiAmount, dueDate: time.Now().Add(time.Duration(dueDate) * time.Hour)}
	}
	return emis[1:]
}

func (b *BusinessLoan) GetOutStanding() (float64, error) {
	b.lock.RLock()
	defer b.lock.RUnlock()
	if b.loanStatus == Closed {
		return 0.00, fmt.Errorf("Loan is closed for loan-id %s ", b.loanId)
	}
	return b.outstandingAmount, nil
}

func (b *BusinessLoan) GetTotalLoanAmount() float64 {
	return b.totalLoanAmount
}

func (b *BusinessLoan) IsDelinquent() bool {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.isDelinquent {
		return b.isDelinquent
	}
	currentYear, currentWeek := time.Now().ISOWeek()
	for i, emi := range b.upcomingInstallments {
		yearOfEmi, weekOfEmi := emi.dueDate.ISOWeek()

		if currentYear == yearOfEmi && weekOfEmi == currentWeek && i >= 1 {
			b.isDelinquent = true
			return b.isDelinquent
		}
	}
	return b.isDelinquent
}

// MakePayment
/** function focuses on deducting repayment based on the week of the payment,
	if the payment week does not match, with the emi week,it won't process. We are making sure that borrower will
	do payment each week or won't do. If they do,it should match with current week of payment.
**/

func (b *BusinessLoan) MakePayment(payment *Payment) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	if b.loanStatus == Closed {
		return fmt.Errorf("Loan is closed for loan-id %s, payment is not accepted ", b.loanId)
	}
	if len(b.upcomingInstallments) == 0 && len(b.delayedInstallments) == 0 {
		return fmt.Errorf("Loan repayment is done, no schduled emi left for loanId %s ", b.loanId)
	}
	// find the week from the payment date and check with the pending emi,
	//if does not match move the pending to delayed list, also update delinquent
	currentYear, weekOfPayment := payment.Date.ISOWeek()

	// once current scheduled emi's are over, repayment is done based on delayed emi's
	if len(b.upcomingInstallments) == 0 && len(b.delayedInstallments) > 0 {
		// check in delayed emi and remove from them one by one
		b.outstandingAmount = b.outstandingAmount - payment.Amount
		b.delayedInstallments = b.delayedInstallments[1:]
		if b.outstandingAmount == 0.00 {
			// mark loan closed
			b.loanStatus = Closed
		}
		return nil
	}

	startYear, startWeek := b.upcomingInstallments[0].dueDate.ISOWeek()
	endYear, endWeek := b.upcomingInstallments[len(b.upcomingInstallments)-1].dueDate.ISOWeek()
	for i, emi := range b.upcomingInstallments {
		yearOfEmi, weekOfEmi := emi.dueDate.ISOWeek()

		// check if payment is from the current week or atleast from future date to make sure,
		//we stay in the boundary limit of 50 week
		if len(b.upcomingInstallments) > 0 && ((startWeek > weekOfPayment && startYear == currentYear) || (weekOfPayment > endWeek && endYear == currentYear)) {
			// can't process this request as it is coming past weeks,
			return fmt.Errorf("payment date is supposed to be within week range of year %d and week %d - year %d and"+
				" week %d ", startYear, startWeek, endYear, endWeek)
		}

		if weekOfEmi == weekOfPayment && yearOfEmi == currentYear {
			// found the emi,
			if payment.Amount != emi.installmentAmount {
				return fmt.Errorf("repayment amout is not matching , not accepting the payment")
			}
			b.outstandingAmount = b.outstandingAmount - payment.Amount
			if b.outstandingAmount == 0.00 {
				// mark loan closed
				b.loanStatus = Closed
			}
			b.upcomingInstallments = b.upcomingInstallments[i+1:]
			break
		} else {
			// check if the borrower missed the 2 continuous payments
			// payment does not belong to first 2 payment of upcoming emi so it means borrower missed two continuous emi
			if i >= 1 {
				b.isDelinquent = true
			}
			b.delayedInstallments = append(b.delayedInstallments, emi)
		}
	}
	b.lastRepaymentAt = payment.Date
	return nil
}

func (b *BusinessLoan) GetLoanStatus() LoanStatus {
	b.lock.RLock()
	defer b.lock.RUnlock()
	return b.loanStatus
}

func (b *BusinessLoan) GetLoanId() string {
	return b.loanId
}

func generateLoanId() string {
	mu.Lock()
	defer mu.Unlock()
	totalBusinessLoan++
	return fmt.Sprintf("BUSINESS-LOAN-ID-%d", totalBusinessLoan)
}

func generateLoanAccountNumber() int64 {
	mu.Lock()
	defer mu.Unlock()
	loanAccountNumberSeries++
	return loanAccountNumberSeries
}
