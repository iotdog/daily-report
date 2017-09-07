package utils

import (
	"encoding/json"
	"net/http"

	"github.com/leesper/holmes"
)

// Jsonify 返回JSON应答
func Jsonify(w http.ResponseWriter, message interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(message); err != nil {
		holmes.Errorln(err)
	}
}
