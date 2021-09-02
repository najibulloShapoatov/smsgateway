package middleware

import (
	"context"
	"errors"
	"net/http"
	"smsc/pkg/log"
	"smsc/pkg/models"
	"strings"
)

//ErrNoAuthentication ...
var ErrNoAuthentication = errors.New("No authentication")
var authenticationContextKey = &contextKey{"authentication context"}

type contextKey struct {
	name string
}

func AuthAPI(nextHandler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		APIKey := r.URL.Query().Get("key")
		if len(strings.TrimSpace(APIKey)) != 64 {
			log.Warn(APIKey, len(strings.TrimSpace(APIKey)))
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		serviceID, err := CheckServiceAuth(r.Context(), APIKey, r)
		if err != nil {
			log.Warn("CheckServiceAuth(APIKey, r)", err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), authenticationContextKey, serviceID)
		r = r.WithContext(ctx)

		nextHandler.ServeHTTP(w, r)
		return
	})
}

//CheckServiceAuth Function
func CheckServiceAuth(ctx context.Context, APIKey string, r *http.Request) (int64, error) {
	customer, err := (&models.Customer{APIKey: APIKey}).Get(ctx)
	if err != nil {
		return 0, errors.New("service not found")
	}

	return customer.ID, nil
}

//Authentication returned authenticated serviceID
func Authentication(ctx context.Context) (int64, error) {
	if value, ok := ctx.Value(authenticationContextKey).(int64); ok {
		return value, nil
	}
	return 0, ErrNoAuthentication
}
