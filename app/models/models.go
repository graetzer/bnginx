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

type Comment struct {
	CommentId   int64
	PostId      int64
	Created     int64
	Title       string
	Body        string
	Email       string
	Approved    bool
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