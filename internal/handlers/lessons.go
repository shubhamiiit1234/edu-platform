package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourname/edu-backend-starter/internal/database"
	"github.com/yourname/edu-backend-starter/internal/models"
)

func ListLessons(c *gin.Context) {
	subjectStr := c.Query("subject_id")

	db, err := database.GetDBInstance()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db not initialized"})
		return
	}

	var rows *sql.Rows
	if subjectStr != "" {
		sid, err := strconv.Atoi(subjectStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subject_id"})
			return
		}
		rows, err = db.Query(
			`SELECT id, subject_id, title, theory, examples, animation_url, created_at
			 FROM lessons WHERE subject_id=$1 ORDER BY id`,
			sid,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error: " + err.Error()})
			return
		}
	} else {
		rows, err = db.Query(
			`SELECT id, subject_id, title, theory, examples, animation_url, created_at
			 FROM lessons ORDER BY id`,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error: " + err.Error()})
			return
		}
	}
	defer rows.Close()

	var lessons []models.Lesson
	for rows.Next() {
		var l models.Lesson
		var examplesRaw []byte
		if err := rows.Scan(&l.ID, &l.SubjectID, &l.Title, &l.Theory, &examplesRaw, &l.AnimationURL, &l.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan error: " + err.Error()})
			return
		}
		if len(examplesRaw) > 0 {
			var ex any
			if err := json.Unmarshal(examplesRaw, &ex); err == nil {
				l.Examples = ex
			}
		}
		lessons = append(lessons, l)
	}

	c.JSON(http.StatusOK, lessons)
}

func GetLesson(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	db, err := database.GetDBInstance()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db not initialized"})
		return
	}

	row := db.QueryRow(
		`SELECT id, subject_id, title, theory, examples, animation_url, created_at
		 FROM lessons WHERE id=$1`, id,
	)

	var l models.Lesson
	var exRaw []byte
	if err := row.Scan(&l.ID, &l.SubjectID, &l.Title, &l.Theory, &exRaw, &l.AnimationURL, &l.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "lesson not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error: " + err.Error()})
		return
	}

	if len(exRaw) > 0 {
		var ex any
		if err := json.Unmarshal(exRaw, &ex); err == nil {
			l.Examples = ex
		}
	}

	c.JSON(http.StatusOK, l)
}
