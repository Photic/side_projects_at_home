package model

import (
	"database/sql/driver"
	"fmt"
)

type TransactionHistory struct {
	ID     int64   `db:"id"`
	LoanID string  `db:"loan_id"`
	Amount float64 `db:"amount"`
	Type   TxType  `db:"type"`
	Date   string  `db:"date"`
}

func (t *TxType) Scan(value interface{}) error {
	if value == nil {
		*t = ""
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*t = TxType(string(v))
		return nil
	case string:
		*t = TxType(v)
		return nil
	default:
		return fmt.Errorf("cannot scan %T into TxType", value)
	}
}

func (t TxType) Value() (driver.Value, error) {
	return string(t), nil
}

type TxType string

const (
	TxTypeDeposit  TxType = "deposit"
	TxTypeWithdraw TxType = "withdraw"
	TxTypeInterest TxType = "interest"
)
