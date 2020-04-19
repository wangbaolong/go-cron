package cron

import (
    "errors"
    "github.com/panjf2000/ants/v2"
    "github.com/wangbaolong/go-cron/robfigcron"
    "math"
    "sync/atomic"
    "time"
)

const statusShowdownNow int32 = 1

// ErrPoolClosed will be returned when submitting task to a closed pool.
var ErrScheduleShutdown = errors.New("this schedule has been shutdown")

type schedule struct {
    queue *DelayQueue
    defParser robfigcron.Parser
    status int32
    pool *ants.Pool
}

// Creates and executes a one-shot action that becomes enabled
// after the given delay.
// name the task name
// delay the time from now to delay execution
// runnable the task to execute
func (s *schedule) Schedule(name string, delay uint64, runnable func()) error {
    if s.isShutdown() {
        return ErrScheduleShutdown
    }
    s.queue.Offer(NewTask(name, int64(delay), 0, runnable))
    return nil
}

// ScheduledWithSpec adds a func to the Cron to be run on the given schedule.
// The spec is parsed using the time zone of this Cron instance as the default.
// name the task name
// spec the cron spec
// runnable the task to execute
func (s *schedule) ScheduleWithSpec(name string, spec string, runnable func()) error {
    if s.isShutdown() {
        return ErrScheduleShutdown
    }
    var scheduled, err = s.defParser.Parse(spec)
    if err != nil {
        return err
    }
    s.queue.Offer(NewTaskWithSchedule(name, scheduled, runnable))
    return nil
}

// Creates and executes a periodic action that becomes enabled first
// after the given initial delay, and subsequently with the
// given delay between the termination of one execution and the
// commencement of the next.  If any execution of the task
// encounters an exception, subsequent executions are suppressed.
// Otherwise, the task will only terminate via cancellation or
// termination of the executor.
// name the task name
// initialDelay the time to delay first execution
// delay the delay between the termination of one
// runnable the task to execute
func (s *schedule) ScheduleWithFixedDelay(name string, initialDelay uint64, delay uint64, runnable func()) error {
    if s.isShutdown() {
        return ErrScheduleShutdown
    }
    s.queue.Offer(NewTask(name, int64(initialDelay), int64(-delay), runnable))
    return nil
}

// Creates and executes a periodic action that becomes enabled first
// after the given initial delay, and subsequently with the given
// period; that is executions will commence after
// initialDelay then initialDelay+period, then
// initialDelay + 2 * period, and so on.
// If any execution of the task
// encounters an exception, subsequent executions are suppressed.
// Otherwise, the task will only terminate via cancellation or
// termination of the executor.  If any execution of this task
// takes longer than its period, then subsequent executions
// may start late, but will not concurrently execute.
func (s *schedule) ScheduleAtFixedRate(name string, initialDelay uint64, period uint64, runnable func()) error {
    if s.isShutdown() {
        return ErrScheduleShutdown
    }
    s.queue.Offer(NewTask(name, int64(initialDelay), int64(period), runnable))
    return nil
}

// uses the provided logger.
func (s *schedule) Logger(log Logger){
    logger = log
}

func (s *schedule) start() {
    go func() {
        for {
            var delayed = s.queue.TakeWithTimeout(60)
            if delayed != nil && s.status != statusShowdownNow {
                s.runWithRecovery(delayed.(*Task))
            }
        }
    }()
}

func (s *schedule) runWithRecovery(task *Task) {
    var err = s.pool.Submit(func() {
        defer func() {
            if r := recover(); r != nil {
                var err = r.(error)
                logger.Error(err, err.Error(), "runWithRecovery recover")
            }
        }()

        if task.schedule != nil {
            task.executeTime = task.schedule.Next(time.Now()).Unix()
            s.queue.Offer(task)
        }
        logger.Info(task.name, " task run start")
        task.Run()
        logger.Info(task.name, " task run end")
        if task.schedule == nil && task.period != 0 {
            if task.period > 0 {
                task.executeTime += task.period
            } else if task.period < 0 {
                task.executeTime = time.Now().Unix() - task.period
            }
            s.queue.Offer(task)
        }
    })
    if err != nil {
        logger.Error(err, err.Error(), "runWithRecovery pool submit")
    }
}


func (s *schedule) ShutdownNow() {
    atomic.AddInt32(&s.status, statusShowdownNow)
    s.pool.Release()
}

func (s *schedule) isShutdown() bool {
    return s.status == statusShowdownNow
}

// ants pool size, default math.MaxInt32
func NewSchedule() *schedule {
    return NewScheduleWithSize(math.MaxInt32)
}

// size ants pool size
func NewScheduleWithSize(size int) *schedule {
    pool, _ := ants.NewPool(size)
    var s = &schedule{queue: NewDelayQueue(),
        defParser: robfigcron.NewParser(
            robfigcron.Second | robfigcron.Minute | robfigcron.Hour | robfigcron.Dom | robfigcron.Month | robfigcron.DowOptional | robfigcron.Descriptor,
        ),
        pool:pool}
    s.start()
    return s
}

