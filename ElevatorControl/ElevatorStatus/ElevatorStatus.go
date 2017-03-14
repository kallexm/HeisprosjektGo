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

var curentOrder ElevatorStructs.Order;
var curentState ElevatorStructs.State;
var unconfirmedOrder []ElevatorStructs.Order;
var curentPosition ElevatorStructs.Position;
var timerDoorchanel chan bool;

func GetState() ElevatorStructs.State {
	return curentState
}

func GetPosition() ElevatorStructs.Position{
	return curentPosition
}

func GetCurentOrder() ElevatorStructs.Order{
	return curentOrder
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
	curentState = ElevatorStructs.Malfunction
}
//Kalles når modulen får en ny ordere fra master. Retunerer retningen heise må kjøre for å fulføre orderen
//retunerer en error dersom heisen allerede har en ordere ettersom den kun kan ha 1. 
//retunerer en boolsk variabel om orderen er øyeblikelig utført. Vi står i riktig etasje
//vill det reelt kunne skje noen gang? Skall det være tilat?
func NewCurentOrder(newCurentOrder ElevatorStructs.Order) (ElevatorStructs.Dir, error, bool){
	//Fjer denne if statmenten, order distrutubotoren skall kunne overskrive over ordere. 
	/*if (curentOrder != Order{}){
		return -2, errors.New("Kan kun ha en ordere om gangen error 005"), false 
	}*/
	if curentState == ElevatorStructs.DoorOpen{
		curentOrder = newCurentOrder
		return ElevatorStructs.DirNon, nil, false
	}
	if (newCurentOrder.Floor == curentPosition.Floor && curentPosition.Dir == ElevatorStructs.DirNon){
		curentState = ElevatorStructs.DoorOpen
		//starter en timer
		go timer.TimerThredTwo(timerDoorchanel,2)
		return ElevatorStructs.DirNon, nil, true
	} else if newCurentOrder.Floor > curentPosition.Floor{
		curentOrder = newCurentOrder
		curentState = ElevatorStructs.Up
		curentPosition.Dir = ElevatorStructs.DirUp
		return ElevatorStructs.DirUp, nil, false
	} else{
		curentOrder = newCurentOrder
		curentState = ElevatorStructs.Down
		curentPosition.Dir = ElevatorStructs.DirDown
		return ElevatorStructs.DirDown, nil, false
	}
}

//Kalles når heisen kommer til en ny etasje. retunerer retningen heise skall ta
//Det skjer heisen har kommet til en etasje den har en ordere på. 
func NewFloor(floor int) (ElevatorStructs.Dir, bool){
	curentPosition.Floor = floor
	if floor == curentOrder.Floor{
		curentPosition.Dir = ElevatorStructs.DirNon
		curentOrder = ElevatorStructs.Order{}
		curentState = ElevatorStructs.DoorOpen
		//starter en timer
		go timer.TimerThredTwo(timerDoorchanel,2)
		return ElevatorStructs.DirNon, true
	}
	return curentPosition.Dir, false
}


//kalles om timeren til door har timet ut. Retunerer den nye retningen til heisen
func DoorTimeOut()ElevatorStructs.Dir{
	if (curentOrder == ElevatorStructs.Order{}){
		curentState = ElevatorStructs.Idel
		curentPosition.Dir = ElevatorStructs.DirNon
		return ElevatorStructs.DirNon
	} else if curentOrder.Floor > curentPosition.Floor{
		curentState = ElevatorStructs.Up
		curentPosition.Dir = ElevatorStructs.DirUp
		return ElevatorStructs.DirUp
	} else {
		curentState = ElevatorStructs.Down
		curentPosition.Dir = ElevatorStructs.DirDown
		return ElevatorStructs.DirDown
	}
}


func InitElevatorStatus(timerDoorch chan bool) {
	timerDoorchanel = timerDoorch
	curentPosition = ElevatorStructs.Position{}
	curentOrder = ElevatorStructs.Order{}
	curentState = ElevatorStructs.Down
	unconfirmedOrder = []ElevatorStructs.Order{}
}



