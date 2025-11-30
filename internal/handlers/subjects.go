package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourname/edu-backend-starter/internal/database"
	"github.com/yourname/edu-backend-starter/internal/models"
)

func ListSubjects(c *gin.Context) {
	classStr := c.Query("class")

	db, err := database.GetDBInstance()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db not initialized"})
		return
	}

	var rows *sql.Rows

	if classStr != "" {
		class, err := strconv.Atoi(classStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid class"})
			return
		}
		rows, err = db.Query(
			`SELECT id, name, class, created_at FROM subjects WHERE class=$1 ORDER BY name`,
			class,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error: " + err.Error()})
			return
		}
	} else {
		rows, err = db.Query(
			`SELECT id, name, class, created_at FROM subjects ORDER BY class, name`,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error: " + err.Error()})
			return
		}
	}
	defer rows.Close()

	var subjects []models.Subject
	for rows.Next() {
		var s models.Subject
		if err := rows.Scan(&s.ID, &s.Name, &s.Class, &s.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan error: " + err.Error()})
			return
		}
		subjects = append(subjects, s)
	}

	c.JSON(http.StatusOK, subjects)
}
