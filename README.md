# captron
[![Go Report Card](https://goreportcard.com/badge/github.com/pravj/captron)](https://goreportcard.com/report/github.com/pravj/captron)
[![GoDoc](https://godoc.org/github.com/pravj/captron?status.svg)](https://godoc.org/github.com/pravj/captron)

An 'apt' cron job service in golang.

###### Improved version of [roylee0704/gron](https://github.com/roylee0704/gron).

## Goals

- Minimalist APIs for scheduling jobs.
- Thread safety.
- Customizable Job Type.
- Customizable Schedule.

## Installation

```sh
$ go get github.com/pravj/captron
```

## Usage
Create `schedule.go`

```go
package main

import (
	"fmt"
	"time"
	"github.com/pravj/captron"
)

func main() {
	c := captron.New()
	c.AddFunc(captron.Every(1*time.Hour), func() {
		fmt.Println("runs every hour.")
	})
	c.Start()
}
```

#### Schedule Parameters

All scheduling is done in the machine's local time zone (as provided by the Go [time package](http://www.golang.org/pkg/time)).


Setup basic periodic schedule with `captron.Every()`.

```go
captron.Every(1*time.Second)
captron.Every(1*time.Minute)
captron.Every(1*time.Hour)
```

Also support `Day`, `Week` by importing `captron/xtime`:
```go
import "github.com/pravj/captron/xtime"

captron.Every(1 * xtime.Day)
captron.Every(1 * xtime.Week)
```

Schedule to run at specific time with `.At(hh:mm)`
```go
captron.Every(30 * xtime.Day).At("00:00")
captron.Every(1 * xtime.Week).At("23:59")
```

#### Custom Job Type
You may define custom job types by implementing `captron.Job` interface: `Run()`.

For example:

```go
type Reminder struct {
	Msg string
}

func (r Reminder) Run() {
  fmt.Println(r.Msg)
}
```

After job has defined, instantiate it and schedule to run in captron.
```go
c := captron.New()
r := Reminder{ "Feed the baby!" }
c.Add(captron.Every(8*time.Hour), r)
c.Start()
```

#### Custom Job Func
You may register `Funcs` to be executed on a given schedule. captron will run them in their own goroutines, asynchronously.

```go
c := captron.New()
c.AddFunc(captron.Every(1*time.Second), func() {
	fmt.Println("runs every second")
})
c.Start()
```


#### Custom Schedule
Schedule is the interface that wraps the basic `Next` method: `Next(p time.Duration) time.Time`

In `captron`, the interface value `Schedule` has the following concrete types:

- **periodicSchedule**. adds time instant t to underlying period p.
- **atSchedule**. reoccurs every period p, at time components(hh:mm).

For more info, checkout `schedule.go`.

### Full Example

```go
package main

import (
	"fmt"
	"github.com/pravj/captron"
	"github.com/pravj/captron/xtime"
)

type PrintJob struct{ Msg string }

func (p PrintJob) Run() {
	fmt.Println(p.Msg)
}

func main() {

	var (
		// schedules
		daily     = captron.Every(1 * xtime.Day)
		weekly    = captron.Every(1 * xtime.Week)
		monthly   = captron.Every(30 * xtime.Day)
		yearly    = captron.Every(365 * xtime.Day)

		// contrived jobs
		purgeTask = func() { fmt.Println("purge aged records") }
		printFoo  = printJob{"Foo"}
		printBar  = printJob{"Bar"}
	)

	c := captron.New()

	c.Add(daily.At("12:30"), printFoo)
	c.AddFunc(weekly, func() { fmt.Println("Every week") })
	c.Start()

	// Jobs may also be added to a running captron
	c.Add(monthly, printBar)
	c.AddFunc(yearly, purgeTask)

	// Stop captron (running jobs are not halted).
	c.Stop()
}
```
