package OrderEvaluator
/*
||	File: OrderEvaluator.go
||
||	Authors:  
||
||	Date: 	 Spring 2017
||	Course:  TTK4145 - Real-time Programming, NTNU
||	
||	Summary of File: 
||		Contains functions used to optimize on which
||		elevator is best suited to handle an order.
||
*/

import
(
	"../OrderQueue"
)



//Cost parameters
const LENGTH_COST 	= 1
const STOP_COST 	= 2
const TURNING_COST 	= 4
const INF 			= 100000


func CalculateOrderAssignment(orders []OrderQueue.Order, elevators map[OrderQueue.Id_t]OrderQueue.Elev) map[OrderQueue.Id_t]*OrderQueue.Order{
	assignCostToOrders(orders, elevators)
	return assignOrdersToElevators(orders,elevators)
}



func assignCostToOrders(orders []OrderQueue.Order, elevators map[OrderQueue.Id_t]OrderQueue.Elev){
	for _, order := range orders{
		for _, elev := range elevators{
			cost := 0
			if order.OrderType == OrderQueue.Command && elev.Id != order.DesignatedElevator{
				order.Cost[elev.Id] = INF
				continue
			} 
			if elev.CurentOrder != 0{
				curentOrder := OrderQueue.GetElevatorCurentOrder(int(elev.Id))
				if curentOrder.OrderType == OrderQueue.Command && requireTurning(elev,order){
					order.Cost[elev.Id] = INF
					continue
				}
			}
			if requireTurning(elev,order){
				cost += TURNING_COST
			}
			if requireStop(elev,order,orders){
				cost += STOP_COST
			}
			lengthCost := order.Floor - elev.CurentPosition.Floor
			if lengthCost < 0{
				cost += lengthCost*-1
			} else {
				cost += lengthCost
			}
			order.Cost[elev.Id] = cost
		}
	} 
}



func requireTurning(elev OrderQueue.Elev, order OrderQueue.Order) bool{
	if elev.CurentPosition.Dir == OrderQueue.DirDown && elev.CurentPosition.Floor < order.Floor{
		return true
	} else if elev.CurentPosition.Dir == OrderQueue.DirUp && elev.CurentPosition.Floor > order.Floor{
		return true
	} else if elev.CurentPosition.Dir == OrderQueue.DirDown && order.OrderType == OrderQueue.Up{
		return true
	} else if elev.CurentPosition.Dir == OrderQueue.DirUp && order.OrderType == OrderQueue.Down{
		return true
	} else {
		return false
	}

}



func requireStop(elev OrderQueue.Elev, order OrderQueue.Order, orders []OrderQueue.Order) bool{
	curentOrder := OrderQueue.GetElevatorCurentOrder(int(elev.Id))
	if (curentOrder.OrderId == 0){
		return false
	}	
	if elev.CurentPosition.Floor < order.Floor && order.Floor < curentOrder.Floor  {
		return true
	} else if elev.CurentPosition.Floor > order.Floor && order.Floor > curentOrder.Floor{
		return true
	} else {
		return false
	}
	return false
}



func assignOrdersToElevators(orders []OrderQueue.Order, elevators map[OrderQueue.Id_t]OrderQueue.Elev) map[OrderQueue.Id_t]*OrderQueue.Order{
	sortedOrderList := make(map[OrderQueue.Id_t][]*OrderQueue.Order)
	for id, _ := range elevators{
		if _, ok := sortedOrderList[id]; !ok{
			sortedOrderList[id] = []*OrderQueue.Order{}
		}
		for i, _ := range orders{
			if len(sortedOrderList[id]) == 0{
				sortedOrderList[id] = []*OrderQueue.Order{&orders[i]}
				continue
			}
			iteratedOrder := &orders[i]
			for j, _ := range sortedOrderList[id]{
				if iteratedOrder.Cost[id] < sortedOrderList[id][j].Cost[id]{
					temp := sortedOrderList[id][j]
					sortedOrderList[id][j] = iteratedOrder
					iteratedOrder = temp 
				}
			}
			sortedOrderList[id] = append(sortedOrderList[id], iteratedOrder)

		}
	}
	
	_, asignedOrders := tryOrderCombo(sortedOrderList)
	for asignedId, _ := range asignedOrders{
		elevators[asignedId] = elevators[asignedId].ChangeCurentOrder(asignedOrders[asignedId].OrderId)
		*asignedOrders[asignedId] = (*asignedOrders[asignedId]).ChangeDesignatedElevator(asignedId)
	}
	return asignedOrders
	
}



func tryOrderCombo(sortedOrderList map[OrderQueue.Id_t][]*OrderQueue.Order) (int, map[OrderQueue.Id_t]*OrderQueue.Order){
	asignedOrders := make(map[OrderQueue.Id_t]*OrderQueue.Order)
	conflictAsignedOrders := make(map[OrderQueue.Id_t]*OrderQueue.Order)
	tempAsignedOrders := make(map[OrderQueue.Id_t]*OrderQueue.Order)
	highestCost := INF
	sum := 0
	nrUnicOrders := 0
	uniceOrder := true
	tempCost := INF
	for firstId, _ := range sortedOrderList{
		for secondId, _ := range sortedOrderList{
			if firstId == secondId{
				continue
			}
			if len(sortedOrderList[secondId]) == 0 || len(sortedOrderList[firstId]) == 0{
				uniceOrder = true
				continue
			}
			if sortedOrderList[firstId][0].OrderId == sortedOrderList[secondId][0].OrderId{
				for id, _ := range sortedOrderList{
					if firstId == id && len(sortedOrderList[firstId]) == 1 && len(sortedOrderList[secondId]) != 1{
						temp := sortedOrderList[secondId]
						sortedOrderList[secondId] = sortedOrderList[secondId][1:]
						tempCost, tempAsignedOrders = tryOrderCombo(sortedOrderList)
						sortedOrderList[id] = temp
					} else if secondId == id && len(sortedOrderList[secondId]) == 1 && len(sortedOrderList[firstId]) != 1{
						temp := sortedOrderList[firstId]
						sortedOrderList[firstId] = sortedOrderList[firstId][1:]
						tempCost, tempAsignedOrders = tryOrderCombo(sortedOrderList)
						sortedOrderList[id] = temp
					} else if len(sortedOrderList[firstId]) == 1 && len(sortedOrderList[secondId]) == 1 {
						temp := sortedOrderList
						delete(sortedOrderList, secondId)
						tempCost, tempAsignedOrders = tryOrderCombo(sortedOrderList)
						sortedOrderList = temp
					} else {
						temp := sortedOrderList[id]
						sortedOrderList[id] = sortedOrderList[id][1:]
						tempCost, tempAsignedOrders = tryOrderCombo(sortedOrderList)
						sortedOrderList[id] = temp
					}
					if tempCost < highestCost{
							highestCost = tempCost
							conflictAsignedOrders = tempAsignedOrders
					}

				}
				uniceOrder = false
			} else {
				uniceOrder = true
			}
		}
		if uniceOrder == true{
			nrUnicOrders += 1
			if len(sortedOrderList[firstId]) != 0{
				sum += sortedOrderList[firstId][0].Cost[firstId]
				asignedOrders[firstId] = sortedOrderList[firstId][0]
			}
		}
	}
	if nrUnicOrders == len(sortedOrderList){	
		return sum, asignedOrders
	}else{
		return highestCost, conflictAsignedOrders
	}
}