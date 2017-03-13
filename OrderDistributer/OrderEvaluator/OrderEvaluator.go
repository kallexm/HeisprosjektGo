package OrderEvaluator

import
(
	"fmt"
	"../OrderQueue"
)



//Cost parameters
const LENGTH_COST = 1
const STOP_COST = 2
const TURNING_COST = 4
const INF = 100000


func CalculateOrderAsignment(orders []OrderQueue.Order, elevators map[OrderQueue.Id_t]OrderQueue.Elev) map[OrderQueue.Id_t]*OrderQueue.Order{
	asigneCostToOrders(orders, elevators)
	return asigneOrdersToElevators(orders,elevators)
}

func asigneCostToOrders(orders []OrderQueue.Order, elevators map[OrderQueue.Id_t]OrderQueue.Elev){
	fmt.Println("Elevators: ", elevators)
	for _, order := range orders{
		for _, elev := range elevators{
			cost := 0
			if order.OrderType == OrderQueue.Comand && elev.Id != order.DesignatedElevator{
				order.Cost[elev.Id] = INF
				continue
			} 
			if (*elev.CurentOrder).OrderType == OrderQueue.Comand && requierTurning(elev,order){
				order.Cost[elev.Id] = INF
				continue
			}
			if requierTurning(elev,order){
				fmt.Println("Requiered turing")
				cost += TURNING_COST
			}
			if requierStop(elev,order){
				fmt.Println("requierStop")
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

func requierTurning(elev OrderQueue.Elev, order OrderQueue.Order) bool{
	if elev.CurentPosition.Dir == OrderQueue.DirDown && elev.CurentPosition.Floor < order.Floor{
		return true
	} else if elev.CurentPosition.Dir == OrderQueue.DirUp && elev.CurentPosition.Floor > order.Floor{
		return true
	// Lit usiker på om disse kan være med i et generelt tilfelle , må vær med i situasjoner hvor du har Comand ordere
	} else if elev.CurentPosition.Dir == OrderQueue.DirDown && order.OrderType == OrderQueue.Up{
		return true
	} else if elev.CurentPosition.Dir == OrderQueue.DirUp && order.OrderType == OrderQueue.Down{
		return true
	} else {
		return false
	}

}

func requierStop(elev OrderQueue.Elev, order OrderQueue.Order) bool{
	if elev.CurentPosition.Floor < order.Floor && order.Floor < elev.CurentOrder.Floor  {
		return true
	} else if elev.CurentOrder.Floor > order.Floor && order.Floor > elev.CurentOrder.Floor{
		return true
	} else {
		return false
	}
}

/*func asigneOrdersToElevators(orders []OrderQueue.Order, elevators map[OrderQueue.Id_t]OrderQueue.Elev) map[OrderQueue.Id_t]*OrderQueue.Order{
	asignedOrders := make(map[OrderQueue.Id_t]*OrderQueue.Order)
	for i, _:= range orders{
		sugestedSignement :=  make(map[OrderQueue.Id_t]*OrderQueue.Order)
		for id, _ := range elevators{
			if _, ok := asignedOrders[id]; !ok{
				asignedOrders[id] = &OrderQueue.Order{Cost: map[OrderQueue.Id_t]int{id: INF}}
			}
			fmt.Println("elevatord id: ", id, "order ", orders[i], "asignedOrder ", asignedOrders[id], "curentOrder ", elevators[id].CurentOrder)
			if orders[i].Cost[id] < (*asignedOrders[id]).Cost[id] && (len(elevators[id].CurentOrder.Cost) == 0 || orders[i].Cost[id] < elevators[id].CurentOrder.Cost[id]){
				sugestedSignement[id] = &orders[i]
				fmt.Println("elev id: ", id, "sugested asignemtn: ", sugestedSignement[id])
			}
			//fmt.Println("Elev id: ", id, "sugested order: ", sugestedSignement)
		}
		if len(sugestedSignement) > 1{
			largestCost := [2]int{0, 0}
			for asignedId, _ := range asignedOrders{
				if asignedOrders[asignedId].Cost[asignedId] > largestCost[1]{
					largestCost[0] = int(asignedId)
					largestCost[1] = (*asignedOrders[asignedId]).Cost[asignedId]
				} 
			}
			asignedOrders[OrderQueue.Id_t(largestCost[0])] = sugestedSignement[OrderQueue.Id_t(largestCost[0])] 
		} else{
			for sugestedId, _ := range sugestedSignement{
				asignedOrders[sugestedId] = sugestedSignement[sugestedId]
			}
		}
	}
	for asignedId, _ := range asignedOrders{
		elevators[asignedId] = elevators[asignedId].ChangeCurentOrder(*asignedOrders[asignedId])
		*asignedOrders[asignedId] = (*asignedOrders[asignedId]).ChangeDesignatedElevator(asignedId)
	}
	return asignedOrders
}*/


func asigneOrdersToElevators(orders []OrderQueue.Order, elevators map[OrderQueue.Id_t]OrderQueue.Elev) map[OrderQueue.Id_t]*OrderQueue.Order{
	sortedOrderList := make(map[OrderQueue.Id_t][]*OrderQueue.Order)
	for id, _ := range elevators{
		if _, ok := sortedOrderList[id]; !ok{
			sortedOrderList[id] = []*OrderQueue.Order{}
			fmt.Println("adde new element")
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
	/*for id, orders := range sortedOrderList{
		fmt.Println("Elevator id: ", id)
		for _, order := range orders{
			fmt.Println("Orders in sortedOrderList", order)
		}
	}*/
	_, asignedOrders := tryCombo(sortedOrderList)
	for asignedId, _ := range asignedOrders{
		elevators[asignedId] = elevators[asignedId].ChangeCurentOrder(asignedOrders[asignedId])
		*asignedOrders[asignedId] = (*asignedOrders[asignedId]).ChangeDesignatedElevator(asignedId)
	}
	return asignedOrders
	
}


func tryCombo(sortedOrderList map[OrderQueue.Id_t][]*OrderQueue.Order) (int, map[OrderQueue.Id_t]*OrderQueue.Order){
	/*for id, orders := range sortedOrderList{
		fmt.Println("Elevator id: ", id)
		for _, order := range orders{
			fmt.Println("Orders in sortedOrderList", order)
		}
	}*/
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
						tempCost, tempAsignedOrders = tryCombo(sortedOrderList)
						sortedOrderList[id] = temp
					} else if secondId == id && len(sortedOrderList[secondId]) == 1 && len(sortedOrderList[firstId]) != 1{
						temp := sortedOrderList[firstId]
						sortedOrderList[firstId] = sortedOrderList[firstId][1:]
						tempCost, tempAsignedOrders = tryCombo(sortedOrderList)
						sortedOrderList[id] = temp
					} else if len(sortedOrderList[firstId]) == 1 && len(sortedOrderList[secondId]) == 1 {
						temp := sortedOrderList
						delete(sortedOrderList, secondId)
						tempCost, tempAsignedOrders = tryCombo(sortedOrderList)
						sortedOrderList = temp
					} else {
						temp := sortedOrderList[id]
						sortedOrderList[id] = sortedOrderList[id][1:]
						tempCost, tempAsignedOrders = tryCombo(sortedOrderList)
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
		/*for _, order := range asignedOrders{
			fmt.Println("asignedOrder: ", order)
		}*/
		return sum, asignedOrders
	}
	/*for _, order := range conflictAsignedOrders{
		fmt.Println("conflictAsignedOrder: ", order)
	}*/
	return highestCost, conflictAsignedOrders
}