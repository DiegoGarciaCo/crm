package handlers

import (
	"bufio"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/DiegoGarciaCo/CRM/internal/database"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/keighl/postmark"
	"github.com/sirupsen/logrus"
)

type apiCfg struct {
	Port             string
	JWTSecret        string
	DB               *database.Queries
	RawDB            *sql.DB
	dev              bool
	logger           *slog.Logger
	S3Client         *s3.Client
	S3Bucket         string
	S3Region         string
	postmarkClient   *postmark.Client
	EmailSecret      []byte
	betterAuthSecret string
	BaseURL          string
	FromEmail        string
}

func New(port, JWTSecret string, db *database.Queries, dbSQL *sql.DB, dev bool, logger *slog.Logger, s3Client *s3.Client, s3Bucket string, s3Region string, postmarkClient *postmark.Client, emailSecret []byte, betterAuthSecret string, baseURL string, fromEmail string) *apiCfg {
	return &apiCfg{
		Port:             port,
		JWTSecret:        JWTSecret,
		DB:               db,
		RawDB:            dbSQL,
		dev:              dev,
		logger:           logger,
		S3Client:         s3Client,
		S3Bucket:         s3Bucket,
		S3Region:         s3Region,
		postmarkClient:   postmarkClient,
		EmailSecret:      emailSecret,
		betterAuthSecret: betterAuthSecret,
		BaseURL:          baseURL,
		FromEmail:        fromEmail,
	}
}

// JSON response helpers
func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")

	// Skip body for no-content statuses
	if code == http.StatusNoContent || code == http.StatusAccepted {
		w.WriteHeader(code)
		return
	}

	// Marshal and write payload for other statuses
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	if _, err := w.Write(data); err != nil {
		log.Printf("Error writing response: %s", err)
		// Can't call WriteHeader again, response already sent
	}
}

// key type for context
type contextKey string

const userIDKey contextKey = "userID"

func VerifySignedCookie(cookieValue, secret string) (string, error) {
	parts := strings.Split(cookieValue, ".")
	if len(parts) != 2 {
		return "", errors.New("invalid cookie format")
	}

	payload := parts[0]
	signature := parts[1]

	// Step 1: URL-decode the signature (handles %2B, %3D, etc.)
	decodedSig, err := url.QueryUnescape(signature)
	if err != nil {
		decodedSig = signature // fallback to raw
	}

	// Step 2: Try standard Base64 decoding
	if sigBytes, err := base64.StdEncoding.DecodeString(decodedSig); err == nil {
		if verifyHMAC(payload, secret, sigBytes) {
			return payload, nil
		}
	}

	// Step 3: Try base64 URL encoding without padding
	if sigBytes, err := base64.RawURLEncoding.DecodeString(decodedSig); err == nil {
		if verifyHMAC(payload, secret, sigBytes) {
			return payload, nil
		}
	}

	return "", errors.New("invalid cookie signature")
}

func verifyHMAC(payload, secret string, signature []byte) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	expected := mac.Sum(nil)
	return hmac.Equal(signature, expected)
}

func HashAPIKey(raw string) string {
	sha := sha256.Sum256([]byte(raw))
	return base64.RawURLEncoding.EncodeToString(sha[:])
}

// --------------------------------------------------------------
// Authentication middleware
// --------------------------------------------------------------

func (cfg *apiCfg) AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/verify" || r.URL.Path == "/api/resend-verification" || r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			if strings.HasPrefix(r.URL.Path, "/webhooks/") {
				// Get X-API-Key header
				apiKey := r.Header.Get("X-API-Key")
				if apiKey == "" {
					respondWithError(w, http.StatusUnauthorized, "Missing API key", nil)
					return
				}

				// Hash the provided API key
				hashedKey := HashAPIKey(apiKey)

				// Check API key in database
				dbKey, err := cfg.DB.GetAPIKeyByHash(r.Context(), hashedKey)
				if err != nil {
					respondWithError(w, http.StatusUnauthorized, "Invalid API key", err)
					return
				}

				// Check if API key is active
				if dbKey.Enabled.Valid && !dbKey.Enabled.Bool {
					respondWithError(w, http.StatusUnauthorized, "Disabled API key", nil)
					return
				}

				// Check if API key expired
				if dbKey.ExpiresAt.Valid && dbKey.ExpiresAt.Time.Before(time.Now()) {
					respondWithError(w, http.StatusUnauthorized, "Expired API key", nil)
					return
				}

				// Add userID to request context
				ctx := context.WithValue(r.Context(), userIDKey, dbKey.UserId.String())
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			cookieName := "__Secure-crm.session_token"
			if cfg.dev == true {
				cookieName = "crm.session_token"
			}

			cookie, err := r.Cookie(cookieName)
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "Missing session cookie", err)
				return
			}

			// Verify the signed cookie
			token, err := VerifySignedCookie(cookie.Value, cfg.betterAuthSecret)
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "Invalid session cookie", err)
				return
			}

			// Query the Better Auth sessions table
			dbToken, err := cfg.DB.CheckSessionByID(r.Context(), token)
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "Invalid session", err)
				return
			}

			// Check if session expired
			if dbToken.ExpiresAt.Before(time.Now()) {
				respondWithError(w, http.StatusUnauthorized, "Session expired", nil)
				return
			}

			// Add userID to request context
			ctx := context.WithValue(r.Context(), userIDKey, dbToken.UserId.String())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// --------------------------------------------------------------
// Logger middleware
// --------------------------------------------------------------

const (
	requestIDKey contextKey = "requestID"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := rw.ResponseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, fmt.Errorf("underlying ResponseWriter does not implement http.Hijacker")
}

func (cfg *apiCfg) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrappedWriter := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Generate and set request ID before calling the next handler
		requestID := uuid.New().String()
		r = r.WithContext(context.WithValue(r.Context(), requestIDKey, requestID))

		// Call the next handler
		next.ServeHTTP(wrappedWriter, r)

		duration := time.Since(start)

		// Handle client IP with load balancer
		clientIP := r.Header.Get("X-Forwarded-For")
		if clientIP == "" {
			clientIP = r.Header.Get("X-Real-IP")
		}
		if clientIP == "" {
			clientIP = r.RemoteAddr
		}

		// Build log fields
		fields := logrus.Fields{
			"request_id": requestID,
			"method":     r.Method,
			"path":       r.URL.Path, // or r.URL.String() for full URL
			"status":     wrappedWriter.statusCode,
			"duration":   duration.Milliseconds(),
			"client_ip":  clientIP,
			"timestamp":  time.Now().UTC().Format(time.RFC3339),
		}

		// Log based on status code
		if wrappedWriter.statusCode >= 400 {
			logrus.WithFields(fields).Error("Request failed")
		} else {
			logrus.WithFields(fields).Info("Request processed")
		}
	})
}

// Helper to retrieve userID from context in handlers
func GetUserUUID(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(userIDKey).(string)
	if !ok {
		return uuid.Nil, errors.New("user ID not found in context")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, err
	}

	return userUUID, nil
}
