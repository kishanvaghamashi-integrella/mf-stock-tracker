package util

import "context"

type contextKey string

const UserIDKey contextKey = "user_id"

// GetUserIDFromContext retrieves the user ID from the request context
func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	val := ctx.Value(UserIDKey)
	userID, ok := val.(int64)
	return userID, ok
}
