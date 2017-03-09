package main

import
(
	"fmt"
	"./ElevatorControl/ElevatorDriver"
	"./ElevatorControl/ElevatorDriver/Elev"
	"./ElevatorControl/ElevatorControlThread"
	"./ElevatorControl/ElevatorStatus"
	//"./NodeCommunication/NodeCommunicationThread"
	//"./OrderDistributer/OrderDistributerThread"
	"./MessageFormat"
	"encoding/json"
)


/*func main() {
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
}*/


func main() {
	main_to_Elev_ch := make(chan []byte)
	Elev_To_main_ch := make(chan []byte)
	fmt.Println("Starting")

	go ElevatorControlThread.ElevatorControlThread(main_to_Elev_ch,Elev_To_main_ch)
	testMessageHeader := MessageFormat.MessageHeader_t{To: MessageFormat.ELEVATOR,ToNodeID: 1,From: MessageFormat.MASTER, FromNodeID: 2,MsgType: MessageFormat.NEW_ORDER_TO_ELEVATOR}
	testOrder := ElevatorStatus.Order{Floor: 1, OrderDir: ElevatorStatus.DirUp}
	data, err := json.Marshal(testOrder)
	if err != nil{
		fmt.Println("Erro in Marshal: ", err)
	}
	testMsg, err := MessageFormat.Encode_msg(testMessageHeader,data)
	if err != nil{
		fmt.Println("Error in Envode_msg: ", err)
	}
	main_to_Elev_ch <- testMsg
	testMessageHeader = MessageFormat.MessageHeader_t{To: MessageFormat.ELEVATOR,ToNodeID: 1,From: MessageFormat.MASTER, FromNodeID: 2,MsgType: MessageFormat.SET_LIGHT}
 	testButton := ElevatorDriver.ButtonPlacement{Floor:1,ButtonType: Elev.Up, Value: 1}
 	data, err = json.Marshal(testButton)
 	if err != nil{
 		fmt.Println("Error in Encoding testButton", err)
 	}
 	testMsg, err = MessageFormat.Encode_msg(testMessageHeader,data)
 	if err != nil{
 		fmt.Println("Error in Encoding data: ", err)
 	}
 	//fmt.Println("GOing inn")
 	main_to_Elev_ch <- testMsg 


}