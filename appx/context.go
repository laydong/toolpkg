package appx

import (
	"github.com/laydong/toolpkg/alarmx"
	"github.com/laydong/toolpkg/datax"
	"github.com/laydong/toolpkg/log"
	"github.com/laydong/toolpkg/tracex"
	"github.com/laydong/toolpkg/utils"
	uuid "github.com/satori/go.uuid"
	"time"
)

// Context is the carrier of request and response
type Context struct {
	*datax.MemoryContext
	*log.LogContext
	*tracex.TraceContext
	*alarmx.AlarmContext
}

// NewDefaultContext 创建 app 默认的context, spanName
func NewDefaultContext(logId string, spanName string) *Context {
	if logId == "" {
		logId = utils.Md5(uuid.NewV4().String())
	}

	tmp := &Context{
		LogContext:    log.NewLogContext(logId),
		TraceContext:  tracex.NewTraceContext(spanName, make(map[string][]string)),
		MemoryContext: datax.NewMemoryContext(),
	}

	return tmp
}

// Deadline returns the time when work done on behalf of this contextx
// should be canceled. Deadline returns ok==false when no deadline is
// set. Successive calls to Deadline return the same results.
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return
}

// Done returns a channel that's closed when work done on behalf of this
// contextx should be canceled. Done may return nil if this contextx can
// never be canceled. Successive calls to Done return the same value.
func (c *Context) Done() <-chan struct{} {
	c.SpanFinish(c.TopSpan)
	return nil
}

// Err returns a non-nil error value after Done is closed,
// successive calls to Err return the same error.
// If Done is not yet closed, Err returns nil.
// If Done is closed, Err returns a non-nil error explaining why:
// Canceled if the contextx was canceled
// or DeadlineExceeded if the contextx's deadline passed.
func (c *Context) Err() error {
	return nil
}

// Value returns the value associated with this contextx for key, or nil
// if no value is associated with key. Successive calls to Value with
// the same key returns the same result.
func (c *Context) Value(key interface{}) interface{} {
	if keyAsString, ok := key.(string); ok {
		val, _ := c.Get(keyAsString)
		return val
	}
	return nil
}
