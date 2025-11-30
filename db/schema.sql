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
