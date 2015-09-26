package controllers

import "github.com/revel/revel"

type Base struct {
	*revel.Controller
}

// ================ Helper functions ================

func (c Base) addTemplateVars() revel.Result {
	c.RenderArgs["user"] = c.connected()
	return nil
}

func (c Base) connected() *User {
	if c.RenderArgs["user"] != nil {
		return c.RenderArgs["user"].(*User)
	}
	if email, ok := c.Session["user"]; ok {
		user := c.getUser(email)
		if user == nil { // Email seems invalid
			delete(c.Session, "user")
		}
		return user
	}
	return nil
}

func (c Base) getUser(email string) *User {
	var user User
	if DB.Where("email = ?", email).First(&user).RecordNotFound() {
		c.Flash.Error("You are not logged in")
		return nil
	}
	return &user
}

func (c Base) getUserById(userId int64) *User {
	var user User
	if DB.First(&user, userId).RecordNotFound() {
		c.Flash.Error("No user with this id")
		return nil
	}
	return &user
}

func (c Base) getPostById(postId int64) *Blogpost {
	var post Blogpost
	if DB.First(&post, postId).RecordNotFound() {
		c.Flash.Error("This Post does not exist")
		return nil
	}
	return &post
}

func (c Base) getPublishedPosts(offset int64) []*Blogpost {
	var posts []*Blogpost
	DB.Where("published").Order("updated_at DESC").Limit(5).Offset(offset).Find(&posts)
	return posts
}
