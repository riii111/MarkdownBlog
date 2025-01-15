package health

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/riii111/markdown-blog-api/pkg/health"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	response := health.Response{
		Status:    "ok!!",
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
