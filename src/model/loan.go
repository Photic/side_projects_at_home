package model

type Loan struct {
	ID      string  `db:"id"`
	Amount  float64 `db:"amount"`
	Initial float64 `db:"initial"`
}
