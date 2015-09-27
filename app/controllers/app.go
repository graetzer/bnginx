package controllers

import (
	"strings"
	"time"

	"github.com/graetzer/bnginx/app/routes"
	"github.com/graetzer/go-recaptcha"
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
	Base
}

// ================ Actions ================

func (c App) Login(email string, password string) revel.Result {
	c.Validation.Required(email)
	c.Validation.Required(password)
	c.Validation.MaxSize(password, 300) // Kinda important
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.App.Index(0))
	}

	user := c.getUser(email)
	if user == nil || !user.CheckPassword(password) {
		c.Flash.Error("Wrong email or password")
		return c.Redirect(routes.App.Index(0))
	}
	c.Session["user"] = user.Email
	c.RenderArgs["User"] = user // Probably not needed
	return c.Redirect(routes.Admin.Index())
}

func (c App) Logout() revel.Result {
	delete(c.Session, "user")
	delete(c.RenderArgs, "user")
	return c.Redirect(routes.App.Index(0))
}

func (c App) Index(offset int64) revel.Result {
	posts := c.getPublishedPosts(offset)
	return c.Render(posts, offset)
}

func (c App) Feed() revel.Result {
	c.RenderArgs["posts"] = c.getPublishedPosts(0)
	c.RenderArgs["time"] = time.Now()
	c.Response.ContentType = "application/rss+xml"
	return c.RenderTemplate("App/Feed.xml")
}

func (c App) Search(query string, offset int64) revel.Result {
	var posts []*Blogpost
	q := "%" + query + "%"
	DB.Where("published AND (body like ? OR title like ?)", q, q).Limit(10).Offset(offset).Find(&posts)
	return c.Render(posts, query, offset)
}

func (c App) Post(postId int64) revel.Result {
	post := c.getPostById(postId)
	if post == nil {
		return c.NotFound("Oh no! I couldn't find this page")
	}
	var comments []Comment
	DB.Model(post).Related(&comments, "PostId")
	recaptchaSiteKey := revel.Config.StringDefault("recaptcha.sitekey", "")
	return c.Render(post, comments, recaptchaSiteKey)
}

func (c App) SaveComment(postId int64, name, title, body string) revel.Result {
	recaptcha_response := c.Params.Get("g-recaptcha-response")

	c.Validation.Required(postId)
	c.Validation.Required(name)
	c.Validation.MaxSize(name, 50)
	//c.Validation.Match(email, regexp.MustCompile(`(\w[-._\w]*\w@\w[-._\w]*\w\.\w{2,3})`))
	c.Validation.Required(title)
	c.Validation.MaxSize(title, 100)
	c.Validation.Required(body)
	c.Validation.MaxSize(body, 500)
	c.Validation.Required(recaptcha_response)

	// Get client IP, optional
	client_ip := c.Request.Header.Get("X-Real-IP")
	if client_ip == "" {
		client_ip = strings.Split(c.Request.RemoteAddr, ":")[0]
	}

	if c.Validation.HasErrors() {
		c.Flash.Error("Could not validate your input")
		c.Validation.Keep()
		c.FlashParams()
	} else if ok := recaptcha.Confirm(client_ip, recaptcha_response); !ok {
		c.Flash.Error("Wrong captcha")
	} else if post := c.getPostById(postId); post != nil {
		comment := Comment{PostId: postId, Name: name, Title: title, Body: body}
		DB.Save(&comment)
		c.Flash.Success("Thanks for commenting")
	}
	return c.Redirect(routes.App.Post(postId))
}

func (c App) Projects(hidden bool) revel.Result {
	var projects []*Project
	if hidden {
		DB.Order("updated_at DESC").Find(&projects)
	} else {
		DB.Where("!hidden").Order("updated_at DESC").Find(&projects)
	}
	return c.Render(projects)
}

func (c App) About() revel.Result {
	var (
		places []Place
		stays []Stay
	)
	DB.Order("started_at DESC").Find(&stays)
	DB.Find(&places)
	return c.Render(places, stays)
}

func (c App) Imprint() revel.Result {
	return c.Render()
}
