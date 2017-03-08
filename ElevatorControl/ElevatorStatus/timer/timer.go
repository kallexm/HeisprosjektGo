package timer

import
(
	"time"
)


func TimerThred(startTimerCh <-chan int, timerFinishedCh chan<- bool) {
	for {
		select{
		case timerTime := <-startTimerCh:
			timer := time.NewTimer(time.Second *time.Duration(timerTime))
			<- timer.C
			timerFinishedCh <- true
		}
	}
}

func TimerThredTwo(timerFinishedCh chan<- bool, time_ int){
	timer := time.NewTimer(time.Second *time.Duration(time_))
	<- timer.C
	timerFinishedCh <- true
}