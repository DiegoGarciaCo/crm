package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/DiegoGarciaCo/CRM/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiCfg) CreateContact(w http.ResponseWriter, r *http.Request) {
	type requestContact struct {
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		Birthdate    string `json:"birthdate"`
		PhoneNumbers []struct {
			Number    string `json:"number"`
			Type      string `json:"type"`
			IsPrimary bool   `json:"is_primary"`
		} `json:"phone_numbers"`
		Emails []struct {
			Email     string `json:"email"`
			Type      string `json:"type"`
			IsPrimary bool   `json:"is_primary"`
		} `json:"emails"`
		Source     string `json:"source"`
		Status     string `json:"status"`
		Address    string `json:"address"`
		City       string `json:"city"`
		State      string `json:"state"`
		Zipcode    string `json:"zipCode"`
		Lender     string `json:"lender"`
		PriceRange string `json:"price_range"`
		Timeframe  string `json:"timeframe"`
	}

	var newContact requestContact
	err := json.NewDecoder(r.Body).Decode(&newContact)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Get ownerUUID from context
	ownerUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No User ID in Context", err)
		return
	}

	// Parse the birthdate string to time.Time
	parsedDate, err := time.Parse("2006-01-02", newContact.Birthdate)
	if err != nil && newContact.Birthdate != "" {
		respondWithError(w, http.StatusBadRequest, "Invalid birthdate format. Use YYYY-MM-DD.", err)
		return
	}

	// Start DB transaction to create contact
	tx, err := cfg.RawDB.BeginTx(r.Context(), nil)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to start database transaction", err)
		return
	}
	defer tx.Rollback()

	qtx := cfg.DB.WithTx(tx)

	// Create contact
	contact, err := qtx.CreateContact(r.Context(), database.CreateContactParams{
		FirstName:  newContact.FirstName,
		LastName:   newContact.LastName,
		Birthdate:  sql.NullTime{Time: parsedDate, Valid: newContact.Birthdate != ""},
		Source:     sql.NullString{String: newContact.Source, Valid: newContact.Source != ""},
		Status:     sql.NullString{String: newContact.Status, Valid: newContact.Status != ""},
		Address:    sql.NullString{String: newContact.Address, Valid: newContact.Address != ""},
		City:       sql.NullString{String: newContact.City, Valid: newContact.City != ""},
		State:      sql.NullString{String: newContact.State, Valid: newContact.State != ""},
		ZipCode:    sql.NullString{String: newContact.Zipcode, Valid: newContact.Zipcode != ""},
		Lender:     sql.NullString{String: newContact.Lender, Valid: newContact.Lender != ""},
		PriceRange: sql.NullString{String: newContact.PriceRange, Valid: newContact.PriceRange != ""},
		Timeframe:  sql.NullString{String: newContact.Timeframe, Valid: newContact.Timeframe != ""},
		OwnerID:    uuid.NullUUID{UUID: ownerUUID, Valid: ownerUUID != uuid.Nil},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create contact", err)
		return
	}

	// Insert phone numbers
	for _, phone := range newContact.PhoneNumbers {
		err = qtx.EnterPhoneNumber(r.Context(), database.EnterPhoneNumberParams{
			ContactID:   uuid.NullUUID{UUID: contact.ID, Valid: contact.ID != uuid.Nil},
			PhoneNumber: phone.Number,
			Type:        sql.NullString{String: phone.Type, Valid: phone.Type != ""},
			IsPrimary:   sql.NullBool{Bool: phone.IsPrimary, Valid: true},
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to add phone number", err)
			return
		}
	}

	// Insert emails
	for _, email := range newContact.Emails {
		err = qtx.EnterEmail(r.Context(), database.EnterEmailParams{
			ContactID:    uuid.NullUUID{UUID: contact.ID, Valid: contact.ID != uuid.Nil},
			EmailAddress: email.Email,
			Type:         sql.NullString{String: email.Type, Valid: email.Type != ""},
			IsPrimary:    sql.NullBool{Bool: email.IsPrimary, Valid: true},
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to add email", err)
			return
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to commit transaction", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, nil)
}

func (cfg *apiCfg) GetContactByID(w http.ResponseWriter, r *http.Request) {
	contactUUID, err := GetUUIDFromUrl("contactID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid contact ID", err)
		return
	}

	contact, err := cfg.DB.GetContactWithDetails(r.Context(), contactUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get contact", err)
		return
	}

	respondWithJSON(w, http.StatusOK, contact)
}

func (cfg *apiCfg) GetAllContacts(w http.ResponseWriter, r *http.Request) {
	ownerUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid owner ID", err)
		return
	}

	// Get query parameters from url for pagination (limit and offset)
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 50
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	contacts, err := cfg.DB.GetAllContacts(r.Context(), database.GetAllContactsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
		OwnerID: uuid.NullUUID{
			UUID:  ownerUUID,
			Valid: true,
		},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get contacts", err)
		return
	}

	if len(contacts) == 0 {
		contacts = []database.GetAllContactsRow{}
		respondWithJSON(w, http.StatusOK, contacts)
		return
	}

	respondWithJSON(w, http.StatusOK, contacts)
}

func (cfg *apiCfg) SearchContacts(w http.ResponseWriter, r *http.Request) {
	ownerUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid owner ID", err)
		return
	}

	// Get search query from URL
	query := r.URL.Query().Get("q")
	if query == "" {
		respondWithError(w, http.StatusBadRequest, "Search query cannot be empty", nil)
		return
	}

	contacts, err := cfg.DB.SearchContacts(r.Context(), database.SearchContactsParams{
		OwnerID:   uuid.NullUUID{UUID: ownerUUID, Valid: true},
		FirstName: query,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to search contacts", err)
		return
	}

	respondWithJSON(w, http.StatusOK, contacts)
}

func (cfg *apiCfg) GetContactsBySmartList(w http.ResponseWriter, r *http.Request) {
	// Get smartListID from url
	smartListUUID, err := GetUUIDFromUrl("smartListID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid smart list ID", err)
		return
	}

	// Get User UUID from context
	userUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No User ID in Context", err)
		return
	}

	// Get Limit and Offset from URL query parameters
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 50
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	contacts, err := cfg.DB.GetContactsBySmartList(r.Context(), database.GetContactsBySmartListParams{
		ID:     smartListUUID,
		Limit:  int32(limit),
		Offset: int32(offset),
		OwnerID: uuid.NullUUID{
			UUID:  userUUID,
			Valid: true,
		},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get contacts by smart list", err)
		return
	}

	if len(contacts) == 0 {
		contacts = []database.GetContactsBySmartListRow{}
		respondWithJSON(w, http.StatusOK, contacts)
		return
	}

	respondWithJSON(w, http.StatusOK, contacts)
}

func (cfg *apiCfg) ImportContacts(w http.ResponseWriter, r *http.Request) {
	cfg.logger.Info("ImportContacts endpoint called")
	type requestContact struct {
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		Birthdate    string `json:"birthdate"`
		PhoneNumbers []struct {
			Number    string `json:"number"`
			Type      string `json:"type"`
			IsPrimary bool   `json:"is_primary"`
		} `json:"phone_numbers"`
		Emails []struct {
			Email     string `json:"email"`
			Type      string `json:"type"`
			IsPrimary bool   `json:"is_primary"`
		} `json:"emails"`
		Source     string   `json:"source"`
		Status     string   `json:"status"`
		Address    string   `json:"address"`
		City       string   `json:"city"`
		State      string   `json:"state"`
		Zipcode    string   `json:"zipCode"`
		Lender     string   `json:"lender"`
		PriceRange string   `json:"price_range"`
		Timeframe  string   `json:"timeframe"`
		Tags       []string `json:"tags"`
	}

	var newContacts []requestContact
	err := json.NewDecoder(r.Body).Decode(&newContacts)
	if err != nil {
		cfg.logger.Error("Failed to decode import contacts payload", "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Get ownerUUID from context
	ownerUUID, err := GetUserUUID(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No User ID in Context", err)
		return
	}

	// Start DB transaction to create contacts
	tx, err := cfg.RawDB.BeginTx(r.Context(), nil)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to start database transaction", err)
		return
	}
	defer tx.Rollback()

	qtx := cfg.DB.WithTx(tx)

	for _, newContact := range newContacts {
		// Collect all string fields in a map
		fields := map[string]string{
			"FirstName":  newContact.FirstName,
			"LastName":   newContact.LastName,
			"Source":     newContact.Source,
			"Status":     newContact.Status,
			"Address":    newContact.Address,
			"City":       newContact.City,
			"State":      newContact.State,
			"Zipcode":    newContact.Zipcode,
			"Lender":     newContact.Lender,
			"PriceRange": newContact.PriceRange,
			"Timeframe":  newContact.Timeframe,
		}

		// Filter fields >= 90 characters
		longFields := map[string]string{}
		for k, v := range fields {
			if len(v) >= 90 {
				longFields[k] = v
			}
		}

		if len(longFields) > 0 {
			cfg.logger.Warn("Contact has long fields (>=90 chars)", "contact", longFields)
		}
		// Parse the birthdate string to time.Time
		parsedDate, err := time.Parse("2006-01-02", newContact.Birthdate)
		if err != nil && newContact.Birthdate != "" {
			cfg.logger.Error("Invalid birthdate format", "error", err)
			respondWithError(w, http.StatusBadRequest, "Invalid birthdate format. Use YYYY-MM-DD.", err)
			return
		}

		// Create contact
		contact, err := qtx.CreateContact(r.Context(), database.CreateContactParams{
			FirstName:  newContact.FirstName,
			LastName:   newContact.LastName,
			Birthdate:  sql.NullTime{Time: parsedDate, Valid: newContact.Birthdate != ""},
			Source:     sql.NullString{String: newContact.Source, Valid: newContact.Source != ""},
			Status:     sql.NullString{String: newContact.Status, Valid: newContact.Status != ""},
			Address:    sql.NullString{String: newContact.Address, Valid: newContact.Address != ""},
			City:       sql.NullString{String: newContact.City, Valid: newContact.City != ""},
			State:      sql.NullString{String: newContact.State, Valid: newContact.State != ""},
			ZipCode:    sql.NullString{String: newContact.Zipcode, Valid: newContact.Zipcode != ""},
			Lender:     sql.NullString{String: newContact.Lender, Valid: newContact.Lender != ""},
			PriceRange: sql.NullString{String: newContact.PriceRange, Valid: newContact.PriceRange != ""},
			Timeframe:  sql.NullString{String: newContact.Timeframe, Valid: newContact.Timeframe != ""},
			OwnerID:    uuid.NullUUID{UUID: ownerUUID, Valid: ownerUUID != uuid.Nil},
		})
		if err != nil {
			cfg.logger.Error("Failed to create contact", "error", err)
			respondWithError(w, http.StatusInternalServerError, "Failed to create contact", err)
			return
		}

		// Insert phone numbers
		for _, phone := range newContact.PhoneNumbers {
			err = qtx.EnterPhoneNumber(r.Context(), database.EnterPhoneNumberParams{
				ContactID:   uuid.NullUUID{UUID: contact.ID, Valid: contact.ID != uuid.Nil},
				PhoneNumber: phone.Number,
				Type:        sql.NullString{String: phone.Type, Valid: phone.Type != ""},
				IsPrimary:   sql.NullBool{Bool: phone.IsPrimary, Valid: true},
			})
			if err != nil {
				cfg.logger.Error("Failed to add phone number", "error", err)
				respondWithError(w, http.StatusInternalServerError, "Failed to add phone number", err)
				return
			}
		}

		// Insert emails
		for _, email := range newContact.Emails {
			err = qtx.EnterEmail(r.Context(), database.EnterEmailParams{
				ContactID:    uuid.NullUUID{UUID: contact.ID, Valid: contact.ID != uuid.Nil},
				EmailAddress: email.Email,
				Type:         sql.NullString{String: email.Type, Valid: email.Type != ""},
				IsPrimary:    sql.NullBool{Bool: email.IsPrimary, Valid: true},
			})
			if err != nil {
				cfg.logger.Error("Failed to add email", "error", err)
				respondWithError(w, http.StatusInternalServerError, "Failed to add email", err)
				return
			}
		}

		// Insert tags
		err = qtx.AssignTagsToContact(r.Context(), database.AssignTagsToContactParams{
			ContactID: contact.ID,
			UserID:    contact.OwnerID,
			Column2:   newContact.Tags,
		})
		if err != nil {
			cfg.logger.Error("Failed to assign tags to contact", "error", err)
			respondWithError(w, http.StatusInternalServerError, "Failed to assign tags to contact", err)
			return
		}

	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to commit transaction", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

// func (cfg *apiCfg) TestImportContacts(w http.ResponseWriter, r *http.Request) {
// 	cfg.logger.Info("ImportContacts endpoint called")
//
// 	type requestPhone struct {
// 		Number    string `json:"number"`
// 		Type      string `json:"type"`
// 		IsPrimary bool   `json:"is_primary"`
// 	}
// 	type requestEmail struct {
// 		Email     string `json:"email"`
// 		Type      string `json:"type"`
// 		IsPrimary bool   `json:"is_primary"`
// 	}
// 	type requestContact struct {
// 		FirstName    string         `json:"first_name"`
// 		LastName     string         `json:"last_name"`
// 		Birthdate    string         `json:"birthdate"`
// 		PhoneNumbers []requestPhone `json:"phone_numbers"`
// 		Emails       []requestEmail `json:"emails"`
// 		Source       string         `json:"source"`
// 		Status       string         `json:"status"`
// 		Address      string         `json:"address"`
// 		City         string         `json:"city"`
// 		State        string         `json:"state"`
// 		Zipcode      string         `json:"zipCode"`
// 		Lender       string         `json:"lender"`
// 		PriceRange   string         `json:"price_range"`
// 		Timeframe    string         `json:"timeframe"`
// 		Tags         []string       `json:"tags"`
// 	}
//
// 	var newContacts []requestContact
// 	if err := json.NewDecoder(r.Body).Decode(&newContacts); err != nil {
// 		cfg.logger.Error("Failed to decode import contacts payload", "error", err)
// 		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
// 		return
// 	}
//
// 	ownerUUID, err := GetUserUUID(r.Context())
// 	if err != nil {
// 		respondWithError(w, http.StatusUnauthorized, "No User ID in Context", err)
// 		return
// 	}
//
// 	tx, err := cfg.RawDB.BeginTx(r.Context(), nil)
// 	if err != nil {
// 		respondWithError(w, http.StatusInternalServerError, "Failed to start DB transaction", err)
// 		return
// 	}
// 	defer tx.Rollback()
// 	qtx := cfg.DB.WithTx(tx)
//
// 	// ---------------------------
// 	// 1. Prepare contacts JSON
// 	// ---------------------------
// 	type contactForDB struct {
// 		FirstName  string  `json:"first_name"`
// 		LastName   string  `json:"last_name"`
// 		Birthdate  *string `json:"birthdate"`
// 		Source     *string `json:"source"`
// 		Status     *string `json:"status"`
// 		Address    *string `json:"address"`
// 		City       *string `json:"city"`
// 		State      *string `json:"state"`
// 		ZipCode    *string `json:"zip_code"`
// 		Lender     *string `json:"lender"`
// 		PriceRange *string `json:"price_range"`
// 		Timeframe  *string `json:"timeframe"`
// 		OwnerID    string  `json:"owner_id"`
// 	}
//
// 	var contactsJSON []contactForDB
// 	for _, c := range newContacts {
// 		var birth *string
// 		if c.Birthdate != "" {
// 			birth = &c.Birthdate
// 		}
// 		toNullString := func(s string) *string {
// 			if s == "" {
// 				return nil
// 			}
// 			return &s
// 		}
// 		contactsJSON = append(contactsJSON, contactForDB{
// 			FirstName:  c.FirstName,
// 			LastName:   c.LastName,
// 			Birthdate:  birth,
// 			Source:     toNullString(c.Source),
// 			Status:     toNullString(c.Status),
// 			Address:    toNullString(c.Address),
// 			City:       toNullString(c.City),
// 			State:      toNullString(c.State),
// 			ZipCode:    toNullString(c.Zipcode),
// 			Lender:     toNullString(c.Lender),
// 			PriceRange: toNullString(c.PriceRange),
// 			Timeframe:  toNullString(c.Timeframe),
// 			OwnerID:    ownerUUID.String(),
// 		})
// 	}
//
// 	contactsJSONBytes, _ := json.Marshal(contactsJSON)
//
// 	// ---------------------------
// 	// 2. Bulk insert contacts
// 	// ---------------------------
// 	newContactIDs, err := qtx.TestBulkInsertContacts(r.Context(), contactsJSONBytes)
// 	if err != nil {
// 		cfg.logger.Error("Failed to bulk insert contacts", "error", err)
// 		respondWithError(w, http.StatusInternalServerError, "Failed to insert contacts", err)
// 		return
// 	}
//
// 	// ---------------------------
// 	// 3. Prepare phone numbers and emails
// 	// ---------------------------
// 	type phoneForDB struct {
// 		ContactID string  `json:"contact_id"`
// 		Number    string  `json:"number"`
// 		Type      *string `json:"type"`
// 		IsPrimary bool    `json:"is_primary"`
// 	}
// 	type emailForDB struct {
// 		ContactID    string  `json:"contact_id"`
// 		EmailAddress string  `json:"email_address"`
// 		Type         *string `json:"type"`
// 		IsPrimary    bool    `json:"is_primary"`
// 	}
// 	type tagForDB struct {
// 		ContactID string `json:"contact_id"`
// 		UserID    string `json:"user_id"`
// 		Tag       string `json:"tag"`
// 	}
//
// 	var allPhones []phoneForDB
// 	var allEmails []emailForDB
// 	var allTags []tagForDB
//
// 	for i, c := range newContacts {
// 		contactID := newContactIDs[i].ID.String()
// 		for _, p := range c.PhoneNumbers {
// 			toNull := func(s string) *string {
// 				if s == "" {
// 					return nil
// 				}
// 				return &s
// 			}
// 			allPhones = append(allPhones, phoneForDB{
// 				ContactID: contactID,
// 				Number:    p.Number,
// 				Type:      toNull(p.Type),
// 				IsPrimary: p.IsPrimary,
// 			})
// 		}
// 		for _, e := range c.Emails {
// 			toNull := func(s string) *string {
// 				if s == "" {
// 					return nil
// 				}
// 				return &s
// 			}
// 			allEmails = append(allEmails, emailForDB{
// 				ContactID:    contactID,
// 				EmailAddress: e.Email,
// 				Type:         toNull(e.Type),
// 				IsPrimary:    e.IsPrimary,
// 			})
// 		}
// 		for _, t := range c.Tags {
// 			allTags = append(allTags, tagForDB{
// 				ContactID: contactID,
// 				UserID:    ownerUUID.String(),
// 				Tag:       t,
// 			})
// 		}
// 	}
//
// 	phonesJSON, _ := json.Marshal(allPhones)
// 	emailsJSON, _ := json.Marshal(allEmails)
// 	tagsJSON, _ := json.Marshal(allTags)
//
// 	// ---------------------------
// 	// 4. Bulk insert phone numbers, emails, tags
// 	// ---------------------------
// 	if len(allPhones) > 0 {
// 		if err := qtx.TestBulkInsertPhoneNumbers(r.Context(), phonesJSON); err != nil {
// 			cfg.logger.Error("Failed to bulk insert phone numbers", "error", err)
// 			respondWithError(w, http.StatusInternalServerError, "Failed to insert phone numbers", err)
// 			return
// 		}
// 	}
// 	if len(allEmails) > 0 {
// 		if err := qtx.TestBulkInsertEmails(r.Context(), emailsJSON); err != nil {
// 			cfg.logger.Error("Failed to bulk insert emails", "error", err)
// 			respondWithError(w, http.StatusInternalServerError, "Failed to insert emails", err)
// 			return
// 		}
// 	}
// 	if len(allTags) > 0 {
// 		if err := qtx.TestBulkAssignTagsToContacts(r.Context(), tagsJSON); err != nil {
// 			cfg.logger.Error("Failed to bulk insert tags", "error", err)
// 			respondWithError(w, http.StatusInternalServerError, "Failed to insert tags", err)
// 			return
// 		}
// 	}
//
// 	if err := tx.Commit(); err != nil {
// 		respondWithError(w, http.StatusInternalServerError, "Failed to commit transaction", err)
// 		return
// 	}
//
// 	respondWithJSON(w, http.StatusNoContent, nil)
// }

func (cfg *apiCfg) UpdateContact(w http.ResponseWriter, r *http.Request) {
	contactUUID, err := GetUUIDFromUrl("contactID", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid contact ID", err)
		return
	}

	type req struct {
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Birthday   string `json:"birthdate"`
		Source     string `json:"source"`
		Status     string `json:"status"`
		Address    string `json:"address"`
		City       string `json:"city"`
		ZipCode    string `json:"zip_code"`
		Lender     string `json:"lender"`
		PriceRange string `json:"price_range"`
		Timeframe  string `json:"timeframe"`
	}

	var updatedData req
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Parse Birthdate
	parsedBirthDate, err := time.Parse("2006-01-02", updatedData.Birthday)
	if err != nil && updatedData.Birthday != "" {
		respondWithError(w, http.StatusBadRequest, "Invalid birthdate format. Use YYYY-MM-DD.", err)
		return
	}

	contact, err := cfg.DB.UpdateContact(r.Context(), database.UpdateContactParams{
		ID:         contactUUID,
		FirstName:  updatedData.FirstName,
		LastName:   updatedData.LastName,
		Birthdate:  sql.NullTime{Time: parsedBirthDate, Valid: updatedData.Birthday != ""},
		Source:     sql.NullString{String: updatedData.Source, Valid: updatedData.Source != ""},
		Status:     sql.NullString{String: updatedData.Status, Valid: updatedData.Status != ""},
		Address:    sql.NullString{String: updatedData.Address, Valid: updatedData.Address != ""},
		City:       sql.NullString{String: updatedData.City, Valid: updatedData.City != ""},
		ZipCode:    sql.NullString{String: updatedData.ZipCode, Valid: updatedData.ZipCode != ""},
		Lender:     sql.NullString{String: updatedData.Lender, Valid: updatedData.Lender != ""},
		PriceRange: sql.NullString{String: updatedData.PriceRange, Valid: updatedData.PriceRange != ""},
		Timeframe:  sql.NullString{String: updatedData.Timeframe, Valid: updatedData.Timeframe != ""},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not update Contact", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, contact)
}
