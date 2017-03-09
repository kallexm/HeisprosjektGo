package main
/*
||	File: main
||
||	Author:  Andreas Hanssen Moltumyr	
||	Partner: Martin Mostad
||	Date: 	 Spring 2017
||	Course:  TTK4145 - Real-time Programming, NTNU
||	
||	Summary of File:
||		Starting point of program.
||		It creates the nessesary channels for communication between the different parts of the program.
||		Gets the nodeID of this node from the standard I/O.
||		It spawns the three goroutines (NodeConnectionManager_thread, OrderDist_thr and ElevatorControl_thread)
||		and gives them the appropriate channels so the goroutines can communicate.
||		Then it waits for the goroutines to finish.
||
*/

import
(
	"fmt"
	"./ElevatorControl/ElevatorControlThread"
	"./NodeCommunication/NodeConnectionManager"
	"./OrderDistributer/OrderDistributerThread"
	
)


func main() {
	OrderDist_to_NodeComm_Ch := make(chan []byte)
	NodeComm_to_OrderDist_Ch := make(chan []byte)
	
	ElevCtrl_to_NodeComm_Ch := make(chan []byte)
	NodeComm_to_ElevCtrl_Ch := make(chan []byte)
	
	OrderDist_exit_Ch := make(chan bool)
	ElevCtrl_exit_Ch := make(chan bool)
	NodeComm_exit_Ch := make(chan bool)
	
	nodeID := getNodeIDfromStdIO()
	
	fmt.Println("Starting main")
	
	
	go NodeConnectionManager.NodeConnectionManager_thread(OrderDist_to_NodeComm_Ch, NodeComm_to_OrderDist_Ch, 
														  ElevCtrl_to_NodeComm_Ch , NodeComm_to_ElevCtrl_Ch ,
														  NodeComm_exit_Ch        , nodeID                  )
														  
	go OrderDistributerThread.OrderDistributer_thread(NodeComm_to_OrderDist_Ch, OrderDist_to_NodeComm_Ch,
													  OrderDist_exit_Ch                                 )
	go ElevatorControlThread.ElevatorControl_thread(NodeComm_to_ElevCtrl_Ch, ElevCtrl_to_NodeComm_Ch,
	                                                ElevCtrl_exit_Ch                                )
	
	
	
	if <- NodeComm_exit_Ch {
		fmt.Println("Network thread exited normaly")
	} else {
		fmt.Println("Notwork thread exited with error")
	}
	
	if <- OrderDist_exit_Ch {
		fmt.Println("Order distributer thread exited normaly")
	} else {
		fmt.Println("Order distributer thread exited with error")
	}
	
	if <- ElevCtrl_exit_Ch {
		fmt.Println("Elevator control thread exited normaly")
	} else {
		fmt.Println("Elevator control thread exited with error")
	}

	fmt.Println("exiting main")
}


func getNodeIDfromStdIO() uint8 {
	var nodeID uint8
	for {
		fmt.Printf("%s", "Enter this node's ID (0-255): ")
		_, err := fmt.Scanln(&nodeID)
		if err == nil {
			break
		}else{
			fmt.Println(err)
		}
	}
	return nodeID
}
	
	