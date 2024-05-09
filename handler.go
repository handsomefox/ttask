package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
)

// Calculate endpoint has to take JSON with following structure: {"a":int,"b":int}
type CalculateRequest struct {
	A *int `json:"a"`
	B *int `json:"b"`
}

// Calculate will return JSON with the a! and b!
type CalculateResponse struct {
	A uint64 `json:"a"`
	B uint64 `json:"b"`
}

// HandleCalculate expects the context to contain the value of the request. It does not read the request body.
func HandleCalculate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, ok := r.Context().Value(CalculateRequestKey).(CalculateRequest)
	if !ok {
		slog.Warn("HandleCalculate must not be called before or without the CalculateMiddleware")
		WriteIncorrectInputError(w)
		return
	}

	a, b := concurrentFactorial(uint64(*body.A), uint64(*body.B))
	if err := writeJSON(w, CalculateResponse{A: a, B: b}, http.StatusOK); err != nil {
		slog.Error("Error when writing json to connection", "err", err)
	}
}

func concurrentFactorial(a, b uint64) (uint64, uint64) {
	var wg sync.WaitGroup

	// calculate factorial of a and b using goroutines
	wg.Add(1)
	go func() {
		defer wg.Done()
		slog.Debug("Starting calculating A")
		a = factorial(a)
		slog.Debug("Finished calculating A")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		slog.Debug("Starting calculating B")
		b = factorial(b)
		slog.Debug("Finished calculating B")
	}()

	wg.Wait()

	return a, b
}

func factorial(n uint64) (result uint64) {
	if n > 0 {
		result = n * factorial(n-1)
		return result
	}
	return 1
}

func writeJSON(w http.ResponseWriter, v any, status int) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
