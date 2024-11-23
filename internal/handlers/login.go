package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gatsu420/ngetes/internal/auth"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	auth := auth.NewAuth()

	_, token, _ := auth.Encode(map[string]interface{}{
		"user_id": "ngetes",
		"exp":     time.Now().Add(20 * time.Second).Unix(),
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
