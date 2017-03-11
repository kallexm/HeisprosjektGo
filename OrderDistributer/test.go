package main 

import(
	"fmt"
	"./OrderQueue"

)


func main(){
	OrderQueue.AddElevator(1)
	OrderQueue.AddElevator(2)
	fmt.Println("Elevators: ", OrderQueue.GetElevators())
	OrderQueue.RemoveElevator(1)
	fmt.Println("Elevators: ", OrderQueue.GetElevators())
	fmt.Println("disabeled elevators: ", OrderQueue.GetDisabeledElevators())
	newOrder := OrderQueue.AddOrder(OrderQueue.Order{Floor:1,OrderType: OrderQueue.Up,DesignatedElevator: 0, Cost: make(map[int]int)})
	newOrder = OrderQueue.AddOrder(OrderQueue.Order{Floor:1,OrderType: OrderQueue.Down,DesignatedElevator: 0, Cost: make(map[int]int)})
	newOrder = OrderQueue.AddOrder(OrderQueue.Order{Floor:1,OrderType: OrderQueue.Up,DesignatedElevator: 0, Cost: make(map[int]int)})
	if newOrder == false{
		fmt.Println("Suckes")
	}
	for i, order := range OrderQueue.GetOrders(){
		fmt.Println("order nr: ", i, order)
	}
}