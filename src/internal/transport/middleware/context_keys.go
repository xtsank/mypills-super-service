package middleware

type contextKey string

const (
	ResponsePayloadKey = contextKey("response_payload")
	ResponseStatusKey  = contextKey("response_status")
)
