package main 

import(
	"fmt"
	"./OrderQueue"
	"./OrderEvaluator"

)


func main(){
	OrderQueue.AddElevator(1)
	OrderQueue.AddElevator(2)
	OrderQueue.ChangeElevatorPosition(1,OrderQueue.Position{1,OrderQueue.DirNon})
	OrderQueue.ChangeElevatorPosition(2,OrderQueue.Position{2,OrderQueue.DirUp})
	elevators := OrderQueue.GetElevators()
	//elevators[OrderQueue.Id_t(2)] = elevators[OrderQueue.Id_t(2)].ChangeCurentOrder(OrderQueue.Order{Floor:4,OrderType: OrderQueue.Comand,DesignatedElevator: OrderQueue.Id_t(2), Cost: make(map[OrderQueue.Id_t]int)})
	_ = OrderQueue.AddOrder(OrderQueue.Order{Floor:4,OrderType: OrderQueue.Comand,DesignatedElevator: OrderQueue.Id_t(2), Cost: make(map[OrderQueue.Id_t]int)})
	fmt.Println("Elevators: ", elevators)
	newOrder := OrderQueue.AddOrder(OrderQueue.Order{Floor:3,OrderType: OrderQueue.Up,DesignatedElevator: 0, Cost: make(map[OrderQueue.Id_t]int)})
	//newOrder = OrderQueue.AddOrder(OrderQueue.Order{Floor:3,OrderType: OrderQueue.Down,DesignatedElevator: 0, Cost: make(map[OrderQueue.Id_t]int)})
	newOrder = OrderQueue.AddOrder(OrderQueue.Order{Floor:3,OrderType: OrderQueue.Down,DesignatedElevator: 0, Cost: make(map[OrderQueue.Id_t]int)})
	if newOrder == false{
		fmt.Println("Suckes")
	}
	for i, order := range OrderQueue.GetOrders(){
		fmt.Println("order nr: ", i, order)
	}
	//OrderEvaluator.AsigneCostToOrders(OrderQueue.GetOrders(), OrderQueue.GetElevators())

	test := OrderEvaluator.CalculateOrderAsignment(OrderQueue.GetOrders(), OrderQueue.GetElevators())
	for id, order := range test{
		fmt.Println("Elve id: ", id, "Order: ", *order)
	} 
	for i, order := range OrderQueue.GetOrders(){
		fmt.Println("order nr: ", i, order)
	}
	for id, elev := range OrderQueue.GetElevators(){
		fmt.Println("elev id: ", id, "curentorder", elev.CurentOrder)
	}
}