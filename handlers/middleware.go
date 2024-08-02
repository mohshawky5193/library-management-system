package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"library_management_system/apperrors"
	"library_management_system/config/jsonconfig"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get(jsonconfig.AuthorizationHeader)
		tokenString = strings.TrimPrefix(tokenString, jsonconfig.Bearer)
		if tokenString == "" {
			panic(&apperrors.MalFormedTokenError{})
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrNoLocation
			}
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			panic(&apperrors.InvalidTokenError{})
		}

		ctx := context.WithValue(r.Context(), jsonconfig.UsernameContextKey, token.Claims.(jwt.MapClaims)[jsonconfig.UsernameClaimKey])
		ctx = context.WithValue(ctx, jsonconfig.RoleContextKey, token.Claims.(jwt.MapClaims)[jsonconfig.RoleClaimKey])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RoleMiddleware(requiredRole string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(jsonconfig.RoleContextKey).(string)
			if !ok || role != requiredRole {
				panic(&apperrors.RoleNotMatchingError{Role: role, RequiredRole: requiredRole})
			}
			next.ServeHTTP(w, r)
		})
	}
}

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				var apiError error
				switch e := err.(type) {
				case error:
					apiError = e
				default:
					apiError = errors.New("unexpected error occurred")
				}
				handleError(apiError, w)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func handleError(err error, w http.ResponseWriter) {
	w.Header().Set(jsonconfig.ContentType, jsonconfig.ApplicationJson)
	switch err.(type) {
	case *apperrors.AlreadyHaveBookError,
		*apperrors.AmountIsZeroError,
		*apperrors.BookNotBorrowedError,
		*apperrors.UsernameAlreadyExistsError,
		*apperrors.BookNotFoundError,
		*apperrors.BookValidationError,
		*apperrors.DeleteBorrowedBookError,
		*apperrors.BookWithSameIDError,
		*apperrors.CredentialsDecodingError:
		w.WriteHeader(http.StatusBadRequest)
	case *apperrors.UnauthorizedUserError,
		*apperrors.UnauthenticatedUserError,
		*apperrors.MalFormedTokenError,
		*apperrors.InvalidTokenError:
		w.WriteHeader(http.StatusUnauthorized)
	case *apperrors.RoleNotMatchingError:
		w.WriteHeader(http.StatusForbidden)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]string{jsonconfig.ErrorJsonKey: err.Error()})
}
