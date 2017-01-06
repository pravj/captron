package main

import (
	"fmt"
	"time"

	"github.com/pravj/captron"
	"github.com/pravj/captron/xtime"
)

type printJob struct{ Msg string }

func (p printJob) Run() {
	fmt.Println(p.Msg)
}

func main() {

	var (
		daily     = captron.Every(1 * xtime.Day)
		weekly    = captron.Every(1 * xtime.Week)
		monthly   = captron.Every(30 * xtime.Day)
		yearly    = captron.Every(365 * xtime.Day)
		purgeTask = func() { fmt.Println("purge unwanted records") }
		printFoo  = printJob{"Foo"}
		printBar  = printJob{"Bar"}
	)

	c := captron.New()

	c.AddFunc(captron.Every(1*time.Hour), func() {
		fmt.Println("Every 1 hour")
	})
	c.Start()

	c.AddFunc(weekly, func() { fmt.Println("Every week") })
	c.Add(daily.At("12:30"), printFoo)
	c.Start()

	// Jobs may also be added to a running Cron
	c.Add(monthly, printBar)
	c.AddFunc(yearly, purgeTask)

	// Stop the scheduler (does not stop any jobs already running).
	defer c.Stop()
}
