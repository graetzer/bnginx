package models

import "time"

type User struct {
    UserId      int64
    Name        string
    Email       string
    Password    string
    IsAdmin       bool
}

func (u User) CheckPassword(inPass string) bool {
	//h := md5.New()
 //   io.WriteString(h, "The fog is getting thicker!")
 //   fmt.Printf("%x", h.Sum(nil))
	
	return inPass == u.Password
}

type Post struct {
    PostId      int64    
    Updated     int64 //Datetime
    Published   bool
    Title       string
    Body        string
    AuthorId    int64
	IsPage      bool
	PageOrder   int16
}

func NewPost(user *User) *Post {
	now := time.Now()
	unix := now.Unix()
	//formatted := now.Format("Jan 2, 2006 at 15:04 (MST)")
	
	post := Post{PostId:0, Updated:unix,
	 Published:false, Title:"New Post", Body:"", AuthorId:user.UserId, IsPage:false, PageOrder:0}
	
	return &post
}

func (p Post) UpdatedTime() time.Time {
	return time.Unix(p.Updated, 0)
}

func (p Post) SetUpdatedTime(t time.Time) {
	p.Updated = t.Unix()
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

func (c Comment) SetCreatedTime(t time.Time) {
	c.Created = t.Unix()
}