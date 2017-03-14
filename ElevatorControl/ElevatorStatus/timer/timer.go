package timer
/*
||	File: ElevatorStatus/timer/timer.go
||
||	Authors: 
||
||	Date: 	 Spring 2017
||	Course:  TTK4145 - Real-time Programming, NTNU
||	
||	Summary of File: 
||		Timer function
||
*/

import
(
	"time"
)


func TimerThread(timerFinishedCh chan<- bool, time_ int){
	timer := time.NewTimer(time.Second *time.Duration(time_))
	<- timer.C
	timerFinishedCh <- true
}