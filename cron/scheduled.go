package cron

import (
    "github.com/wangbaolong/go-cron/robfigcron"
    "time"
)

type scheduled struct {
    queue *DelayQueue
    defParser robfigcron.Parser
}

// Creates and executes a one-shot action that becomes enabled
// after the given delay.
// name the task name
// delay the time from now to delay execution
// runnable the task to execute
func (s *scheduled) Scheduled(name string, delay uint64, runnable func()) {
    s.queue.Offer(NewTask(name, int64(delay), 0, runnable))
}

// ScheduledWithSpec adds a func to the Cron to be run on the given schedule.
// The spec is parsed using the time zone of this Cron instance as the default.
// name the task name
// spec the cron spec
// runnable the task to execute
func (s *scheduled) ScheduledWithSpec(name string, spec string, runnable func()) error {
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
func (s *scheduled) ScheduleWithFixedDelay(name string, initialDelay uint64, delay uint64, runnable func()) {
    s.queue.Offer(NewTask(name, int64(initialDelay), int64(-delay), runnable))
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
func (s *scheduled) ScheduleAtFixedRate(name string, initialDelay uint64, period uint64, runnable func()) {
    s.queue.Offer(NewTask(name, int64(initialDelay), int64(period), runnable))
}

// uses the provided logger.
func (s *scheduled) Logger(log Logger){
    logger = log
}

func (s *scheduled) start() {
    go func() {
        for {
            var delayed = s.queue.TakeWithTimeout(60)
            if delayed != nil {
                go s.runWithRecovery(delayed.(*Task))
            }
        }
    }()
}

func (s *scheduled) runWithRecovery(task *Task) {
    defer func() {
        if r := recover(); r != nil {
            // TODO
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
}

func (s *scheduled) shutdown() {
    // TODO
}

func (s *scheduled) shutdownNow() {
    // TODO
}

func NewScheduled() *scheduled {
    var s = &scheduled{queue: NewDelayQueue(),
        defParser: robfigcron.NewParser(
            robfigcron.Second | robfigcron.Minute | robfigcron.Hour | robfigcron.Dom | robfigcron.Month | robfigcron.DowOptional | robfigcron.Descriptor,
            )}
    s.start()
    return s
}
