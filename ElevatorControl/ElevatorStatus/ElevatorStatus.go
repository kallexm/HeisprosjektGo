package ElevatorStatus

import
(
	"fmt"
	"../ElevatorDriver/" 
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
	floor int 
	orderDir dir
}

type position struct{
	floor int
	dir dir
}

var curentOrder Order;
var state state;
var unconfirmedOrder []Order;
var position position;
var timerDoorchanel chan bool;
var timerNewUnconfirmedOrderchanel chan bool;

func GetState() state {
	return state
}

func GetPosition() prosition{
	return position
}

func GetCurentOrder() Order{
	return curentOrder
}

//Den skall kalles når det blir trykket på en knap i heisen. Orderen skall lagres her til det blir 
//Bekreftet fra master at orderen er håntert riktig.
func NewUnconfirmedOrder(newUnconfirmedOrder ElevatorDriverMain.ButtonPlacement){
	temOrder = {floor:newUnconfirmedOrder.floor,dir:newUnconfirmedOrder.ButtonType}
	unconfirmedOrder.append(temOrder)
	timerNewUnconfirmedOrderchanel <- true
}
//kalles om heisen sluter å fungerer
func SetStateMalfunction(){
	state = malfunction
}
//Kalles når modulen får en ny ordere fra master. Retunerer retningen heise må kjøre for å fulføre orderen
//retunerer en error dersom heisen allerede har en ordere ettersom den kun kan ha 1. 
//retunerer en boolsk variabel om orderen er øyeblikelig utført. Vi står i riktig etasje
//vill det reelt kunne skje noen gang? Skall det være tilat?
func NewCurentOrder(newCurentOrder Order) (dir, error, bool){
	if curentOrder == Order{}{
		return nil, errors.New("Kan kun ha en ordere om gangen error 005"), false 
	}
	if status == doorOpen{
		curentOrder = newCurentOrder
		return dirNon, nil, false
	}
	if newCurentOrder.floor == position.floor && position.dir == dirNon{
		state = doorOpen
		//starter en timer
		timerDoorchanel <-true
		return dirNon, nil, true
	} else if newCurentOrder.floor > position.floor{
		curentOrder = newCurentOrder
		state = up
		return dirUp, nil, false
	} else{
		curentOrder = newCurentOrder
		state = dirDown
		return dirDown, nil, false
	}
}

//Kalles når heisen kommer til en ny etasje. retunerer retningen heise skall ta
//Det skjer heisen har kommet til en etasje den har en ordere på. 
func NewFloor(floor int) dir{
	position.floor = floor
	if floor == curentOrder.floor{
		position.dir = dirNon
		curentOrder = Order{}
		state = doorOpen
		//starter en timer
		timerDoorchanel <- true
		return dirNon, true
	}
	return position.dir, false
}


//kalles om timeren til door har timet ut. Retunerer den nye retningen til heisen
func doorTimeOut()dir{
	if curentOrder == Order{}{
		state = idel
		return dirNon
	} else if curentOrder.dir == dirUp{
		state = dirUp
		position.dir = dirUp
		return dirUp
	} else if curentOrder.dir == dirDown{
		state = dirDown
		position.dir = dirDown
		return dirDown
	} else if curentOrder.floor > position.floor{
		state = dirUp
		position.dir = dirUp
		return dirUp
	} else {
		state = dirDown
		position.dir = dirDown
		return dirDown
	}
}


func initElevatorStatus(timerDoorch chan bool, timerNewUnconfirmedOrderch chan bool) {
	timerDoorchanel = timerDoorchanel
	timerNewUnconfirmedOrderchanel = timerNewUnconfirmedOrderch
	position = position{}
	curentOrder = Order{}
	state = down
	unconfirmedOrder = Order{}
}



