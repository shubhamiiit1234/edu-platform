package routes

import (
	"net/http"

	"edu-learning-platform/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {

	// Health
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok!"}`))
	})

	// Auth
	r.Post("/signup", handlers.Signup)
	r.Post("/login", handlers.Login)

	// Subjects
	r.Get("/subjects", handlers.ListSubjects)

	// Lessons
	r.Get("/lessons", handlers.ListLessons)
	r.Get("/lessons/{id}", handlers.GetLesson)

	// Practice questions
	r.Get("/practice_questions", handlers.ListPracticeQuestions)

	// Aptitude
	r.Get("/aptitude/questions", handlers.ListAptitudeQuestions)
	r.Post("/aptitude/submit", handlers.SubmitAptitude)

	r.Get("/quizzes", handlers.ListQuizzes)
	r.Get("/quizzes/{quizId}/questions", handlers.ListQuizQuestions)
	r.Post("/quizzes/{quizId}/start", handlers.StartQuizAttempt)
	r.Post("/quizzes/{quizId}/submit", handlers.SubmitQuiz)

}
