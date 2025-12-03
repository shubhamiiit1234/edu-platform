package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"edu-learning-platform/internal/database"
	"edu-learning-platform/internal/models"
)

func ListSubjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	classStr := r.URL.Query().Get("class")

	db, err := database.GetDBInstance()
	if err != nil {
		http.Error(w, `{"error":"db not initialized"}`, http.StatusInternalServerError)
		return
	}

	var rows *sql.Rows

	if classStr != "" {
		class, err := strconv.Atoi(classStr)
		if err != nil {
			http.Error(w, `{"error":"invalid class"}`, http.StatusBadRequest)
			return
		}

		rows, err = db.Query(`
			SELECT id, name, class, created_at 
			FROM subjects 
			WHERE class = $1 
			ORDER BY name
		`, class)

		if err != nil {
			http.Error(w, `{"error":"db error: `+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

	} else {
		rows, err = db.Query(`
			SELECT id, name, class, created_at 
			FROM subjects 
			ORDER BY class, name
		`)
		if err != nil {
			http.Error(w, `{"error":"db error: `+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
	}

	defer rows.Close()

	var subjects []models.Subject

	for rows.Next() {
		var s models.Subject
		if err := rows.Scan(&s.ID, &s.Name, &s.Class, &s.CreatedAt); err != nil {
			http.Error(w, `{"error":"scan error: `+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		subjects = append(subjects, s)
	}

	if err := json.NewEncoder(w).Encode(subjects); err != nil {
		http.Error(w, `{"error":"encoding error"}`, http.StatusInternalServerError)
	}
}
