package loanengine

import "time"

type Installment struct {
	installmentAmount float64
	dueDate           time.Time
	paidOn            time.Time
}

type Payment struct {
	Amount float64   `json:"loanAmount"`
	Date   time.Time `json:"paymentDate"`
	Mode   string    // default on-line
}
