## cron lib for go

A support cron expression、 fixed-rate、fixed-delay cron lib

## example

```go
    var scheduled = cron.NewScheduled()
    scheduled.Schedule("Scheduled", 300, func() {
            logger.Info("A task exec")
    })

    scheduled.ScheduleAtFixedRate("ScheduleAtFixedRate", 5, 5, func() {
        logger.Info("B task running spend time 10s")
        time.Sleep(time.Second * 10)
    })

    scheduled.ScheduleWithFixedDelay("ScheduleWithFixedDelay", 5, 5, func() {
        logger.Info("C task running spend time 10s")
        time.Sleep(time.Second * 10)
    })

     _ = schedule.ScheduleWithSpec("specTask", "*/20 * * * * *", func() {
            logger.Info("D task running spend time 2min")
            time.Sleep(time.Minute * 2)
        })
    
        _ = schedule.ScheduleWithSpec("specTask", "@every 10s", func() {
            logger.Info("D task running spend time 2min")
            time.Sleep(time.Minute * 2)
        })
```