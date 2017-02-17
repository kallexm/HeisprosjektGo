package main

import
(
	"fmt"
	//"./ElevatorControlThread/ElevatorControlThread"
	"./NodeCommunication/NodeCommunicationThread"
	"./OrderDistributer/OrderDistributerThread"
	
)



func main() {
	OD_to_NC_Ch := make(chan []byte)
	NC_to_OD_Ch := make(chan []byte)
	OD_exit_Ch := make(chan bool)
	NC_exit_Ch := make(chan bool)
	
	fmt.Println("Starting main")
	
	go NodeCommunicationThread.NC_thr(OD_to_NC_Ch, NC_to_OD_Ch, NC_exit_Ch)
	go OrderDistributerThread.OD_thr(NC_to_OD_Ch, OD_to_NC_Ch, OD_exit_Ch)
	
	
	
	if <- NC_exit_Ch {
		fmt.Println("Network thread exited normaly")
	} else {
		fmt.Println("Notwork thread exited with error")
	}
	
	if <- OD_exit_Ch {
		fmt.Println("Order distributer thread exited normaly")
	} else {
		fmt.Println("Order distributer thread exited with error")
	}

	fmt.Println("exiting main")
}