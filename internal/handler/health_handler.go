package handler

import "net/http"

func (h *LedgerHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("up"))
}
