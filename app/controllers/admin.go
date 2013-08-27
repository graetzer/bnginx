package controllers

import (
    "github.com/robfig/revel"
    "bngnix/app/models"
	"bngnix/app/routes"
	"time"
)

type Admin struct {
	*revel.Controller
    App
}

func (c Admin) checkUser() revel.Result {
    if user := c.connected(); user == nil {
        c.Flash.Error("Please log in first")
        return c.Redirect(routes.App.Index())
    }
    return nil
}

func (c Admin) Index() revel.Result {
	var posts []*models.Post
	_, err := c.Txn.Select(&posts, "SELECT * FROM Post")
	if err != nil {
		revel.ERROR.Panic(err)
	}
	return c.Render(posts)
}

func (c Admin) CreateUser(username string, name string, email string) revel.Result {
	user := c.getUser(username)	
	if user != nil {
		c.Flash.Error("User already exists")
		return c.Redirect(routes.Admin.Index())
	}
	
    user = &models.User{UserId:0, Name:name, Email:email, Username:username, Password:"default"}
    err := c.Txn.Insert(user)
	if err != nil {
		revel.ERROR.Panic(err)
	}
	c.Flash.Success("Set password to default")
    return c.Redirect(routes.Admin.Index())
}

func (c Admin) CreatePost() revel.Result {
	post := models.NewPost(c.connected())// Creates
	err := c.Txn.Insert(post)
	if err != nil {
		revel.WARN.Println(err)
	}
	return c.Redirect(routes.Admin.EditPost(post.PostId))
}

func (c Admin) EditPost(postId int64) revel.Result {
	obj, err := c.Txn.Get(models.Post{}, postId)
	if err != nil {
		revel.ERROR.Panic(err)
	}
	if obj == nil {
		c.Flash.Error("No result for this id")
		return c.Redirect(routes.Admin.Index())
	}
	
	post := obj.(*models.Post)
	return c.Render(post)
}

func (c Admin) SavePost(post *models.Post) revel.Result {
	post.SetUpdatedTime(time.Now())
    _, err := c.Txn.Update(post)
	if err != nil {
		revel.ERROR.Panic(err)
	}
	return c.Redirect(routes.Admin.Index())
}