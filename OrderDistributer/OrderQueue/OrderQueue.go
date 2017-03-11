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
	Id int
	CurentOrder Order
	CurentPosition Position
	CurentInternalOrders []Order

}

type Order struct{
	Floor int
	OrderType OrderType_t
	DesignatedElevator int
	Cost map[int]int

}

type Position struct{
	Floor int
	dir Dir_t
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
	newElev := Elev{Id:id,CurentOrder: Order{},CurentPosition: Position{}, CurentInternalOrders: make([]Order,0)}
	elevators[Id_t(id)] = newElev
}

func AddOrder(newOrder Order) bool{
	for _, order := range orders{
		if newOrder.Floor == order.Floor && newOrder.OrderType == order.OrderType{
			return false
		}
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

