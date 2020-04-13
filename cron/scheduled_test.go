package cron

import (
    "sync"
    "testing"
    "time"
)

func TestScheduled(t *testing.T) {
    var waitGroup sync.WaitGroup
    waitGroup.Add(1)
    var scheduled = NewScheduled()
    scheduled.Scheduled("Scheduled", 300, func() {
            logger.Info("A task exec")
            waitGroup.Done()
    })

    scheduled.ScheduleAtFixedRate("ScheduleAtFixedRate", 5, 5, func() {
        logger.Info("B task running spend time 10s")
        time.Sleep(time.Second * 10)
    })

    scheduled.ScheduleWithFixedDelay("ScheduleWithFixedDelay", 5, 5, func() {
        logger.Info("C task running spend time 10s")
        time.Sleep(time.Second * 10)
    })

    _ = scheduled.ScheduledWithSpec("specTask", "@every 10s", func() {
        logger.Info("D task running spend time 2min")
        time.Sleep(time.Minute * 2)
    })
    waitGroup.Wait()
}
