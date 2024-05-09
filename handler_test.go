package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func Test_CalculateHandler(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		h := CalculateMiddleware(HandleCalculate)
		h(w, r, httprouter.Params{})
	}

	tests := []struct {
		name         string
		request      *http.Request
		wantResponse string
		wantStatus   int
	}{
		{
			name:         "A is nil",
			request:      mustCreateRequest(t, http.MethodPost, "localhost:8989/calculate", CalculateRequest{A: nil, B: new(int)}),
			wantResponse: string(mustMarshalJSON(t, map[string]string{"error": "Incorrect input"})) + "\n",
			wantStatus:   http.StatusBadRequest,
		},
		{
			name:         "B is nil",
			request:      mustCreateRequest(t, http.MethodPost, "localhost:8989/calculate", CalculateRequest{A: new(int), B: nil}),
			wantResponse: string(mustMarshalJSON(t, map[string]string{"error": "Incorrect input"})) + "\n",
			wantStatus:   http.StatusBadRequest,
		},
		{
			name:         "A is -1",
			request:      mustCreateRequest(t, http.MethodPost, "localhost:8989/calculate", CalculateRequest{A: intPtrValue(t, -1), B: new(int)}),
			wantResponse: string(mustMarshalJSON(t, map[string]string{"error": "Incorrect input"})) + "\n",
			wantStatus:   http.StatusBadRequest,
		},
		{
			name:         "B is -1",
			request:      mustCreateRequest(t, http.MethodPost, "localhost:8989/calculate", CalculateRequest{A: new(int), B: intPtrValue(t, -1)}),
			wantResponse: string(mustMarshalJSON(t, map[string]string{"error": "Incorrect input"})) + "\n",
			wantStatus:   http.StatusBadRequest,
		},
		{
			name:         "A 2, B 3",
			request:      mustCreateRequest(t, http.MethodPost, "localhost:8989/calculate", CalculateRequest{A: intPtrValue(t, 2), B: intPtrValue(t, 3)}),
			wantResponse: string(mustMarshalJSON(t, CalculateResponse{A: 2, B: 6})) + "\n",
			wantStatus:   http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			handler(rr, tt.request)

			resp := rr.Result()
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("HandleCalculate() error = %v", err)
				return
			}

			bodyStr := string(body)
			if !reflect.DeepEqual(bodyStr, tt.wantResponse) {
				t.Errorf("HandleCalculate() = %v, want %v", bodyStr, tt.wantResponse)
			}

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("HandleCalculate() status = %d, want = %d", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}

func Test_factorial(t *testing.T) {
	tests := []struct {
		name       string
		n          uint64
		wantResult uint64
	}{
		{"Factorial of 0", 0, 1},
		{"Factorial of 1", 1, 1},
		{"Factorial of 2", 2, 2},
		{"Factorial of 3", 3, 6},
		{"Factorial of 4", 4, 24},
		{"Factorial of 5", 5, 120},
		{"Factorial of 10", 10, 3628800},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := factorial(tt.n); gotResult != tt.wantResult {
				t.Errorf("factorial() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func mustCreateRequest(t *testing.T, method, url string, body any) *http.Request {
	t.Helper()

	encoded, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
		return nil
	}

	ctx := context.TODO()
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(encoded))
	if err != nil {
		t.Fatal(err)
		return nil
	}

	return req
}

func mustMarshalJSON(t *testing.T, v any) []byte {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func intPtrValue(t *testing.T, v int) *int {
	t.Helper()
	i := new(int)
	*i = v
	return i
}
