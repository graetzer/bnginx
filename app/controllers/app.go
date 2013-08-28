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
    if email, ok := c.Session["user"]; ok {
        return c.getUser(email)
    }
    return nil
}

func (c App) getUser(email string) *models.User {
    users, err := c.Txn.Select(models.User{}, `SELECT * FROM User WHERE Email = ?`, email)
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

func (c App) Login(email string, password string) revel.Result {
    if user := c.connected(); user != nil {
        c.Flash.Error("Already logged in as %s", user.Name)
        return c.Redirect(routes.App.Index())
    }
	
    user := c.getUser(email)
    if user == nil || user.Password != password {
        c.Flash.Error("Wrong email or password")
        return c.Redirect(routes.App.Index())
    }
    c.Session["user"] = user.Email
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
		revel.ERROR.Panic(err.Error())
	}
	return c.Render(posts)
}

func (c App) Search(query string) revel.Result {
	var posts []*models.Post
	q := "%"+query+"%"
	_, err := c.Txn.Select(&posts, `SELECT * FROM Post WHERE Published AND (Body like ? OR Title like ?)`, q, q)
	if err != nil {
		revel.ERROR.Panic(err.Error())
	}
	return c.Render(posts, query)
}

func (c App) Page(postId int64) revel.Result {
	var post *models.Post
	_, err := c.Txn.Get(&post, postId)
	if err != nil {
		c.Flash.Error("Wrong post ID")
		return c.Redirect(routes.App.Index())
	}
	
	user := c.connected()
	if post == nil || !post.Published && (user == nil || !user.Admin) {
		c.Flash.Error("Wrong post ID")
		return c.Redirect(routes.App.Index())
	}
	
	input := []byte(post.Body)
	output := blackfriday.MarkdownCommon(input)
	
	return c.RenderText(string(output))
}
