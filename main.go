package main

import (
	"log"
	"net/http"
)

func handleRequests() {
	// better to add the prefix /customer/{customerId}/loans/{loanId} to endpoints,
	//but for simplicity going ahead without it.
	http.Handle("/loans/balance", http.HandlerFunc(getOutstandingAmount))
	http.Handle("/loans/payment", http.HandlerFunc(makePayment))
	http.Handle("/loans/delinquentStatus", http.HandlerFunc(isDelinquent))
	http.Handle("/customer/create", http.HandlerFunc(createCustomer))
	http.Handle("/loans/create", http.HandlerFunc(createLoanProfile))
	log.Fatal(http.ListenAndServe(":8989", nil))
}

func main() {
	handleRequests()
}
