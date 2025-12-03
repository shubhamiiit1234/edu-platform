package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"edu-learning-platform/internal/database"
	"edu-learning-platform/internal/models"
)

func ListPracticeQuestions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	lessonStr := r.URL.Query().Get("lesson_id")
	if lessonStr == "" {
		http.Error(w, `{"error":"lesson_id is required"}`, http.StatusBadRequest)
		return
	}

	lessonID, err := strconv.Atoi(lessonStr)
	if err != nil {
		http.Error(w, `{"error":"invalid lesson_id"}`, http.StatusBadRequest)
		return
	}

	db, err := database.GetDBInstance()
	if err != nil {
		http.Error(w, `{"error":"db not initialized"}`, http.StatusInternalServerError)
		return
	}

	rows, err := db.Query(`
		SELECT id, lesson_id, question, options, correct_option, difficulty, created_at
		FROM practice_questions
		WHERE lesson_id = $1
		ORDER BY id
	`, lessonID)

	if err != nil {
		http.Error(w, `{"error":"db error: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var qs []models.PracticeQuestion

	for rows.Next() {
		var q models.PracticeQuestion
		var optionsRaw []byte

		if err := rows.Scan(
			&q.ID,
			&q.LessonID,
			&q.Question,
			&optionsRaw,
			&q.CorrectIndex,
			&q.Difficulty,
			&q.CreatedAt,
		); err != nil {
			http.Error(w, `{"error":"scan error: `+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		// Parse JSON options
		var options []string
		if len(optionsRaw) > 0 {
			_ = json.Unmarshal(optionsRaw, &options)
		}
		q.Options = options

		// 🚨 IMPORTANT:
		// Normally we should NOT send correct index to clients
		// but you kept it for testing, so I leave it as is.

		qs = append(qs, q)
	}

	json.NewEncoder(w).Encode(qs)
}
