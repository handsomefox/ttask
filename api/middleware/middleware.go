package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/handsomefox/ttask/internal/shared"
	"github.com/handsomefox/ttask/pkg/types"
	"github.com/julienschmidt/httprouter"
)

type ContextKey string

const (
	CalculateRequestKey ContextKey = "calculate-request-body"
)

var ErrCalculateInvalidInput = errors.New("middleware: invalid input in calculate")

// Calculate validates the request and sets the context.Value(CalculateRequestKey) to the CalculateRequest struct.
func Calculate(next httprouter.Handle) httprouter.Handle {
	// Create middleware
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var body types.CalculateRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			slog.Error("CalculateMiddleware#Decode", "err", err)
			shared.WriteIncorrectInputError(w)
			return
		}

		if err := validateCalculateRequest(body); err != nil {
			slog.Info("CalculateMiddleware#validateCalculateRequest", "err", err)
			shared.WriteIncorrectInputError(w)
			return
		}

		// Store the decoded value inside the context.
		// This is done to avoid decoding the body twice (which means the body needs to be copied, since you can't read it twice).
		ctx := context.WithValue(r.Context(), CalculateRequestKey, body)
		r = r.WithContext(ctx)

		next(w, r, p)
	}
}

func validateCalculateRequest(body types.CalculateRequest) error {
	// which will check if a and b exists
	if body.A == nil {
		return fmt.Errorf("%w: A is nil", ErrCalculateInvalidInput)
	}
	if body.B == nil {
		return fmt.Errorf("%w: B is nil", ErrCalculateInvalidInput)
	}
	// and they are non-negative int
	if *body.A < 0 {
		return fmt.Errorf("%w: A is a negative number", ErrCalculateInvalidInput)
	}
	if *body.B < 0 {
		return fmt.Errorf("%w: B is a negative number", ErrCalculateInvalidInput)
	}

	return nil
}
