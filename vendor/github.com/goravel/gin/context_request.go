package gin

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
	contractsfilesystem "github.com/goravel/framework/contracts/filesystem"
	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/log"
	contractsession "github.com/goravel/framework/contracts/session"
	contractsvalidate "github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/filesystem"
	"github.com/goravel/framework/support/json"
	"github.com/goravel/framework/support/str"
	"github.com/goravel/framework/validation"
	"github.com/spf13/cast"
)

var contextRequestPool = sync.Pool{New: func() any {
	return &ContextRequest{
		log:        LogFacade,
		validation: ValidationFacade,
	}
}}

type ContextRequest struct {
	ctx        *Context
	instance   *gin.Context
	httpBody   map[string]any
	log        log.Log
	validation contractsvalidate.Validation
}

func NewContextRequest(ctx *Context, log log.Log, validation contractsvalidate.Validation) contractshttp.ContextRequest {
	request := contextRequestPool.Get().(*ContextRequest)
	httpBody, err := getHttpBody(ctx)
	if err != nil {
		log.Error(fmt.Sprintf("%+v", err))
	}
	request.ctx = ctx
	request.instance = ctx.instance
	request.httpBody = httpBody
	request.log = log
	request.validation = validation
	return request
}

func (r *ContextRequest) Abort(code ...int) {
	realCode := contractshttp.DefaultAbortStatus
	if len(code) > 0 {
		realCode = code[0]
	}

	r.instance.AbortWithStatus(realCode)
}

// DEPRECATED: Use Abort instead.
func (r *ContextRequest) AbortWithStatus(code int) {
	r.instance.AbortWithStatus(code)
}

// DEPRECATED: Use Response().Json().Abort() instead.
func (r *ContextRequest) AbortWithStatusJson(code int, jsonObj any) {
	r.instance.AbortWithStatusJSON(code, jsonObj)
}

func (r *ContextRequest) All() map[string]any {
	var (
		dataMap  = make(map[string]any)
		queryMap = make(map[string]any)
	)

	for key, query := range r.instance.Request.URL.Query() {
		queryMap[key] = strings.Join(query, ",")
	}

	for _, param := range r.instance.Params {
		dataMap[param.Key] = param.Value
	}
	for k, v := range queryMap {
		dataMap[k] = v
	}
	for k, v := range r.httpBody {
		dataMap[k] = v
	}

	return dataMap
}

func (r *ContextRequest) Bind(obj any) error {
	return r.instance.ShouldBind(obj)
}

func (r *ContextRequest) BindQuery(obj any) error {
	return r.instance.ShouldBindQuery(obj)
}

func (r *ContextRequest) Cookie(key string, defaultValue ...string) string {
	cookie, err := r.instance.Cookie(key)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}

		return ""
	}

	return cookie
}

func (r *ContextRequest) Form(key string, defaultValue ...string) string {
	if len(defaultValue) == 0 {
		return r.instance.PostForm(key)
	}

	return r.instance.DefaultPostForm(key, defaultValue[0])
}

func (r *ContextRequest) File(name string) (contractsfilesystem.File, error) {
	file, err := r.instance.FormFile(name)
	if err != nil {
		return nil, err
	}

	return filesystem.NewFileFromRequest(file)
}

func (r *ContextRequest) Files(name string) ([]contractsfilesystem.File, error) {
	form, err := r.instance.MultipartForm()
	if err != nil {
		return nil, err
	}

	if files, ok := form.File[name]; ok && len(files) > 0 {
		var result []contractsfilesystem.File
		for i := range files {
			var file contractsfilesystem.File
			file, err = filesystem.NewFileFromRequest(files[i])
			if err != nil {
				return nil, err
			}
			result = append(result, file)
		}

		return result, nil
	}

	return nil, http.ErrMissingFile
}

func (r *ContextRequest) FullUrl() string {
	prefix := "https://"
	if r.instance.Request.TLS == nil {
		prefix = "http://"
	}

	if r.instance.Request.Host == "" {
		return ""
	}

	return prefix + r.instance.Request.Host + r.instance.Request.RequestURI
}

func (r *ContextRequest) Header(key string, defaultValue ...string) string {
	header := r.instance.GetHeader(key)
	if header != "" {
		return header
	}

	if len(defaultValue) == 0 {
		return ""
	}

	return defaultValue[0]
}

func (r *ContextRequest) Headers() http.Header {
	return r.instance.Request.Header
}

func (r *ContextRequest) Host() string {
	return r.instance.Request.Host
}

func (r *ContextRequest) HasSession() bool {
	_, ok := r.ctx.Value(sessionKey).(contractsession.Session)
	return ok
}

func (r *ContextRequest) Json(key string, defaultValue ...string) string {
	var data map[string]any
	if err := r.Bind(&data); err != nil {
		if len(defaultValue) == 0 {
			return ""
		} else {
			return defaultValue[0]
		}
	}

	if value, exist := data[key]; exist {
		return cast.ToString(value)
	}

	if len(defaultValue) == 0 {
		return ""
	}

	return defaultValue[0]
}

func (r *ContextRequest) Method() string {
	return r.instance.Request.Method
}

func (r *ContextRequest) Name() string {
	return r.Info().Name
}

func (r *ContextRequest) Next() {
	r.instance.Next()
}

func (r *ContextRequest) Query(key string, defaultValue ...string) string {
	if len(defaultValue) > 0 {
		return r.instance.DefaultQuery(key, defaultValue[0])
	}

	return r.instance.Query(key)
}

func (r *ContextRequest) QueryInt(key string, defaultValue ...int) int {
	if val, ok := r.instance.GetQuery(key); ok {
		return cast.ToInt(val)
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return 0
}

func (r *ContextRequest) QueryInt64(key string, defaultValue ...int64) int64 {
	if val, ok := r.instance.GetQuery(key); ok {
		return cast.ToInt64(val)
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return 0
}

func (r *ContextRequest) QueryBool(key string, defaultValue ...bool) bool {
	if value, ok := r.instance.GetQuery(key); ok {
		return stringToBool(value)
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return false
}

func (r *ContextRequest) QueryArray(key string) []string {
	return r.instance.QueryArray(key)
}

func (r *ContextRequest) QueryMap(key string) map[string]string {
	return r.instance.QueryMap(key)
}

func (r *ContextRequest) Queries() map[string]string {
	queries := make(map[string]string)

	for key, query := range r.instance.Request.URL.Query() {
		queries[key] = strings.Join(query, ",")
	}

	return queries
}

func (r *ContextRequest) Origin() *http.Request {
	return r.instance.Request
}

func (r *ContextRequest) OriginPath() string {
	return colonToBracket(r.instance.FullPath())
}

func (r *ContextRequest) Path() string {
	return r.instance.Request.URL.Path
}

func (r *ContextRequest) Info() contractshttp.Info {
	methodToInfo, exist := routes[r.OriginPath()]
	if !exist {
		return contractshttp.Info{}
	}

	method := r.Method()

	methodsToTry := []string{
		method,
		contractshttp.MethodAny,
		contractshttp.MethodResource,
	}

	if method == contractshttp.MethodGet || method == contractshttp.MethodHead {
		methodsToTry = append([]string{contractshttp.MethodGet + "|" + contractshttp.MethodHead}, methodsToTry...)
	}

	for _, tryMethod := range methodsToTry {
		if info, exist := methodToInfo[tryMethod]; exist {
			info.Method = r.Method()
			return info
		}
	}

	return contractshttp.Info{
		Method: r.Method(),
		Path:   r.OriginPath(),
	}
}

func (r *ContextRequest) Input(key string, defaultValue ...string) string {
	valueFromHttpBody := r.getValueFromHttpBody(key)
	if valueFromHttpBody != nil {
		switch reflect.ValueOf(valueFromHttpBody).Kind() {
		case reflect.Map:
			valueFromHttpBodyObByte, err := json.Marshal(valueFromHttpBody)
			if err != nil {
				return ""
			}

			return string(valueFromHttpBodyObByte)
		case reflect.Slice:
			return strings.Join(cast.ToStringSlice(valueFromHttpBody), ",")
		default:
			return cast.ToString(valueFromHttpBody)
		}
	}

	if value, exist := r.instance.GetQuery(key); exist {
		return value
	}

	if value, exist := r.instance.Params.Get(key); exist {
		return value
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return ""
}

func (r *ContextRequest) InputArray(key string, defaultValue ...[]string) []string {
	if valueFromHttpBody := r.getValueFromHttpBody(key); valueFromHttpBody != nil {
		if value := cast.ToStringSlice(valueFromHttpBody); value == nil {
			return []string{}
		} else {
			return value
		}
	}

	if value, exist := r.instance.GetQueryArray(key); exist {
		if len(value) == 1 && value[0] == "" {
			return []string{}
		}

		return value
	}

	if value, exist := r.instance.Params.Get(key); exist {
		return str.Of(value).Split(",")
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return []string{}
}

func (r *ContextRequest) InputMap(key string, defaultValue ...map[string]any) map[string]any {
	if valueFromHttpBody := r.getValueFromHttpBody(key); valueFromHttpBody != nil {
		return cast.ToStringMap(valueFromHttpBody)
	}

	if _, exist := r.instance.GetQuery(key); exist {
		valueStr := r.instance.QueryMap(key)
		var value = make(map[string]any)
		for k, v := range valueStr {
			value[k] = v
		}

		return value
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return map[string]any{}
}

func (r *ContextRequest) InputMapArray(key string, defaultValue ...[]map[string]any) []map[string]any {
	if valueFromHttpBody := r.getValueFromHttpBody(key); valueFromHttpBody != nil {
		var result = make([]map[string]any, 0)
		for _, item := range cast.ToSlice(valueFromHttpBody) {
			res, err := cast.ToStringMapE(item)
			if err != nil {
				return []map[string]any{}
			}
			result = append(result, res)
		}

		if len(result) == 0 {
			for _, item := range cast.ToStringSlice(valueFromHttpBody) {
				res, err := cast.ToStringMapE(item)
				if err != nil {
					return []map[string]any{}
				}
				result = append(result, res)
			}
		}

		return result
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return []map[string]any{}
}

func (r *ContextRequest) InputInt(key string, defaultValue ...int) int {
	value := r.Input(key)
	if value == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return cast.ToInt(value)
}

func (r *ContextRequest) InputInt64(key string, defaultValue ...int64) int64 {
	value := r.Input(key)
	if value == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return cast.ToInt64(value)
}

func (r *ContextRequest) InputBool(key string, defaultValue ...bool) bool {
	value := r.Input(key)
	if value == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return stringToBool(value)
}

func (r *ContextRequest) Ip() string {
	return r.instance.ClientIP()
}

func (r *ContextRequest) Route(key string) string {
	return r.param(key)
}

func (r *ContextRequest) RouteInt(key string) int {
	val := r.param(key)

	return cast.ToInt(val)
}

func (r *ContextRequest) RouteInt64(key string) int64 {
	val := r.param(key)

	return cast.ToInt64(val)
}

func (r *ContextRequest) Session() contractsession.Session {
	s, ok := r.ctx.Value(sessionKey).(contractsession.Session)
	if !ok {
		return nil
	}
	return s
}

func (r *ContextRequest) SetSession(session contractsession.Session) contractshttp.ContextRequest {
	r.ctx.WithValue(sessionKey, session)

	return r
}

func (r *ContextRequest) Url() string {
	return r.instance.Request.RequestURI
}

func (r *ContextRequest) Validate(rules map[string]string, options ...contractsvalidate.Option) (contractsvalidate.Validator, error) {
	if len(rules) == 0 {
		return nil, errors.New("rules can't be empty")
	}

	options = append(options, validation.Rules(rules), validation.CustomRules(r.validation.Rules()), validation.CustomFilters(r.validation.Filters()))

	dataFace, err := validate.FromRequest(r.ctx.Request().Origin())
	if err != nil {
		return nil, err
	}

	for key, query := range r.instance.Request.URL.Query() {
		if _, exist := dataFace.Get(key); !exist {
			if _, err := dataFace.Set(key, strings.Join(query, ",")); err != nil {
				return nil, err
			}
		}
	}

	for _, param := range r.instance.Params {
		if _, exist := dataFace.Get(param.Key); !exist {
			if _, err := dataFace.Set(param.Key, param.Value); err != nil {
				return nil, err
			}
		}
	}

	return r.validation.Make(r.ctx, dataFace, rules, options...)
}

func (r *ContextRequest) ValidateRequest(request contractshttp.FormRequest) (contractsvalidate.Errors, error) {
	if err := request.Authorize(r.ctx); err != nil {
		return nil, err
	}

	var options []contractsvalidate.Option
	if requestWithFilters, ok := request.(contractshttp.FormRequestWithFilters); ok {
		options = append(options, validation.Filters(requestWithFilters.Filters(r.ctx)))
	}
	if requestWithMessage, ok := request.(contractshttp.FormRequestWithMessages); ok {
		options = append(options, validation.Messages(requestWithMessage.Messages(r.ctx)))
	}
	if requestWithAttributes, ok := request.(contractshttp.FormRequestWithAttributes); ok {
		options = append(options, validation.Attributes(requestWithAttributes.Attributes(r.ctx)))
	}
	if prepareForValidation, ok := request.(contractshttp.FormRequestWithPrepareForValidation); ok {
		options = append(options, validation.PrepareForValidation(func(ctx context.Context, data contractsvalidate.Data) error {
			httpCtx, ok := ctx.(contractshttp.Context)
			if !ok {
				httpCtx = r.ctx
			}

			return prepareForValidation.PrepareForValidation(httpCtx, data)
		}))
	}

	validator, err := r.Validate(request.Rules(r.ctx), options...)
	if err != nil {
		return nil, err
	}

	if err := validator.Bind(request); err != nil {
		return nil, err
	}

	return validator.Errors(), nil
}

func (r *ContextRequest) getValueFromHttpBody(key string) any {
	if r.httpBody == nil {
		return nil
	}

	var current any
	current = r.httpBody
	keys := strings.Split(key, ".")
	for _, k := range keys {
		currentValue := reflect.ValueOf(current)
		switch currentValue.Kind() {
		case reflect.Map:
			if value := currentValue.MapIndex(reflect.ValueOf(k)); value.IsValid() {
				current = value.Interface()
			} else {
				if value := currentValue.MapIndex(reflect.ValueOf(k + "[]")); value.IsValid() {
					current = value.Interface()
				} else {
					return nil
				}
			}
		case reflect.Slice:
			if number, err := strconv.Atoi(k); err == nil {
				return cast.ToStringSlice(current)[number]
			}

			return nil
		default:
		}
	}

	return current
}

func (r *ContextRequest) param(key string) string {
	if val := r.instance.Param(key); val != "" {
		return val
	}

	for _, param := range r.instance.Params {
		if suffix, exist := strings.CutPrefix(param.Key, key); exist && strings.HasSuffix(param.Value, suffix) {
			return strings.TrimSuffix(param.Value, suffix)
		}
	}

	return ""
}

func getHttpBody(ctx *Context) (map[string]any, error) {
	request := ctx.instance.Request
	if request == nil || request.Body == nil || request.ContentLength == 0 {
		return nil, nil
	}

	contentType := ctx.instance.ContentType()
	data := make(map[string]any)
	if contentType == "application/json" {
		bodyBytes, err := io.ReadAll(request.Body)
		_ = request.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("retrieve json error: %v", err)
		}

		if len(bodyBytes) > 0 {
			if err := json.Unmarshal(bodyBytes, &data); err != nil {
				return nil, fmt.Errorf("decode json [%v] error: %v", string(bodyBytes), err)
			}

			request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
	}

	if contentType == "multipart/form-data" {
		if request.PostForm == nil {
			const defaultMemory = 32 << 20
			if err := request.ParseMultipartForm(defaultMemory); err != nil {
				return nil, fmt.Errorf("parse multipart form error: %v", err)
			}
		}
		for k, v := range request.PostForm {
			if len(v) > 1 {
				data[k] = v
			} else if len(v) == 1 {
				data[k] = v[0]
			}
		}
		for k, v := range request.MultipartForm.File {
			if len(v) > 1 {
				data[k] = v
			} else if len(v) == 1 {
				data[k] = v[0]
			}
		}
	}

	if contentType == "application/x-www-form-urlencoded" {
		if request.PostForm == nil {
			if err := request.ParseForm(); err != nil {
				return nil, fmt.Errorf("parse form error: %v", err)
			}
		}
		for k, v := range request.PostForm {
			if len(v) > 1 {
				data[k] = v
			} else if len(v) == 1 {
				data[k] = v[0]
			}
		}
	}

	return data, nil
}

func stringToBool(value string) bool {
	return value == "1" || value == "true" || value == "on" || value == "yes"
}
