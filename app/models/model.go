package models

import (
	"crypto/sha256"
	"encoding/base64"
	"time"
)
// User model
type User struct {
	ID       int64
	Name     string
	Email    string
	Password string
	IsAdmin  bool

	Blogposts []Blogpost // One-To-Many relationship (has many)
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

// CheckPassword hashes and compares the passwords
func (u User) CheckPassword(in string) bool {
	return hashPassword(in) == u.Password
}

// SetPassword hashes and sets the password
func (u *User) SetPassword(in string) {
	u.Password = hashPassword(in)
}

// Blogpost contains the text and metadata of a post
type Blogpost struct {
	ID        int64 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Title     string
	Body      string
	Published bool
	UserID    int64 // Foreign key of User
	Comments []Comment // One-To-Many relationship (has many)
}

// Just the first 200 characters (+"..."), may not be very good
// considering that this is supposed to contain valid markdown
func (p Blogpost) Summary() string {
	if len(p.Body) > 200 {
		return p.Body[0:200] + "..."
	}
	return p.Body
}

// Comment contains user comments
type Comment struct {
	ID     int64
	CreatedAt time.Time

	PostID    int64 // Foreign key of BlogPost
	Name      string
	Title     string
	Body      string
	Approved  bool
}

// Contains info about my pet projects
type Project struct {
	ID        int64
	UpdatedAt time.Time

	Title       string
	Description string
	CoverUrl    string
	RepoUrl			string
	Tags        string
	Hidden      bool
}

// Stay documents where I was at a point
type Stay struct {
	ID     int64 `json:"id"`
	PlaceID    int64 // Foreign key of Place
	StartedAt time.Time `json:"startedAt"`
	EndedAt   time.Time `json:"endedAt"`
	Url string `json:"url"`
}

// Place is somwhere in the world
type Place struct {
	ID     int64 `json:"id"`

	Name  string `json:"name"`
	CoverUrl string `json:"coverUrl"`
	Latitude   float32 `json:"lat"`
	Longitude   float32 `json:"lng"`
	Stays []Stay // One-To-Many relationship (has many)
}
