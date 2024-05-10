package handler

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/handsomefox/ttask/api/middleware"
	"github.com/handsomefox/ttask/internal/shared"
	"github.com/handsomefox/ttask/pkg/types"
	"github.com/julienschmidt/httprouter"
)

// Calculate expects the context to contain the value of the request. It does not read the request body.
func Calculate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, ok := r.Context().Value(middleware.CalculateRequestKey).(types.CalculateRequest)
	if !ok {
		slog.Warn("HandleCalculate must not be called before or without the CalculateMiddleware")
		shared.WriteIncorrectInputError(w)
		return
	}

	a, b := concurrentFactorial(uint64(*body.A), uint64(*body.B))
	if err := shared.WriteJSON(w, types.CalculateResponse{A: a, B: b}, http.StatusOK); err != nil {
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
