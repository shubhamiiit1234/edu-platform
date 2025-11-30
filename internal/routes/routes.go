package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yourname/edu-backend-starter/internal/handlers"
)

func RegisterRoutes(r *gin.Engine) {
	// Health
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Auth
	r.POST("/signup", handlers.Signup)
	r.POST("/login", handlers.Login)

	// Subjects
	r.GET("/subjects", handlers.ListSubjects)

	// Lessons
	r.GET("/lessons", handlers.ListLessons)
	r.GET("/lessons/:id", handlers.GetLesson)

	// Practice questions
	r.GET("/practice_questions", handlers.ListPracticeQuestions)

	// Aptitude
	r.GET("/aptitude/questions", handlers.ListAptitudeQuestions)
	r.POST("/aptitude/submit", handlers.SubmitAptitude)
}
