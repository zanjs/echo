package echo

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/logger"

	"bytes"

	netContext "golang.org/x/net/context"
	"golang.org/x/net/websocket"
)

type (
	// Context represents context for the current request. It holds request and
	// response objects, path parameters, data and registered handler.
	Context interface {
		netContext.Context
		Request() engine.Request
		Response() engine.Response
		Socket() *websocket.Conn
		Path() string
		P(int) string
		Param(string) string
		Query(string) string
		Form(string) string
		Set(string, interface{})
		Get(string) interface{}
		Bind(interface{}) error
		Render(int, string, interface{}) error
		HTML(int, string) error
		String(int, string) error
		JSON(int, interface{}) error
		JSONBlob(int, []byte) error
		JSONP(int, string, interface{}) error
		XML(int, interface{}) error
		XMLBlob(int, []byte) error
		Attachment(string) error
		NoContent(int) error
		Redirect(int, string) error
		Error(err error)
		Handle(Context) error
		Logger() logger.Logger
		Object() *context

		SetFunc(string, interface{})
		GetFunc(string) interface{}
		Funcs() map[string]interface{}
		Reset(engine.Request, engine.Response)
		Fetch(string, interface{}) ([]byte, error)
		SetRenderer(Renderer)
	}

	context struct {
		request  engine.Request
		response engine.Response
		socket   *websocket.Conn
		path     string
		pnames   []string
		pvalues  []string
		store    store
		handler  Handler
		echo     *Echo
		funcs    map[string]interface{}
		renderer Renderer
	}

	store map[string]interface{}
)

const (
	indexPage = "index.html"
)

// NewContext creates a Context object.
func NewContext(req engine.Request, res engine.Response, e *Echo) Context {
	return &context{
		request:  req,
		response: res,
		echo:     e,
		pvalues:  make([]string, *e.maxParam),
		store:    make(store),
		handler:  notFoundHandler,
		funcs:    make(map[string]interface{}),
	}
}

func (c *context) Handle(ctx Context) error {
	return c.handler.Handle(ctx)
}

func (c *context) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c *context) Done() <-chan struct{} {
	return nil
}

func (c *context) Err() error {
	return nil
}

func (c *context) Value(key interface{}) interface{} {
	return nil
}

// Request returns *http.Request.
func (c *context) Request() engine.Request {
	return c.request
}

// Response returns *Response.
func (c *context) Response() engine.Response {
	return c.response
}

// Socket returns *websocket.Conn.
func (c *context) Socket() *websocket.Conn {
	return c.socket
}

// Path returns the registered path for the handler.
func (c *context) Path() string {
	return c.path
}

// P returns path parameter by index.
func (c *context) P(i int) (value string) {
	l := len(c.pnames)
	if i < l {
		value = c.pvalues[i]
	}
	return
}

// Param returns path parameter by name.
func (c *context) Param(name string) (value string) {
	l := len(c.pnames)
	for i, n := range c.pnames {
		if n == name && i < l {
			value = c.pvalues[i]
			break
		}
	}
	return
}

// Query returns query parameter by name.
func (c *context) Query(name string) string {
	return c.request.URL().QueryValue(name)
}

// Form returns form parameter by name.
func (c *context) Form(name string) string {
	return c.request.FormValue(name)
}

// Get retrieves data from the context.
func (c *context) Get(key string) interface{} {
	return c.store[key]
}

// Set saves data in the context.
func (c *context) Set(key string, val interface{}) {
	if c.store == nil {
		c.store = make(store)
	}
	c.store[key] = val
}

// Bind binds the request body into specified type `i`. The default binder does
// it based on Content-Type header.
func (c *context) Bind(i interface{}) error {
	return c.echo.binder.Bind(i, c)
}

// Render renders a template with data and sends a text/html response with status
// code. Templates can be registered using `Echo.SetRenderer()`.
func (c *context) Render(code int, name string, data interface{}) (err error) {
	b, err := c.Fetch(name, data)
	if err != nil {
		return
	}
	c.response.Header().Set(ContentType, TextHTMLCharsetUTF8)
	c.response.WriteHeader(code)
	c.response.Write(b)
	return
}

// HTML sends an HTTP response with status code.
func (c *context) HTML(code int, html string) (err error) {
	c.response.Header().Set(ContentType, TextHTMLCharsetUTF8)
	c.response.WriteHeader(code)
	c.response.Write([]byte(html))
	return
}

// String sends a string response with status code.
func (c *context) String(code int, s string) (err error) {
	c.response.Header().Set(ContentType, TextPlainCharsetUTF8)
	c.response.WriteHeader(code)
	c.response.Write([]byte(s))
	return
}

// JSON sends a JSON response with status code.
func (c *context) JSON(code int, i interface{}) (err error) {
	b, err := json.Marshal(i)
	if c.echo.Debug() {
		b, err = json.MarshalIndent(i, "", "  ")
	}
	if err != nil {
		return err
	}
	return c.JSONBlob(code, b)
}

// JSONBlob sends a JSON blob response with status code.
func (c *context) JSONBlob(code int, b []byte) (err error) {
	c.response.Header().Set(ContentType, ApplicationJSONCharsetUTF8)
	c.response.WriteHeader(code)
	c.response.Write(b)
	return
}

// JSONP sends a JSONP response with status code. It uses `callback` to construct
// the JSONP payload.
func (c *context) JSONP(code int, callback string, i interface{}) (err error) {
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}
	c.response.Header().Set(ContentType, ApplicationJavaScriptCharsetUTF8)
	c.response.WriteHeader(code)
	c.response.Write([]byte(callback + "("))
	c.response.Write(b)
	c.response.Write([]byte(");"))
	return
}

// XML sends an XML response with status code.
func (c *context) XML(code int, i interface{}) (err error) {
	b, err := xml.Marshal(i)
	if c.echo.Debug() {
		b, err = xml.MarshalIndent(i, "", "  ")
	}
	if err != nil {
		return err
	}
	return c.XMLBlob(code, b)
}

// XMLBlob sends a XML blob response with status code.
func (c *context) XMLBlob(code int, b []byte) (err error) {
	c.response.Header().Set(ContentType, ApplicationXMLCharsetUTF8)
	c.response.WriteHeader(code)
	c.response.Write([]byte(xml.Header))
	c.response.Write(b)
	return
}

// Attachment sends specified file as an attachment to the client.
func (c *context) Attachment(file string) (err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	_, name := filepath.Split(file)
	c.response.Header().Set(ContentDisposition, "attachment; filename="+name)
	c.response.Header().Set(ContentType, c.detectContentType(file))
	c.response.WriteHeader(http.StatusOK)
	_, err = io.Copy(c.response, f)
	return
}

// NoContent sends a response with no body and a status code.
func (c *context) NoContent(code int) error {
	c.response.WriteHeader(code)
	return nil
}

// Redirect redirects the request with status code.
func (c *context) Redirect(code int, url string) error {
	if code < http.StatusMultipleChoices || code > http.StatusTemporaryRedirect {
		return ErrInvalidRedirectCode
	}
	c.response.Redirect(url, code)
	return nil
}

// Error invokes the registered HTTP error handler. Generally used by middleware.
func (c *context) Error(err error) {
	c.echo.httpErrorHandler(err, c)
}

// Logger returns the `Logger` instance.
func (c *context) Logger() logger.Logger {
	return c.echo.logger
}

// Object returns the `context` object.
func (c *context) Object() *context {
	return c
}

func (c *context) detectContentType(name string) (t string) {
	if t = mime.TypeByExtension(filepath.Ext(name)); t == "" {
		t = OctetStream
	}
	return
}

func (c *context) reset(req engine.Request, res engine.Response) {
	c.request = req
	c.response = res
	c.store = nil
	c.funcs = make(map[string]interface{})
	c.renderer = nil
}

// Echo returns the `Echo` instance.
func (c *context) Echo() *Echo {
	return c.echo
}

func (c *context) Reset(req engine.Request, res engine.Response) {
	c.reset(req, res)
}

func (c *context) GetFunc(key string) interface{} {
	return c.funcs[key]
}

func (c *context) SetFunc(key string, val interface{}) {
	c.funcs[key] = val
}

func (c *context) Funcs() map[string]interface{} {
	return c.funcs
}

func (c *context) Fetch(name string, data interface{}) (b []byte, err error) {
	if c.renderer == nil {
		if c.echo.renderer == nil {
			return nil, ErrRendererNotRegistered
		}
		c.renderer = c.echo.renderer
	}
	buf := new(bytes.Buffer)
	err = c.renderer.Render(buf, name, data, c.funcs)
	if err != nil {
		return
	}
	b = buf.Bytes()
	return
}

// SetRenderer registers an HTML template renderer.
func (c *context) SetRenderer(r Renderer) {
	c.renderer = r
}
