package context

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	msg "github.com/guilhermealegre/go-clean-arch-core-lib/pagination"
)

// IContext interface
type IContext interface {
	context.Context
	FullPath() string
	Next()
	Set(key string, value any)
	Get(key string) (value any, exists bool)
	MustGet(key string) any
	GetString(key string) (s string)
	GetBool(key string) (b bool)
	GetInt(key string) (i int)
	GetInt64(key string) (i64 int64)
	GetUint(key string) (ui uint)
	GetUint64(key string) (ui64 uint64)
	GetFloat64(key string) (f64 float64)
	GetTime(key string) (t time.Time)
	GetDuration(key string) (d time.Duration)
	GetStringSlice(key string) (ss []string)
	GetStringMap(key string) (sm map[string]any)
	GetStringMapString(key string) (sms map[string]string)
	GetStringMapStringSlice(key string) (smss map[string][]string)
	Param(key string) string
	AddParam(key, value string)
	Query(key string) (value string)
	DefaultQuery(key, defaultValue string) string
	GetQuery(key string) (string, bool)
	QueryArray(key string) (values []string)
	GetQueryArray(key string) (values []string, ok bool)
	QueryMap(key string) (dicts map[string]string)
	GetQueryMap(key string) (map[string]string, bool)
	PostForm(key string) (value string)
	DefaultPostForm(key, defaultValue string) string
	GetPostForm(key string) (string, bool)
	PostFormArray(key string) (values []string)
	GetPostFormArray(key string) (values []string, ok bool)
	PostFormMap(key string) (dicts map[string]string)
	GetPostFormMap(key string) (map[string]string, bool)
	FormFile(name string) (*multipart.FileHeader, error)
	MultipartForm() (*multipart.Form, error)
	SaveUploadedFile(file *multipart.FileHeader, dst string) error
	Bind(obj any) error
	BindJSON(obj any) error
	BindQuery(obj any) error
	BindHeader(obj any) error
	BindUri(obj any) error
	MustBindWith(obj any, b binding.Binding) error
	ShouldBind(obj any) error
	ShouldBindJSON(obj any) error
	ShouldBindQuery(obj any) error
	ShouldBindHeader(obj any) error
	ShouldBindUri(obj any) error
	ShouldBindWith(obj any, b binding.Binding) error
	ShouldBindBodyWith(obj any, bb binding.BindingBody) (err error)
	ClientIP() string
	RemoteIP() string
	ContentType() string
	Status(code int)
	Header(key, value string)
	GetHeader(key string) string
	GetRawData() ([]byte, error)
	SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool)
	SetSameSite(http.SameSite)
	Cookie(name string) (string, error)
	IndentedJSON(code int, obj any)
	JSONP(code int, obj any)
	JSON(code int, obj any)
	String(code int, format string, values ...any)
	Redirect(code int, location string)
	Data(code int, contentType string, data []byte)
	DataFromReader(code int, contentLength int64, contentType string, reader io.Reader, extraHeaders map[string]string)
	File(filepath string)
	FileFromFS(filepath string, fs http.FileSystem)
	FileAttachment(filepath, filename string)
	Stream(step func(w io.Writer) bool) bool
	SetAccepted(formats ...string)
	Values(key any) any
	Params() gin.Params
	Keys() map[string]any
	Request() *http.Request
	Response() gin.ResponseWriter
	SetBody([]byte)
	GetBody() []byte
	Abort()
	AddMeta(meta any) IContext
	AddPagination(pagination *msg.Pagination) IContext
	GetMeta() any
	GetPagination() *msg.Pagination
	FromGrpc(ctx context.Context) IContext
	ToGrpc() context.Context
	RequestContext() context.Context
}
