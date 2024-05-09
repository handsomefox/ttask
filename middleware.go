package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type ContextKey string

const (
	CalculateRequestKey ContextKey = "calculate-request-body"
)

var ErrCalculateInvalidInput = errors.New("middleware: invalid input in calculate")

// in case of failure return JSON { "error":"Incorrect input"} with error status code 400 Bad Request
func WriteIncorrectInputError(w http.ResponseWriter) {
	if err := writeJSON(w, map[string]string{"error": "Incorrect input"}, http.StatusBadRequest); err != nil {
		slog.Error("Error when writing incorrect input to stream", "err", err)
	}
}

// CalculateMiddleware validates the request and sets the context.Value(CalculateRequestKey) to the CalculateRequest struct.
func CalculateMiddleware(next httprouter.Handle) httprouter.Handle {
	// Create middleware
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var body CalculateRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			slog.Error("CalculateMiddleware#Decode", "err", err)
			WriteIncorrectInputError(w)
			return
		}

		if err := validateCalculateRequest(body); err != nil {
			slog.Info("CalculateMiddleware#validateCalculateRequest", "err", err)
			WriteIncorrectInputError(w)
			return
		}

		// Store the decoded value inside the context.
		// This is done to avoid decoding the body twice (which means the body needs to be copied, since you can't read it twice).
		ctx := context.WithValue(r.Context(), CalculateRequestKey, body)
		r = r.WithContext(ctx)

		next(w, r, p)
	}
}

func validateCalculateRequest(body CalculateRequest) error {
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
