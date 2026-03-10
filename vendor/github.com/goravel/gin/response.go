package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	contractshttp "github.com/goravel/framework/contracts/http"
)

type DataResponse struct {
	code        int
	contentType string
	data        []byte
	instance    *gin.Context
}

func (r *DataResponse) Render() error {
	r.instance.Data(r.code, r.contentType, r.data)

	return nil
}

func (r *DataResponse) Abort() error {
	r.instance.Abort()

	return r.Render()
}

type DownloadResponse struct {
	filename string
	filepath string
	instance *gin.Context
}

func (r *DownloadResponse) Render() error {
	r.instance.FileAttachment(r.filepath, r.filename)

	return nil
}

type FileResponse struct {
	filepath string
	instance *gin.Context
}

func (r *FileResponse) Render() error {
	r.instance.File(r.filepath)

	return nil
}

type JsonResponse struct {
	code     int
	obj      any
	instance *gin.Context
}

func (r *JsonResponse) Render() error {
	r.instance.JSON(r.code, r.obj)

	return nil
}

func (r *JsonResponse) Abort() error {
	r.instance.Abort()

	return r.Render()
}

type NoContentResponse struct {
	code     int
	instance *gin.Context
}

func (r *NoContentResponse) Render() error {
	r.instance.Status(r.code)

	return nil
}

func (r *NoContentResponse) Abort() error {
	r.instance.AbortWithStatus(r.code)

	return nil
}

type RedirectResponse struct {
	code     int
	location string
	instance *gin.Context
}

func (r *RedirectResponse) Render() error {
	r.instance.Redirect(r.code, r.location)

	return nil
}

func (r *RedirectResponse) Abort() error {
	r.instance.Abort()

	return r.Render()
}

type StringResponse struct {
	code     int
	format   string
	instance *gin.Context
	values   []any
}

func (r *StringResponse) Render() error {
	r.instance.String(r.code, r.format, r.values...)

	return nil
}

func (r *StringResponse) Abort() error {
	r.instance.Abort()

	return r.Render()
}

type HtmlResponse struct {
	data     any
	instance *gin.Context
	view     string
}

func (r *HtmlResponse) Render() error {
	r.instance.HTML(http.StatusOK, r.view, r.data)

	return nil
}

type StreamResponse struct {
	code     int
	instance *gin.Context
	writer   func(w contractshttp.StreamWriter) error
}

func (r *StreamResponse) Render() error {
	r.instance.Status(r.code)

	w := NewStreamWriter(r.instance)

	ctx := r.instance.Request.Context()
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			return r.writer(w)
		}
	}
}
