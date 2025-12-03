package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"edu-learning-platform/internal/database"
	"edu-learning-platform/internal/models"

	"github.com/go-chi/chi/v5"
)

func ListQuizzes(w http.ResponseWriter, r *http.Request) {
	db, err := database.GetDBInstance()
	if err != nil {
		http.Error(w, `{"error":"DB not initialized"}`, http.StatusInternalServerError)
		return
	}

	rows, err := db.Query(`
		SELECT id, title, description, lesson_id, total_marks, time_limit, created_at
		FROM quizzes
		ORDER BY id
	`)
	if err != nil {
		http.Error(w, `{"error":"DB error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var list []models.Quiz

	for rows.Next() {
		var q models.Quiz
		if err := rows.Scan(&q.ID, &q.Title, &q.Description, &q.LessonID, &q.TotalMarks, &q.TimeLimit, &q.CreatedAt); err != nil {
			http.Error(w, `{"error":"scan error"}`, http.StatusInternalServerError)
			return
		}
		list = append(list, q)
	}

	json.NewEncoder(w).Encode(list)
}

func ListQuizQuestions(w http.ResponseWriter, r *http.Request) {
	quizIDStr := chi.URLParam(r, "quizId")
	quizID, err := strconv.Atoi(quizIDStr)
	if err != nil {
		http.Error(w, `{"error":"invalid quiz id"}`, http.StatusBadRequest)
		return
	}

	db, err := database.GetDBInstance()
	if err != nil {
		http.Error(w, `{"error":"DB not initialized"}`, http.StatusInternalServerError)
		return
	}

	rows, err := db.Query(`
		SELECT qq.id, qq.quiz_id, pq.id, pq.question
		FROM quiz_questions qq
		JOIN practice_questions pq ON pq.id = qq.question_id
		WHERE qq.quiz_id = $1
	`, quizID)

	if err != nil {
		http.Error(w, `{"error":"DB error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var list []models.QuizQuestion

	for rows.Next() {
		var q models.QuizQuestion
		if err := rows.Scan(&q.ID, &q.QuizID, &q.QuestionID, &q.Question); err != nil {
			http.Error(w, `{"error":"scan error"}`, http.StatusInternalServerError)
			return
		}
		list = append(list, q)
	}

	json.NewEncoder(w).Encode(list)
}

func StartQuizAttempt(w http.ResponseWriter, r *http.Request) {
	quizIDStr := chi.URLParam(r, "quizId")
	userIDStr := r.URL.Query().Get("user_id")

	quizID, err := strconv.Atoi(quizIDStr)
	userID, err2 := strconv.Atoi(userIDStr)

	if err != nil || err2 != nil {
		http.Error(w, `{"error":"invalid params"}`, http.StatusBadRequest)
		return
	}

	db, err := database.GetDBInstance()
	if err != nil {
		http.Error(w, `{"error":"DB not initialized"}`, http.StatusInternalServerError)
		return
	}

	var attemptID int
	err = db.QueryRow(`
		INSERT INTO quiz_attempts (user_id, quiz_id)
		VALUES ($1, $2) RETURNING id
	`, userID, quizID).Scan(&attemptID)

	if err != nil {
		http.Error(w, `{"error":"failed to start attempt"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"attempt_id": attemptID,
	})
}

func SubmitQuiz(w http.ResponseWriter, r *http.Request) {
	quizIDStr := chi.URLParam(r, "quizId")
	quizID, err := strconv.Atoi(quizIDStr)
	if err != nil {
		http.Error(w, `{"error":"invalid quiz id"}`, http.StatusBadRequest)
		return
	}

	var req models.SubmitQuizRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	db, err := database.GetDBInstance()
	if err != nil {
		http.Error(w, `{"error":"DB not initialized"}`, http.StatusInternalServerError)
		return
	}

	// Start attempt
	var attemptID int
	err = db.QueryRow(`
		INSERT INTO quiz_attempts (user_id, quiz_id)
		VALUES ($1, $2) RETURNING id
	`, req.UserID, quizID).Scan(&attemptID)

	if err != nil {
		http.Error(w, `{"error":"failed to create attempt"}`, http.StatusInternalServerError)
		return
	}

	// Load correct answers
	rows, err := db.Query(`
		SELECT pq.id, pq.correct_option
		FROM quiz_questions qq
		JOIN practice_questions pq ON pq.id = qq.question_id
		WHERE qq.quiz_id = $1
	`, quizID)

	if err != nil {
		http.Error(w, `{"error":"DB error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	score := 0

	for rows.Next() {
		var qid, correct int
		if err := rows.Scan(&qid, &correct); err != nil {
			http.Error(w, `{"error":"scan error"}`, http.StatusInternalServerError)
			return
		}

		selected := req.Answers[qid]
		isCorrect := selected == correct
		if isCorrect {
			score++
		}

		db.Exec(`
			INSERT INTO quiz_attempt_answers (attempt_id, question_id, selected_option, correct_option, is_correct)
			VALUES ($1, $2, $3, $4, $5)
		`, attemptID, qid, selected, correct, isCorrect)
	}

	// Update score
	db.Exec(`UPDATE quiz_attempts SET score=$1, finished_at=NOW() WHERE id=$2`, score, attemptID)

	json.NewEncoder(w).Encode(map[string]any{
		"attempt_id": attemptID,
		"score":      score,
	})
}
