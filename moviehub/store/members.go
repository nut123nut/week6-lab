import (
    "context"
    "errors"
    "time"

    "github.com/jackc/pgx/v5"

    "github.com/jackc/pgx/v5/pgconn"
)

type Member struct {
    ID           int64      `json:"id"`
    Email        string     `json:"email"`
    DisplayName  string     `json:"displayName"`
    Role         string     `json:"role"`
    CreatedAt    *time.Time `json:"createdAt,omitempty"` // มีเฉพาะตอน list ให้ admin ดู
    PasswordHash string     `json:"-"` // ไม่มีวันออกไปกับ JSON
}

// ... CreateMember และ MemberByEmail ของเดิม ไม่แตะ ...

// MemberByID ใช้กับ GET /api/me: เลขมาจากตั๋วที่ตรวจแล้ว
func (s *Store) MemberByID(ctx context.Context, id int64) (Member, error) {
    var m Member
    err := s.Pool.QueryRow(ctx,
        `SELECT id, email, display_name, role, created_at
         FROM members WHERE id = $1`, id).
        Scan(&m.ID, &m.Email, &m.DisplayName, &m.Role, &m.CreatedAt)
    if errors.Is(err, pgx.ErrNoRows) {
        return Member{}, ErrNotFound
    }
    return m, err
}

// ListMembers รายชื่อสมาชิกทั้งหมดสำหรับ admin เรียงคนใหม่ก่อน
func (s *Store) ListMembers(ctx context.Context) ([]Member, error) {
    rows, err := s.Pool.Query(ctx,
        `SELECT id, email, display_name, role, created_at
         FROM members ORDER BY created_at DESC, id DESC`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    out := []Member{}
    for rows.Next() {
        var m Member
        if err := rows.Scan(&m.ID, &m.Email, &m.DisplayName, &m.Role, &m.CreatedAt); err != nil {
            return nil, err
        }
        out = append(out, m)
    }
    return out, rows.Err()
}

// UpdateMemberRole เปลี่ยนบทบาท (member <-> admin) ไม่พบคนตอบ ErrNotFound
func (s *Store) UpdateMemberRole(ctx context.Context, id int64, role string) (Member, error) {
    var m Member
    err := s.Pool.QueryRow(ctx,
        `UPDATE members SET role = $2, updated_at = now()
         WHERE id = $1
         RETURNING id, email, display_name, role, created_at`, id, role).
        Scan(&m.ID, &m.Email, &m.DisplayName, &m.Role, &m.CreatedAt)
    if errors.Is(err, pgx.ErrNoRows) {
        return Member{}, ErrNotFound
    }
    return m, err
}