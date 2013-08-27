package controllers

import (
    "github.com/robfig/revel"
    "bngnix/app/models"
	"bngnix/app/routes"
	"github.com/russross/blackfriday"
)

type App struct {
    GorpController
}

// ================ Helper functions ================
func (c App) AddUser() revel.Result {
    if user := c.connected(); user != nil {
         c.RenderArgs["user"] = user
    }
    return nil
}

func (c App) connected() *models.User {
    if c.RenderArgs["user"] != nil {
        return c.RenderArgs["user"].(*models.User)
    }
    if username, ok := c.Session["user"]; ok {
        return c.getUser(username)
    }
    return nil
}

func (c App) getUser(username string) *models.User {
    users, err := c.Txn.Select(models.User{}, `SELECT * FROM User WHERE Username = ?`, username)
    if err != nil {
		//c.Flash.Error("Database error: "+err.Error())
		//return nil
        panic(err)
    }
    if len(users) == 0 {
        return nil
    }
    return users[0].(*models.User)
}

// ================ Actions ================

func (c App) Login(username string, password string) revel.Result {
    if user := c.connected(); user != nil {
        c.Flash.Error("Already logged in as %s", user.Username)
        return c.Redirect(routes.App.Index())
    }
    user := c.getUser(username)
    if user == nil || user.Password != password {
        c.Flash.Error("Wrong username or password")
        return c.Redirect(routes.App.Index())
    }
    c.Session["user"] = user.Username
    c.RenderArgs["User"] = user // Probably not needed
    return c.Redirect(routes.Admin.Index())
}

func (c App) Logout() revel.Result {
	delete(c.Session, "user")
	delete(c.RenderArgs, "user")
	return c.Redirect(routes.App.Index())
}

func (c App) Index() revel.Result {
	var posts []*models.Post
	_, err := c.Txn.Select(&posts, `SELECT * FROM Post WHERE Published`)
	if err != nil {
		revel.ERROR.Fatal(err.Error())
	}
	return c.Render(posts)
}

func (c App) Page(postId int64) revel.Result {
	obj, err := c.Txn.Get(models.Post{}, postId)
	post := obj.(*models.Post)
	if err != nil {
		c.Flash.Error("Wrong post ID")
		return c.Redirect(routes.App.Index())
	}
	
	user := c.connected()
	if !post.Published && (user == nil || !user.Admin) {
		c.Flash.Error("Wrong post ID")
		return c.Redirect(routes.App.Index())
	}
	
	input := []byte(post.Body)
	output := blackfriday.MarkdownCommon(input)
	
	return c.RenderText(string(output))
}
