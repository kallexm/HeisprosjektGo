package OrderQueue

import
(
	"fmt"
)

type Dir_t int
const(
	DirDown = iota -1 
	DirNon
	DirUp


)
//Byttes plass på Comand og Down, ver OBS! på feil som kan følge av dette.
type OrderType_t int
const(
	Down = iota -1 
	Comand 
	Up
)

type Id_t int

type State_t int
const(
	Idel = iota
	DoorOpen 
	StateUp
	StateDown
	Malfunction
)

type Elev struct{
	Id Id_t
	CurentOrder *Order
	CurentPosition Position
	CurentInternalOrders []*Order
	ElevatorStatus State_t

}

func (elev Elev) changePosition(position Position) Elev{
	elev.CurentPosition = position
	return elev
	
}

func (elev Elev) AddInternalOrder(order *Order) Elev{
	elev.CurentInternalOrders = append(elev.CurentInternalOrders, order)
	return elev
}

func (elev Elev) ChangeCurentOrder(order *Order) Elev{
	elev.CurentOrder = order
	return elev
}

type Order struct{
	Floor int
	OrderType OrderType_t
	DesignatedElevator Id_t
	Cost map[Id_t]int
	OrderId Id_t
}

func (order Order) ChangeDesignatedElevator(id Id_t) Order{
	order.DesignatedElevator = id
	return order
}

type Position struct{
	Floor int
	Dir Dir_t
}


var elevators map[Id_t]Elev
var disabeledElevators map[Id_t]Elev
var orders []Order
var orderIdNr int

func init(){
	elevators = make(map[Id_t]Elev)
	disabeledElevators = make(map[Id_t]Elev)
	orders = make([]Order,0)
	orderIdNr = 0
}

func AddElevator(id int){
	fmt.Println("New helevator id: ", id)
	if _, ok := disabeledElevators[Id_t(id)]; ok{
		elevators[Id_t(id)] = disabeledElevators[Id_t(id)]
		delete(disabeledElevators,Id_t(id))
		return
	}
	newElev := Elev{Id:Id_t(id),CurentOrder: &Order{},CurentPosition: Position{}, CurentInternalOrders: make([]*Order,0)}
	elevators[Id_t(id)] = newElev
}
//Comand orderes must have a designatedElevator id representing the elevator it originated from
func AddOrder(newOrder Order) bool{
	newOrder.OrderId = Id_t(orderIdNr)
	for _, order := range orders{
		if newOrder.OrderId == order.OrderId{
			return false
		}
	} 
	orders = append(orders,newOrder)
	if newOrder.OrderType == Comand{
		elevators[newOrder.DesignatedElevator] = elevators[newOrder.DesignatedElevator].AddInternalOrder(&orders[len(orders)-1])
	} else{
		orders[len(orders)-1] = orders[len(orders)-1].ChangeDesignatedElevator(Id_t(0)) 
	}
	orderIdNr += 1
	return true
}
// Må legge til funksjonalitet for at interne ordere til heisen skall tass ut av listen over ordere. 
func RemoveElevator(id int){
	newDisabelElev := elevators[Id_t(id)]
	disabeledElevators[Id_t(id)] = newDisabelElev
	delete(elevators, Id_t(id))

}

func RemoveOrder(remOrder Order){
	for i, order := range orders{
		if remOrder.Floor == order.Floor && remOrder.OrderType == order.OrderType{
			orders = append(orders[:i], orders[i+1:]...)
		}
	}
}



func GetElevators() map[Id_t]Elev{
	return elevators
}

func GetDisabeledElevators() map[Id_t]Elev{
	return disabeledElevators
}

func GetOrders() []Order{
	return orders
}

func GetOrderIdNr() int{
	return orderIdNr
}

func ChangeElevatorPosition(id int, position Position){
	elevators[Id_t(id)] = elevators[Id_t(id)].changePosition(position)
}

func OrderCompleet( id int) {
	RemoveOrder(*(elevators[Id_t(id)].CurentOrder))
	elevators[Id_t(id)] = elevators[Id_t(id)].ChangeCurentOrder(&Order{})
}

