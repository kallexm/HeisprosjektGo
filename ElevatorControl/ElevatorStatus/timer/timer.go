package timer

import
(
	"time"
)


func TimerThread(timerFinishedCh chan<- bool, time_ int){
	timer := time.NewTimer(time.Second *time.Duration(time_))
	<- timer.C
	timerFinishedCh <- true
}