package main

import (
	"log"
	"sync"
	"time"
)

type TokenBucket struct {
	rate         int64
	capacity     int64
	tokens       int64
	lastTokenSec int64

	lock sync.Mutex
}

func (l *TokenBucket) Allow() bool {
	l.lock.Lock()
	defer l.lock.Unlock()

	now := time.Now().UnixNano()
	l.tokens = l.tokens + (now-l.lastTokenSec)*l.rate
	if l.tokens > l.capacity {
		l.tokens = l.capacity
	}
	l.lastTokenSec = now
	if l.tokens > 0 {
		return true
	}
	return false
}

func (l *TokenBucket) Set(r, c int64) {
	l.rate = r
	l.capacity = c
	l.tokens = 0
	l.lastTokenSec = time.Now().Unix()
}

func main() {
	var wg sync.WaitGroup
	var tb TokenBucket

	tb.Set(1, 2)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		log.Println("create a response: ", i)

		go func(i int) {
			if tb.Allow() {
				log.Println("request a response: ", i)
			}
			wg.Done()
		}(i)

		time.Sleep(time.Microsecond)
	}

	wg.Wait()
}
