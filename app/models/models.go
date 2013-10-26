package models

import (
	"time"
	"crypto/sha256"
	"encoding/base64"
	"github.com/coopernurse/gorp"
)

type User struct {
    UserId      int64
    Name        string
    Email       string
    Password    string
    IsAdmin       bool
}

func HashPassword(in string) string {
	hash := sha256.New()
	buffer := []byte(in)
	for i := 0; i < 100; i++ {
		hash.Write(buffer) 
		buffer = hash.Sum(nil)
	}
	return base64.StdEncoding.EncodeToString(buffer)
}

func (u User) CheckPassword(in string) bool {
	return HashPassword(in) == u.Password
}

type Post struct {
    PostId      int64    
    Updated     int64
    Published   bool
    Title       string
    Body        string
    AuthorId    int64
	IsPage      bool
	PageOrder   int16
	
	// Transient
	User		*User	`db:"-"`
}

func NewPost(user *User) *Post {
	post := Post{PostId:0, Updated:time.Now().Unix(), Published:false,
	 Title:"New Post", Body:"", AuthorId:user.UserId, IsPage:false, PageOrder:0}
	return &post
}

func (p Post) UpdatedTime() time.Time {
	return time.Unix(p.Updated, 0)
}

func (p *Post) SetUpdatedTime(t time.Time) {
	p.Updated = t.Unix()
}

// Just the first 200 characters (+"..."), may not be very good
// considering that this is supposed to contain valid markdown
func (p Post) Summary() string {
	if len(p.Body) > 200 {
		return p.Body[0:200] + "..."
	}
	return p.Body
}

// Executed when gorp fetches an instance of this struct
// Fixes gorp's lack of support for relations
func (p *Post) PostGet(s gorp.SqlExecutor) error {
	obj, err := s.Get(User{}, p.AuthorId)
	if (obj != nil) {
		p.User = obj.(*User)
	}
	return err
}

type Comment struct {
	CommentId   int64
	PostId      int64
	Created     int64
	Email       string
	Name        string
	Title       string
	Body        string
	Approved    bool
}

func NewComment() *Comment {
	return &Comment{Created:time.Now().Unix(), Approved:false}
}

func (c Comment) CreatedTime() time.Time {
	return time.Unix(c.Created, 0)
}

func (c *Comment) SetCreatedTime(t time.Time) {
	c.Created = t.Unix()
}