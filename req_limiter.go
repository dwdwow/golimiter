package golimiter

import (
	"context"
	"time"
)

type ReqLimiter struct {
	l int
	d time.Duration
	c chan struct{}
}

// ReqLimiter implements a request rate limiter that allows up to a specified limit
// of concurrent requests within a given time interval.
//
// The limiter uses a buffered channel to track active requests. When a request arrives,
// it attempts to send a value to the channel. If the channel is full (limit reached),
// the request blocks until capacity becomes available. After the specified interval,
// the request is automatically removed from the channel.
func NewLimiter(interval time.Duration, limit int) *ReqLimiter {
	return &ReqLimiter{
		l: limit,
		d: interval,
		c: make(chan struct{}, limit),
	}
}

func (l *ReqLimiter) Wait(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	case l.c <- struct{}{}:
		go func() {
			time.Sleep(l.d)
			<-l.c
		}()
		return
	}
}
