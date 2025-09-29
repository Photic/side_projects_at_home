package router

import (
	"net/http"

	"side_projects_at_home/src/views"
)

func GETDashboardPage() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		views.ConditionalRender(writer, request, views.Dashboard())
	}
}
