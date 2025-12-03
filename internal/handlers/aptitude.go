package handlers

import (
	"encoding/json"
	"net/http"

	"edu-learning-platform/internal/database"
	"edu-learning-platform/internal/models"
)

func ListAptitudeQuestions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, err := database.GetDBInstance()
	if err != nil {
		http.Error(w, `{"error":"db not initialized"}`, http.StatusInternalServerError)
		return
	}

	rows, err := db.Query(`
		SELECT id, category, question, options, correct_option, points, created_at
		FROM aptitude_questions ORDER BY id
	`)
	if err != nil {
		http.Error(w, `{"error":"db error: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var list []models.AptitudeQuestion

	for rows.Next() {
		var q models.AptitudeQuestion
		var optsRaw []byte

		if err := rows.Scan(
			&q.ID,
			&q.Category,
			&q.Question,
			&optsRaw,
			&q.CorrectOption,
			&q.Points,
			&q.CreatedAt,
		); err != nil {
			http.Error(w, `{"error":"scan error: `+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		// Parse JSONB options
		var options []string
		if len(optsRaw) > 0 {
			_ = json.Unmarshal(optsRaw, &options)
		}
		q.Options = options

		// Hide correct option
		q.CorrectOption = -1

		list = append(list, q)
	}

	json.NewEncoder(w).Encode(list)
}

type AptitudeSubmitRequest struct {
	UserID  int         `json:"user_id" binding:"required"`
	Answers map[int]int `json:"answers" binding:"required"` // question_id -> selectedOptionIndex
}

func SubmitAptitude(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req AptitudeSubmitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	db, err := database.GetDBInstance()
	if err != nil {
		http.Error(w, `{"error":"db not initialized"}`, http.StatusInternalServerError)
		return
	}

	// Load ALL aptitude questions (fine for small dataset)
	rows, err := db.Query(`
		SELECT id, category, correct_option, points 
		FROM aptitude_questions
	`)
	if err != nil {
		http.Error(w, `{"error":"db error: `+err.Error()+`"}`, http.StatusInternalServerError)
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
			http.Error(w, `{"error":"scan error: `+err.Error()+`"}`, http.StatusInternalServerError)
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

	// Simple recommendation logic { TODO: Will be changed!! }
	rec := "General"
	if scores["math"] >= scores["science"] && scores["math"] >= scores["creativity"] {
		rec = "Engineering"
	} else if scores["science"] > scores["math"] && scores["science"] >= scores["creativity"] {
		rec = "Medical"
	} else if scores["creativity"] > scores["math"] && scores["creativity"] > scores["science"] {
		rec = "Arts / Design"
	}

	var resultID int
	err = db.QueryRow(`
		INSERT INTO aptitude_results
			(user_id, math_score, science_score, creativity_score, memory_score, recommended_stream)
		VALUES
			($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, req.UserID, scores["math"], scores["science"], scores["creativity"], scores["memory"], rec).
		Scan(&resultID)

	if err != nil {
		http.Error(w, `{"error":"failed to save result: `+err.Error()+`"}`, http.StatusInternalServerError)
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

	json.NewEncoder(w).Encode(map[string]any{
		"result": res,
	})
}
