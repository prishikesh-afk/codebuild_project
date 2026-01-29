package responses

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Status      int
	Data        any
	Message     string
	IsPlainText bool
	Err         error
}

func Write(w http.ResponseWriter, resp APIResponse) {
	w.WriteHeader(resp.Status)

	if resp.IsPlainText {
		if resp.Err != nil {
			w.Write([]byte(resp.Err.Error()))
			return
		}
		w.Write([]byte(resp.Message))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	out := make(map[string]any)

	if resp.Data != nil {
		out["data"] = resp.Data
	}
	if resp.Message != "" {
		out["message"] = resp.Message
	}
	if resp.Err != nil {
		out["error"] = resp.Err.Error()
	}

	json.NewEncoder(w).Encode(out)
}
