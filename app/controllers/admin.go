package controllers

import (
	"github.com/graetzer/bnginx/app/routes"
	"github.com/revel/revel"

	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
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
		posts []*Post
		users []*User
	)
	DB.Order("page_order DESC, updated_at DESC").Find(&posts)
	DB.Find(&users)
	return c.Render(posts, users)
}

// ==================== Handle Users ====================

func (c Admin) EditUser(email string) revel.Result {
	profile := new(User)
	profile.Id = -1
	if email != "create" {
		profile = c.getUser(email)
		if profile == nil {
			return c.Redirect(routes.Admin.Index())
		}
	}

	return c.Render(profile)
}

func (c Admin) SaveUser(userId int64, name, email, oldPassword, password string) revel.Result {
	c.Validation.Required(userId)
	c.Validation.Required(email)
	c.Validation.MinSize(password, 8)
	c.Validation.MaxSize(password, 300)
	//c.Validation.Match(username, regexp.MustCompile("^\\w*$"))
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		if userId == -1 {
			return c.Redirect(routes.Admin.Index())
		} else {
			return c.Redirect(routes.Admin.EditUser(email))
		}
	}

	var user *User
	u := c.connected()
	if userId <= 0 {
		user = new(User)
		if c.getUser(email) != nil {
			c.Flash.Error("This is email address is already used")
			return c.Redirect(routes.Admin.EditUser("create"))
		}
		delete(c.Flash.Out, "error")
	} else {
		if u.Id != userId && !u.IsAdmin {
			c.Flash.Error("You have no permission to edit this profile")
			return c.Redirect(routes.Admin.EditUser("create"))
		}

		user = c.getUserById(userId)
		if user == nil || !user.CheckPassword(oldPassword) {
			c.Flash.Error("Db corruption or the password is incorrect")
			return c.Redirect(routes.Admin.EditUser(user.Email))
		}
	}
	user.Name = name
	user.Email = email
	user.SetPassword(password)
	DB.Save(user)

	return c.Redirect(routes.Admin.Index())
}

func (c Admin) DeleteUser(email string) revel.Result {

	if user := c.getUser(email); user != nil {
		u := c.connected()
		if u.Id == user.Id || u.IsAdmin {
			DB.Delete(user)
			c.Flash.Success("Deleted user")
		}
	}
	return c.Redirect(routes.Admin.Index())
}

// ==================== Handle Posts ====================

func (c Admin) EditPost(postId int64) revel.Result {
	post := new(Post)
	if postId > 0 {
		post = c.getPostById(postId)
		if post == nil {
			return c.Redirect(routes.Admin.Index())
		}
	} else {
		post.Title = "A new Blogpost"
		post.Body = "### Start with something\n\nE.g.\n\n1. Make a List\n2. Of Interesting\n3. Things"
	}
	return c.Render(post)
}

func (c Admin) SavePost(postId int64, published bool, title, body string, isPage bool, pageOrder int16) revel.Result {
	var post *Post
	u := c.connected()
	if postId <= 0 { // Create a new one
		post = &Post{UserId: u.Id}
	} else {
		post = c.getPostById(postId)
		if post == nil {
			return c.Redirect(routes.Admin.Index())
		}
		if u.Id != post.UserId && !u.IsAdmin {
			c.Flash.Error("You have not the permission to save this post")
			return c.Redirect(routes.Admin.Index())
		}
	}
	post.Published = published
	post.Title = title
	post.Body = body
	post.IsPage = isPage
	post.PageOrder = pageOrder
	DB.Save(post)

	return c.Redirect(routes.Admin.Index())
}

func (c Admin) DeletePost(postId int64) revel.Result {
	post := c.getPostById(postId)
	if post == nil {
		return c.Redirect(routes.Admin.Index())
	}

	u := c.connected()
	if u.Id == post.UserId || u.IsAdmin {
		DB.Where("PostId = ?", post.Id).Delete(&Comment{})
		DB.Delete(post)
		c.Flash.Success("Deleted post")
	}
	return c.Redirect(routes.Admin.Index())
}

// ==================== Handle Uploads ====================

func (c Admin) Media() revel.Result {
	// TODO configure that
	basePath := filepath.Join(revel.BasePath, filepath.FromSlash("public/uploads/"))

	fs, err := ioutil.ReadDir(basePath)
	if err != nil {
		revel.ERROR.Println(err)
		return c.RenderError(err)
	}

	// Remove hidden files
	files := make([]os.FileInfo, 0, len(fs))
	for _, f := range fs {
		if !strings.HasPrefix(f.Name(), ".") {
			files = append(files, f)
		}
	}

	uploadPrefix := time.Now().Format("2006_01_02_")
	return c.Render(files, uploadPrefix)
}

func (c Admin) Upload() revel.Result {
	basePath := filepath.Join(revel.BasePath, filepath.FromSlash("public/uploads/"))

	for _, fInfo := range c.Params.Files["file"] {

		fi, err := fInfo.Open()
		if err != nil {
			return c.RenderError(err)
		}
		defer fi.Close()

		var uploadPrefix = time.Now().Format("2006_01_02_")
		full := filepath.Join(basePath, uploadPrefix+filepath.Base(fInfo.Filename))

		fo, err := os.Create(full)
		if err != nil {
			return c.RenderError(err)
		}
		defer fo.Close()

		if _, err = io.Copy(fo, fi); err != nil {
			return c.RenderError(err)
		}
		return c.RenderJson(struct{ Message string }{"Success"})
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

// ==================== Handle Comments ====================

func (c Admin) Comments() revel.Result {
	var comments []*Comment
	DB.Order("created_at DESC").Find(&comments)
	return c.Render(comments)
}

func (c Admin) UpdateComment(commentId int64, approved bool) revel.Result {
	var comment Comment
	if !DB.First(&comment, commentId).RecordNotFound() {
		comment.Approved = approved
		DB.Save(&comment)
	}
	return c.Redirect(routes.Admin.Comments())
}

func (c Admin) DeleteComment(commentId int64) revel.Result {
	var comment Comment
	if !DB.First(&comment, commentId).RecordNotFound() {
		DB.Delete(&comment)
	}
	return c.Redirect(routes.Admin.Comments())
}
