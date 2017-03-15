package ElevatorStatus
/*
||	File: ElevatorStatus.go
||
||	Authors:  
||
||	Date: 	 Spring 2017
||	Course:  TTK4145 - Real-time Programming, NTNU
||	
||	Summary of File: 
||		
||
*/

import
(
	//"../ElevatorDriver/"
	"../ElevatorStructs"
	
	"./timer"
)



var currentOrder 		ElevatorStructs.Order;
var currentState 		ElevatorStructs.State;
var unconfirmedOrder 	[]ElevatorStructs.Order;
var currentPosition 	ElevatorStructs.Position;
var timerDoorchannel 	chan bool;




func GetState() ElevatorStructs.State {
	return currentState
}



func GetPosition() ElevatorStructs.Position{
	return currentPosition
}



func GetCurentOrder() ElevatorStructs.Order{
	return currentOrder
}



func GetUnconfirmedOrder() (ElevatorStructs.Order, bool){
	if len(unconfirmedOrder) == 0{
		return ElevatorStructs.Order{}, false
	}
	return unconfirmedOrder[0], true
}



func NewUnconfirmedOrder(newUnconfirmedOrder ElevatorStructs.ButtonPlacement){
	tempOrder := ElevatorStructs.Order{Floor:newUnconfirmedOrder.Floor,OrderDir:ElevatorStructs.Dir(newUnconfirmedOrder.ButtonType)}
	unconfirmedOrder = append(unconfirmedOrder,tempOrder)
}



func RemoveUnconfirmedOrder(){
	unconfirmedOrder = unconfirmedOrder[1:]
}



func SetStateMalfunction(){
	currentState = ElevatorStructs.StateMalfunction
}



func NewCurrentOrder(newCurrentOrder ElevatorStructs.Order) (ElevatorStructs.Dir, error, bool){
	if currentState == ElevatorStructs.StateDoorOpen{
		currentOrder = newCurrentOrder
		return ElevatorStructs.DirNone, nil, false
	}
	if (newCurrentOrder.Floor == currentPosition.Floor && currentPosition.Dir == ElevatorStructs.DirNone){
		currentState = ElevatorStructs.StateDoorOpen

		go timer.TimerThread(timerDoorchannel,2)
		return ElevatorStructs.DirNone, nil, true
	} else if newCurrentOrder.Floor > currentPosition.Floor{
		currentOrder 		= newCurrentOrder
		currentState 		= ElevatorStructs.StateUp
		currentPosition.Dir = ElevatorStructs.DirUp
		return ElevatorStructs.DirUp, nil, false
	} else{
		currentOrder 		= newCurrentOrder
		currentState 		= ElevatorStructs.StateDown
		currentPosition.Dir = ElevatorStructs.DirDown
		return ElevatorStructs.DirDown, nil, false
	}
}



func NewFloor(floor int) (ElevatorStructs.Dir, bool){
	currentPosition.Floor = floor
	if floor == currentOrder.Floor{
		currentPosition.Dir = ElevatorStructs.DirNone
		currentOrder 		= ElevatorStructs.Order{}
		currentState 		= ElevatorStructs.StateDoorOpen

		go timer.TimerThread(timerDoorchannel,2)
		return ElevatorStructs.DirNone, true
	}
	return currentPosition.Dir, false
}



func DoorTimeOut()ElevatorStructs.Dir{
	if (currentOrder == ElevatorStructs.Order{}){
		currentState 		= ElevatorStructs.StateIdle
		currentPosition.Dir = ElevatorStructs.DirNone
		return ElevatorStructs.DirNone
	} else if currentOrder.Floor > currentPosition.Floor{
		currentState 		= ElevatorStructs.StateUp
		currentPosition.Dir = ElevatorStructs.DirUp
		return ElevatorStructs.DirUp
	} else {
		currentState 		= ElevatorStructs.StateDown
		currentPosition.Dir = ElevatorStructs.DirDown
		return ElevatorStructs.DirDown
	}
}




func InitElevatorStatus(timerDoorch chan bool) {
	timerDoorchannel 	= timerDoorch
	currentPosition 	= ElevatorStructs.Position{}
	currentOrder 		= ElevatorStructs.Order{}
	currentState 		= ElevatorStructs.StateDown
	unconfirmedOrder 	= []ElevatorStructs.Order{}
}