package controllers

import (
	"github.com/revel/revel"
	
	"os"
	fpath "path/filepath"
	"strings"
	"syscall"
)

type Uploads struct {
	*revel.Controller
}

func (c Uploads) ServeUpload(filepath string) revel.Result {

	basePathPrefix := fpath.Join(DataBaseDir(), "uploads")
	fname := fpath.Join(basePathPrefix, fpath.FromSlash(filepath))
	// Verify the request file path is within the application's scope of access
	if !strings.HasPrefix(fname, basePathPrefix) {
		c.Log.Warn("Attempted to read file outside of base path: %s", fname)
		return c.NotFound("")
	}

	// Verify file path is accessible
	finfo, err := os.Stat(fname)
	if err != nil {
		if os.IsNotExist(err) || err.(*os.PathError).Err == syscall.ENOTDIR {
			c.Log.Warn("File not found (%s): %s ", fname, err)
			return c.NotFound("File not found")
		}
		c.Log.Error("Error trying to get fileinfo for '%s': %s", fname, err)
		return c.RenderError(err)
	}

	// Disallow directory listing
	if finfo.Mode().IsDir() {
		c.Log.Warn("Attempted directory listing of %s", fname)
		return c.Forbidden("Directory listing not allowed")
	}

	// Open request file path
	file, err := os.Open(fname)
	if err != nil {
		if os.IsNotExist(err) {
			c.Log.Warn("File not found (%s): %s ", fname, err)
			return c.NotFound("File not found")
		}
		c.Log.Error("Error opening '%s': %s", fname, err)
		return c.RenderError(err)
	}
	return c.RenderFile(file, revel.Inline)
}
