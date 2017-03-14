package ElevatorStructs

import(
	//"../ElevatorDriver/Elev"
	"../ElevatorDriver/simulator/client"
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
	ButtonType 	Elev.ButtonType
	Value 		int
}

//Data types from Status
type State int
const(
	Idel 		= iota
	DoorOpen 
	Up
	Down
	Malfunction
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