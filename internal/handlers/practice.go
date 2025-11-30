package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourname/edu-backend-starter/internal/database"
	"github.com/yourname/edu-backend-starter/internal/models"
)

func ListPracticeQuestions(c *gin.Context) {
	lessonStr := c.Query("lesson_id")
	if lessonStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lesson_id is required"})
		return
	}
	lessonID, err := strconv.Atoi(lessonStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lesson_id"})
		return
	}

	db, err := database.GetDBInstance()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db not initialized"})
		return
	}

	rows, err := db.Query(
		`SELECT id, lesson_id, question, options, correct_option, difficulty, created_at
		 FROM practice_questions WHERE lesson_id=$1 ORDER BY id`,
		lessonID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error: " + err.Error()})
		return
	}
	defer rows.Close()

	var qs []models.PracticeQuestion
	for rows.Next() {
		var q models.PracticeQuestion
		var optionsRaw []byte
		if err := rows.Scan(&q.ID, &q.LessonID, &q.Question, &optionsRaw, &q.CorrectIndex, &q.Difficulty, &q.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan error: " + err.Error()})
			return
		}
		var opts []string
		if len(optionsRaw) > 0 {
			_ = json.Unmarshal(optionsRaw, &opts)
		}
		q.Options = opts

		// Do NOT send correctIndex to client normally; here we send it for now for testing
		qs = append(qs, q)
	}

	c.JSON(http.StatusOK, qs)
}
