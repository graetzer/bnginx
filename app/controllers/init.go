package controllers

import (
	"github.com/robfig/revel"
	"github.com/russross/blackfriday"
	"html/template"
)

func init() {
	revel.OnAppStart(DBInit)
	revel.InterceptMethod((*GorpController).Begin, revel.BEFORE)
	revel.InterceptMethod((*App).AddUser, revel.BEFORE)
	revel.InterceptMethod((*Admin).checkUser, revel.BEFORE)
	revel.InterceptMethod((*GorpController).Commit, revel.AFTER)
	revel.InterceptMethod((*GorpController).Rollback, revel.PANIC)
	
	revel.TemplateFuncs["markdown"] = func(input string) template.HTML { 
		return template.HTML(string(blackfriday.MarkdownCommon([]byte(input))))
	}
}
