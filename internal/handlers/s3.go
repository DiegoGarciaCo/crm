package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/DiegoGarciaCo/CRM/internal/database"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (cfg *apiCfg) UploadProfilePicture(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userUUID, err := GetUserUUID(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse multipart form with a max memory of 32MB
	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "File too big", http.StatusBadRequest)
		return
	}

	// Ensure exactly one file is uploaded
	files := r.MultipartForm.File["file"]
	if len(files) == 0 || len(files) > 1 {
		http.Error(w, "Exactly one image file is required", http.StatusBadRequest)
		return
	}
	fileHeader := files[0]

	file, err := fileHeader.Open()
	if err != nil {
		http.Error(w, "Could not open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	contentType := fileHeader.Header.Get("Content-Type")

	// Validate that the file is an image
	if !strings.HasPrefix(contentType, "image/") {
		http.Error(w, "File is not an image", http.StatusBadRequest)
		return
	}

	tx, err := cfg.RawDB.BeginTx(r.Context(), nil)
	if err != nil {
		http.Error(w, "Could not begin transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	qtx := cfg.DB.WithTx(tx)

	// Check if user has an existing profile picture
	image, err := qtx.GetUserImage(r.Context(), userUUID)
	if err != sql.ErrNoRows && err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not check existing image", err)
		return
	}
	// If so, delete it from S3
	if image.Valid {
		err = DeleteFromS3(cfg, r.Context(), image.String)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Could not delete existing image", err)
			return
		}
	}

	// Upload the new image to S3
	key := GetAssetPath(contentType, "thumbnails")

	_, err = cfg.S3Client.PutObject(r.Context(), &s3.PutObjectInput{
		Bucket:      aws.String(cfg.S3Bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not upload image to S3", err)
		return
	}

	imageURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", cfg.S3Bucket, cfg.S3Region, key)

	// Update user's profile picture URL in the database
	err = qtx.UpdateUserImage(r.Context(), database.UpdateUserImageParams{
		ID:    userUUID,
		Image: sql.NullString{String: imageURL, Valid: true},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not update user image", err)
		return
	}

	err = tx.Commit()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not commit transaction", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
