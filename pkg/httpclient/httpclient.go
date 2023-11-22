package httpclient

type RequestMethodBuilder[T any] interface {
	Post(path string) RequestBuilder[T]
	Get(path string) RequestBuilder[T]
	Patch(path string) RequestBuilder[T]
	Put(path string) RequestBuilder[T]
	Delete(path string) RequestBuilder[T]
	Authorization(token string) RequestMethodBuilder[T]
}

type RequestBuilder[T any] interface {
	Body(any) RequestBuilder[T]
	Query(key, value string) RequestBuilder[T]
	Param(key string, value any) RequestBuilder[T]
	StatusCode(int) RequestBuilder[T]
	Exec() (T, error)
}
