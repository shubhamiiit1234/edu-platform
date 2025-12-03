package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"edu-learning-platform/internal/database"
	"edu-learning-platform/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type signupRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Class    int    `json:"class" binding:"required"`
}

func Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req signupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	db, err := database.GetDBInstance()
	if err != nil {
		http.Error(w, `{"error":"db not initialized"}`, http.StatusInternalServerError)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"error":"failed to hash password"}`, http.StatusInternalServerError)
		return
	}

	var id int
	err = db.QueryRow(`
		INSERT INTO users (name, email, password_hash, class)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, req.Name, req.Email, string(hash), req.Class).Scan(&id)

	if err != nil {
		http.Error(w, `{"error":"could not create user: `+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	resp := map[string]any{
		"id":    id,
		"name":  req.Name,
		"email": req.Email,
		"class": req.Class,
	}

	json.NewEncoder(w).Encode(resp)
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	db, err := database.GetDBInstance()
	if err != nil {
		http.Error(w, `{"error":"db not initialized"}`, http.StatusInternalServerError)
		return
	}

	var u models.User
	err = db.QueryRow(`
		SELECT id, name, email, password_hash, class, created_at
		FROM users
		WHERE email = $1
	`, req.Email).Scan(
		&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Class, &u.CreatedAt,
	)

	if err == sql.ErrNoRows {
		http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, `{"error":"db error: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// Validate password
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)) != nil {
		http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	// Response
	resp := map[string]any{
		"user": map[string]any{
			"id":    u.ID,
			"name":  u.Name,
			"email": u.Email,
			"class": u.Class,
		},
	}

	json.NewEncoder(w).Encode(resp)
}
