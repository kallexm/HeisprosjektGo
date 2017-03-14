package ElevatorStatus

import
(
	//"fmt"
	//"../ElevatorDriver/" 
	//"errors"
	"./timer"
	"../ElevatorStructs"
)


//Ikke noe av koden er testet. Planen for module er at den skall håntere tilstanden til heisen
//Den får informasjon om når det har bli trykket på knapper, kommer til ny etasje, kommer ny orderer 
//fra master eller timer for dør åpen er ferdig. Den kan også retunere alle status situasjonene


/*type state int
const(
	idel = iota
	doorOpen 
	up
	down
	malfunction
)*/

//dirComand virker som et dårlig navn. Diskuter med andreas formatet på ordere. 
//Hvordan skall det abstraheres bort. 
//Er det mulig å gjøre dattatypen mere gjenerel så den kan brukkes med Position

/*ype Dir int
const(
	DirDown = iota -1 
	DirNon
	DirUp
)*/


/*type Position struct{
	Floor int
	Dir Dir
}*/

var currentOrder ElevatorStructs.Order;
var currentState ElevatorStructs.State;
var unconfirmedOrder []ElevatorStructs.Order;
var currentPosition ElevatorStructs.Position;
var timerDoorchannel chan bool;

func GetState() ElevatorStructs.State {
	return currentState
}

func GetPosition() ElevatorStructs.Position{
	return currentPosition
}

func GetCurentOrder() ElevatorStructs.Order{
	return currentOrder
}

func GetUnconfirmedOrder() ElevatorStructs.Order{
	return unconfirmedOrder[0]
}

//Den skall kalles når det blir trykket på en knap i heisen. Orderen skall lagres her til det blir 
//Bekreftet fra master at orderen er håntert riktig.
func NewUnconfirmedOrder(newUnconfirmedOrder ElevatorStructs.ButtonPlacement){
	tempOrder := ElevatorStructs.Order{Floor:newUnconfirmedOrder.Floor,OrderDir:ElevatorStructs.Dir(newUnconfirmedOrder.ButtonType)}
	unconfirmedOrder = append(unconfirmedOrder,tempOrder)
}

func removeUnconfirmedOrder(){
	unconfirmedOrder = unconfirmedOrder[1:]
}
//kalles om heisen sluter å fungerer
func SetStateMalfunction(){
	currentState = ElevatorStructs.Malfunction
}
//Kalles når modulen får en ny ordere fra master. Retunerer retningen heise må kjøre for å fulføre orderen
//retunerer en error dersom heisen allerede har en ordere ettersom den kun kan ha 1. 
//retunerer en boolsk variabel om orderen er øyeblikelig utført. Vi står i riktig etasje
//vill det reelt kunne skje noen gang? Skall det være tilat?
func NewCurentOrder(newCurrentOrder ElevatorStructs.Order) (ElevatorStructs.Dir, error, bool){
	//Fjer denne if statmenten, order distrutubotoren skall kunne overskrive over ordere. 
	/*if (currentOrder != Order{}){
		return -2, errors.New("Kan kun ha en ordere om gangen error 005"), false 
	}*/
	if currentState == ElevatorStructs.DoorOpen{
		currentOrder = newCurrentOrder
		return ElevatorStructs.DirNon, nil, false
	}
	if (newCurrentOrder.Floor == currentPosition.Floor && currentPosition.Dir == ElevatorStructs.DirNon){
		currentState = ElevatorStructs.DoorOpen
		//starter en timer
		go timer.TimerThreadTwo(timerDoorchannel,2)
		return ElevatorStructs.DirNon, nil, true
	} else if newCurrentOrder.Floor > currentPosition.Floor{
		currentOrder 		= newCurrentOrder
		currentState 		= ElevatorStructs.Up
		currentPosition.Dir = ElevatorStructs.DirUp
		return ElevatorStructs.DirUp, nil, false
	} else{
		currentOrder 		= newCurrentOrder
		currentState 		= ElevatorStructs.Down
		currentPosition.Dir = ElevatorStructs.DirDown
		return ElevatorStructs.DirDown, nil, false
	}
}

//Kalles når heisen kommer til en ny etasje. retunerer retningen heise skal ta
//Det skjer heisen har kommet til en etasje den har en ordere på. 
func NewFloor(floor int) (ElevatorStructs.Dir, bool){
	currentPosition.Floor = floor
	if floor == currentOrder.Floor{
		currentPosition.Dir = ElevatorStructs.DirNon
		currentOrder 		= ElevatorStructs.Order{}
		currentState 		= ElevatorStructs.DoorOpen
		//starter en timer
		go timer.TimerThreadTwo(timerDoorchannel,2)
		return ElevatorStructs.DirNon, true
	}
	return currentPosition.Dir, false
}


//kalles om timeren til door har timet ut. Returnerer den nye retningen til heisen
func DoorTimeOut()ElevatorStructs.Dir{
	if (currentOrder == ElevatorStructs.Order{}){
		currentState = ElevatorStructs.Idel
		currentPosition.Dir = ElevatorStructs.DirNon
		return ElevatorStructs.DirNon
	} else if currentOrder.Floor > currentPosition.Floor{
		currentState = ElevatorStructs.Up
		currentPosition.Dir = ElevatorStructs.DirUp
		return ElevatorStructs.DirUp
	} else {
		currentState = ElevatorStructs.Down
		currentPosition.Dir = ElevatorStructs.DirDown
		return ElevatorStructs.DirDown
	}
}


func InitElevatorStatus(timerDoorch chan bool) {
	timerDoorchannel = timerDoorch
	currentPosition = ElevatorStructs.Position{}
	currentOrder = ElevatorStructs.Order{}
	currentState = ElevatorStructs.Down
	unconfirmedOrder = []ElevatorStructs.Order{}
}



