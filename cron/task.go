package cron

import (
    "github.com/wangbaolong/go-cron/robfigcron"
    "time"
)

type Task struct {
    name        string
    executeTime time.Time
    period      time.Duration
    schedule    robfigcron.Schedule
    runnable    func()
}

func (p *Task) GetCompareValue() int64 {
    return p.executeTime.Unix()
}

func (p *Task) GetDelay() time.Duration {
    return p.executeTime.Sub(time.Now())
}

func (p *Task) GetExecuteTime() time.Time {
    return p.executeTime
}

func (p *Task) Run() {
    (p.runnable)()
}

// delayTime units seconds
func NewTask(name string, delayTime time.Duration, period time.Duration, runnable func()) *Task {
    var executeTime = time.Now().Add(delayTime)
    return &Task{name: name, executeTime: executeTime, period: period, runnable: runnable}
}

// delayTime units seconds
func NewTaskWithSchedule(name string, schedule robfigcron.Schedule, runnable func()) *Task {
    var executeTime = schedule.Next(time.Now())
    return &Task{name: name, executeTime:executeTime, schedule:schedule, runnable:runnable}
}

