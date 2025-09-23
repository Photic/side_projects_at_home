package control

import (
	"database/sql"
	"os"
	"path/filepath"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/phuslu/log"

	"side_projects_at_home/src/model"
)

type Sqlite struct {
	pool *sqlx.DB
}

var (
	sqliteOnce sync.Once
	sqlitePool *sqlx.DB
)

func getSqlitePath() (string, error) {
	configDir, err := os.UserConfigDir()

	if err != nil {
		log.Error().Msgf("%s", err)
		return "", err
	}

	dbDir := filepath.Join(configDir, "side_projects")

	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(dbDir, "database.sqlite"), nil
}

func poolConnector() (*sqlx.DB, error) {
	var err error
	sqliteOnce.Do(func() {
		var dbPath string
		dbPath, err = getSqlitePath()
		if err != nil {
			return
		}
		sqlitePool, err = sqlx.Open("sqlite3", dbPath)

		if err != nil {
			log.Error().Msgf("%s", err)
			return
		}

		// Set pool options if needed, e.g. MaxOpenConns, MaxIdleConns

		// Create tables
		schemaLoans := `
			CREATE TABLE IF NOT EXISTS loans (
				id TEXT PRIMARY KEY,
				amount REAL NOT NULL,
				payed REAL NOT NULL,
				initial REAL NOT NULL
			);
		`
		schemaHistory := `
			CREATE TABLE IF NOT EXISTS transaction_history (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				loan_id TEXT NOT NULL,
				amount REAL NOT NULL,
				type TEXT NOT NULL,
				date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				FOREIGN KEY (loan_id) REFERENCES loans(id)
			);
		`
		sqlitePool.MustExec(schemaLoans)
		sqlitePool.MustExec(schemaHistory)
		// SQLite doesn't support triggers in the same way as Postgres, so skip trigger for initial update
	})
	return sqlitePool, err
}

func SqliteConnector() (*Sqlite, error) {
	pool, err := poolConnector()

	if err != nil {
		log.Error().Msgf("%s", err)
	}

	if err != nil {
		return nil, err
	}

	return &Sqlite{pool: pool}, nil
}

func (s *Sqlite) InsertAmount(id string, amount float64, tx_type model.TxType) error {
	var currentAmount, payed float64

	row := s.pool.QueryRowx("SELECT amount, payed FROM loans WHERE id = ?", id)
	err := row.Scan(&currentAmount, &payed)

	if err != nil {
		log.Error().Msgf("%s", err)
	}

	if err == sql.ErrNoRows {
		// Insert new loan
		_, err = s.insertIntoLoans(id, amount, amount, 0.0)
		if err != nil {
			log.Error().Msg("Could not insert into loans")
			return err
		}

		// We have no row initially ofcourse.
		currentAmount = 0.0
		payed = 0.0

		if err := s.updateLoan(id, currentAmount+amount, 0.0); err != nil {
			return err
		}

		if err := s.insertIntoHistory(id, amount, tx_type); err != nil {
			return err
		}

		return nil
	} else if err != nil {
		return err
	}

	if err := s.updateLoan(id, currentAmount+amount, payed+amount); err != nil {
		return err
	}

	if err := s.insertIntoHistory(id, amount, tx_type); err != nil {
		return err
	}

	return nil
}

func (s *Sqlite) GetTransactionsPage(offset, limit int64) ([]model.TransactionHistory, error) {
	rows := []model.TransactionHistory{}

	err := s.pool.Select(&rows, `
		SELECT id, loan_id, amount, date FROM transaction_history
		ORDER BY date DESC
		LIMIT ? OFFSET ?`, limit, offset)

	if err != nil {
		log.Error().Msgf("%s", err)
	}

	return rows, err
}

func (s *Sqlite) GetLatestTransactions(id string) ([]model.TransactionHistory, error) {
	rows := []model.TransactionHistory{}

	err := s.pool.Select(&rows, `
		SELECT id, loan_id, amount, date FROM transaction_history
		WHERE loan_id = ?
		ORDER BY date DESC
		LIMIT 5`, id)

	return rows, err
}

func (s *Sqlite) GetLoan(id string) (model.Loan, error) {
	var loan model.Loan
	err := s.pool.Get(&loan, "SELECT id, amount, payed, initial FROM loans WHERE id = ?", id)

	if err != nil {
		log.Error().Msgf("%s", err)
	}

	if err == sql.ErrNoRows {
		return model.Loan{ID: "", Amount: 0, Payed: 0, Initial: 0}, nil
	}

	return loan, err
}

func (s *Sqlite) insertIntoLoans(id string, amount, payed, initial float64) (sql.Result, error) {
	obj, err := s.pool.Exec("INSERT OR REPLACE INTO loans (id, amount, payed, initial) VALUES (?, ?, ?, ?)", id, amount, payed, initial)

	if err != nil {
		log.Error().Msgf("%s", err)
	}

	return obj, err
}

func (s *Sqlite) updateLoan(id string, amount float64, payed float64) error {
	_, err := s.pool.Exec("UPDATE loans SET amount = ?, payed = ? WHERE id = ?", amount, payed, id)

	if err != nil {
		log.Error().Msgf("%s", err)
	}

	return err
}

func (s *Sqlite) insertIntoHistory(id string, amount float64, tx_type model.TxType) error {
	_, err := s.pool.Exec("INSERT INTO transaction_history (loan_id, amount, type) VALUES (?, ?, ?)", id, amount, tx_type)

	if err != nil {
		log.Error().Msgf("%s", err)
	}

	return err
}
