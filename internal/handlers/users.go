package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourname/edu-backend-starter/internal/database"
	"github.com/yourname/edu-backend-starter/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type signupRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Class    int    `json:"class" binding:"required"`
}

func Signup(c *gin.Context) {
	var req signupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, err := database.GetDBInstance()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db not initialized"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	var id int
	err = db.QueryRow(
		`INSERT INTO users (name, email, password_hash, class)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		req.Name, req.Email, string(hash), req.Class,
	).Scan(&id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not create user: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    id,
		"name":  req.Name,
		"email": req.Email,
		"class": req.Class,
	})
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, err := database.GetDBInstance()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db not initialized"})
		return
	}

	var u models.User
	err = db.QueryRow(
		`SELECT id, name, email, password_hash, class, created_at
		 FROM users WHERE email=$1`, req.Email,
	).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Class, &u.CreatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error: " + err.Error()})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// For now no JWT, just return user basic data
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    u.ID,
			"name":  u.Name,
			"email": u.Email,
			"class": u.Class,
		},
	})
}
