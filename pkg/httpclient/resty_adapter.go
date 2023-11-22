package httpclient

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
)

func Resty[T any](c *resty.Client) RequestMethodBuilder[T] {
	return &requestConfig[T]{
		client: c,
		query:  url.Values{},
		params: make(map[string]string),
	}
}

type requestConfig[T any] struct {
	client     *resty.Client
	body       any
	query      url.Values
	params     map[string]string
	statusCode int
	method     string
	path       string
	authToken  string
}

func (rc *requestConfig[T]) Body(body any) RequestBuilder[T] {
	rc.body = body
	return rc
}
func (rc *requestConfig[T]) Query(key, value string) RequestBuilder[T] {
	rc.query.Add(key, value)
	return rc
}

func (rc *requestConfig[T]) Param(key string, value any) RequestBuilder[T] {
	rc.params[key] = fmt.Sprintf("%v", value)
	return rc
}

func (rc *requestConfig[T]) StatusCode(code int) RequestBuilder[T] {
	rc.statusCode = code
	return rc
}

func (rc *requestConfig[T]) Post(path string) RequestBuilder[T] {
	rc.path = path
	rc.method = http.MethodPost
	return rc
}

func (rc *requestConfig[T]) Get(path string) RequestBuilder[T] {
	rc.path = path
	rc.method = http.MethodGet
	return rc
}

func (rc *requestConfig[T]) Patch(path string) RequestBuilder[T] {
	rc.path = path
	rc.method = http.MethodPatch
	return rc
}

func (rc *requestConfig[T]) Put(path string) RequestBuilder[T] {
	rc.path = path
	rc.method = http.MethodPut
	return rc
}

func (rc *requestConfig[T]) Delete(path string) RequestBuilder[T] {
	rc.path = path
	rc.method = http.MethodDelete
	return rc
}

func (rc *requestConfig[T]) Authorization(token string) RequestMethodBuilder[T] {
	rc.authToken = token
	return rc
}
func (rc *requestConfig[T]) Exec() (T, error) {
	return execRequest[T](rc)
}

func execRequest[T any](config *requestConfig[T]) (T, error) {
	var (
		result        T
		err           error
		restyResponse *resty.Response
	)

	restyRequest := config.client.R()
	if config.body != nil {
		restyRequest = restyRequest.SetBody(config.body)
	}

	if len(config.params) > 0 {
		for key, value := range config.params {
			config.path = strings.Replace(config.path, fmt.Sprintf("{%s}", key), value, 1)
		}
	}

	if len(config.query) > 0 {
		config.path = fmt.Sprintf("%s?%s", config.path, config.query.Encode())
	}

	if config.authToken != "" {
		restyRequest = restyRequest.SetAuthToken(config.authToken)
	}

	restyRequest = restyRequest.SetResult(result)

	switch config.method {
	case http.MethodGet:
		restyResponse, err = restyRequest.Get(config.path)
	case http.MethodPost:
		restyResponse, err = restyRequest.Post(config.path)
	case http.MethodPatch:
		restyResponse, err = restyRequest.Patch(config.path)
	case http.MethodPut:
		restyResponse, err = restyRequest.Put(config.path)
	case http.MethodDelete:
		restyResponse, err = restyRequest.Delete(config.path)
	default:
		return result, fmt.Errorf("wrong http method: got %s", config.method)
	}

	if err != nil {
		return result, fmt.Errorf("error on http client: %w", err)
	}

	if code := restyResponse.StatusCode(); code != config.statusCode {
		return result, fmt.Errorf("http response must return %d: got %d", config.statusCode, code)
	}

	resultPtr, ok := restyResponse.Result().(*T)

	if !ok {
		return result, fmt.Errorf("http response body must be of type %T", result)
	}

	return *resultPtr, nil
}
