package ElevatorStructs

import(
	//"../ElevatorDriver/Elev"
	"../ElevatorDriver/simulator/client"
)

//Struckt originaly from Driver
type ButtonPlacement struct{
	Floor int
	ButtonType Elev.ButtonType
	Value int 
}

//Data types from Status
type State int
const(
	Idel = iota
	DoorOpen 
	Up
	Down
	Malfunction
)

type Dir int
const(
	DirDown = iota -1 
	DirNon
	DirUp
)


//Struckt from Status
type Order struct{
	Floor int 
	OrderDir Dir
}

type Position struct{
	Floor int
	Dir Dir
}

//STruckt from ControlThred
type OrderCompletStruck struct{
	OrderComplet bool
}