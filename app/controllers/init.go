package controllers

import (
	"html/template"

	"github.com/dpapathanasiou/go-recaptcha"
	"github.com/russross/blackfriday"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/revel/revel"
)

func init() {
	revel.OnAppStart(AppInit)
	revel.InterceptMethod((*App).addTemplateVars, revel.BEFORE)
	revel.InterceptMethod((*Admin).checkUser, revel.BEFORE)

	revel.TemplateFuncs["markdown"] = func(input string) template.HTML {
		return template.HTML(string(blackfriday.MarkdownCommon([]byte(input))))
	}

	revel.TemplateFuncs["markdownSave"] = func(input string) template.HTML {
		return template.HTML(string(MarkdownSave([]byte(input))))
	}

	revel.TemplateFuncs["sub"] = func(a, b int64) int64 {
		return a - b
	}

	revel.TemplateFuncs["add"] = func(a, b int64) int64 {
		return a + b
	}

	revel.TemplateFuncs["username"] = func(userId int64) string {
		var user User
		if DB.First(&user, userId).RecordNotFound() {
			return ""
		}
		return user.Name
	}

	revel.TemplateFuncs["commentCount"] = func(post *Post) int64 {
		var result int64
		if post == nil {
			DB.Model(Comment{}).Where("NOT approved").Count(&result)
		} else {
			DB.Model(Comment{}).Where("post_id = ?", post.Id).Count(&result)
		}
		return result
	}
}

var DB *gorm.DB

func AppInit() {
	db, err := gorm.Open("sqlite3", "bnginx.db")
	if err == nil {
		db.CreateTable(&User{})
		db.CreateTable(&Post{})
		db.CreateTable(&Comment{})

		var user User

		if db.First(&user).RecordNotFound() {
			user = User{Name: "Simon", Email: "simon@graetzer.org", IsAdmin: true}
			user.SetPassword("default")
			db.Save(&user)
		}
		DB = &db
	}

	if private, found := revel.Config.String("recaptcha.private"); found {
		recaptcha.Init(private)
	}
}

func MarkdownSave(input []byte) []byte {
	// set up the HTML renderer
	htmlFlags := 0
	htmlFlags |= blackfriday.HTML_USE_XHTML
	htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
	htmlFlags |= blackfriday.HTML_SKIP_HTML
	htmlFlags |= blackfriday.HTML_SKIP_STYLE
	htmlFlags |= blackfriday.HTML_SKIP_IMAGES
	htmlFlags |= blackfriday.HTML_SKIP_SCRIPT
	htmlFlags |= blackfriday.HTML_SAFELINK
	renderer := blackfriday.HtmlRenderer(htmlFlags, "", "")

	// set up the parser
	extensions := 0
	extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= blackfriday.EXTENSION_TABLES
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_STRIKETHROUGH
	extensions |= blackfriday.EXTENSION_SPACE_HEADERS

	return blackfriday.Markdown(input, renderer, extensions)
}
