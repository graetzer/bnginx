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

func (u User) CheckPassword(in string) bool {
	return hashPassword(in) == u.Password
}

func (u *User) SetPassword(in string) {
	u.Password = hashPassword(in)
}

type Blogpost struct {
	Id        int64
	CreatedAt time.Time
	UpdatedAt time.Time

	Title     string
	Body      string
	Published bool
	UserId    int64 // Foreign key of User
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

type Comment struct {
	Id     int64
	CreatedAt time.Time

	PostId    int64 // Foreign key of BlogPost
	Name      string
	Title     string
	Body      string
	Approved  bool
}

type Project struct {
	Id        int64
	UpdatedAt time.Time

	Title       string
	Description string
	CoverUrl    string
	RepoUrl			string
	Tags        string
}

type Stay struct {
	Id     int64 `json:"id"`
	PlaceId    int64 // Foreign key of Place
	StartedAt time.Time `json:"startedAt"`
	EndedAt   time.Time `json:"endedAt"`
	Url string `json:"url"`
}

type Place struct {
	Id     int64 `json:"id"`

	Name  string `json:"name"`
	CoverUrl string `json:"coverUrl"`
	Latitude   float32 `json:"lat"`
	Longitude   float32 `json:"lng"`
	Stays []Stay // One-To-Many relationship (has many)
}
