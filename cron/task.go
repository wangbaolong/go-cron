package cron

import (
    "github.com/wangbaolong/go-cron/robfigcron"
    "time"
)

type Task struct {
    name        string
    executeTime int64
    period      int64
    schedule    robfigcron.Schedule
    runnable    func()
}

func (p *Task) GetCompareValue() int64 {
    return p.executeTime
}

func (p *Task) GetDelay() int64 {
    return p.executeTime - time.Now().Unix()
}

func (p *Task) GetExecuteTime() int64 {
    return p.executeTime
}

func (p *Task) Run() {
    (p.runnable)()
}

// delayTime units seconds
func NewTask(name string, delayTime int64, period int64, runnable func()) *Task {
    var executeTime = time.Now().Unix() + delayTime
    return &Task{name: name, executeTime: executeTime, period: period, runnable: runnable}
}

// delayTime units seconds
func NewTaskWithSchedule(name string, schedule robfigcron.Schedule, runnable func()) *Task {
    var executeTime = schedule.Next(time.Now())
    return &Task{name: name, executeTime:executeTime.Unix(), schedule:schedule, runnable:runnable}
}

