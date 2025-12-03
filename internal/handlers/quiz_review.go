package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"edu-learning-platform/internal/database"

	"github.com/go-chi/chi/v5"
)

func GetQuizReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	attemptIDStr := chi.URLParam(r, "attemptId")
	attemptID, err := strconv.Atoi(attemptIDStr)
	if err != nil {
		http.Error(w, `{"error":"invalid attempt id"}`, http.StatusBadRequest)
		return
	}

	db, err := database.GetDBInstance()
	if err != nil {
		http.Error(w, `{"error":"db not initialized"}`, http.StatusInternalServerError)
		return
	}

	// 1️⃣ Load attempt metadata
	var quizID, userID, score int
	err = db.QueryRow(`
		SELECT quiz_id, user_id, score
		FROM quiz_attempts
		WHERE id = $1
	`, attemptID).Scan(&quizID, &userID, &score)

	if err != nil {
		http.Error(w, `{"error":"attempt not found"}`, http.StatusNotFound)
		return
	}

	// 2️⃣ Join answers + questions
	rows, err := db.Query(`
		SELECT 
			qa.question_id,
			pq.question,
			pq.options,
			qa.selected_option,
			qa.correct_option,
			qa.is_correct
		FROM quiz_attempt_answers qa
		JOIN practice_questions pq ON pq.id = qa.question_id
		WHERE qa.attempt_id = $1
	`, attemptID)

	if err != nil {
		http.Error(w, `{"error":"db error while fetching answers"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type ReviewQuestion struct {
		QuestionID     int      `json:"question_id"`
		Question       string   `json:"question"`
		Options        []string `json:"options"`
		SelectedOption int      `json:"selected_option"`
		CorrectOption  int      `json:"correct_option"`
		IsCorrect      bool     `json:"is_correct"`
	}

	var reviewQuestions []ReviewQuestion

	for rows.Next() {
		var q ReviewQuestion
		var optionsRaw []byte

		if err := rows.Scan(
			&q.QuestionID,
			&q.Question,
			&optionsRaw,
			&q.SelectedOption,
			&q.CorrectOption,
			&q.IsCorrect,
		); err != nil {
			http.Error(w, `{"error":"scan error"}`, http.StatusInternalServerError)
			return
		}

		var opts []string
		_ = json.Unmarshal(optionsRaw, &opts)
		q.Options = opts

		// Add to list
		reviewQuestions = append(reviewQuestions, q)
	}

	// 3️⃣ Build final JSON response
	resp := map[string]interface{}{
		"attempt_id": attemptID,
		"quiz_id":    quizID,
		"user_id":    userID,
		"score":      score,
		"questions":  reviewQuestions,
	}

	json.NewEncoder(w).Encode(resp)
}
