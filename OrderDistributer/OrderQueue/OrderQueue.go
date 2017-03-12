package OrderQueue

import
(
	//"fmt"
)

type Dir_t int
const(
	DirDown = iota -1 
	DirNon
	DirUp


)

type OrderType_t int
const(
	Up = iota 
	Down
	Comand
)

type Id_t int

type Elev struct{
	Id Id_t
	CurentOrder Order
	CurentPosition Position
	CurentInternalOrders []Order

}

func (elev Elev) changePosition(position Position) Elev{
	elev.CurentPosition = position
	return elev
	
}

func (elev Elev) AddInternalOrder(order Order) Elev{
	elev.CurentInternalOrders = append(elev.CurentInternalOrders, order)
	return elev
}

func (elev Elev) ChangeCurentOrder(order Order) Elev{
	elev.CurentOrder = order
	return elev
}

type Order struct{
	Floor int
	OrderType OrderType_t
	DesignatedElevator Id_t
	Cost map[Id_t]int
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

func init(){
	elevators = make(map[Id_t]Elev)
	disabeledElevators = make(map[Id_t]Elev)
	orders = make([]Order,0)
}

func AddElevator(id int){
	if _, ok := disabeledElevators[Id_t(id)]; ok{
		delete(disabeledElevators,Id_t(id))
		elevators[Id_t(id)] = disabeledElevators[Id_t(id)]
		return
	}
	newElev := Elev{Id:Id_t(id),CurentOrder: Order{},CurentPosition: Position{}, CurentInternalOrders: make([]Order,0)}
	elevators[Id_t(id)] = newElev
}
//Comand orderes must have a designatedElevator id representing the elevator it originated from
func AddOrder(newOrder Order) bool{
	for _, order := range orders{
		if newOrder.Floor == order.Floor && newOrder.OrderType == order.OrderType{
			return false
		}
	}
	if newOrder.OrderType == Comand{
		//elevators[newOrder.DesignatedElevator].CurentInternalOrders = append(elevators[newOrder.DesignatedElevator].CurentInternalOrders,newOrder)
		elevators[newOrder.DesignatedElevator] = elevators[newOrder.DesignatedElevator].AddInternalOrder(newOrder)
	}  
	orders = append(orders,newOrder)
	return true
}

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


func ChangeElevatorPosition(id int, position Position){
	elevators[Id_t(id)] = elevators[Id_t(id)].changePosition(position)
}

