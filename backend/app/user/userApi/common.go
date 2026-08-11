package userApi

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sealsee/web-base/public/IOFile"
	"github.com/sealsee/web-base/public/context"
	"github.com/sealsee/web-base/public/cst/common"
	"github.com/sealsee/web-base/public/errs"
	"github.com/sealsee/web-base/public/route"
	"github.com/sealsee/web-base/public/utils/fileUtils"
	"github.com/sealsee/web-base/public/web"
)

type CommonController struct{}

func init() {
	route.AddGroup("/common", "公共").
		AddPOSTHandel("/upload_public", fileUploadPublic, true, "文件上传-共享").
		AddPOSTHandel("/upload_private", fileUploadPrivate, true, "文件上传-私有").
		AddPOSTHandel("/download", fileDownload, true, "文件下载")
}

func fileUploadPublic(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	form, err := c.MultipartForm()
	if err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	files := form.File["file"]
	oper := ctx.GetUserId()
	// 上传文件-公有
	urls, originalNames, err := _upload(files, oper, false)
	if err != nil {
		json.SetErr(err.Error()).Render()
		return
	}
	json.AddData("url", strings.Join(urls, ",")).
		AddData("originalName", strings.Join(originalNames, ",")).
		Render()
}

func fileUploadPrivate(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	form, err := c.MultipartForm()
	if err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	files := form.File["file"]
	oper := ctx.GetUserId()
	// 上传文件-私有
	urls, originalNames, err := _upload(files, oper, true)
	if err != nil {
		json.SetErr(err.Error()).Render()
		return
	}
	json.AddData("url", strings.Join(urls, ",")).
		AddData("originalName", strings.Join(originalNames, ",")).
		Render()
}

func fileDownload(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	oper := ctx.GetUserId()
	urlpath := ctx.PostForm("url")
	if urlpath == "" || oper < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	data, err := IOFile.GetConfig().Download(urlpath)
	if err != nil {
		json.SetErrs(errs.FILE_DOWNLOAD_ERR).Render()
		return
	}
	fileName := ctx.PostForm("fileName")
	if fileName == "" {
		fileName = time.Now().Format("20060102030405")
	}
	fileName += filepath.Ext(urlpath)
	fileName = url.QueryEscape(fileName)
	c.Header("Content-Type", "application/octet-stream;charset=UTF-8")
	c.Header("Content-Length", strconv.Itoa(len(data)))
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=utf-8''%s", fileName))
	c.Data(http.StatusOK, "application/octet-stream;charset=UTF-8", data)
}

func _upload(files []*multipart.FileHeader, oper int64, isPrivate bool) (urls []string, originalNames []string, err error) {
	if len(files) < 1 {
		return nil, nil, errors.New(errs.PARAM_INVALID_ERR[1])
	}
	// 校验个数
	if len(files) > common.UploadMaxNum {
		return nil, nil, errors.New(errs.FILE_UPLOAD_MAX_NUM[1])
	}
	// 校验大小
	for _, file := range files {
		if file.Size > common.UploadMaxSize {
			return nil, nil, errors.New(errs.FILE_UPLOAD_MAX_SIZE[1])
		}
	}
	for _, file := range files {
		open, err := file.Open()
		if err != nil {
			return nil, nil, errors.New(errs.FILE_UPLOAD_ERR[1])
		}
		defer open.Close()
		extension := fileUtils.GetExtension(file)
		url, err := IOFile.GetConfig().Upload(open, strconv.FormatInt(oper, 10), extension, isPrivate)
		if err != nil {
			return nil, nil, errors.New(errs.FILE_UPLOAD_ERR[1])
		}
		urls = append(urls, url)
		originalNames = append(originalNames, file.Filename)
	}
	return
}
