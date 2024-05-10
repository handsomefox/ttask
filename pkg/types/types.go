// types is a package with types that are shared in middleware and handler packages to avoid dependency cycles.
package types

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
