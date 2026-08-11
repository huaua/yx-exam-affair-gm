package webres

import (
	"embed"
	"errors"
	"io/fs"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type WebRes struct {
	fs   embed.FS
	path string
}

func (r *WebRes) Open(name string) (fs.File, error) {
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return nil, errors.New("http: invalid character in file path")
	}
	fullName := filepath.Join(r.path, filepath.FromSlash(path.Clean("/"+name)))
	fullName = filepath.ToSlash(fullName)
	// zap.L().Info(fmt.Sprintf("请求资源路径 => %v", fullName))
	file, err := r.fs.Open(fullName)
	return file, err
}

func InitWebResource(engine *gin.Engine) *gin.Engine {
	engine.StaticFS("/assets", http.FS(&WebRes{fs: Assets, path: "html/assets"}))
	engine.StaticFS("/static", http.FS(&WebRes{fs: Static, path: "html/static"}))
	return engine
}
