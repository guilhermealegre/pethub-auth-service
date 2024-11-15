package context

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	contextDomain "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain/context"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	msg "github.com/guilhermealegre/go-clean-arch-core-lib/pagination"
	"github.com/stretchr/testify/mock"
)

func NewContextMock() *ContextMock {
	return &ContextMock{}
}

type ContextMock struct {
	mock.Mock
}

func (c *ContextMock) FullPath() string {
	args := c.Called()
	return args.Get(0).(string)
}

func (c *ContextMock) Next() {
	c.Called()
}

func (c *ContextMock) Set(key string, value any) {
	c.Called(key, value)
}

func (c *ContextMock) Get(key string) (value any, exists bool) {
	args := c.Called(key)
	return args.Get(0), args.Get(1).(bool)
}

func (c *ContextMock) MustGet(key string) any {
	args := c.Called(key)
	return args.Get(0)
}

func (c *ContextMock) GetString(key string) (s string) {
	args := c.Called(key)
	return args.Get(0).(string)
}

func (c *ContextMock) GetBool(key string) (b bool) {
	args := c.Called(key)
	return args.Get(0).(bool)
}

func (c *ContextMock) GetInt(key string) (i int) {
	args := c.Called(key)
	return args.Get(0).(int)
}

func (c *ContextMock) GetInt64(key string) (i64 int64) {
	args := c.Called(key)
	return args.Get(0).(int64)
}

func (c *ContextMock) GetUint(key string) (ui uint) {
	args := c.Called(key)
	return args.Get(0).(uint)
}

func (c *ContextMock) GetUint64(key string) (ui64 uint64) {
	args := c.Called(key)
	return args.Get(0).(uint64)
}

func (c *ContextMock) GetFloat64(key string) (f64 float64) {
	args := c.Called(key)
	return args.Get(0).(float64)
}

func (c *ContextMock) GetTime(key string) (t time.Time) {
	args := c.Called(key)
	return args.Get(0).(time.Time)
}

func (c *ContextMock) GetDuration(key string) (d time.Duration) {
	args := c.Called(key)
	return args.Get(0).(time.Duration)
}

func (c *ContextMock) GetStringSlice(key string) (ss []string) {
	args := c.Called(key)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]string)
}

func (c *ContextMock) GetStringMap(key string) (sm map[string]any) {
	args := c.Called(key)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(map[string]any)
}

func (c *ContextMock) GetStringMapString(key string) (sms map[string]string) {
	args := c.Called(key)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(map[string]string)
}

func (c *ContextMock) GetStringMapStringSlice(key string) (smss map[string][]string) {
	args := c.Called(key)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(map[string][]string)
}

func (c *ContextMock) Param(key string) string {
	args := c.Called(key)
	return args.Get(0).(string)
}

func (c *ContextMock) AddParam(key, value string) {
	c.Called(key, value)
}

func (c *ContextMock) Query(key string) (value string) {
	args := c.Called(key)
	return args.Get(0).(string)
}

func (c *ContextMock) DefaultQuery(key, defaultValue string) string {
	args := c.Called(key, defaultValue)
	return args.Get(0).(string)
}

func (c *ContextMock) GetQuery(key string) (string, bool) {
	args := c.Called(key)
	return args.Get(0).(string), args.Get(1).(bool)
}

func (c *ContextMock) QueryArray(key string) (values []string) {
	args := c.Called(key)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]string)
}

func (c *ContextMock) GetQueryArray(key string) (values []string, ok bool) {
	args := c.Called(key)
	if args.Get(0) == nil {
		return nil, args.Get(1).(bool)
	}
	return args.Get(0).([]string), args.Get(1).(bool)
}

func (c *ContextMock) QueryMap(key string) (dicts map[string]string) {
	args := c.Called(key)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(map[string]string)
}
func (c *ContextMock) GetQueryMap(key string) (map[string]string, bool) {
	args := c.Called(key)
	if args.Get(0) == nil {
		return nil, args.Get(1).(bool)
	}
	return args.Get(0).(map[string]string), args.Get(1).(bool)
}

func (c *ContextMock) PostForm(key string) (value string) {
	args := c.Called(key)
	return args.Get(0).(string)
}

func (c *ContextMock) DefaultPostForm(key, defaultValue string) string {
	args := c.Called(key, defaultValue)
	return args.Get(0).(string)
}

func (c *ContextMock) GetPostForm(key string) (string, bool) {
	args := c.Called(key)
	return args.Get(0).(string), args.Get(1).(bool)
}

func (c *ContextMock) PostFormArray(key string) (values []string) {
	args := c.Called(key)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]string)
}

func (c *ContextMock) GetPostFormArray(key string) (values []string, ok bool) {
	args := c.Called(key)
	if args.Get(0) == nil {
		return nil, args.Get(1).(bool)
	}
	return args.Get(0).([]string), args.Get(1).(bool)
}

func (c *ContextMock) PostFormMap(key string) (dicts map[string]string) {
	args := c.Called(key)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(map[string]string)
}

func (c *ContextMock) GetPostFormMap(key string) (map[string]string, bool) {
	args := c.Called(key)
	if args.Get(0) == nil {
		return nil, args.Get(1).(bool)
	}
	return args.Get(0).(map[string]string), args.Get(1).(bool)
}

func (c *ContextMock) FormFile(name string) (*multipart.FileHeader, error) {
	args := c.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*multipart.FileHeader), args.Error(1)
}

func (c *ContextMock) MultipartForm() (*multipart.Form, error) {
	args := c.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*multipart.Form), args.Error(1)
}

func (c *ContextMock) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	args := c.Called(file, dst)
	return args.Error(0)
}

func (c *ContextMock) Bind(obj any) error {
	args := c.Called(obj)
	return args.Error(0)
}

func (c *ContextMock) BindJSON(obj any) error {
	args := c.Called(obj)
	return args.Error(0)
}

func (c *ContextMock) BindQuery(obj any) error {
	args := c.Called(obj)
	return args.Error(0)
}

func (c *ContextMock) BindHeader(obj any) error {
	args := c.Called(obj)
	return args.Error(0)
}

func (c *ContextMock) BindUri(obj any) error {
	args := c.Called(obj)
	return args.Error(0)
}

func (c *ContextMock) MustBindWith(obj any, b binding.Binding) error {
	args := c.Called(obj, b)
	return args.Error(0)
}

func (c *ContextMock) ShouldBind(obj any) error {
	args := c.Called(obj)
	return args.Error(0)
}

func (c *ContextMock) ShouldBindJSON(obj any) error {
	args := c.Called(obj)
	return args.Error(0)
}

func (c *ContextMock) ShouldBindQuery(obj any) error {
	args := c.Called(obj)
	return args.Error(0)
}

func (c *ContextMock) ShouldBindHeader(obj any) error {
	args := c.Called(obj)
	return args.Error(0)
}

func (c *ContextMock) ShouldBindUri(obj any) error {
	args := c.Called(obj)
	return args.Error(0)
}

func (c *ContextMock) ShouldBindWith(obj any, b binding.Binding) error {
	args := c.Called(obj, b)
	return args.Error(0)
}

func (c *ContextMock) ShouldBindBodyWith(obj any, bb binding.BindingBody) (err error) {
	args := c.Called(obj, bb)
	return args.Error(0)
}

func (c *ContextMock) ClientIP() string {
	args := c.Called()
	return args.Get(0).(string)
}

func (c *ContextMock) RemoteIP() string {
	args := c.Called()
	return args.Get(0).(string)
}

func (c *ContextMock) ContentType() string {
	args := c.Called()
	return args.Get(0).(string)
}

func (c *ContextMock) Status(code int) {
	c.Called(code)
}

func (c *ContextMock) Header(key, value string) {
	c.Called(key, value)
}

func (c *ContextMock) GetHeader(key string) string {
	args := c.Called(key)
	return args.Get(0).(string)
}

func (c *ContextMock) GetRawData() ([]byte, error) {
	args := c.Called()
	return args.Get(0).([]byte), args.Error(1)
}

func (c *ContextMock) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	c.Called(name, value, maxAge, path, domain, secure, httpOnly)
}

func (c *ContextMock) Cookie(name string) (string, error) {
	args := c.Called(name)
	return args.Get(0).(string), args.Error(1)
}

func (c *ContextMock) IndentedJSON(code int, obj any) {
	c.Called(code, obj)
}

func (c *ContextMock) JSONP(code int, obj any) {
	c.Called(code, obj)
}

func (c *ContextMock) JSON(code int, obj any) {
	c.Called(code, obj)
}

func (c *ContextMock) String(code int, format string, values ...any) {
	c.Called(code, format, append([]interface{}{code, format}, values...))
}

func (c *ContextMock) Redirect(code int, location string) {
	c.Called(code, location)
}

func (c *ContextMock) Data(code int, contentType string, data []byte) {
	c.Called(code, contentType, data)
}

func (c *ContextMock) DataFromReader(code int, contentLength int64, contentType string, reader io.Reader, extraHeaders map[string]string) {
	c.Called(code, contentLength, contentType, reader, extraHeaders)
}

func (c *ContextMock) File(filepath string) {
	c.Called(filepath)
}

func (c *ContextMock) FileFromFS(filepath string, fs http.FileSystem) {
	c.Called(filepath, fs)
}

func (c *ContextMock) FileAttachment(filepath, filename string) {
	c.Called(filepath, filename)
}

func (c *ContextMock) Stream(step func(w io.Writer) bool) bool {
	args := c.Called(step)
	return args.Get(0).(bool)
}

func (c *ContextMock) SetAccepted(formats ...string) {
	var list []interface{}
	for _, value := range formats {
		list = append(list, value)
	}
	c.Called(list...)
}

func (c *ContextMock) Values(key any) any {
	args := c.Called(key)
	return args.Get(0)
}

func (c *ContextMock) Params() gin.Params {
	args := c.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(gin.Params)
}

func (c *ContextMock) Keys() map[string]any {
	args := c.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(map[string]any)
}

func (c *ContextMock) Request() *http.Request {
	args := c.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*http.Request)
}

func (c *ContextMock) Response() gin.ResponseWriter {
	args := c.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(gin.ResponseWriter)
}

func (c *ContextMock) Deadline() (deadline time.Time, ok bool) {
	args := c.Called()
	return args.Get(0).(time.Time), args.Get(1).(bool)
}

func (c *ContextMock) Done() <-chan struct{} {
	args := c.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(<-chan struct{})
}

func (c *ContextMock) Err() error {
	args := c.Called()
	return args.Error(0)
}

func (c *ContextMock) Value(key any) any {
	args := c.Called(key)
	return args.Get(0)
}

func (c *ContextMock) SetSameSite(site http.SameSite) {
	c.Called(site)
}

func (c *ContextMock) SetBody(b []byte) {
	c.Called(b)
}

func (c *ContextMock) GetBody() []byte {
	args := c.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]byte)
}

func (c *ContextMock) AddMeta(meta interface{}) contextDomain.IContext {
	args := c.Called(meta)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(contextDomain.IContext)
}

func (c *ContextMock) AddPagination(pagination *msg.Pagination) contextDomain.IContext {
	args := c.Called(pagination)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(contextDomain.IContext)
}

func (c *ContextMock) GetMeta() any {
	args := c.Called()
	return args.Get(0)
}

func (c *ContextMock) GetPagination() *msg.Pagination {
	args := c.Called()
	return args.Get(0).(*msg.Pagination)
}

func (c *ContextMock) Abort() {
	_ = c.Called()
}

func (c *ContextMock) FromGrpc(ctx context.Context) contextDomain.IContext {
	args := c.Called(ctx)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(contextDomain.IContext)
}

func (c *ContextMock) ToGrpc() context.Context {
	args := c.Called()
	return args.Get(0).(context.Context)
}

func (c *ContextMock) RequestContext() context.Context {
	args := c.Called()
	return args.Get(0).(context.Context)
}
