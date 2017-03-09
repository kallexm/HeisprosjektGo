package main

import
(
	"fmt"
	//"./ElevatorControl/ElevatorDriver"
	//"./ElevatorControl/ElevatorDriver/Elev"
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


/*func main() {
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


}*/

func main() {
	Elev_To_main_ch := make(chan []byte)
	main_to_Elev_ch := make(chan []byte)
	mutex_Ec_Ch := make(chan bool, 1)
	fmt.Println("Starting")

	go ElevatorControlThread.ElevatorControlThred(main_to_Elev_ch, Elev_To_main_ch,mutex_Ec_Ch)
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
	mutex_Ec_Ch <- true
	for {
		msg := <- Elev_To_main_ch
		msgHead, data, err := MessageFormat.Decode_msg(msg)
		fmt.Println("msgHead MsgType: ", msgHead.MsgType)
		if err != nil{
			fmt.Println("Error in decoding: ", err)
		}
		if msgHead.MsgType == MessageFormat.NEW_ELEVATOR_REQUEST{
			var newOrder ElevatorStatus.Order
			if err := json.Unmarshal(data, &newOrder); err != nil{
				fmt.Println("Error in Unmarshal: ", err)
			}
			fmt.Println("New Order: ", newOrder)


		} else if msgHead.MsgType == MessageFormat.ELEVATOR_STATUS_DATA{
			var position ElevatorStatus.Position
			if err := json.Unmarshal(data, &position); err != nil{
				fmt.Println("Error in MsgType SET_LITGHT", err)
			}
			//setLightCh <- button
			fmt.Println("New position: ", position)
		}
	} 

}