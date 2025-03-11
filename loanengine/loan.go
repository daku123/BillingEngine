package loanengine

type Loan interface {
	GetOutStanding() (float64, error)
	MakePayment(payment *Payment) error
	IsDelinquent() bool
	GetLoanStatus() LoanStatus
	GetLoanId() string
	GetTotalLoanAmount() float64
}
