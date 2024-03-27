package dao

const (
	insertUser = `
		INSERT INTO users (name, email, hashed_password, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id`
)
