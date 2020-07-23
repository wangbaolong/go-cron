package cron

import (
    "container/heap"
    "time"
)

type Delayed interface {
	Comparable
	GetDelay() time.Duration
}

type DelayQueue struct {
	queue *PriorityQueue
	reader chan Delayed
	writer chan Delayed
}

func NewDelayQueue() *DelayQueue {
    queue := DelayQueue{
        queue: &PriorityQueue{},
        reader: make(chan Delayed, 0),
        writer: make(chan Delayed, 0),
    }
    go queue.start()
	return &queue
}

func (d *DelayQueue) Release() {
    d.queue = nil
    close(d.reader)
    close(d.writer)
}

// Add delay task
func (d *DelayQueue) Offer(ele Delayed) {
    d.writer <- ele
}


// Gets a deferred task; no task blocks a timeout for a long time，if timeout return nil
// If there is a task in the queue and it has run out of time, the task is returned immediately
// Timeout units are seconds
func (d *DelayQueue) TakeWithTimeout(timeout time.Duration) Delayed {
    timer := time.NewTimer(timeout)
    select {
    case <-timer.C:
        return nil
    case newDelay := <-d.reader:
        return newDelay
    }
}

func (d *DelayQueue) start() {
    for {
        var first Delayed
        if d.queue.Len() > 0 {
            first = (*d.queue)[0].(Delayed)
        }
        if first != nil {
            delay := first.GetDelay()
            if delay <= 0 {
                delay := heap.Pop(d.queue).(Delayed)
                d.reader <- delay
            } else {
                timer := time.NewTimer(delay)
                select {
                case <-timer.C:
                case newDelay := <- d.writer:
                    if newDelay != nil {
                        heap.Push(d.queue, newDelay)
                    }
                }
            }
        } else {
            newDelay := <- d.writer
            if newDelay != nil {
                heap.Push(d.queue, newDelay)
            }
        }
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

