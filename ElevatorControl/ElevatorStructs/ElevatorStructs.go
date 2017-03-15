package ElevatorStructs
/*
||	File: ElevatorStructs.go
||
||	Authors:  
||
||	Date: 	 Spring 2017
||	Course:  TTK4145 - Real-time Programming, NTNU
||	
||	Summary of File: 
||		Contains different elevator data structures
||		used by the elevator module and the order manager module
||
*/



type ButtonType int
const(
	Down ButtonType = iota - 1
	Command
	Up
	Door
)


type ButtonPlacement struct{
	Floor 		int
	ButtonType 	ButtonType
	Value 		int
}


type State int
const(
	StateIdle 		= iota
	StateDoorOpen 
	StateUp
	StateDown
	StateMalfunction
)

type Dir int
const(
	DirDown = iota -1 
	DirNone
	DirUp
)


type Order struct{
	Floor 	 int 
	OrderDir Dir
}


type Position struct{
	Floor 	 int
	Dir 	 Dir
}


type OrderCompleteStruct struct{
	OrderComplete bool
}