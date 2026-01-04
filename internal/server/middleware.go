package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/dormitory-life/gateway/internal/constants"
	rmodel "github.com/dormitory-life/gateway/internal/server/request_models"
	"github.com/golang-jwt/jwt/v5"
)

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("request started",
			slog.String("method", r.Method),
			slog.String("url", r.URL.Path),
			slog.Any("date", time.Now().UTC()))

		next.ServeHTTP(w, r)

		s.logger.Info("request completed",
			slog.String("method", r.Method),
			slog.String("url", r.URL.Path),
			slog.Any("date", time.Now().UTC()))
	})
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			s.logger.Error("error getting header", slog.String("error", ErrEmptyAuthHeader.Error()))
			writeErrorResponse(w, constants.ErrBadRequest, http.StatusUnauthorized, "empty request auth header")
			return
		}

		token, ok := strings.CutPrefix(header, "Bearer ")
		if !ok {
			s.logger.Error("error getting token string", slog.String("error", ErrInternal.Error()))
			writeErrorResponse(w, constants.ErrBadRequest, http.StatusUnauthorized, "empty request auth header")
			return
		}

		claims, err := s.validateToken(token)
		if err != nil {
			s.logger.Error("error validating token", slog.String("error", err.Error()))
			writeErrorResponse(w, err, http.StatusUnauthorized)
			return
		}

		r.Header.Set("X-User-ID", claims.UserId)
		r.Header.Set("X-Dormitory-ID", claims.DormitoryId)

		s.logger.Debug("extracted ids from token", slog.String("userId", claims.UserId), slog.String("dormitoryId", claims.DormitoryId))

		next.ServeHTTP(w, r)
	})
}

func (s *Server) validateToken(tokenString string) (*rmodel.Claims, error) {

	claims := &rmodel.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse token: %v", ErrInternal, err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("%w: invalid token", ErrInvalidToken)
	}

	if claims.Type != constants.AccessPrivateKey {
		return nil, fmt.Errorf("%w: invalid token type", ErrInvalidToken)
	}

	if claims.UserId == "" {
		return nil, fmt.Errorf("%w: user_id is required", ErrBadRequest)
	}

	if claims.DormitoryId == "" {
		return nil, fmt.Errorf("%w: dormitory_id is required", ErrBadRequest)
	}

	return claims, nil
}
