package cron

import (
    "container/heap"
    "sync"
    "time"
)

type Delayed interface {
	Comparable
	GetDelay() int64
}

type DelayQueue struct {
	queue *PriorityQueue
	lock  sync.Mutex
    timerMap map[*time.Timer]*time.Timer
}

func NewDelayQueue() *DelayQueue {
	return &DelayQueue{queue: &PriorityQueue{}, timerMap: make(map[*time.Timer]*time.Timer)}
}

// Add delay task
func (d *DelayQueue) Offer(ele Delayed) {
    log("DelayQueue Offer Exec")
	d.lock.Lock()
	defer d.lock.Unlock()
    heap.Push(d.queue, ele)
	first := (*d.queue)[0].(Delayed)
	if first == ele {
        d.signal()
    }
}


// Gets a deferred task; no task blocks a timeout for a long time，if timeout return nil
// If there is a task in the queue and it has run out of time, the task is returned immediately
// Timeout units are seconds
func (d *DelayQueue) TakeWithTimeout(timeout int64) Delayed {
    log("DelayQueue TakeWithTimeout Exec")
	d.lock.Lock()
	defer d.lock.Unlock()
	for {
		var first Delayed
		if d.queue.Len() > 0 {
			first = (*d.queue)[0].(Delayed)
		}
		if first == nil {
			if timeout <= 0 {
				return nil
			} else {
				timeout = d.waitWithTimeout(timeout)
			}
		} else {
			delay := first.GetDelay()
			if delay <= 0 {
				return heap.Pop(d.queue).(Delayed)
			} else if timeout <= 0 {
				return nil
			}
			first = nil
			if timeout < delay {
				timeout = d.waitWithTimeout(timeout)
			} else {
				timeout = d.waitWithTimeout(delay)
			}
            log("DelayQueue TakeWithTimeout surplus timeout:", timeout)
		}
	}
}


// wait and release lock and set timeout wake up
func (d *DelayQueue) waitWithTimeout(timeout int64) int64 {
    log("DelayQueue waitWithTimeout:", timeout)
    timer := time.NewTimer(time.Duration(timeout) * time.Second)
    d.timerMap[timer] = timer
    d.lock.Unlock()
    absTimeout := time.Now().Unix() + timeout
    select {
    case <-timer.C:
        d.lock.Lock()
        delete(d.timerMap, timer)
        return absTimeout - time.Now().Unix()
    }
}

// Just iterate out a timer to wake up directly
func (d *DelayQueue) signal() {
    for _, val := range d.timerMap {
        val.Reset(0)
        break
    }
}

// Gets the number of tasks in the queue
func (d *DelayQueue) Len() int {
	return d.queue.Len()
}

func log(msg string, keysAndValues ...interface{}) {
    if false {
        logger.Info(msg, keysAndValues)
    }
}

