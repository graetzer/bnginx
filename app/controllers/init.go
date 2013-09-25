package controllers

import (
	"github.com/robfig/revel"
	"github.com/graetzer/bnginx/app/models"
	"github.com/russross/blackfriday"
	"github.com/dpapathanasiou/go-recaptcha"
	"html/template"
)

func init() {
	revel.OnAppStart(DBInit)
	revel.OnAppStart(AppInit)
	
	revel.InterceptMethod((*GorpController).Begin, revel.BEFORE)
	revel.InterceptMethod((*GorpController).Commit, revel.AFTER)
	revel.InterceptMethod((*GorpController).Rollback, revel.PANIC)
	
	revel.InterceptMethod((*App).addUser, revel.BEFORE)
	revel.InterceptMethod((*App).addGlobalPages, revel.BEFORE)
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
	
	revel.TemplateFuncs["commentCount"] = func(post *models.Post) (result int64) {
		var err error
		if post == nil {
			result, err = Dbm.SelectInt("SELECT count(*) FROM Comment WHERE NOT Approved")
		} else {
			result, err = Dbm.SelectInt("SELECT count(*) FROM Comment WHERE PostId = ?", post.PostId)
		}
		if err != nil {revel.ERROR.Panic(err)}
		return result
	}
	
	revel.TemplateFuncs["recaptchaKey"] = func () string {
		return revel.Config.StringDefault("recaptcha.public", "")
	}
}

func AppInit() {	
	if private, found := revel.Config.String("recaptcha.private"); found {
		recaptcha.Init (private)
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