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
        return c.Redirect(routes.App.Index(0))
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

func (c Admin) EditUser (userId int64) revel.Result {
	user := new (models.User)
	user.UserId = -1
	if  c.Params.Get("userId") != "create" {
		obj, err := c.Txn.Get(models.User{}, userId)
		if err != nil {
			revel.ERROR.Panic(err)
		} else if user == nil {
			c.Flash.Error("No result for this id")
			return c.Redirect(routes.Admin.Index())
		}
		user = obj.(*models.User)
	}
	return c.Render(user)
}

func (c Admin) SaveUser(name string, email string) revel.Result {
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
	post := new(models.Post)
	post.PostId = -1
    if c.Params.Get("postId") != "create" {
		post = c.getPost(postId)
		if post == nil {
			return c.Redirect(routes.Admin.Index())
		}
	}
	
	return c.Render(post)
}

func (c Admin) SavePost(postId int64, published bool, title, body string, isPage bool, pageOrder int16) revel.Result {
	var post *models.Post
	u := c.connected()
	
	if postId == -1 {
		post = models.NewPost(c.connected())
		post.AuthorId = u.UserId
	} else {
		post = c.getPost(postId)
		if post == nil {
			return c.Redirect(routes.Admin.Index())
		}
		
		if u.UserId != post.AuthorId && !u.Admin {
			c.Flash.Error("You have no permission to edit this post")
			return c.Redirect(routes.Admin.Index())
		}
	}
	post.SetUpdatedTime(time.Now())
	post.Published = published
	post.Title = title
	post.Body = body
	post.IsPage = isPage
	post.PageOrder = pageOrder
	
	var err error
	if (postId == -1) {
		err = c.Txn.Insert(post)
	} else {
		_, err = c.Txn.Update(post)
	}
	if err != nil {
		revel.ERROR.Panic(err)
	}
	
	return c.Redirect(routes.Admin.Index())
}

func (c Admin) DeletePost(postId int64) revel.Result {
	post := c.getPost(postId)
	if post == nil {
		return c.Redirect(routes.Admin.Index())
	}

	u := c.connected()
	if u.UserId == postId || u.Admin {
		_, err := c.Txn.Delete(post.PostId)
		if err == nil {
			revel.ERROR.Panic(err)
		}
		c.Flash.Success("Deleted post")
	}
	return c.Redirect(routes.Admin.Index())
}
