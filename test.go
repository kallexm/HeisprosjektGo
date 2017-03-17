package main 

import(
	"fmt"
	"./OrderDistributer/OrderEvaluator"
	"./OrderDistributer/OrderQueue"
)


func main(){
	//elev1 := OrderQueue.Elev{Id: OrderQueue.Id_t(1), CurenOrder: OrderQueue.Id_t(0), CurentPosition: OrderQueue.Position{}, CurentInternalOrders: []OrderQueue.Id_t{}}
	//elev1 := OrderQueue.Elev{Id: OrderQueue.Id_t(2),CurenOrder: OrderQueue.Id_t(0), CurentPosition: OrderQueue.Position{}, CurentInternalOrders: []OrderQueue.Id_t{}}
	OrderQueue.AddElevator(1)
	OrderQueue.AddElevator(2)
	elevators := OrderQueue.GetElevators()
	elevators[1] = elevators[1].ChangePosition(OrderQueue.Position{Floor: 3, Dir: OrderQueue.DirNone})
	elevators[2] = elevators[2].ChangePosition(OrderQueue.Position{Floor: 4, Dir: OrderQueue.DirNone})
	fmt.Println("Elevators: ", elevators)
	order1 := OrderQueue.Order{Floor: 1, OrderType: OrderQueue.Command, DesignatedElevator: OrderQueue.Id_t(1)}
	order2 := OrderQueue.Order{Floor: 2, OrderType: OrderQueue.Command, DesignatedElevator: OrderQueue.Id_t(1)}
	OrderQueue.AddOrder(order1)
	OrderQueue.AddOrder(order2)
	orders := OrderQueue.GetOrders()
	fmt.Println("Orders: ", orders)
	asignedOrders := OrderEvaluator.CalculateOrderAssignment(orders, elevators)
	for id, _ := range asignedOrders{
		fmt.Println("Id: ", id, "order: ", *asignedOrders[id])
	}
	elevators = OrderQueue.GetElevators()
	fmt.Println("Elevators: ", elevators)
	orders = OrderQueue.GetOrders()
	fmt.Println("Orders: ", orders)
	OrderQueue.OrderComplete(1)
	orders = OrderQueue.GetOrders()
	fmt.Println("Orders: ", orders)
	asignedOrders = OrderEvaluator.CalculateOrderAssignment(orders, elevators)
	fmt.Println("Length: ", len(asignedOrders))
	for id, _ := range asignedOrders{
		fmt.Println("Id: ", id, "order: ", *asignedOrders[id])
	}


}