package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/keighl/postmark"
)

func splitList(s string) []string {
	if s == "" {
		return []string{}
	}
	items := strings.Split(s, ";")
	for i := range items {
		items[i] = strings.TrimSpace(items[i])
	}
	return items
}

func GetUUIDFromUrl(pathValue string, r *http.Request) (uuid.UUID, error) {
	// Get the ID from the URL Path value
	id := r.PathValue(pathValue)
	if id == "" {
		log.Printf("%s not found in URL", pathValue)
		return uuid.Nil, fmt.Errorf("%s not found in URL", pathValue)
	}

	// Parse the ID
	idUUID, err := uuid.Parse(id)
	if err != nil {
		log.Printf("Error parsing %s: %v", pathValue, err)
		return uuid.Nil, fmt.Errorf("invalid %s: %v", pathValue, err)
	}

	return idUUID, nil
}

func (cfg *apiCfg) GenerateEmailToken(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(30 * time.Minute).Unix(),
		"typ":   "email_verification",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(cfg.EmailSecret)
}

func (cfg *apiCfg) SendVerificationEmail(to, token string) error {
	verifyURL := cfg.BaseURL + "/api/verify?token=" + token

	_, err := cfg.postmarkClient.SendEmail(postmark.Email{
		From:     "diego.garcia@soldbyghost.com",
		To:       to,
		Subject:  "Verify your email",
		HtmlBody: fmt.Sprintf("<p>Click here to verify:</p><p><a href='%s'>Verify Email</a></p>", verifyURL),
		TextBody: "Verify your email: " + verifyURL,
	})

	return err
}

func MediaTypeToExt(mediaType string) string {
	parts := strings.Split(mediaType, "/")
	if len(parts) != 2 {
		return ".bin"
	}
	return "." + parts[1]
}

func GetAssetPath(mediaType, folder string) string {
	base := make([]byte, 32)
	_, err := rand.Read(base)
	if err != nil {
		panic("failed to generate random bytes")
	}
	id := base64.RawURLEncoding.EncodeToString(base)

	ext := MediaTypeToExt(mediaType)
	return fmt.Sprintf("%s/%s%s", folder, id, ext)
}

func CleanupS3(cfg *apiCfg, ctx context.Context, keys []string) {
	if len(keys) == 0 {
		return
	}

	for _, key := range keys {
		_, err := cfg.S3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(cfg.S3Bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			log.Printf("Failed to delete %s: %s", key, err)
		}
	}
}

func DeleteFromS3(cfg *apiCfg, ctx context.Context, imageURL string) error {
	baseURL := "https://" + cfg.S3Bucket + ".s3." + cfg.S3Region + ".amazonaws.com/"
	key := strings.TrimPrefix(imageURL, baseURL)

	if key == imageURL {
		log.Printf("Invalid S3 URL format: %s", imageURL)
		return fmt.Errorf("Invalid S3 URL format: %s", imageURL)
	}

	_, err := cfg.S3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(cfg.S3Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Printf("Failed to delete %s: %s", key, err)
	}

	return nil
}
