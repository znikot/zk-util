package kttp

import (
	"bytes"
	"crypto/tls"
	"io"
	"mime"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/goccy/go-json"

	"github.com/anaskhan96/soup"
	"github.com/google/uuid"
	"github.com/spf13/cast"
)

type Form url.Values
type MultiPartForm url.Values

type Option func(tp *http.Transport)

// set connect timeout
func WithConnectTimeout(timeout time.Duration) Option {
	return func(tp *http.Transport) {
		httpDialer.Timeout = timeout
		tp.DialContext = httpDialer.DialContext
	}
}

// set keepalive
func WithKeepAlive(interval time.Duration) Option {
	return func(tp *http.Transport) {
		httpDialer.KeepAlive = interval
		tp.DialContext = httpDialer.DialContext
	}
}

// set keepalives
func WithKeepAlives(kp bool) Option {
	return func(tp *http.Transport) {
		tp.DisableKeepAlives = !kp
	}
}

// set MaxIdleConnsPerHost
func WithMaxIdleConnsPerHost(mi int) Option {
	return func(tp *http.Transport) {
		tp.MaxIdleConnsPerHost = mi
	}
}

// set MaxIdleConns
func WithMaxIdleConns(mc int) Option {
	return func(tp *http.Transport) {
		tp.MaxIdleConns = mc
	}
}

// set IdleConnTimeout
func WithIdleConnTimeout(timeout time.Duration) Option {
	return func(tp *http.Transport) {
		tp.IdleConnTimeout = timeout
	}
}

// set ExpectContinueTimeout
func WithExpectContinueTimeout(timeout time.Duration) Option {
	return func(tp *http.Transport) {
		tp.ExpectContinueTimeout = timeout
	}
}

// set TLSHandshakeTimeout
func WithTLSHandshakeTimeout(timeout time.Duration) Option {
	return func(tp *http.Transport) {
		tp.TLSHandshakeTimeout = timeout
	}
}

// set ResponseHeaderTimeout
func WithResponseHeaderTimeout(timeout time.Duration) Option {
	return func(tp *http.Transport) {
		tp.ResponseHeaderTimeout = timeout
	}
}

// set InsecureSkipVerify
func WithInsecureSkipVerify(skip bool) Option {
	return func(tp *http.Transport) {
		if tp.TLSClientConfig == nil {
			tp.TLSClientConfig = &tls.Config{InsecureSkipVerify: skip}
		} else {
			tp.TLSClientConfig.InsecureSkipVerify = skip
		}
	}
}

// http net dialer
var httpDialer *net.Dialer = &net.Dialer{}

// http client
var httpClient *http.Client = &http.Client{}

// set client options
func TransportOptions(opts ...Option) {
	transport := &http.Transport{DialContext: httpDialer.DialContext}

	// set options
	for _, opt := range opts {
		opt(transport)
	}

	// httpClient.Transport = transport
	httpClient = &http.Client{Transport: transport}
}

type Request struct {
	raw    *http.Request
	client *http.Client

	header http.Header
}

// create new http request with default http client.
func NewRequest(url string, vars PathVar, body any) *Request {
	r := &Request{
		header: make(http.Header),
	}

	url = FillPathVariables(url, vars)

	r.raw, _ = http.NewRequest("", url, r.body(body))

	return r.WithClient(httpClient)
}

// set request header
func (r *Request) SetHeader(name, value string) *Request {
	r.header.Set(name, value)
	return r
}

// add request header
func (r *Request) AddHeader(name, value string) *Request {
	r.header.Add(name, value)
	return r
}

// add cookie
func (r *Request) AddCookie(cookie *http.Cookie) *Request {
	r.raw.AddCookie(cookie)
	return r
}

// delete header
func (r *Request) DelHeader(name string) *Request {
	r.header.Del(name)
	return r
}

// set header with function
func (r *Request) SetHeaderFunc(hfunc func(h http.Header)) *Request {
	if hfunc != nil {
		hfunc(r.raw.Header)
	}

	return r
}

// set http client
func (r *Request) WithClient(client *http.Client) *Request {
	r.client = client
	return r
}

// set body for post, put
func (r *Request) body(v any) io.Reader {
	if v == nil {
		return nil
	}

	switch v := v.(type) {
	case io.Reader:
		// r.raw.Body = v
		return v
	case []byte:
		return r.bytesBody(v)
	case string:
		return r.bytesBody([]byte(v))
	case Form:
		return r.formBody(v)
	case MultiPartForm:
		return r.multiFormBody(v)
	default:
		return r.jsonBody(v)
	}
}

// bytes body
func (r *Request) bytesBody(bodyBytes []byte) io.Reader {
	if l := len(bodyBytes); l > 0 {
		return bytes.NewBuffer(bodyBytes)
		// reader := bytes.NewReader(bodyBytes)
		// r.raw.Body = io.NopCloser(reader)
		// r.SetHeader("Content-Length", strconv.Itoa(l))
	}
	return nil
}

// json body
func (r *Request) jsonBody(v any) io.Reader {
	// convert to json
	bodyBytes, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	// r.bytesBody(bodyBytes)
	// set content type to json
	r.SetHeader("Content-Type", "application/json; charset=utf-8")

	return r.bytesBody(bodyBytes)
}

// form body
func (r *Request) formBody(form Form) io.Reader {
	params := url.Values{}
	for k, v := range form {
		for _, vi := range v {
			params.Add(k, cast.ToString(vi))
		}
	}
	r.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	return r.bytesBody([]byte(params.Encode()))
}

// multipart form body
func (r *Request) multiFormBody(form MultiPartForm) io.Reader {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	// add file part
	for k, v := range form {
		if k == "[file]" {
			part, _ := writer.CreateFormFile("file", "file_encode")
			part.Write([]byte(v[0]))
		}
	}
	//Adds additional parameters
	for k, v := range form {
		if k != "[file]" {
			for _, vi := range v {
				writer.WriteField(k, cast.ToString(vi))
			}
		}
	}
	writer.Close()

	// r.raw.Body = io.NopCloser(body)

	r.AddHeader("Content-Type", writer.FormDataContentType())

	return body
}

// do get
func (r *Request) Get() (*Response, error) {
	r.raw.Method = http.MethodGet
	r.raw.Header = r.header
	return r.do()
}

// do post
func (r *Request) Post() (*Response, error) {
	r.raw.Method = http.MethodPost
	r.raw.Header = r.header
	return r.do()
}

// do put
func (r *Request) Put() (*Response, error) {
	r.raw.Method = http.MethodPut
	r.raw.Header = r.header
	return r.do()
}

// do patch
func (r *Request) Patch() (*Response, error) {
	r.raw.Method = http.MethodPatch
	r.raw.Header = r.header
	return r.do()
}

// do delete
func (r *Request) Delete() (*Response, error) {
	r.raw.Method = http.MethodDelete
	r.raw.Header = r.header
	return r.do()
}

func (r *Request) do() (*Response, error) {
	resp, err := r.client.Do(r.raw)
	return &Response{raw: resp}, err
}

type Response struct {
	raw *http.Response
}

// close raw body
func (r *Response) Close() {
	r.raw.Body.Close()
}

// get status code
func (r *Response) StatusCode() int {
	return r.raw.StatusCode
}

// get status text
func (r *Response) Status() string {
	return r.raw.Status
}

// get header
func (r *Response) GetHeader(name string) string {
	return r.raw.Header.Get(name)
}

// get cookies
func (r *Response) GetCookies(name string) []*http.Cookie {
	return r.raw.Cookies()
}

// as reader
func (r *Response) AsReader() io.ReadCloser {
	return r.raw.Body
}

// as bytes
func (r *Response) AsBytes() ([]byte, error) {
	return io.ReadAll(r.raw.Body)
}

// as string
func (r *Response) AsString() (string, error) {
	buf, err := io.ReadAll(r.raw.Body)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

// unmarshal http response data to interface{}. this method will close raw response body.
//
//	resultObj := Result{}
//	if err := r.AsJson(&resultObj); err!= nil {
//		fmt.Printf("read response fail: %v\n", err)
//	} else {
//		fmt.Printf("read response success: %v\n", resultObj)
//	}
func (r *Response) AsJson(v any) error {
	return json.NewDecoder(r.raw.Body).Decode(v)
}

// as html dom. this method will close raw response body.
//
// use github.com/anaskhor6/soup to parse html.
func (r *Response) AsDom() (*soup.Root, error) {
	str, err := r.AsString()
	if err != nil {
		return nil, err
	}
	root := soup.HTMLParse(str)
	return &root, nil
}

// save http response to file. this method will close raw response body.
func (r *Response) AsFile(location, name string) error {
	// if len(name)==0, use uuid for name
	if len(name) == 0 {
		name = ExtractFileName(r.raw.Header)
		if len(name) == 0 {
			name = strings.ReplaceAll(uuid.New().String(), "-", "")
		}
	}

	// if no extension in name
	if path.Ext(name) == "" {
		// get content type
		contentType := r.raw.Header.Get("Content-Type")
		// get extension by type
		extensions, err := mime.ExtensionsByType(contentType)
		if err != nil {
			// ignore
		} else if len(extensions) > 0 {
			// use first
			name += extensions[0]
		}
	}

	// open file
	file, err := os.OpenFile(path.Join(location, name), os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	defer file.Close()

	// copy data to file
	_, err = io.Copy(file, r.raw.Body)

	return err
}
