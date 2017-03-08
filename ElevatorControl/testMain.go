package main

import
(
	"fmt"
	//"./ElevatorStatus"
	"./ElevatorStatus/timer"
	//"time"
)


func main() {
	finishTimerCh := make(chan bool)
	fmt.Println("foo")
	go timer.TimerThredTwo(finishTimerCh, 2)
	<- finishTimerCh
	fmt.Println("bar")
}


/*func main() {
	startTimerCh := make(chan int)
	finishTimerCh := make(chan bool)
	go timer.TimerThred(startTimerCh,finishTimerCh)
	ElevatorStatus.InitElevatorStatus(startTimerCh)
	dir, orderCompleat := ElevatorStatus.NewFloor(3)
	fmt.Println("dir: ", dir, "orderCompleat: ", orderCompleat)
	testOrder := ElevatorStatus.Order{Floor:2,OrderDir:-1}
	if dir, err, orderCompleat := ElevatorStatus.NewCurentOrder(testOrder); err != nil{
		fmt.Println("Error: ", err)
	} else{
		fmt.Println("dir: ",dir, "orderCompleat", orderCompleat)
	}
	order := ElevatorStatus.GetCurentOrder()
	fmt.Println("Order floor: ", order.Floor, "Order OrderDir", order.OrderDir)
	pos := ElevatorStatus.GetPosition()
	fmt.Println("Pos floor", pos.Floor, "Pos dir", pos.Dir)
	dir, orderCompleat = ElevatorStatus.NewFloor(2)
	testOrder = ElevatorStatus.Order{Floor:3,OrderDir: 1}
	if dir, err, orderCompleat := ElevatorStatus.NewCurentOrder(testOrder); err != nil{
		fmt.Println("Error: ", err)
	} else {
		fmt.Println("dir: ",dir, "orderCompleat", orderCompleat)
	}
	<- finishTimerCh
	dir = ElevatorStatus.DoorTimeOut()
	fmt.Println("New dir: ", dir)
	pos = ElevatorStatus.GetPosition()
	fmt.Println("Pos floor", pos.Floor, "Pos dir", pos.Dir)
	status := ElevatorStatus.GetState()
	fmt.Println("stat: ", status)
	
	fmt.Println("sukses")
}*/


/*func main() {
	startTimerCh := make(chan int)
	finishTimerCh := make(chan bool)
	go timer.TimerThred(startTimerCh,finishTimerCh)
	ElevatorStatus.InitElevatorStatus(startTimerCh)
	fmt.Println("Hei")
	dir, orderCompleat := ElevatorStatus.NewFloor(1)
	fmt.Println("dir: ", dir, "orderCompleat: ", orderCompleat)
}*/


/*temp := true
	for temp{
		select{
		case finished := <- finishTimerCh:
			fmt.Println("we did it!", finished)
			temp = false 
		default:
			//fmt.Println("Not working")
		}
	}*/