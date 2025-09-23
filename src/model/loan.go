package model

type Loan struct {
	ID      string  `db:"id"`
	Amount  float64 `db:"amount"`
	Payed   float64 `db:"payed"`
	Initial float64 `db:"initial"`
}
