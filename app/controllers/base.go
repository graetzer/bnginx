package controllers

import (
	"bnginx/app/models"

	"github.com/revel/revel"
)

type Base struct {
	*revel.Controller
}

// ================ Helper functions ================

func (c Base) addTemplateVars() revel.Result {
	c.ViewArgs["user"] = c.connected()
	return nil
}

func (c Base) connected() *models.User {
	if c.ViewArgs["user"] != nil {
		return c.ViewArgs["user"].(*models.User)
	}
	if email, ok := c.Session["user"]; ok {
		user := c.getUser(email.(string))
		if user == nil { // Email seems invalid
			delete(c.Session, "user")
		}
		return user
	}
	return nil
}

func (c Base) getUser(email string) *models.User {
	var user models.User
	if DB.Where("email = ?", email).First(&user).RecordNotFound() {
		c.Flash.Error("You are not logged in")
		return nil
	}
	return &user
}

func (c Base) getUserByID(userID int64) *models.User {
	var user models.User
	if DB.First(&user, userID).RecordNotFound() {
		c.Flash.Error("No user with this id")
		return nil
	}
	return &user
}

func (c Base) getPostByID(postID int64) *models.Blogpost {
	var post models.Blogpost
	if DB.First(&post, postID).RecordNotFound() {
		c.Flash.Error("Post does not exist")
		return nil
	}
	return &post
}

func (c Base) getPublishedPosts(offset int, limit int) []*models.Blogpost {
	var posts []*models.Blogpost
	DB.Where("published").Order("created_at DESC").Limit(limit).Offset(offset).Find(&posts)
	return posts
}
