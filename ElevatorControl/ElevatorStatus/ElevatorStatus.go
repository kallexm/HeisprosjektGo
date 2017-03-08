package ElevatorStatus

import
(
	//"fmt"
	//"../ElevatorDriver/" 
	"errors"
)


//Ikke noe av koden er testet. Planen for module er at den skall håntere tilstanden til heisen
//Den får informasjon om når det har bli trykket på knapper, kommer til ny etasje, kommer ny orderer 
//fra master eller timer for dør åpen er ferdig. Den kan også retunere alle status situasjonene


type state int
const(
	idel = iota
	doorOpen 
	up
	down
	malfunction
)

//dirComand virker som et dårlig navn. Diskuter med andreas formatet på ordere. 
//Hvordan skall det abstraheres bort. 
//Er det mulig å gjøre dattatypen mere gjenerel så den kan brukkes med position

type dir int
const(
	dirDown = iota -1 
	dirNon
	dirUp
)


type Order struct{
	Floor int 
	OrderDir dir
}

type position struct{
	Floor int
	Dir dir
}

var curentOrder Order;
var curentState state;
var unconfirmedOrder []Order;
var curentPosition position;
var timerDoorchanel chan int;

func GetState() state {
	return curentState
}

func GetPosition() position{
	return curentPosition
}

func GetCurentOrder() Order{
	return curentOrder
}

func GetUnconfirmedOrder() Order{
	return unconfirmedOrder[0]
}

//Den skall kalles når det blir trykket på en knap i heisen. Orderen skall lagres her til det blir 
//Bekreftet fra master at orderen er håntert riktig.
func NewUnconfirmedOrder(newUnconfirmedOrder ElevatorDriver.ButtonPlacement){
	tempOrder := Order{floor:newUnconfirmedOrder.Floor,orderDir:int(newUnconfirmedOrder.ButtonType)}
	unconfirmedOrder = append(unconfirmedOrder,tempOrder)
	timerNewUnconfirmedOrderchanel <- true
}

func removeUnconfirmedOrder(){
	unconfirmedOrder = unconfirmedOrder[1:]
}
//kalles om heisen sluter å fungerer
func SetStateMalfunction(){
	curentState = malfunction
}
//Kalles når modulen får en ny ordere fra master. Retunerer retningen heise må kjøre for å fulføre orderen
//retunerer en error dersom heisen allerede har en ordere ettersom den kun kan ha 1. 
//retunerer en boolsk variabel om orderen er øyeblikelig utført. Vi står i riktig etasje
//vill det reelt kunne skje noen gang? Skall det være tilat?
func NewCurentOrder(newCurentOrder Order) (dir, error, bool){
	if (curentOrder != Order{}){
		return -2, errors.New("Kan kun ha en ordere om gangen error 005"), false 
	}
	if curentState == doorOpen{
		curentOrder = newCurentOrder
		return dirNon, nil, false
	}
	if (newCurentOrder.Floor == curentPosition.Floor && curentPosition.Dir == dirNon){
		curentState = doorOpen
		//starter en timer
		go timer.TimerThredTwo(timerDoorchanel,2)
		return dirNon, nil, true
	} else if newCurentOrder.Floor > curentPosition.Floor{
		curentOrder = newCurentOrder
		curentState = up
		curentPosition.Dir = dirUp
		return dirUp, nil, false
	} else{
		curentOrder = newCurentOrder
		curentState = down
		curentPosition.Dir = dirDown
		return dirDown, nil, false
	}
}

//Kalles når heisen kommer til en ny etasje. retunerer retningen heise skall ta
//Det skjer heisen har kommet til en etasje den har en ordere på. 
func NewFloor(floor int) (dir, bool){
	curentPosition.Floor = floor
	if floor == curentOrder.Floor{
		curentPosition.Dir = dirNon
		curentOrder = Order{}
		curentState = doorOpen
		//starter en timer
		go timer.TimerThredTwo(timerDoorchanel,2)
		return dirNon, true
	}
	return curentPosition.Dir, false
}


//kalles om timeren til door har timet ut. Retunerer den nye retningen til heisen
func DoorTimeOut()dir{
	if (curentOrder == Order{}){
		curentState = idel
		return dirNon
	} else if curentOrder.Floor > curentPosition.Floor{
		curentState = up
		curentPosition.Dir = dirUp
		return dirUp
	} else {
		curentState = down
		curentPosition.Dir = dirDown
		return dirDown
	}
}


func InitElevatorStatus(timerDoorch chan int) {
	timerDoorchanel = timerDoorch
	curentPosition = position{}
	curentOrder = Order{}
	curentState = down
	unconfirmedOrder = []Order{}
}



