package limiter

import (
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"
)

type Client struct {
	limiter          *rate.Limiter
	LastSeenAtUnixNs int64
}

type IPLimiter struct {
	Limiters typedSyncMap[string, *Client]
	rps      rate.Limit
	burst    int
}

func NewIPLimiter(rps rate.Limit, burst int) *IPLimiter {
	il := &IPLimiter{
		rps:   rps,
		burst: burst,
	}

	go func() {
		for {
			time.Sleep(time.Minute)
			now := time.Now().UnixNano()
			fiveMinutes := int64(5 * time.Minute)

			il.Limiters.Range(func(ip string, c *Client) bool {
				lastSeen := atomic.LoadInt64(&c.LastSeenAtUnixNs)
				if now-lastSeen > fiveMinutes {
					il.Limiters.Delete(ip)
				}
				return true
			})
		}
	}()
	return il
}

func (il *IPLimiter) Get(ip string) *rate.Limiter {
	now := time.Now().UnixNano()
	if c, ok := il.Limiters.Load(ip); ok {
		atomic.StoreInt64(&c.LastSeenAtUnixNs, now)
		il.Limiters.m.Store(ip, c)
		return c.limiter
	}

	l := &Client{
		limiter:          rate.NewLimiter(il.rps, il.burst),
		LastSeenAtUnixNs: now,
	}
	actual, _ := il.Limiters.LoadOrStore(ip, l)
	return actual.limiter
}
