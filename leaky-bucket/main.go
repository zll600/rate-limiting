package main

import (
	"log"
	"math"
	"sync"
	"time"
)

type LeakyBucket struct {
	rate       float64
	capacity   float64
	water      float64
	lastLeakMs int64

	lock sync.Mutex
}

func (l *LeakyBucket) Allow() bool {
	l.lock.Lock()
	defer l.lock.Unlock()

	curTime := time.Now().UnixNano()
	eclipse := float64(curTime-l.lastLeakMs) * l.rate / 1000
	l.water = l.water - eclipse
	l.water = math.Max(0, l.water)
	l.lastLeakMs = curTime
	if (l.water + 1) < l.capacity {
		l.water++
		return true
	}

	return false
}

func (l *LeakyBucket) Set(r, c float64) {
	l.rate = r
	l.capacity = c
	l.water = 0
	l.lastLeakMs = time.Now().UnixNano() / 1e6
}

func main() {
	var wg sync.WaitGroup
	var lb LeakyBucket

	lb.Set(1.0, 3.0)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		log.Println("create a request: ", i)

		go func(i int) {
			if lb.Allow() {
				log.Println("response a request: ", i)
			}
			wg.Done()
		}(i)

		// time.Sleep(2 * time.Millisecond)
	}
	wg.Wait()
}
