package main

import (
	"log"
	"sync"
	"time"
)

type Counter struct {
	rate  int
	begin time.Time
	cycle time.Duration
	count int
	lock  sync.Mutex
}

func (l *Counter) Allow() bool {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.count == l.rate-1 {
		curTime := time.Now()
		if curTime.Sub(l.begin) >= l.cycle {
			l.Reset(curTime)
			return true
		} else {
			return false
		}
	} else {
		l.count++
		return true
	}
}

func (l *Counter) Set(r int, cycle time.Duration) {
	l.rate = r
	l.begin = time.Now()
	l.cycle = cycle
	l.count = 0
}

func (l *Counter) Reset(t time.Time) {
	l.begin = t
	l.count = 0
}

func main() {
	var wg sync.WaitGroup
	var lr Counter

	lr.Set(3, time.Second)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		log.Println("create a request: ", i)

		go func(i int) {
			if lr.Allow() {
				log.Println("response a request: ", i)
			}
			wg.Done()
		}(i)

		time.Sleep(200 * time.Millisecond)
	}
	wg.Wait()
}
