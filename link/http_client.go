package link

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/laydong/toolpkg"
	"github.com/laydong/toolpkg/logx"
	"github.com/laydong/toolpkg/utils"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"
)

var b3Headers = []string{toolpkg.XtraceKey, toolpkg.RequestIdKey, "x-request-id", "x-b3-traceid", "x-b3-spanid", "x-b3-parentspanid", "x-b3-sampled", "x-b3-flags", "x-ot-span-context", "x-huayu-traffic-tag"}

// CONST String map of options
var CONST = map[string]int{
	"OPT_AUTOREFERER":       toolpkg.OPT_AUTOREFERER,
	"OPT_FOLLOWLOCATION":    toolpkg.OPT_FOLLOWLOCATION,
	"OPT_CONNECTTIMEOUT":    toolpkg.OPT_CONNECTTIMEOUT,
	"OPT_CONNECTTIMEOUT_MS": toolpkg.OPT_CONNECTTIMEOUT_MS,
	"OPT_MAXREDIRS":         toolpkg.OPT_MAXREDIRS,
	"OPT_PROXYTYPE":         toolpkg.OPT_PROXYTYPE,
	"OPT_TIMEOUT":           toolpkg.OPT_TIMEOUT,
	"OPT_TIMEOUT_MS":        toolpkg.OPT_TIMEOUT_MS,
	"OPT_COOKIEJAR":         toolpkg.OPT_COOKIEJAR,
	"OPT_INTERFACE":         toolpkg.OPT_INTERFACE,
	"OPT_PROXY":             toolpkg.OPT_PROXY,
	"OPT_REFERER":           toolpkg.OPT_REFERER,
	"OPT_USERAGENT":         toolpkg.OPT_USERAGENT,
	"OPT_REDIRECT_POLICY":   toolpkg.OPT_REDIRECT_POLICY,
	"OPT_PROXY_FUNC":        toolpkg.OPT_PROXY_FUNC,
	"OPT_DEBUG":             toolpkg.OPT_DEBUG,
	"OPT_UNSAFE_TLS":        toolpkg.OPT_UNSAFE_TLS,
	"OPT_CONTEXT":           toolpkg.OPT_CONTEXT,
}

// Default options for any clients.
var defaultOptions = map[int]interface{}{
	toolpkg.OPT_FOLLOWLOCATION: true,
	toolpkg.OPT_MAXREDIRS:      10,
	toolpkg.OPT_AUTOREFERER:    true,
	toolpkg.OPT_USERAGENT:      "client",
	toolpkg.OPT_COOKIEJAR:      true,
	toolpkg.OPT_DEBUG:          false,
}

// These options affect transport, transport may not be reused if you change any
// of these options during a request.
var transportOptions = []int{
	toolpkg.OPT_CONNECTTIMEOUT,
	toolpkg.OPT_CONNECTTIMEOUT_MS,
	toolpkg.OPT_PROXYTYPE,
	toolpkg.OPT_TIMEOUT,
	toolpkg.OPT_TIMEOUT_MS,
	toolpkg.OPT_INTERFACE,
	toolpkg.OPT_PROXY,
	toolpkg.OPT_PROXY_FUNC,
	toolpkg.OPT_UNSAFE_TLS,
}

// These options affect cookie jar, jar may not be reused if you change any of
// these options during a request.
var jarOptions = []int{
	toolpkg.OPT_COOKIEJAR,
}

func DefaultClient() *HttpClient {
	return &HttpClient{
		reuseTransport: true,
		reuseJar:       true,
		lock:           new(sync.Mutex),
	}
}

// Response Thin wrapper of http.Response(can also be used as http.Response).
type Response struct {
	*http.Response
}

// ReadAll Read response body into a byte slice.
func (client *Response) ReadAll() ([]byte, error) {
	var reader io.ReadCloser
	var err error
	switch client.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(client.Body)
		if err != nil {
			return nil, err
		}
	default:
		reader = client.Body
	}

	defer reader.Close()
	return ioutil.ReadAll(reader)
}

// ToString Read response body into string.
func (client *Response) ToString() (string, error) {
	bytes, err := client.ReadAll()
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// Prepare a request.
func prepareRequest(method string, url_ string, headers map[string]string,
	body io.Reader, options map[int]interface{}) (*http.Request, error) {
	req, err := http.NewRequest(method, url_, body)

	if err != nil {
		return nil, err
	}

	// OPT_REFERER
	if referer, ok := options[toolpkg.OPT_REFERER]; ok {
		if refererStr, ok := referer.(string); ok {
			req.Header.Set("Referer", refererStr)
		}
	}

	// OPT_USERAGENT
	if useragent, ok := options[toolpkg.OPT_USERAGENT]; ok {
		if useragentStr, ok := useragent.(string); ok {
			req.Header.Set("User-Agent", useragentStr)
		}
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return req, nil
}

// Prepare a transport.
//
// Handles timemout, proxy and maybe other transport related options here.
func prepareTransport(options map[int]interface{}) (http.RoundTripper, error) {
	transport := &http.Transport{}

	var connectTimeout time.Duration

	if connectTimeoutMS_, ok := options[toolpkg.OPT_CONNECTTIMEOUT_MS]; ok {
		if connectTimeoutMS, ok := connectTimeoutMS_.(int); ok {
			connectTimeout = time.Duration(connectTimeoutMS) * time.Millisecond
		} else {
			return nil, fmt.Errorf("OPT_CONNECTTIMEOUT_MS must be int")
		}
	} else if connectTimeout_, ok := options[toolpkg.OPT_CONNECTTIMEOUT]; ok {
		if connectTimeout, ok = connectTimeout_.(time.Duration); !ok {
			if connectTimeoutS, ok := connectTimeout_.(int); ok {
				connectTimeout = time.Duration(connectTimeoutS) * time.Second
			} else {
				return nil, fmt.Errorf("OPT_CONNECTTIMEOUT must be int or time.Duration")
			}
		}
	}

	var timeout time.Duration

	if timeoutMS_, ok := options[toolpkg.OPT_TIMEOUT_MS]; ok {
		if timeoutMS, ok := timeoutMS_.(int); ok {
			timeout = time.Duration(timeoutMS) * time.Millisecond
		} else {
			return nil, fmt.Errorf("OPT_TIMEOUT_MS must be int")
		}
	} else if timeout_, ok := options[toolpkg.OPT_TIMEOUT]; ok {
		if timeout, ok = timeout_.(time.Duration); !ok {
			if timeoutS, ok := timeout_.(int); ok {
				timeout = time.Duration(timeoutS) * time.Second
			} else {
				return nil, fmt.Errorf("OPT_TIMEOUT must be int or time.Duration")
			}
		}
	}

	// fix connect timeout(important, or it might cause a long time wait during
	//connection)
	if timeout > 0 && (connectTimeout > timeout || connectTimeout == 0) {
		connectTimeout = timeout
	}

	transport.Dial = func(network, addr string) (net.Conn, error) {
		var conn net.Conn
		var err error
		if connectTimeout > 0 {
			conn, err = net.DialTimeout(network, addr, connectTimeout)
			if err != nil {
				return nil, err
			}
		} else {
			conn, err = net.Dial(network, addr)
			if err != nil {
				return nil, err
			}
		}

		if timeout > 0 {
			conn.SetDeadline(time.Now().Add(timeout))
		}

		return conn, nil
	}

	// proxy
	if proxyFunc_, ok := options[toolpkg.OPT_PROXY_FUNC]; ok {
		if proxyFunc, ok := proxyFunc_.(func(*http.Request) (int, string, error)); ok {
			transport.Proxy = func(req *http.Request) (*url.URL, error) {
				proxyType, u_, err := proxyFunc(req)
				if err != nil {
					return nil, err
				}

				if proxyType != toolpkg.PROXY_HTTP {
					return nil, fmt.Errorf("only PROXY_HTTP is currently supported")
				}

				u_ = "http://" + u_

				u, err := url.Parse(u_)

				if err != nil {
					return nil, err
				}

				return u, nil
			}
		} else {
			return nil, fmt.Errorf("OPT_PROXY_FUNC is not a desired function")
		}
	} else {
		var proxytype int
		if proxytype_, ok := options[toolpkg.OPT_PROXYTYPE]; ok {
			if proxytype, ok = proxytype_.(int); !ok || proxytype != toolpkg.PROXY_HTTP {
				return nil, fmt.Errorf("OPT_PROXYTYPE must be int, and only PROXY_HTTP is currently supported")
			}
		}

		var proxy string
		if proxy_, ok := options[toolpkg.OPT_PROXY]; ok {
			if proxy, ok = proxy_.(string); !ok {
				return nil, fmt.Errorf("OPT_PROXY must be string")
			}

			if !strings.Contains(proxy, "://") {
				proxy = "http://" + proxy
			}
			proxyUrl, err := url.Parse(proxy)
			if err != nil {
				return nil, err
			}
			transport.Proxy = http.ProxyURL(proxyUrl)
		}
	}

	// TLS
	if unsafe_tls_, found := options[toolpkg.OPT_UNSAFE_TLS]; found {
		var unsafe_tls, _ = unsafe_tls_.(bool)
		var tls_config = transport.TLSClientConfig
		if tls_config == nil {
			tls_config = &tls.Config{}
			transport.TLSClientConfig = tls_config
		}
		tls_config.InsecureSkipVerify = unsafe_tls
	}

	return transport, nil
}

// Prepare a redirect policy.
func prepareRedirect(options map[int]interface{}) (func(req *http.Request, via []*http.Request) error, error) {
	var redirectPolicy func(req *http.Request, via []*http.Request) error

	if redirectPolicy_, ok := options[toolpkg.OPT_REDIRECT_POLICY]; ok {
		if redirectPolicy, ok = redirectPolicy_.(func(*http.Request, []*http.Request) error); !ok {
			return nil, fmt.Errorf("OPT_REDIRECT_POLICY is not a desired function")
		}
	} else {
		var followlocation bool
		if followlocation_, ok := options[toolpkg.OPT_FOLLOWLOCATION]; ok {
			if followlocation, ok = followlocation_.(bool); !ok {
				return nil, fmt.Errorf("OPT_FOLLOWLOCATION must be bool")
			}
		}

		var maxredirs int
		if maxredirs_, ok := options[toolpkg.OPT_MAXREDIRS]; ok {
			if maxredirs, ok = maxredirs_.(int); !ok {
				return nil, fmt.Errorf("OPT_MAXREDIRS must be int")
			}
		}

		redirectPolicy = func(req *http.Request, via []*http.Request) error {
			// no follow
			if !followlocation || maxredirs <= 0 {
				return &Error{
					Code:    3,
					Message: fmt.Sprintf("redirect not allowed"),
				}
			}

			if len(via) >= maxredirs {
				return &Error{
					Code:    3,
					Message: fmt.Sprintf("stopped after %d redirects", len(via)),
				}
			}

			last := via[len(via)-1]
			// keep necessary headers
			if useragent := last.Header.Get("User-Agent"); useragent != "" {
				req.Header.Set("User-Agent", useragent)
			}

			return nil
		}
	}

	return redirectPolicy, nil
}

// Prepare a cookie jar.
func prepareJar(options map[int]interface{}) (http.CookieJar, error) {
	var jar http.CookieJar
	var err error
	if optCookieJar_, ok := options[toolpkg.OPT_COOKIEJAR]; ok {
		// is bool
		if optCookieJar, ok := optCookieJar_.(bool); ok {
			// default jar
			if optCookieJar {
				jar, err = cookiejar.New(nil)
				if err != nil {
					return nil, err
				}
			}
		} else if optCookieJar, ok := optCookieJar_.(http.CookieJar); ok {
			jar = optCookieJar
		} else {
			return nil, fmt.Errorf("invalid cookiejar")
		}
	}

	return jar, nil
}

// NewHttpClient Create an HTTP client.
func NewHttpClient() *HttpClient {
	c := &HttpClient{
		reuseTransport: true,
		reuseJar:       true,
		lock:           new(sync.Mutex),
	}

	return c
}

// HttpClient Powerful and easy to use HTTP client.
type HttpClient struct {
	// Default options of client client.
	options map[int]interface{}

	// Default headers of client client.
	Headers map[string]string

	// Options of current request.
	oneTimeOptions map[int]interface{}

	// Headers of current request.
	oneTimeHeaders map[string]string

	// Cookies of current request.
	oneTimeCookies []*http.Cookie

	// Global transport of client client, might be shared between different
	// requests.
	transport http.RoundTripper

	// Global cookie jar of client client, might be shared between different
	// requests.
	jar http.CookieJar

	// Whether current request should reuse the transport or not.
	reuseTransport bool

	// Whether current request should reuse the cookie jar or not.
	reuseJar bool

	// Make requests of one client concurrent safe.
	lock *sync.Mutex

	// context 暂时使用gin.Context
	cxt *gin.Context

	// reqLog struct
	reqLog RequestLog

	withLock bool
}

type RequestLog struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	IP      []string          `json:"ip"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
	Query   interface{}       `json:"query"`
	Body    interface{}       `json:"body"`
}

// Custom error
type Error struct {
	Code    int
	Message string
}

// Implement the error interface
func (err Error) Error() string {
	return fmt.Sprintf("httpclient #%d: %s", err.Code, err.Message)
}

// Map is a mixed structure with options and headers
type Map map[interface{}]interface{}

// Parse the Map, return options and headers
func parseMap(m Map) (map[int]interface{}, map[string]string) {
	var options = make(map[int]interface{})
	var headers = make(map[string]string)

	if m == nil {
		return options, headers
	}

	for k, v := range m {
		// integer is option
		if kInt, ok := k.(int); ok {
			// don't need to validate
			options[kInt] = v
		} else if kString, ok := k.(string); ok {
			kStringUpper := strings.ToUpper(kString)
			if kInt, ok := CONST[kStringUpper]; ok {
				options[kInt] = v
			} else {
				// it should be header, but we still need to validate it's type
				if vString, ok := v.(string); ok {
					headers[kString] = vString
				}
			}
		}
	}

	return options, headers
}

// Is opt in options?
func hasOption(opt int, options []int) bool {
	for _, v := range options {
		if opt == v {
			return true
		}
	}

	return false
}

// Defaults Set default options and headers.
func (client *HttpClient) Defaults(defaults Map) *HttpClient {
	options, headers := parseMap(defaults)

	// merge options
	if client.options == nil {
		client.options = options
	} else {
		for k, v := range options {
			client.options[k] = v
		}
	}

	// merge headers
	if client.Headers == nil {
		client.Headers = headers
	} else {
		for k, v := range headers {
			client.Headers[k] = v
		}
	}

	return client
}

// Begin marks the begining of a request, it's necessary for concurrent
// requests.
func (client *HttpClient) Begin() *HttpClient {
	client.lock.Lock()
	client.withLock = true

	return client
}

// Reset the client state so that other requests can begin.
func (client *HttpClient) reset() {
	client.oneTimeOptions = nil
	client.oneTimeHeaders = nil
	client.oneTimeCookies = nil
	client.reuseTransport = true
	client.reuseJar = true

	// nil means the Begin has not been called, asume requests are not
	// concurrent.
	if client.withLock {
		client.withLock = false
		client.lock.Unlock()
	}
}

// WithOption Temporarily specify an option of the current request.
func (client *HttpClient) WithOption(k int, v interface{}) *HttpClient {
	if client.oneTimeOptions == nil {
		client.oneTimeOptions = make(map[int]interface{})
	}
	client.oneTimeOptions[k] = v

	// Conditions we cann't reuse the transport.
	if hasOption(k, transportOptions) {
		client.reuseTransport = false
	}

	// Conditions we cann't reuse the cookie jar.
	if hasOption(k, jarOptions) {
		client.reuseJar = false
	}

	return client
}

// WithOptions Temporarily specify multiple options of the current request.
func (client *HttpClient) WithOptions(m Map) *HttpClient {
	options, _ := parseMap(m)
	for k, v := range options {
		client.WithOption(k, v)
	}

	return client
}

// WithHeader Temporarily specify a header of the current request.
func (client *HttpClient) WithHeader(k string, v string) *HttpClient {
	if client.oneTimeHeaders == nil {
		client.oneTimeHeaders = make(map[string]string)
	}
	client.oneTimeHeaders[k] = v

	return client
}

// WithHeaders Temporarily specify multiple headers of the current request.
func (client *HttpClient) WithHeaders(m map[string]string) *HttpClient {
	for k, v := range m {
		client.WithHeader(k, v)
	}

	return client
}

func (client *HttpClient) WithITrace(ctx *gin.Context) *HttpClient {
	for _, k := range b3Headers {
		if v := ctx.GetHeader(k); v != "" {
			client.WithHeader(k, v)
		}
	}

	client.WithHeader("Content-Type", "application/json")
	client.cxt = ctx
	return client
}

func (client *HttpClient) WithCommonHeader(appName, appSecretKey string) *HttpClient {
	now := time.Now().Unix()
	appSign := utils.Md5(fmt.Sprintf("%s%d", appSecretKey, now))
	client.WithHeaders(map[string]string{
		"APP-NAME":  appName,
		"APP-SIGN":  appSign,
		"TIMESTAMP": fmt.Sprintf("%d", now),
	})
	return client
}

// Do Start a request, and get the response.
func (client *HttpClient) Do(method string, url string, heads map[string]string, body io.Reader) (*Response, error) {
	options := utils.MergeOptions(defaultOptions, client.options, client.oneTimeOptions)
	headers := utils.MergeHeaders(client.Headers, heads, client.oneTimeHeaders)
	cookies := client.oneTimeCookies

	var transport http.RoundTripper
	var jar http.CookieJar
	var err error

	// transport
	if client.transport == nil || !client.reuseTransport {
		transport, err = prepareTransport(options)
		if err != nil {
			client.reset()
			return nil, err
		}

		if client.reuseTransport {
			client.transport = transport
		}
	} else {
		transport = client.transport
	}

	// jar
	if client.jar == nil || !client.reuseJar {
		jar, err = prepareJar(options)
		if err != nil {
			client.reset()
			return nil, err
		}

		if client.reuseJar {
			client.jar = jar
		}
	} else {
		jar = client.jar
	}

	// release lock
	client.reset()

	redirect, err := prepareRedirect(options)
	if err != nil {
		return nil, err
	}

	c := &http.Client{
		Transport:     transport,
		CheckRedirect: redirect,
		Jar:           jar,
	}

	req, err := prepareRequest(method, url, headers, body, options)
	if err != nil {
		return nil, err
	}
	if req.Body != nil {
		errReadBody, _ := ioutil.ReadAll(req.Body)
		req.Body = ioutil.NopCloser(bytes.NewBuffer(errReadBody))
		if method == "GET" {
			_ = json.NewDecoder(bytes.NewBuffer(errReadBody)).Decode(&client.reqLog.Query)
		} else {
			_ = json.NewDecoder(bytes.NewBuffer(errReadBody)).Decode(&client.reqLog.Body)
		}
	}

	client.reqLog.Headers = headers
	client.reqLog.URL = url
	client.reqLog.Method = method

	if debugEnabled, ok := options[toolpkg.OPT_DEBUG]; ok {
		if debugEnabled.(bool) {
			dump, err := httputil.DumpRequestOut(req, true)
			if err == nil {
				fmt.Printf("%s\n", dump)
			}
		}
	}

	if jar != nil {
		jar.SetCookies(req.URL, cookies)
	} else {
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
	}

	if ctx, ok := options[toolpkg.OPT_CONTEXT]; ok {
		if c, ok := ctx.(context.Context); ok {
			req = req.WithContext(c)
		}
	}

	begin := time.Now()

	res, err := c.Do(req)

	var rspLog interface{}

	if err != nil {
		rspLog = err.Error()
	} else {
		errReadResBody, _ := ioutil.ReadAll(res.Body)
		res.Body = ioutil.NopCloser(bytes.NewBuffer(errReadResBody))
		_ = json.NewDecoder(bytes.NewBuffer(errReadResBody)).Decode(&rspLog)
	}

	end := time.Now()
	elapsed := end.Sub(begin)
	runTime := fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)
	logx.InfoSdk(client.cxt, "sdk_log",
		zap.Any("request", client.reqLog),
		zap.Any("respon", rspLog),
		zap.Any("start_time", float64(begin.UnixNano())/1e9),
		zap.Any("end_time", float64(end.UnixNano())/1e9),
		zap.String("runtime", runTime),
	)
	return &Response{res}, err
}

// Head The HEAD request
func (client *HttpClient) Head(url string) (*Response, error) {
	return client.Do("HEAD", url, nil, nil)
}

// Get The GET request
func (client *HttpClient) Get(url string, params ...interface{}) (*Response, error) {
	for _, p := range params {
		url = utils.AddParams(url, utils.ToUrlValues(p))
	}

	return client.Do("GET", url, nil, nil)
}

// Delete The DELETE request
func (client *HttpClient) Delete(url string, params ...interface{}) (*Response, error) {
	for _, p := range params {
		url = utils.AddParams(url, utils.ToUrlValues(p))
	}

	return client.Do("DELETE", url, nil, nil)
}

// Post The POST request
//
// With multipart set to true, the request will be encoded as
// "multipart/form-data".
//
// If any of the params key starts with "@", it is considered as a form file
// (similar to CURL but different).
func (client *HttpClient) Post(url string, params interface{}) (*Response, error) {
	t := utils.CheckParamsType(params)
	if t == 2 {
		return client.Do("POST", url, nil, utils.ToReader(params))
	}

	paramsValues := utils.ToUrlValues(params)
	// Post with files should be sent as multipart.
	if utils.CheckParamFile(paramsValues) {
		return client.PostMultipart(url, params)
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	body := strings.NewReader(paramsValues.Encode())

	return client.Do("POST", url, headers, body)
}

// PostMultipart Post with the request encoded as "multipart/form-data".
func (client *HttpClient) PostMultipart(url string, params interface{}) (*Response, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	paramsValues := utils.ToUrlValues(params)
	// check files
	for k, v := range paramsValues {
		for _, vv := range v {
			// is file
			if k[0] == '@' {
				err := utils.AddFormFile(writer, k[1:], vv)
				if err != nil {
					return nil, err
				}
			} else {
				_ = writer.WriteField(k, vv)
			}
		}
	}
	headers := make(map[string]string)

	headers["Content-Type"] = writer.FormDataContentType()
	err := writer.Close()
	if err != nil {
		return nil, err
	}

	return client.Do("POST", url, headers, body)
}

func (client *HttpClient) sendJson(method string, url string, data interface{}) (*Response, error) {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"

	var body []byte
	switch t := data.(type) {
	case []byte:
		body = t
	case string:
		body = []byte(t)
	default:
		var err error
		body, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	}

	return client.Do(method, url, headers, bytes.NewReader(body))
}

func (client *HttpClient) PostJson(url string, data interface{}) (*Response, error) {
	return client.sendJson("POST", url, data)
}

// Put The PUT request
func (client *HttpClient) Put(url string, body io.Reader) (*Response, error) {
	return client.Do("PUT", url, nil, body)
}

// PutJson Put json data
func (client *HttpClient) PutJson(url string, data interface{}) (*Response, error) {
	return client.sendJson("PUT", url, data)
}
