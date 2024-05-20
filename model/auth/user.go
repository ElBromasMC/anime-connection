package auth

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"
)

const SessionName = "auth"

type AuthKey struct{}

type UserRole string

const (
	AdminRole    UserRole = "ADMIN"
	NormalRole   UserRole = "NORMAL"
	RecorderRole UserRole = "RECORDER"
)

type User struct {
	Id        uuid.UUID `form:"-"`
	Name      string    `form:"name"`
	Email     string    `form:"email"`
	Password  string    `form:"password"`
	Role      UserRole  `form:"-"`
	CreatedAt time.Time `form:"-"`
}

func GetUser(ctx context.Context) (User, bool) {
	u, ok := ctx.Value(AuthKey{}).(User)
	return u, ok
}
