# BillingEngine

# Assumptions
1. There is no prepayment(anyway mentioned in the question prompt).
2. Processing missed repayments at the end when all current repayments are done by borrower.
3. Solution works for specific type of loan(ex, BusinessLoan), We can make it dynamic
4. Once delinquent status is changed to true, the only way to change it back to false is close the loan after 
   repayment.
5. Not considered the foreclose functionality.
6. hardcoded 50 weeks loan.

# How to Run

1. There are 2 ways to run it, first via CURL and second via test file

# Curl commands to run
1. First create the customer.
2. Create the loan-profile 
3. Now can call either of the method, payment(),outstanding(),delinquent status via curl

## create customer -curl
`curl --location 'localhost:8989/customer/create' \
--header 'Content-Type: application/json' \
--request "POST" \
--data '{
"firstName": "vishu",
"lastName": "Sharma"
}'`


## create loan profile

`curl --location 'localhost:8989/loans/create' \
--header 'Content-Type: application/json' \
--request "POST" \
--data '{
"customerId": "cust-1",
"loanAmount": 50000.00
}'`


## delinquet status
`curl --location 'localhost:8989/loans/delinquentStatus?custId=cust-1&loanId=BUSINESS-LOAN-ID-1' \
--header 'Content-Type: application/json' \
--request "GET"`


## balance-check
`curl --location 'localhost:8989/loans/balance?custId=cust-1&loanId=BUSINESS-LOAN-ID-1' \
--header 'Content-Type: application/json' \
--request "GET"`

## payment

`curl --location 'localhost:8989/loans/payment?custId=cust-1&loanId=BUSINESS-LOAN-ID-1' \
--header 'Content-Type: application/json' \
--request "POST" \
--data '{
"loanAmount": 1100.00,
"paymentDate": "2025-03-10T22:00:00Z"
}'`

** Run the go api server by changing directory to BillingEngine and run
`go run .`

# How to Run Test file
1. Just configure loan amount,customer-id(hard-coded in test) and loan-id(hardcoded in test)
2. Once we run the test, it prints outstanding and delinquent status after each payment.

** Run either `go test ./...` or use ide to run the test.