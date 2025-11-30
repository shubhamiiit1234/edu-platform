package handlers

// import (
//     "net/http"
//     "strconv"
//     "time"

//     "github.com/gin-gonic/gin"
//     "github.com/yourname/edu-backend-starter/internal/models"
// )

// // In-memory stores (prototype)
// var lessons = models.SeedLessons()
// var aptitudeQuestions = models.SeedAptitudeQuestions()
// var users = map[string]models.User{}

// func Health(c *gin.Context) {
//     c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now()})
// }

// // --- Auth (mocked) ---
// type signupReq struct {
//     Name  string `json:"name" binding:"required"`
//     Email string `json:"email" binding:"required,email"`
//     Class int    `json:"class" binding:"required"`
// }

// func Signup(c *gin.Context) {
//     var req signupReq
//     if err := c.ShouldBindJSON(&req); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }
//     // Mock: create user with email as key
//     users[req.Email] = models.User{
//         ID:    req.Email,
//         Name:  req.Name,
//         Email: req.Email,
//         Class: req.Class,
//     }
//     c.JSON(http.StatusOK, gin.H{"message": "user created (mock)", "user": users[req.Email]})
// }

// type loginReq struct {
//     Email string `json:"email" binding:"required,email"`
// }

// func Login(c *gin.Context) {
//     var req loginReq
//     if err := c.ShouldBindJSON(&req); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }
//     u, ok := users[req.Email]
//     if !ok {
//         c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found (mock). signup first"})
//         return
//     }
//     // Mock token
//     c.JSON(http.StatusOK, gin.H{"message": "login success (mock)", "token": "mock-token", "user": u})
// }

// // --- Lessons ---
// func ListLessons(c *gin.Context) {
//     c.JSON(http.StatusOK, lessons)
// }

// func GetLesson(c *gin.Context) {
//     idStr := c.Param("id")
//     id, err := strconv.Atoi(idStr)
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
//         return
//     }
//     for _, l := range lessons {
//         if l.ID == id {
//             c.JSON(http.StatusOK, l)
//             return
//         }
//     }
//     c.JSON(http.StatusNotFound, gin.H{"error": "lesson not found"})
// }

// // --- Aptitude ---
// func ListAptitudeQuestions(c *gin.Context) {
//     // Return without correct answers
//     out := []models.AptitudeQuestion{}
//     for _, q := range aptitudeQuestions {
//         qCopy := q
//         qCopy.CorrectOption = -1
//         out = append(out, qCopy)
//     }
//     c.JSON(http.StatusOK, out)
// }

// type AptitudeSubmitReq struct {
//     Email   string            `json:"email" binding:"required,email"`
//     Answers map[int]int       `json:"answers"` // question_id -> selected_option
// }

// func SubmitAptitude(c *gin.Context) {
//     var req AptitudeSubmitReq
//     if err := c.ShouldBindJSON(&req); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }
//     // Score calculation
//     scores := map[string]int{
//         "math":      0,
//         "science":   0,
//         "creativity": 0,
//         "memory":    0,
//     }
//     for _, q := range aptitudeQuestions {
//         sel, ok := req.Answers[q.ID]
//         if !ok {
//             continue
//         }
//         if sel == q.CorrectOption {
//             scores[q.Category] += q.Points
//         }
//     }

//     // Simple recommendation logic (prototype)
//     // Compare math vs science vs creativity to decide recommended stream
//     rec := "Balanced"
//     if scores["math"] > scores["science"] && scores["math"] > scores["creativity"] {
//         rec = "Engineering (PCM)"
//     } else if scores["science"] > scores["math"] && scores["science"] > scores["creativity"] {
//         rec = "Medical (PCB)"
//     } else if scores["creativity"] > scores["math"] && scores["creativity"] > scores["science"] {
//         rec = "Design/Arts"
//     }

//     result := models.AptitudeResult{
//         UserEmail:            req.Email,
//         MathScore:            scores["math"],
//         ScienceScore:         scores["science"],
//         CreativityScore:      scores["creativity"],
//         MemoryScore:          scores["memory"],
//         RecommendedStream:    rec,
//     }

//     c.JSON(http.StatusOK, gin.H{"result": result})
// }
