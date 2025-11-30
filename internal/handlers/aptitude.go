package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourname/edu-backend-starter/internal/database"
	"github.com/yourname/edu-backend-starter/internal/models"
)

func ListAptitudeQuestions(c *gin.Context) {
	db, err := database.GetDBInstance()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db not initialized"})
		return
	}

	rows, err := db.Query(
		`SELECT id, category, question, options, correct_option, points, created_at
		 FROM aptitude_questions ORDER BY id`,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error: " + err.Error()})
		return
	}
	defer rows.Close()

	var list []models.AptitudeQuestion
	for rows.Next() {
		var q models.AptitudeQuestion
		var optsRaw []byte
		if err := rows.Scan(&q.ID, &q.Category, &q.Question, &optsRaw, &q.CorrectOption, &q.Points, &q.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan error: " + err.Error()})
			return
		}
		var opts []string
		if len(optsRaw) > 0 {
			_ = json.Unmarshal(optsRaw, &opts)
		}
		q.Options = opts
		// hide correctOption from client
		q.CorrectOption = -1
		list = append(list, q)
	}

	c.JSON(http.StatusOK, list)
}

type AptitudeSubmitRequest struct {
	UserID  int         `json:"user_id" binding:"required"`
	Answers map[int]int `json:"answers" binding:"required"` // question_id -> selectedOptionIndex
}

func SubmitAptitude(c *gin.Context) {
	var req AptitudeSubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, err := database.GetDBInstance()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db not initialized"})
		return
	}

	// Load all questions (for now – fine for small set)
	rows, err := db.Query(
		`SELECT id, category, correct_option, points FROM aptitude_questions`,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error: " + err.Error()})
		return
	}
	defer rows.Close()

	scores := map[string]int{
		"math":       0,
		"science":    0,
		"creativity": 0,
		"memory":     0,
	}

	for rows.Next() {
		var id, correct, points int
		var category string
		if err := rows.Scan(&id, &category, &correct, &points); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan error: " + err.Error()})
			return
		}

		selected, ok := req.Answers[id]
		if !ok {
			continue
		}
		if selected == correct {
			scores[category] += points
		}
	}

	// Very simple recommendation
	rec := "General"
	if scores["math"] >= scores["science"] && scores["math"] >= scores["creativity"] {
		rec = "Engineering"
	} else if scores["science"] > scores["math"] && scores["science"] >= scores["creativity"] {
		rec = "Medical"
	} else if scores["creativity"] > scores["math"] && scores["creativity"] > scores["science"] {
		rec = "Arts / Design"
	}

	var resultID int
	err = db.QueryRow(
		`INSERT INTO aptitude_results
		 (user_id, math_score, science_score, creativity_score, memory_score, recommended_stream)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		req.UserID,
		scores["math"],
		scores["science"],
		scores["creativity"],
		scores["memory"],
		rec,
	).Scan(&resultID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save result: " + err.Error()})
		return
	}

	res := models.AptitudeResult{
		ID:                resultID,
		UserID:            req.UserID,
		MathScore:         scores["math"],
		ScienceScore:      scores["science"],
		CreativityScore:   scores["creativity"],
		MemoryScore:       scores["memory"],
		RecommendedStream: rec,
	}

	c.JSON(http.StatusOK, gin.H{"result": res})
}
