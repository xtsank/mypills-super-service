package middleware

type contextKey string

const (
	ResponsePayloadKey = contextKey("response_payload")
	ResponseStatusKey  = contextKey("response_status")
	UserIDKey          = contextKey("user_id")
	IsAdminKey         = contextKey("is_admin")
)
