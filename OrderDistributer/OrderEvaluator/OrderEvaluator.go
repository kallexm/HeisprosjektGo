package OrderEvaluator

import
(
	"fmt"
	"../OrderQueue"
)


//Cost parameters
const LENGTH_COST = 1
const STOP_COST = 2
const DIRECTION_CHANGE_COST = 4



func calculateOrderAsignment() map[OrderQueue.Id_t]OrderQueue.Order{
	asigneCostToOrders()
	return asigneOrdersToElevators()
}

func asigneCostToOrders(){
	orders := OrderQueue.GetOrders()
	elevators := OrderQueue.GetElevators()
	for i, order := range orders{
		for j, elv := range elevators{
			//Do cost calculations
			// if internal order and change in direction == inf cost
			// if internal order and order_dir is in oposit direction == inf cost
			// if curent order inbetween elevator and order add STOP_COST
			// 
		}
	} 
}


func asigneOrdersToElevators() map[OrderQueue.Id_t]OrderQueue.Order{
	//Må Optimaliser summen av kosten på de tre orderen som blir valgt ut. 
}
