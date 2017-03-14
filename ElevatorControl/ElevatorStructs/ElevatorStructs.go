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

import(
	//"../ElevatorDriver/Elev"
	//"../ElevatorDriver/simulator/client"
)


//Data types from Elev
type ButtonType int
const(
	Up ButtonType = iota
	Comand
	Down
	Door
)

//Struckt originaly from Driver
type ButtonPlacement struct{
	Floor 		int
	ButtonType 	ButtonType
	Value 		int
}

//Data types from Status
type State int
const(
	StateIdel 		= iota
	StateDoorOpen 
	StateUp
	StateDown
	StateMalfunction
)

type Dir int
const(
	DirDown 	= iota -1 
	DirNon
	DirUp
)


//Struct from Status
type Order struct{
	Floor 	 int 
	OrderDir Dir
}

type Position struct{
	Floor 	 int
	Dir 	 Dir
}

//Struct from ControlThread
type OrderCompletStruck struct{
	OrderComplet bool
}