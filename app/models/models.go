package models

import "time"

type User struct {
    UserId      int64
    Name        string
    Email       string
    Username    string
    Password    string
    Admin       bool
}

type Post struct {
    PostId      int64    
    Created     int64 // Datetime
    Updated     int64 //Datetime
    Published   bool
    Title       string
    Body        string
    AuthorId    int64
}

func NewPost(user *User) *Post {
	now := time.Now()
	unix := now.Unix()
	//formatted := now.Format("Jan 2, 2006 at 15:04 (MST)")
	
	post := Post{PostId:0, Created:unix, Updated:unix,
	 Published:false, Title:"New Post", Body:"", AuthorId:user.UserId}
	
	return &post
}

func (p Post) CreatedTime() time.Time {
	return time.Unix(p.Created, 0)
}

func (p Post) SetCreatedTime(t time.Time) {
	p.Created = t.Unix()
}

func (p Post) UpdatedTime() time.Time {
	return time.Unix(p.Updated, 0)
}

func (p Post) SetUpdatedTime(t time.Time) {
	p.Updated = t.Unix()
}