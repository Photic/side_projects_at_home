package main

import (
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/phuslu/log"
	"github.com/robfig/cron/v3"

	"side_projects_at_home/src/control"
	"side_projects_at_home/src/model"
	"side_projects_at_home/src/router"
)

var (
	upgrader = websocket.Upgrader{}
)

func scheduledJob() {
	sqlite, err := control.SqliteConnector()

	if err != nil {
		log.Printf("failed to init sqlite: %v", err)
	}

	// If the returned value has a Close() error method, call it when done.
	if closer, ok := interface{}(sqlite).(interface{ Close() error }); ok {
		defer closer.Close()
	}

	loanID := "house_loan"

	log.Info().Msg("scheduled job started")

	loan, _ := sqlite.GetLoan(loanID)

	add_interest := (loan.Amount * 0.035) / 365.25

	_ = sqlite.InsertAmount(loanID, add_interest, model.TxTypeInterest)
	log.Printf("Added Interest: %f to LoanID: %s", add_interest, loanID)
}

func middleware(handler_func http.HandlerFunc, want string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != want {
			w.Header().Set("Allow", want)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler_func(w, r)
	}
}

func main() {
	log.DefaultLogger = log.Logger{
		Level:      log.DebugLevel,
		Caller:     1,
		TimeField:  "date",
		TimeFormat: "2006-01-02",
		Writer:     &log.IOWriter{Writer: os.Stdout},
	}

	sqlite, err := control.SqliteConnector()

	if err != nil {
		log.Fatal().Msgf("failed to init sqlite: %v", err)
	}

	// Schedular
	cron_job := cron.New()
	if _, err := cron_job.AddFunc("0 3 * * *", scheduledJob); err != nil {
		log.Fatal().Msgf("failed to add cron job: %v", err)
	}
	cron_job.Start()
	// Schedular End

	http.HandleFunc("/", middleware(router.GETLoanPage(sqlite), http.MethodGet))
	http.HandleFunc("/loan/update", middleware(router.PUTUpdateLoan(sqlite), http.MethodPut))

	// Live reload websocket endpoint
	http.HandleFunc("/live-reload-ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		// Keep the connection open until the server restarts (process ends)
		select {}
	})

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "assets/favicon.ico")
	})

	// Serve static files from the "assets" directory at "/assets/"
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	log.Info().Msg("Listening on 0.0.0.0:8080")
	err = http.ListenAndServe("0.0.0.0:8080", nil)

	if err != nil {
		log.Error().Msgf("Could not start server, error: %s", err)
	}
}
