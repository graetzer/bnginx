package controllers

import (
	"crypto/sha256"
	"encoding/base64"
	"time"
)

type User struct {
	Id       int64
	Name     string
	Email    string
	Password string
	IsAdmin  bool

	Posts []Post
}

func hashPassword(in string) string {
	hash := sha256.New()
	buffer := []byte(in)
	for i := 0; i < 1000; i++ {
		hash.Write(buffer)
		buffer = hash.Sum(nil)
	}
	return base64.StdEncoding.EncodeToString(buffer)
}

func (u User) CheckPassword(in string) bool {
	return hashPassword(in) == u.Password
}

func (u *User) SetPassword(in string) {
	u.Password = hashPassword(in)
}

type Post struct {
	Id     int64
	UserId int64

	CreatedAt time.Time
	UpdatedAt time.Time
	Published bool
	Title     string
	Body      string

	IsPage    bool
	PageOrder int16

	Comments []Comment
}

// Just the first 200 characters (+"..."), may not be very good
// considering that this is supposed to contain valid markdown
func (p Post) Summary() string {
	if len(p.Body) > 200 {
		return p.Body[0:200] + "..."
	}
	return p.Body
}

type Comment struct {
	Id     int64
	PostId int64

	CreatedAt time.Time
	Name      string
	Title     string
	Body      string
	Approved  bool
}
