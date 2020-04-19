## cron lib for go

A support cron expression、 fixed-rate、fixed-delay cron lib

## example

```go
import "github.com/wangbaolong/go-cron/cron"
```

```go
    

    var schedule = cron.NewSchedule()
    schedule.Schedule("Scheduled", 300, func() {
                logger.Info("A task exec")
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

    // schedule.ShutdownNow()
    
```