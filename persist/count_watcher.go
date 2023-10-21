package persist

import (
	"github.com/weloe/token-go/model"
	"sync/atomic"
)

type Counter struct {
	count int64
}

func (c *Counter) Increment() {
	atomic.AddInt64(&c.count, 1)
}

func (c *Counter) Decrement() {
	atomic.AddInt64(&c.count, -1)
}

func (c *Counter) Count() int64 {
	return atomic.LoadInt64(&c.count)
}

type CountWatcher struct {
	counter *Counter
}

func NewCountWatcher() *CountWatcher {
	return &CountWatcher{counter: &Counter{count: 0}}
}

func (c *CountWatcher) GetLoginCounts() int64 {
	return c.counter.Count()
}

func (c *CountWatcher) Login(loginType string, id interface{}, tokenValue string, loginModel *model.Login) {
	c.counter.Increment()
}

func (c *CountWatcher) Logout(loginType string, id interface{}, tokenValue string) {
	c.counter.Decrement()
}

func (c *CountWatcher) Kickout(loginType string, id interface{}, tokenValue string) {
	c.counter.Decrement()
}

func (c *CountWatcher) Replace(loginType string, id interface{}, tokenValue string) {
	c.counter.Decrement()
}

func (c *CountWatcher) Ban(loginType string, id interface{}, service string, level int, time int64) {

}

func (c *CountWatcher) UnBan(loginType string, id interface{}, service string) {

}

func (c *CountWatcher) RefreshToken(tokenValue string, id interface{}, timeout int64) {

}

func (c *CountWatcher) OpenSafe(loginType string, token string, service string, time int64) {

}

func (c *CountWatcher) CloseSafe(loginType string, token string, service string) {

}
