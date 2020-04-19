package cron

import (
    "sync"
    "testing"
    "time"
)

func TestScheduled(t *testing.T) {
    var waitGroup sync.WaitGroup
    waitGroup.Add(2)
    var schedule = NewSchedule()
    schedule.Schedule("Scheduled", 180, func() {
           logger.Info("A task exec")
           waitGroup.Done()
    })

    schedule.ScheduleAtFixedRate("ScheduleAtFixedRate", 5, 5, func() {
        logger.Info("B task running spend time 10s")
        time.Sleep(time.Second * 10)
    })

    schedule.ScheduleWithFixedDelay("ScheduleWithFixedDelay", 5, 5, func() {
        logger.Info("C task running spend time 10s")
        time.Sleep(time.Second * 10)
    })

    _ = schedule.ScheduleWithSpec("specTask", "*/20 * * * * *", func() {
        logger.Info("D task running spend time 2min")
        time.Sleep(time.Minute * 2)
    })

    _ = schedule.ScheduleWithSpec("specTask", "@every 10s", func() {
        logger.Info("E task running spend time 2min")
        time.Sleep(time.Minute * 2)
    })
    _ = schedule.Schedule("Shutdown", 200, func() {
        schedule.ShutdownNow()
        logger.Info("schedule ShutdownNow")
        waitGroup.Done()
    })

    waitGroup.Wait()
}
