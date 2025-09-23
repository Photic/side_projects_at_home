package router

import (
	"net/http"
	"strconv"

	"github.com/phuslu/log"

	"side_projects_at_home/src/control"
	"side_projects_at_home/src/views"
)

func renderLoanPage(writer http.ResponseWriter, request *http.Request, sqlite *control.Sqlite, loanId string) {
	loan, _ := sqlite.GetLoan(loanId)
	transaction_history, _ := sqlite.GetLatestTransactions(loanId)

	views.ConditionalRender(writer, request, views.BankLoan(loan, transaction_history))
}

func GETLoanPage(sqlite *control.Sqlite) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		renderLoanPage(writer, request, sqlite, "house_loan")
	}
}

func PUTUpdateLoan(sqlite *control.Sqlite) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		loanId := "house_loan"

		amountStr := request.FormValue("amount")
		action := request.FormValue("action")

		log.Debug().Msgf("POST Amount: %s, as Action: %s", amountStr, action)

		amount, err := strconv.ParseFloat(amountStr, 64)

		if err == nil {
			if action == "withdraw" {
				amount = -amount
			}

			_ = sqlite.InsertAmount(loanId, amount)
		}

		renderLoanPage(writer, request, sqlite, loanId)
	}

}
