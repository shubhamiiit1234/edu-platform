package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Class        int       `json:"class"`
	CreatedAt    time.Time `json:"created_at"`
}

type Subject struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Class     int       `json:"class"`
	CreatedAt time.Time `json:"created_at"`
}

type Lesson struct {
	ID           int       `json:"id"`
	SubjectID    int       `json:"subject_id"`
	Title        string    `json:"title"`
	Theory       string    `json:"theory"`
	Examples     any       `json:"examples"`      // JSONB
	AnimationURL *string   `json:"animation_url"` // nullable
	CreatedAt    time.Time `json:"created_at"`
}

type PracticeQuestion struct {
	ID           int       `json:"id"`
	LessonID     int       `json:"lesson_id"`
	Question     string    `json:"question"`
	Options      []string  `json:"options"`
	CorrectIndex int       `json:"correct_option"`
	Difficulty   int       `json:"difficulty"`
	CreatedAt    time.Time `json:"created_at"`
}

type AptitudeQuestion struct {
	ID            int       `json:"id"`
	Category      string    `json:"category"`
	Question      string    `json:"question"`
	Options       []string  `json:"options"`
	CorrectOption int       `json:"-"` // never sent to client
	Points        int       `json:"points"`
	CreatedAt     time.Time `json:"created_at"`
}

type AptitudeResult struct {
	ID                int       `json:"id"`
	UserID            int       `json:"user_id"`
	MathScore         int       `json:"math_score"`
	ScienceScore      int       `json:"science_score"`
	CreativityScore   int       `json:"creativity_score"`
	MemoryScore       int       `json:"memory_score"`
	RecommendedStream string    `json:"recommended_stream"`
	CreatedAt         time.Time `json:"created_at"`
}

type Quiz struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	LessonID    int       `json:"lesson_id"`
	TotalMarks  int       `json:"total_marks"`
	TimeLimit   int       `json:"time_limit"`
	CreatedAt   time.Time `json:"created_at"`
}

type QuizQuestion struct {
	ID         int    `json:"id"`
	QuizID     int    `json:"quiz_id"`
	QuestionID int    `json:"question_id"`
	Question   string `json:"question,omitempty"`
}

type QuizAttempt struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	QuizID     int       `json:"quiz_id"`
	Score      int       `json:"score"`
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
}

type QuizAnswer struct {
	ID             int       `json:"id"`
	AttemptID      int       `json:"attempt_id"`
	QuestionID     int       `json:"question_id"`
	SelectedOption int       `json:"selected_option"`
	CorrectOption  int       `json:"correct_option"`
	IsCorrect      bool      `json:"is_correct"`
	CreatedAt      time.Time `json:"created_at"`
}

type SubmitQuizRequest struct {
	UserID  int         `json:"user_id"`
	Answers map[int]int `json:"answers"` // questionID → selectedOption
}
