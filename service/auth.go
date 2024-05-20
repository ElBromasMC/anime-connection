package service

import (
	"alc/model/auth"
	"context"
	"net/http"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type Auth struct {
	db *pgxpool.Pool
}

func NewAuthService(db *pgxpool.Pool) Auth {
	return Auth{
		db: db,
	}
}

// User management

func (us Auth) GetUser(id uuid.UUID) (auth.User, error) {
	var user auth.User
	sql := `SELECT user_id, name, email, role, created_at FROM users WHERE user_id = $1`
	if err := us.db.QueryRow(context.Background(), sql, id).
		Scan(&user.Id, &user.Name, &user.Email, &user.Role, &user.CreatedAt); err != nil {
		return auth.User{}, echo.NewHTTPError(http.StatusNotFound, "Usuario no encontrado")
	}
	return user, nil
}

func (us Auth) InsertUser(u auth.User, hpass []byte) error {
	if _, err := us.db.Exec(context.Background(), `INSERT INTO users (name, email, hashed_password)
VALUES ($1, $2, $3)`, u.Name, u.Email, string(hpass)); err != nil {
		// TODO: Test and handle unique email condition
		return echo.NewHTTPError(http.StatusConflict, "Ya existe una cuenta con el email proporcionado")
	}
	return nil
}

func (us Auth) DeleteUser(id uuid.UUID) error {
	sql := `DELETE FROM users WHERE user_id = $1`
	c, err := us.db.Exec(context.Background(), sql, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	if c.RowsAffected() != 1 {
		return echo.NewHTTPError(http.StatusNotFound, "Usuario no encontrado")
	}
	return nil
}

func (us Auth) GetUserIdAndHpassByEmail(email string) (uuid.UUID, []byte, error) {
	var id uuid.UUID
	var hpass string

	if err := us.db.QueryRow(context.Background(), `SELECT user_id, hashed_password
FROM users WHERE email = $1`, email).Scan(&id, &hpass); err != nil {
		return uuid.UUID{}, []byte{}, echo.NewHTTPError(http.StatusUnauthorized, "Email no encontrado")
	}
	return id, []byte(hpass), nil
}

func (us Auth) GetRecorderUsers() ([]auth.User, error) {
	sql := `SELECT user_id, name, email, role, created_at FROM users WHERE role = $1`
	rows, err := us.db.Query(context.Background(), sql, auth.RecorderRole)
	if err != nil {
		return []auth.User{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (auth.User, error) {
		var user auth.User
		if err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Role, &user.CreatedAt); err != nil {
			return auth.User{}, err
		}
		return user, nil
	})
	if err != nil {
		return []auth.User{}, echo.NewHTTPError(http.StatusInternalServerError)
	}

	return users, nil
}

func (us Auth) InsertRecorderUser(u auth.User, hpass []byte) error {
	if _, err := us.db.Exec(context.Background(), `INSERT INTO users (name, email, hashed_password, role)
VALUES ($1, $2, $3, $4)`, u.Name, u.Email, string(hpass), auth.RecorderRole); err != nil {
		// TODO: Test and handle unique email condition
		return echo.NewHTTPError(http.StatusConflict, "Ya existe una cuenta con el email proporcionado")
	}
	return nil
}

// Session management

func (us Auth) GetUserBySession(sessionId uuid.UUID) (auth.User, error) {
	var u auth.User
	sql := `SELECT u.user_id, u.name, u.email, u.role
	FROM users AS u
	JOIN sessions AS s
	ON u.user_id = s.user_id
	WHERE s.session_id = $1`
	if err := us.db.QueryRow(context.Background(), sql, sessionId).Scan(&u.Id, &u.Name, &u.Email, &u.Role); err != nil {
		return auth.User{}, echo.NewHTTPError(http.StatusUnauthorized, "Sesión inválida")
	}
	return u, nil
}

func (us Auth) InsertSession(userId uuid.UUID) (uuid.UUID, error) {
	var session uuid.UUID
	if err := us.db.QueryRow(context.Background(), `INSERT INTO sessions (user_id)
VALUES ($1) RETURNING session_id`, userId).Scan(&session); err != nil {
		return uuid.UUID{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	return session, nil
}

func (us Auth) DeleteSession(id uuid.UUID) error {
	sql := `DELETE FROM sessions WHERE session_id = $1`
	c, err := us.db.Exec(context.Background(), sql, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	if c.RowsAffected() != 1 {
		return echo.NewHTTPError(http.StatusNotFound, "Sesión no encontrada")
	}
	return nil
}
