package model

type TransactionHistory struct {
	ID     int64   `db:"id"`
	LoanID string  `db:"loan_id"`
	Amount float64 `db:"amount"`
	Date   string  `db:"date"`
}
