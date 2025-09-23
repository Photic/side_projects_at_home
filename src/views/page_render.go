package views

import (
	"net/http"
	"os"

	"github.com/a-h/templ"
)

func ConditionalRender(writer http.ResponseWriter, request *http.Request, component templ.Component) {
	devMode := os.Getenv("DEV_MODE") == "1"

	writer.Header().Set("Content-Type", "text/html; charset=utf-8")

	if request.Header.Get("hx-request") == "true" {
		_ = component.Render(request.Context(), writer)
	} else {
		_ = Layout(component, devMode).Render(request.Context(), writer)
	}
}
