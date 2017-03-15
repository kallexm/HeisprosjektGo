package OrderQueue
/*
||	File: OrderQueue.go
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
	"fmt"
)

type Dir_t int
const(
	DirDown = iota -1 
	DirNone
	DirUp
)

type OrderType_t int
const(
	Down = iota -1 
	Command 
	Up
)

type Id_t int



type Position struct{
	Floor int
	Dir Dir_t
}


type Order struct{
	Floor 				int
	OrderType 			OrderType_t
	DesignatedElevator 	Id_t
	Cost 				map[Id_t]int
	OrderId 			Id_t
}

func (order Order) ChangeDesignatedElevator(id Id_t) Order{
	order.DesignatedElevator = id
	return order
}





type Elev struct{
	Id Id_t
	CurentOrder Id_t
	CurentPosition Position
	CurentInternalOrders []Id_t 
}

func (elev Elev) changePosition(position Position) Elev{
	elev.CurentPosition = position
	return elev
	
}

func (elev Elev) AddInternalOrder(order Id_t) Elev{
	elev.CurentInternalOrders = append(elev.CurentInternalOrders, order)
	return elev
}

func (elev Elev) ChangeCurentOrder(order Id_t) Elev{
	elev.CurentOrder = order
	return elev
}







var elevators 			map[Id_t]Elev
var disabeledElevators 	map[Id_t]Elev
var orders 				[]Order
var orderIdNr			int






func init(){
	elevators = make(map[Id_t]Elev)
	disabeledElevators = make(map[Id_t]Elev)
	orders = make([]Order,0)
	orderIdNr = 1
}

func AddElevator(id int){
	fmt.Println("New helevator id: ", id)
	if _, ok := disabeledElevators[Id_t(id)]; ok{
		elevators[Id_t(id)] = disabeledElevators[Id_t(id)]
		delete(disabeledElevators,Id_t(id))
		return
	}
	if _, ok := elevators[Id_t(id)]; ok {
		return
	} 
	newElev := Elev{Id:Id_t(id),CurentOrder: Id_t(0),CurentPosition: Position{}, CurentInternalOrders: make([]Id_t,0)}
	elevators[Id_t(id)] = newElev
}



func AddOrder(newOrder Order) bool{
	newOrder.OrderId = Id_t(orderIdNr)
	for _, order := range orders{
		if areOrdersEqual(newOrder,order) {
			return false
		}
	} 
	orders = append(orders,newOrder)
	if newOrder.OrderType == Command{
		elevators[newOrder.DesignatedElevator] = elevators[newOrder.DesignatedElevator].AddInternalOrder(orders[len(orders)-1].OrderId)
	} else{
		orders[len(orders)-1] = orders[len(orders)-1].ChangeDesignatedElevator(Id_t(0)) 
	}
	orderIdNr += 1
	return true
}


func areOrdersEqual(newOrder Order, order Order) bool{
	return newOrder.Floor == order.Floor && newOrder.OrderType == order.OrderType && newOrder.DesignatedElevator == order.DesignatedElevator
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



func GetOrderIdNr() int{
	return orderIdNr
}



func ChangeElevatorPosition(id int, position Position){
	elevators[Id_t(id)] = elevators[Id_t(id)].changePosition(position)
}



func OrderComplete( id int) {
	orderCompleet := GetElevatorCurentOrder(id)
	for _, order := range orders{
		if order.OrderId == orderCompleet.OrderId{
			RemoveOrder(order)
		}
	}
	elevators[Id_t(id)] = elevators[Id_t(id)].ChangeCurentOrder(Id_t(0))
}



func BackupWrite(newElevators map[Id_t]Elev, newDisabeledElevators map[Id_t]Elev, newOrders []Order, newOrderIdNr int) {
	elevators 			= newElevators
	disabeledElevators 	= newDisabeledElevators
	orders 				= newOrders
	orderIdNr 			= newOrderIdNr
}



func MergeOrderFromSlave(elevatorsFromSlave map[Id_t]Elev, disabeledElevatorsFromSlave map[Id_t]Elev, ordersFromSlave []Order) {
	for _, newOrder := range ordersFromSlave {
		_ = AddOrder(newOrder)
	}
	for _, elevFromSlave := range elevatorsFromSlave {
		for _, elevInMaster := range elevators {
			if elevFromSlave.Id != elevInMaster.Id {
				elevFromSlave.CurentInternalOrders = []Id_t{}
				elevators[elevFromSlave.Id] = elevFromSlave
			}
		}
	}
	for _, disabledElevFromSlave := range disabeledElevatorsFromSlave {
		for _, disabledElevInMaster := range disabeledElevators {
			if disabledElevFromSlave.Id != disabledElevInMaster.Id {
				disabledElevFromSlave.CurentInternalOrders = []Id_t{}
				disabeledElevators[disabledElevFromSlave.Id] = disabledElevFromSlave
			}
		}
	}
	for _, singleOrder := range orders {
		if singleOrder.OrderType == Command {
			for _, elev := range elevators {
				if singleOrder.DesignatedElevator == elev.Id {
					elevators[singleOrder.OrderId].AddInternalOrder(singleOrder.OrderId)
				}
			}
		}
	}
}



func GetElevatorCurentOrder(id int) Order{
	if int(elevators[Id_t(id)].CurentOrder) == 0 {
		return Order{}
	}
	for i, _ := range orders{
		if orders[i].OrderId == elevators[Id_t(id)].CurentOrder{
			return orders[i]
		}
	} 
	return Order{}
}