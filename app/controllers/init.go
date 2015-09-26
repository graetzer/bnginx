package controllers

import (
	"html/template"
	"os"
	"path/filepath"
	"runtime"
	"encoding/json"

	"github.com/graetzer/go-recaptcha"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/revel/revel"
)

func init() {
	revel.OnAppStart(AppInit)
	revel.InterceptMethod((*Base).addTemplateVars, revel.BEFORE)
	revel.InterceptMethod((*Admin).checkUser, revel.BEFORE)

	revel.TemplateFuncs["markdown"] = func(input string) template.HTML {
		return template.HTML(string(blackfriday.MarkdownCommon([]byte(input))))
	}
	// Use for user supplied content, sanitizes html input and javascript
	revel.TemplateFuncs["markdownSave"] = func(input string) template.HTML {
		return template.HTML(string(SecureMarkdown([]byte(input))))
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

	revel.TemplateFuncs["commentCount"] = func(post *Blogpost) int64 {
		var result int64
		if post == nil {
			DB.Model(Comment{}).Where("NOT approved").Count(&result)
		} else {
			DB.Model(Comment{}).Where("post_id = ? AND approved", post.Id).Count(&result)
		}
		return result
	}

	revel.TemplateFuncs["place"] = func(stay *Stay) *Place {
		var place Place
		if DB.First (&place, stay.PlaceId).RecordNotFound() {
			return nil
		}
		return &place
	}

	revel.TemplateFuncs["json"] = func(v interface{}) string {
		if v == nil {
			return "(nil)"
		}
		bytes, err := json.Marshal(v)
		if err != nil {
			return ""//fmt.Println("error:", err)
		}
		return string(bytes)
	}
}

var (
	DB     *gorm.DB
	policy *bluemonday.Policy
)

func AppInit() {
	// Init HTML sanitizer
	policy = bluemonday.UGCPolicy()

	os.MkdirAll(DataBaseDir(), 0777)
	db, err := gorm.Open("sqlite3", filepath.Join(DataBaseDir(), "sqlite_bnginx.db"))
	if err == nil {
		db.LogMode(revel.DevMode)
		db.CreateTable(&User{})
		db.CreateTable(&Blogpost{})
		db.CreateTable(&Comment{})
		db.CreateTable(&Project{})
		db.CreateTable(&Place{})
		db.CreateTable(&Stay{})

		//db.AutoMigrate(&Location{}, )

		var user User // Add default user
		if db.First(&user).RecordNotFound() {
			user = User{Name: "Simon", Email: "simon@graetzer.org", IsAdmin: true}
			user.SetPassword("default")
			db.Save(&user)
		}
		DB = &db
	}

	if secret, found := revel.Config.String("recaptcha.secret"); found {
		recaptcha.Init(secret)
	}
}

func SecureMarkdown(input []byte) []byte {
	// set up the HTML renderer
	htmlFlags := 0
	htmlFlags |= blackfriday.HTML_USE_XHTML
	htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
	htmlFlags |= blackfriday.HTML_SKIP_HTML
	htmlFlags |= blackfriday.HTML_SKIP_STYLE
	htmlFlags |= blackfriday.HTML_SKIP_IMAGES
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

	unsafe := blackfriday.Markdown(input, renderer, extensions)
	return policy.SanitizeBytes(unsafe)
}

func DataBaseDir() string {
	base, found := revel.Config.String("databasedir")
	if !found {
		if runtime.GOOS == "windows" {
			base = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
			if base == "" {
				base = os.Getenv("USERPROFILE")
			}
		} else {
			base = os.Getenv("HOME")
		}
		base = filepath.Join(base, "/bnginx-data")
	}
	revel.INFO.Println("Using basepath " + base)
	return base
}
