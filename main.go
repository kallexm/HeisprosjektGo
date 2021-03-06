package main
/*
||	File: main
||
||	Authors:  
||
||	Date: 	 Spring 2017
||	Course:  TTK4145 - Real-time Programming, NTNU
||	
||	Summary of File:
||		Starting point of program.
||		It creates the nessesary channels for communication between the different parts of the program.
||		Gets the nodeID of this node from the standard I/O.
||		It spawns three goroutines (NodeConnectionManager.Thread, OrderDistributerThread.Thread and ElevatorControlThread.Thread)
||		and gives them the appropriate channels so the goroutines can communicate.
||		Then it waits for the goroutines to finish.
||
*/

import
(
	"./ElevatorControl/ElevatorControlThread"
	"./NodeCommunication/NodeConnectionManager"
	"./OrderDistributer/OrderDistributerThread"

	
	"fmt"
)


func main() {
	OrderDist_to_NodeComm_Ch 	:= make(chan []byte)
	NodeComm_to_OrderDist_Ch 	:= make(chan []byte)
	OrderDist_NodeComm_Mutex_Ch := make(chan bool)
	
	ElevCtrl_to_NodeComm_Ch 	:= make(chan []byte)
	NodeComm_to_ElevCtrl_Ch 	:= make(chan []byte)
	ElevCtrl_NodeComm_Mutex_Ch	:= make(chan bool)
	
	OrderDist_exit_Ch	:= make(chan bool)
	ElevCtrl_exit_Ch	:= make(chan bool)
	NodeComm_exit_Ch	:= make(chan bool)
	
	fmt.Println("Starting Elevator...")
	nodeID := getNodeIDfromStdIO()
	
	
	go NodeConnectionManager.Thread(	
		OrderDist_to_NodeComm_Ch	,
		NodeComm_to_OrderDist_Ch	,
		OrderDist_NodeComm_Mutex_Ch	,
		ElevCtrl_to_NodeComm_Ch		,
		NodeComm_to_ElevCtrl_Ch		,
		ElevCtrl_NodeComm_Mutex_Ch	,
		NodeComm_exit_Ch			,
		nodeID 						)
														  
	go OrderDistributerThread.Thread(	
		NodeComm_to_OrderDist_Ch	,
	 	OrderDist_to_NodeComm_Ch	,
	 	OrderDist_NodeComm_Mutex_Ch ,
	 	OrderDist_exit_Ch       	,
	 	nodeID 						)

	go ElevatorControlThread.Thread(
		NodeComm_to_ElevCtrl_Ch		,
		ElevCtrl_to_NodeComm_Ch		,
		ElevCtrl_NodeComm_Mutex_Ch	,
    	ElevCtrl_exit_Ch        	)
	
	
	if <- NodeComm_exit_Ch {
		fmt.Println("Network thread exited normaly")
	}
	if <- OrderDist_exit_Ch {
		fmt.Println("Order distributer thread exited normaly")
	}
	if <- ElevCtrl_exit_Ch {
		fmt.Println("Elevator control thread exited normaly")
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