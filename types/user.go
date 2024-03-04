package types

import (
	"errors"
	"time"

	"github.com/happilymarrieddad/old-world/api3/internal/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"golang.org/x/crypto/bcrypt"
)

type UserType string

const (
	AdminUserType    UserType = "Admin"
	StandardUserType UserType = "Standard User"
)

type CreateUser struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `validate:"required" json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"confirm_password"`
}

func (u CreateUser) GetPasswordHash() (hash string, err error) {
	if len(u.Password) == 0 || len(u.PasswordConfirm) == 0 {
		return "", errors.New("invalid password and|or confirm")
	}

	if u.Password != u.PasswordConfirm {
		return "", errors.New("invalid password and|or confirm")
	}

	bts, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	return string(bts), err
}

type User struct {
	ID         string     `json:"id"`
	FirstName  string     `validate:"required" json:"first_name"`
	LastName   string     `validate:"required" json:"last_name"`
	Email      string     `validate:"required" json:"email"`
	Password   string     `json:"-"`
	AddressID  int64      `json:"address_id"`
	BirthMonth int64      `json:"birth_month"`
	BirthDay   int64      `json:"birth_day"`
	BirthYear  int64      `json:"birth_year"`
	UserType   UserType   `json:"type"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

func (u *User) PasswordMatches(psw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(psw)) == nil
}

func UserFromNode(node dbtype.Node) *User {
	u := &User{
		ID:        node.Props["id"].(string),
		FirstName: node.Props["first_name"].(string),
		LastName:  node.Props["last_name"].(string),
		Email:     node.Props["email"].(string),
		Password:  node.Props["password"].(string),
	}

	timeRaw, ok := node.Props["created_at"].(int64)
	if ok {
		u.CreatedAt = time.Unix(timeRaw, 0)
	}

	timeRaw, ok = node.Props["updated_at"].(int64)
	if ok {
		u.UpdatedAt = utils.Ref(time.Unix(timeRaw, 0))
	}

	return u
}
