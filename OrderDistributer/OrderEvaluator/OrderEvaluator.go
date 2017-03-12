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
			if elev.CurentOrder.OrderType == OrderQueue.Comand && requierTurning(elev,order){
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
			//Do cost calculations
			// if internal order and change in direction == inf cost
			// if internal order and order_dir is in oposit direction == inf cost
			// if curent order inbetween elevator and order add STOP_COST
			// 
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

func asigneOrdersToElevators(orders []OrderQueue.Order, elevators map[OrderQueue.Id_t]OrderQueue.Elev) map[OrderQueue.Id_t]*OrderQueue.Order{
	asignedOrders := make(map[OrderQueue.Id_t]*OrderQueue.Order)
	for i, _:= range orders{
		sugestedSignement :=  make(map[OrderQueue.Id_t]*OrderQueue.Order)
		for id, _ := range elevators{
			if _, ok := asignedOrders[id]; !ok{
				asignedOrders[id] = &OrderQueue.Order{Cost: map[OrderQueue.Id_t]int{id: INF}}
			}
			if orders[i].Cost[id] < (*asignedOrders[id]).Cost[id] && (len(elevators[id].CurentOrder.Cost) == 0 || orders[i].Cost[id] < elevators[id].CurentOrder.Cost[id]){
				sugestedSignement[id] = &orders[i]
			}
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
}
