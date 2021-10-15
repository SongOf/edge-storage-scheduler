package timer

import (
	"k8s.io/klog/v2"
	"time"
)

// Timer is to call a function periodically.
type Timer struct {
	Function func()
	Duration time.Duration
	Times    int
	Shutdown chan string
}

// Start start a timer.
func (t *Timer) Start() {
	ticker := time.NewTicker(t.Duration)
	if t.Times > 0 {
		for i := 0; i < t.Times; i++ {
			select {
			case <-ticker.C:
				t.Function()
			case <-t.Shutdown:
				klog.Info("timer is exit")
				return
			}
		}
	} else {
		for {
			select {
			case <-ticker.C:
				t.Function()
			case <-t.Shutdown:
				klog.Info("timer is exit")
				return
			}
		}
	}
}

func (t *Timer) Terminated() {
	t.Shutdown <- "Done"
}
