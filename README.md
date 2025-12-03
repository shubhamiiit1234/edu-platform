Backend (Go + Chi + PostgreSQL)

This backend powers the Learning & Aptitude Platform.
It is built using Go (with Chi router), PostgreSQL as the database, and comes with Docker & Docker Compose support.

**1. The backend provides APIs for:**
   1. User Signup & Login
   2. Subjects
   3. Lessons
   4. Practice Questions
   5. Aptitude Questions
   6. Aptitude Test Submission & Recommendation

**2. Requirements:**
   1. Go 1.24+
   2. PostgreSQL 14+

**3. Database Tables:**
   1. Users
   2. Subjects
   3. Lessons
   4. Practice Questions
   5. Aptitude Questions
   6. Aptitude Results
   7. Badges
   8. User Badges
  
**4. Running/Rebuilding Backend With Docker Compose:**
   docker compose up --build

**5. API Testing:**
   1. Health check:
      GET /health
   
   2. Signup:
      POST /signup
      {
        "name": "Shubham",
        "email": "test@test.com",
        "password": "123456",
        "class": 10
      }
   
   3. Login:
      POST /login
      {
        "email": "test@test.com",
        "password": "123456"
      }
   
   4. Subjects:
      GET /subjects
      GET /subjects?class=6
   
   5. Lessons:
      GET /lessons
      GET /lessons?subject_id=2
      GET /lessons/5
   
   6. Aptitude test:
      GET /aptitude/questions
      POST /aptitude/submit

**6. Upcoming Features:**
   1. JWT Authentication
   2. Reward system
   3. Daily challenges
   4. Engineering–Medical–Arts recommendation engine
   5. Admin dashboard for adding lessons/questions
