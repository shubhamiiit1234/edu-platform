package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"edu-learning-platform/internal/database"
	"edu-learning-platform/internal/models"

	"github.com/go-chi/chi/v5"
)

func ListLessons(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	subjectStr := r.URL.Query().Get("subject_id")

	db, err := database.GetDBInstance()
	if err != nil {
		http.Error(w, `{"error":"db not initialized"}`, http.StatusInternalServerError)
		return
	}

	var rows *sql.Rows

	// Filter by subject_id if given
	if subjectStr != "" {
		sid, err := strconv.Atoi(subjectStr)
		if err != nil {
			http.Error(w, `{"error":"invalid subject_id"}`, http.StatusBadRequest)
			return
		}

		rows, err = db.Query(`
			SELECT id, subject_id, title, theory, examples, animation_url, created_at
			FROM lessons
			WHERE subject_id = $1
			ORDER BY id
		`, sid)

		if err != nil {
			http.Error(w, `{"error":"db error: `+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

	} else {
		// No filter → return all lessons
		rows, err = db.Query(`
			SELECT id, subject_id, title, theory, examples, animation_url, created_at
			FROM lessons
			ORDER BY id
		`)

		if err != nil {
			http.Error(w, `{"error":"db error: `+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
	}

	defer rows.Close()

	var lessons []models.Lesson

	for rows.Next() {
		var l models.Lesson
		var examplesRaw []byte

		if err := rows.Scan(
			&l.ID, &l.SubjectID, &l.Title, &l.Theory,
			&examplesRaw, &l.AnimationURL, &l.CreatedAt,
		); err != nil {
			http.Error(w, `{"error":"scan error: `+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		// Parse JSON examples
		if len(examplesRaw) > 0 {
			var ex any
			if json.Unmarshal(examplesRaw, &ex) == nil {
				l.Examples = ex
			}
		}

		lessons = append(lessons, l)
	}

	// Send final JSON response
	json.NewEncoder(w).Encode(lessons)
}

func GetLesson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	db, err := database.GetDBInstance()
	if err != nil {
		http.Error(w, `{"error":"db not initialized"}`, http.StatusInternalServerError)
		return
	}

	row := db.QueryRow(`
		SELECT id, subject_id, title, theory, examples, animation_url, created_at
		FROM lessons
		WHERE id = $1
	`, id)

	var l models.Lesson
	var exRaw []byte

	err = row.Scan(
		&l.ID, &l.SubjectID, &l.Title, &l.Theory,
		&exRaw, &l.AnimationURL, &l.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"lesson not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"db error: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	if len(exRaw) > 0 {
		var ex any
		if json.Unmarshal(exRaw, &ex) == nil {
			l.Examples = ex
		}
	}

	json.NewEncoder(w).Encode(l)
}
