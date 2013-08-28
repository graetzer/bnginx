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

func (c Admin) CreateUser(name string, email string) revel.Result {
	if !c.connected().Admin {
		return c.Redirect(routes.Admin.Index())
	}
	
	user := c.getUser(email)	
	if user != nil {
		c.Flash.Error("User already exists")
		return c.Redirect(routes.Admin.Index())
	}
	
    user = &models.User{UserId:0, Name:name, Email:email, Password:"default", Admin:false}
    err := c.Txn.Insert(user)
	if err != nil {
		revel.ERROR.Panic(err)
	}
	c.Flash.Success("Set password to default")
    return c.Redirect(routes.Admin.Index())
}


func (c Admin) EditPost(postId int64) revel.Result {
	var post *models.Post = new(models.Post)
    if postId > 0 {
		_, err := c.Txn.Get(&post, postId)
		if err != nil {
			revel.ERROR.Panic(err)
		} else if post == nil {
			c.Flash.Error("No result for this id")
			return c.Redirect(routes.Admin.Index())
		}
	}
	
	return c.Render(post)
}

func (c Admin) SavePost(postId int64, published bool, title, body string) revel.Result {
	var post *models.Post
	if postId <= 0 {
		post = models.NewPost(c.connected())
	} else {
		_, err := c.Txn.Get(&post, postId)
		if err != nil {
			revel.ERROR.Panic(err)
		} else if post == nil {
			c.Flash.Error("No result for this id")
			return c.Redirect(routes.Admin.Index())
		}
		
		if c.connected().UserId != post.AuthorId && !c.connected().Admin {
			c.Flash.Error("You have no permission to edit this post")
			return c.Redirect(routes.Admin.Index())
		}
	}
	post.SetUpdatedTime(time.Now())
	post.Published = published
	post.Title = title
	post.Body = body
	
	var err error
	if (post.PostId <= 0) {
		err = c.Txn.Insert(post)
	} else {
		_, err = c.Txn.Update(post)
	}
	if err != nil {
		revel.ERROR.Panic(err)
	}
	
	return c.Redirect(routes.Admin.Index())
}