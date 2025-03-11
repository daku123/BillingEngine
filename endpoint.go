package main

import (
	"BillingEngine/loanengine"
	"encoding/json"
	"net/http"
	"time"
)

type OutstandingResponse struct {
	Balance float64 `json:"outstandingPayment"`
	LoanId  string  `json:"loanId"`
	CustId  string  `json:"customerId"`
}
type DelinquentResponse struct {
	IsDelinquent bool `json:"isDelinquent"`
}
type CreateLoanRequest struct {
	CustId     string  `json:"customerId"`
	LoanAmount float64 `json:"loanAmount"`
}

type CreateCustomerRequest struct {
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	ContactNumber int64  `json:"contactNumber"`
}

func getOutstandingAmount(w http.ResponseWriter, r *http.Request) {
	custId := r.URL.Query().Get("custId")
	loanId := r.URL.Query().Get("loanId")
	if len(custId) == 0 || len(loanId) == 0 {
		http.Error(w, "Pass customerId and loanId", http.StatusBadRequest)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	bankInstance := GetBankInstance()
	balance, err := bankInstance.getOutstandingBalance(custId, loanId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(OutstandingResponse{Balance: balance, CustId: custId, LoanId: loanId})
}

func makePayment(w http.ResponseWriter, r *http.Request) {
	custId := r.URL.Query().Get("custId")
	loanId := r.URL.Query().Get("loanId")
	if len(custId) == 0 || len(loanId) == 0 {
		http.Error(w, "Pass customerId and loanId", http.StatusBadRequest)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	bankInstance := GetBankInstance()

	var payment loanengine.Payment
	json.NewDecoder(r.Body).Decode(&payment)
	//payment.Date = time.Now()
	err := bankInstance.makePayment(custId, loanId, &payment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	remaining, _ := bankInstance.getOutstandingBalance(custId, loanId)
	response := map[string]interface{}{
		"customerId":        custId,
		"loanId":            loanId,
		"outstandingAmount": remaining,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func isDelinquent(w http.ResponseWriter, r *http.Request) {
	custId := r.URL.Query().Get("custId")
	loanId := r.URL.Query().Get("loanId")
	if len(custId) == 0 || len(loanId) == 0 {
		http.Error(w, "Pass customerId and loanId", http.StatusBadRequest)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	bankInstance := GetBankInstance()
	isDelinquent, err := bankInstance.getDelinquentStatus(custId, loanId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(DelinquentResponse{IsDelinquent: isDelinquent})
}

func createLoanProfile(w http.ResponseWriter, r *http.Request) {
	// check if customer exists, else create customer first then loanProfile
	// if customer exists, most likely loan too, create another loan
	var createReq CreateLoanRequest
	json.NewDecoder(r.Body).Decode(&createReq)
	if len(createReq.CustId) == 0 || createReq.LoanAmount == 0 {
		http.Error(w, "Pass customerId and loanAmount", http.StatusBadRequest)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method, expecting POST method call", http.StatusMethodNotAllowed)
		return
	}
	bankInstance := GetBankInstance()

	if bankInstance.customers[createReq.CustId] == nil {
		http.Error(w, "Create Customer First", http.StatusBadRequest)
		return
	}
	customerProfile := bankInstance.customers[createReq.CustId]
	// TODO :- this is an assumption that loan is type of business,
	// TODO :- we can enhance this in future but keeping as it is right now.
	loanProfile := loanengine.NewBusinessLoan(createReq.LoanAmount)
	customerProfile.SetLoanProfile(loanProfile)

	remaining, _ := loanProfile.GetOutStanding()
	response := map[string]interface{}{
		"customerId":        customerProfile.GetCustomerId(),
		"loanId":            loanProfile.GetLoanId(),
		"totalLoan":         loanProfile.GetTotalLoanAmount(),
		"loanStatus":        loanProfile.GetLoanStatus(),
		"outstandingAmount": remaining,
		"firstEmiDate":      time.Now().Add(7 * 24 * time.Hour),
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func createCustomer(w http.ResponseWriter, r *http.Request) {
	var createReq CreateCustomerRequest
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method, expecting POST method call", http.StatusMethodNotAllowed)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&createReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(createReq.FirstName) == 0 || len(createReq.LastName) == 0 {
		http.Error(w, "Pass firstName and LastName ", http.StatusBadRequest)
		return
	}
	customer := GetBankInstance().createCustomer(&createReq)

	response := map[string]interface{}{
		"customerId": customer.GetCustomerId(),
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
