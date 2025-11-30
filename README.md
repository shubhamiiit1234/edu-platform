# Go Backend Starter (Gin)

This is a minimal Go backend starter project for your education platform MVP.
It uses in-memory storage for quick prototyping and example handlers for:
- Signup / Login (mocked)
- Lessons (list + detail)
- Aptitude questions (list)
- Aptitude submission (simple scoring + recommendation)

**Notes**
- This is intentionally simple to help you start quickly.
- Replace in-memory stores with PostgreSQL and proper auth (JWT / OAuth) for production.
- See `db/schema.sql` for suggested database schema.

## How to run (locally)
1. Install Go (>=1.20 recommended).
2. From the project root:
   ```
   go mod tidy
   go run ./cmd/server
   ```
3. Server will run on `:8080`

## Endpoints
- `GET /health` → health check
- `POST /signup` → create user (mock)
- `POST /login` → login (mock)
- `GET /lessons` → list lessons
- `GET /lessons/:id` → lesson detail
- `GET /aptitude/questions` → aptitude questions
- `POST /aptitude/submit` → submit answers, get score & recommendation

## Next steps
- Add PostgreSQL integration (use `db/schema.sql`)
- Add JWT-based auth and middleware
- Add persistence for users, lessons, results
- Add tests
