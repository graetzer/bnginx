package controllers

import (
    "github.com/robfig/revel"
    "bnginx/app/models"
	"bnginx/app/routes"
	"time"
	"path/filepath"
	"os"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
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
	var (
		posts []*models.Post
		users []*models.User
	)
	_, err := c.Txn.Select(&posts, "SELECT * FROM Post")
	if err != nil {
		revel.ERROR.Panic(err)
	}
	_, err = c.Txn.Select(&users, "SELECT * FROM User")
	if err != nil {revel.ERROR.Panic(err)}
	return c.Render(posts, users)
}

// ==================== Handle Users ====================

func (c Admin) EditUser (email string) revel.Result {
	profile := new(models.User)
	profile.UserId = -1
    if email != "create" {
		profile = c.getUser(email)
		if profile == nil {
			return c.Redirect(routes.Admin.Index())
		}
	}
	
	return c.Render(profile)
}

func (c Admin) SaveUser(userId int64, name string, email string, password string) revel.Result {
	c.Validation.Required(userId)
	c.Validation.Required(email)
	c.Validation.MinSize(password, 8)
	//c.Validation.Match(username, regexp.MustCompile("^\\w*$"))
	if c.Validation.HasErrors() {
		// Store the validation errors in the flash context and redirect.
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Admin.Index())
	}
	
	var user *models.User
	u := c.connected()
	if userId == -1 {
		user = new(models.User)
		if c.getUser(email) != nil {
			c.Flash.Error("This is email address is already used")
			return c.Redirect(routes.Admin.Index())
		}
		delete(c.Flash.Out, "error")
	} else {
		if u.UserId != userId && !u.IsAdmin {
			c.Flash.Error("You have not the permission to save this post")
			return c.Redirect(routes.Admin.Index())
		}
		
		user = c.getUserById(userId)
		if user == nil {
			return c.Redirect(routes.Admin.Index())
		}
	}
	user.Name = name
	user.Email = email
	user.Password = password
	
	var err error
	if (userId == -1) {
		err = c.Txn.Insert(user)
	} else {
		_, err = c.Txn.Update(user)
	}
	if err != nil {
		revel.ERROR.Panic(err)
	}
	
	return c.Redirect(routes.Admin.Index())
}

func (c Admin) DeleteUser(email string) revel.Result {
	user := c.getUser(email)
	if user == nil {
		return c.Redirect(routes.Admin.Index())
	}

	u := c.connected()
	if u.UserId == user.UserId || u.IsAdmin {
		_, err := c.Txn.Delete(user)
		if err == nil {
			revel.ERROR.Panic(err)
		}
		c.Flash.Success("Deleted user")
	}
	return c.Redirect(routes.Admin.Index())
}

// ==================== Handle Posts ====================

func (c Admin) EditPost(postId int64) revel.Result {
	post := new(models.Post)
	post.PostId = -1
    if c.Params.Get("postId") != "create" {
		post = c.getPostById(postId)
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
		post = c.getPostById(postId)
		if post == nil {
			return c.Redirect(routes.Admin.Index())
		}
		
		if u.UserId != post.AuthorId && !u.IsAdmin {
			c.Flash.Error("You have not the permission to save this post")
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
	post := c.getPostById(postId)
	if post == nil {
		return c.Redirect(routes.Admin.Index())
	}

	u := c.connected()
	if u.UserId == postId || u.IsAdmin {
		_, err := c.Txn.Delete(post)
		if err == nil {
			revel.ERROR.Panic(err)
		}
		c.Flash.Success("Deleted post")
	}
	return c.Redirect(routes.Admin.Index())
}

// ==================== Handle Uploads ====================

func (c Admin) Upload() revel.Result {
	// TODO configure that
	basePath := filepath.Join(revel.BasePath, filepath.FromSlash("public/uploads/"))
	
	fs, err := ioutil.ReadDir(basePath)
	if err != nil {revel.ERROR.Panic(err)}
	files := make([]os.FileInfo, 0, len(fs))
	for _, f := range fs {
		if ! strings.HasPrefix(f.Name(), ".") {
			files = append(files, f)
		} 
	}
	
	return c.Render(files)
}

func (c Admin) SaveUpload() revel.Result {
	basePath := filepath.Join(revel.BasePath, filepath.FromSlash("public/uploads/"))
	
	for _, fInfo := range c.Params.Files["file"] {
		
		fi, err := fInfo.Open()
		if err != nil { return c.RenderError(err) }
		defer fi.Close()
		
		time := strconv.FormatInt(time.Now().Unix(), 10)
		full := filepath.Join(basePath, time + "_" + filepath.Base(fInfo.Filename))
		
		fo, err := os.Create(full)
		if err != nil { return c.RenderError(err) }
		defer fo.Close()
		
		if _, err = io.Copy(fo, fi); err != nil { 
			return c.RenderError(err) 
		}
		return c.RenderJson(struct {Message string}{"Success"})
    }
	return c.Redirect(routes.Admin.Upload())
}

func (c Admin) DeleteUpload(filename string) revel.Result {
	if c.connected().IsAdmin {
		basePath := filepath.Join(revel.BasePath, filepath.FromSlash("public/uploads/"))
		full := filepath.Join(basePath, filepath.Base(filename))
		err := os.Remove(full)
		if err != nil {
			c.Flash.Error(err.Error())
		}
		
	}
	return c.Redirect(routes.Admin.Upload())
}

// ==================== Handle Uploads ====================

func (c Admin) Comments() revel.Result {
	var comments []*models.Comment
	_, err := c.Txn.Select(&comments, "SELECT * FROM Comment ORDER BY Created DESC")
	if err != nil {revel.ERROR.Panic(err)}
	return c.Render(comments)
}

func (c Admin) UpdateComment (commentId int64, approved bool) revel.Result {
	c.Validation.Required(commentId)
	c.Validation.Required(approved)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Admin.Comments())
	}
	
	obj, err := c.Txn.Get(models.Comment{}, commentId)
	if err != nil {revel.ERROR.Panic(err)}
	if obj != nil {
		comment := obj.(*models.Comment)
		comment.Approved = approved
		_, err := c.Txn.Update(comment)
		if err != nil {
			revel.ERROR.Panic(err)
		}
	}
	
	return c.Redirect(routes.Admin.Comments())
}

func (c Admin) DeleteComment (commentId int64) revel.Result {
	obj, err := c.Txn.Get(models.Comment{}, commentId)
	if err != nil {revel.ERROR.Panic(err)}
	if obj != nil {
		comment := obj.(*models.Comment)
		_, err := c.Txn.Delete(comment)
		if err != nil {
			revel.ERROR.Panic(err)
		}
	}
	
	return c.Redirect(routes.Admin.Comments())
}
