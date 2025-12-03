-- Suggested PostgreSQL schema (starter)

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255),
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash TEXT,
  class INT,
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE subjects (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100),
  class INT
);

CREATE TABLE lessons (
  id SERIAL PRIMARY KEY,
  subject_id INT REFERENCES subjects(id),
  title VARCHAR(255),
  theory TEXT,
  examples JSONB,
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE practice_questions (
  id SERIAL PRIMARY KEY,
  lesson_id INT REFERENCES lessons(id),
  question TEXT,
  options JSONB,
  correct_option INT,
  difficulty INT
);

CREATE TABLE aptitude_questions (
  id SERIAL PRIMARY KEY,
  category VARCHAR(50),
  question TEXT,
  options JSONB,
  correct_option INT,
  points INT
);

CREATE TABLE aptitude_results (
  id SERIAL PRIMARY KEY,
  user_id INT REFERENCES users(id),
  math_score INT,
  science_score INT,
  creativity_score INT,
  memory_score INT,
  recommended_stream VARCHAR(100),
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE badges (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100),
  description TEXT
);

CREATE TABLE user_badges (
  user_id INT REFERENCES users(id) ON DELETE CASCADE,
  badge_id INT REFERENCES badges(id) ON DELETE CASCADE,
  earned_at TIMESTAMP DEFAULT NOW(),
  PRIMARY KEY (user_id, badge_id)
);


CREATE TABLE quizzes (
  id SERIAL PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  lesson_id INT REFERENCES lessons(id),
  total_marks INT DEFAULT 0,
  time_limit INT, -- in seconds (e.g., 600 = 10 mins)
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE quiz_questions (
  id SERIAL PRIMARY KEY,
  quiz_id INT REFERENCES quizzes(id) ON DELETE CASCADE,
  question_id INT REFERENCES practice_questions(id) ON DELETE CASCADE
);

CREATE TABLE quiz_attempts (
  id SERIAL PRIMARY KEY,
  user_id INT REFERENCES users(id),
  quiz_id INT REFERENCES quizzes(id),
  score INT DEFAULT 0,
  started_at TIMESTAMP DEFAULT NOW(),
  finished_at TIMESTAMP
);

CREATE TABLE quiz_answers (
  id SERIAL PRIMARY KEY,
  attempt_id INT REFERENCES quiz_attempts(id) ON DELETE CASCADE,
  question_id INT REFERENCES practice_questions(id),
  selected_option INT,
  correct_option INT,
  is_correct BOOLEAN,
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE user_progress (
  user_id INT REFERENCES users(id) ON DELETE CASCADE,
  lesson_id INT REFERENCES lessons(id) ON DELETE CASCADE,
  completed BOOLEAN DEFAULT FALSE,
  completed_at TIMESTAMP,
  PRIMARY KEY (user_id, lesson_id)
);