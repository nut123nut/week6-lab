package store

import (
    "context"
    "errors"

    "github.com/jackc/pgx/v5/pgconn"
)

var ErrEmailTaken = errors.New("email taken")

type Member struct {
    ID           int64  `json:"id"`
    Email        string `json:"email"`
    DisplayName  string `json:"displayName"`
    Role         string `json:"role"`
    PasswordHash string `json:"-"` // ไม่มีวันออกไปกับ JSON
}

func (s *Store) CreateMember(ctx context.Context, email, hash, name string) (Member, error) {
    m := Member{Email: email, DisplayName: name}
    err := s.Pool.QueryRow(ctx,
        `INSERT INTO members (email, password_hash, display_name)
         VALUES ($1, $2, $3)
         RETURNING id, role`, email, hash, name).Scan(&m.ID, &m.Role)
    var pgErr *pgconn.PgError
    if errors.As(err, &pgErr) && pgErr.Code == "23505" {
        return Member{}, ErrEmailTaken
    }
    return m, err
}