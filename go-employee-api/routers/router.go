package routers

import (
	"net/http"
	"strings"

	"go-employee-api/handlers"
	"go-employee-api/responses"
)

type MethodHandler func(*http.Request, map[string]string) responses.APIResponse

var routes = map[string]map[string]MethodHandler{
	"/employees": {
		http.MethodGet:  handlers.GetEmployees,
		http.MethodPost: handlers.CreateEmployee,
	},
	"/employees/{name}": {
		http.MethodGet:    handlers.GetEmployeeByName,
		http.MethodPut:    handlers.UpdateEmployee,
		http.MethodDelete: handlers.DeleteEmployee,
	},
}

func MainRouter(w http.ResponseWriter, r *http.Request) {
	for route, methods := range routes {
		params, ok := matchRoute(route, r.URL.Path)
		if !ok {
			continue
		}

		handler, exists := methods[r.Method]
		if !exists {
			responses.Write(w, responses.APIResponse{
				Status:      http.StatusMethodNotAllowed,
				Message:     "Method not allowed",
				IsPlainText: true,
			})
			return
		}

		responses.Write(w, handler(r, params))
		return
	}

	responses.Write(w, responses.APIResponse{
		Status:      http.StatusNotFound,
		Message:     "Route not found",
		IsPlainText: true,
	})
}

func matchRoute(route, path string) (map[string]string, bool) {
	rp := strings.Split(strings.Trim(route, "/"), "/")
	pp := strings.Split(strings.Trim(path, "/"), "/")

	if len(rp) != len(pp) {
		return nil, false
	}

	params := map[string]string{}
	for i := range rp {
		if strings.HasPrefix(rp[i], "{") && strings.HasSuffix(rp[i], "}") {
			params[rp[i][1:len(rp[i])-1]] = pp[i]
		} else if rp[i] != pp[i] {
			return nil, false
		}
	}
	return params, true
}
