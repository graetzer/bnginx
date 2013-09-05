package controllers

import (
	"github.com/robfig/revel"
	"github.com/russross/blackfriday"
	"html/template"
	"bnginx/app/models"
)

func init() {
	revel.OnAppStart(DBInit)
	revel.InterceptMethod((*GorpController).Begin, revel.BEFORE)
	revel.InterceptMethod((*GorpController).Commit, revel.AFTER)
	revel.InterceptMethod((*GorpController).Rollback, revel.PANIC)
	revel.InterceptMethod((*App).addUser, revel.BEFORE)
	revel.InterceptMethod((*App).addGlobalPages, revel.BEFORE)
	revel.InterceptMethod((*Admin).checkUser, revel.BEFORE)

	
	revel.TemplateFuncs["markdown"] = func(input string) template.HTML { 
		return template.HTML(string(blackfriday.MarkdownCommon([]byte(input))))
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
}
