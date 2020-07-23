package cron

import (
    "errors"
    "sync"
    "testing"
    "time"
)

func TestScheduled(t *testing.T) {
    var waitGroup sync.WaitGroup
    waitGroup.Add(2)
    var schedule = NewSchedule()
    schedule.Schedule("Scheduled", time.Second * 100, func() {
           logger.Info("A task exec 100s")
           waitGroup.Done()
    })

    schedule.ScheduleAtFixedRate("ScheduleAtFixedRate", time.Second * 5, time.Second * 5, func() {
        defer func() {
            if r := recover(); r != nil {
                logger.Info("Recover -------------")
            }
        }()
        logger.Info("B task running spend time 5  5")
        time.Sleep(time.Second * 10)
        panic(errors.New("测试用 panic "))
    })

    schedule.ScheduleWithFixedDelay("ScheduleWithFixedDelay", time.Second * 5, time.Second * 5, func() {
        logger.Info("C task running spend time 5 5 ")
        time.Sleep(time.Second * 10)
    })

    _ = schedule.ScheduleWithSpec("specTask", "*/20 * * * * *", func() {
        logger.Info("D task running spend time */20")
        time.Sleep(time.Minute * 2)
    })

    _ = schedule.ScheduleWithSpec("specTask", "@every 10s", func() {
        logger.Info("E task running spend time 10s")
        time.Sleep(time.Minute * 2)
    })
    _ = schedule.Schedule("Shutdown", time.Second * 200, func() {
        logger.Info("schedule ShutdownNow")
        waitGroup.Done()
    })

    waitGroup.Wait()

    schedule.ShutdownNow()
}
