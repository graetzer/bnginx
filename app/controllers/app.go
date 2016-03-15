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

// ============== Interceptor ==============
func (c App) addCacheHeaders() revel.Result {
	c.Response.Out.Header().Set("Cache-Control", "public, max-age=14400")
	return nil
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

// Index serves the frontpage including the last stay
func (c App) Index(offset int) revel.Result {
	posts := c.getPublishedPosts(offset)
    var place Place 
    var stay Stay
	DB.Order("started_at DESC").First(&stay)
    if DB.Model(&stay).Related(&place).RecordNotFound() {
        return c.Render(posts, offset)
    }
    return c.Render(place, posts, offset)
}

// Feed serves the current RSS feed
func (c App) Feed() revel.Result {
	c.RenderArgs["posts"] = c.getPublishedPosts(0)
	c.RenderArgs["time"] = time.Now()
	c.Response.ContentType = "application/rss+xml"
	return c.RenderTemplate("App/Feed.xml")
}

// Search serrves matching blogposts
func (c App) Search(query string, offset int) revel.Result {
	var posts []*Blogpost
	q := "%" + query + "%"
	DB.Where("published AND (body like ? OR title like ?)", q, q).Limit(10).Offset(offset).Find(&posts)
	return c.Render(posts, query, offset)
}

// Post renders the post identified by postID
func (c App) Post(postID int64) revel.Result {
	post := c.getPostByID(postID)
	if post == nil {
		return c.NotFound("Oh no! I couldn't find this page")
	}
	var comments []Comment
	DB.Where(&Comment{PostID:postID, Approved:true}).Find(&comments)
	recaptchaSiteKey := revel.Config.StringDefault("recaptcha.sitekey", "")
	return c.Render(post, comments, recaptchaSiteKey)
}

// SaveComment validates and stores user comments
func (c App) SaveComment(postID int64, name, title, body string) revel.Result {
	recaptchaResponse := c.Params.Get("g-recaptcha-response")

	c.Validation.Required(postID)
	c.Validation.Required(name)
	c.Validation.MaxSize(name, 50)
	//c.Validation.Match(email, regexp.MustCompile(`(\w[-._\w]*\w@\w[-._\w]*\w\.\w{2,3})`))
	c.Validation.Required(title)
	c.Validation.MaxSize(title, 100)
	c.Validation.Required(body)
	c.Validation.MaxSize(body, 500)
	c.Validation.Required(recaptchaResponse)

	// Get client IP, optional
	clientIP := c.Request.Header.Get("X-Real-IP")
	if clientIP == "" {
		clientIP = strings.Split(c.Request.RemoteAddr, ":")[0]
	}

	if c.Validation.HasErrors() {
		c.Flash.Error("Could not validate your input")
		c.Validation.Keep()
		c.FlashParams()
	} else if ok := recaptcha.Confirm(clientIP, recaptchaResponse); !ok {
		c.Flash.Error("Wrong captcha")
	} else if post := c.getPostByID(postID); post != nil {
		comment := Comment{PostID: postID, Name: name, Title: title, Body: body}
		DB.Save(&comment)
		c.Flash.Success("Thanks for commenting")
	}
	return c.Redirect(routes.App.Post(postID))
}

func (c App) Projects(hidden bool) revel.Result {
	var projects []*Project
	if hidden {
		DB.Order("updated_at DESC").Find(&projects)
	} else {
		DB.Where("NOT hidden").Order("updated_at DESC").Find(&projects)
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
