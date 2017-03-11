package main

import
(
	"fmt"
	//"./ElevatorControl/ElevatorDriver"
	//"./ElevatorControl/ElevatorDriver/Elev"
	"./ElevatorControl/ElevatorControlThread"
	"./ElevatorControl/ElevatorStatus"
	//"./ElevatorControl/ElevatorDriver"
	//"./ElevatorControl/ElevatorDriver/simulator/client"
	//"./NodeCommunication/NodeCommunicationThread"
	//"./OrderDistributer/OrderDistributerThread"
	"./MessageFormat"
	"encoding/json"
)

func main() {
	Elev_To_main_ch := make(chan []byte)
	main_to_Elev_ch := make(chan []byte)
	mutex_Ec_Ch := make(chan bool, 1)
	fmt.Println("Starting")
	mutex_Ec_Ch <- true


	go ElevatorControlThread.ElevatorControlThred(main_to_Elev_ch, Elev_To_main_ch,mutex_Ec_Ch)
	<-Elev_To_main_ch
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


/*func main() {
	Elev_To_main_ch := make(chan []byte)
	main_to_Elev_ch := make(chan []byte)
	mutex_Ec_Ch := make(chan bool, 1)
	fmt.Println("Starting")
	order := ElevatorStatus.Order{}
	go ElevatorControlThread.ElevatorControlThred(main_to_Elev_ch, Elev_To_main_ch,mutex_Ec_Ch)
	mutex_Ec_Ch <- true
	floorVar := 1
	buttonTypeVar := 1
	for{
		select{
		case msg := <- Elev_To_main_ch: 
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
		case <- mutex_Ec_Ch:
			if (order != ElevatorStatus.Order{}){
				testMessageHeader := MessageFormat.MessageHeader_t{To: MessageFormat.ELEVATOR,ToNodeID: 1,From: MessageFormat.MASTER, FromNodeID: 2,MsgType: MessageFormat.NEW_ORDER_TO_ELEVATOR}
				testOrder := order
				data, err := json.Marshal(testOrder)
				if err != nil{
					fmt.Println("Erro in Marshal: ", err)
				}
				msg, err := MessageFormat.Encode_msg(testMessageHeader,data)
				if err != nil{
					fmt.Println("Error in Encode_msg: ", err)
				}
				main_to_Elev_ch <- msg
			}
			lightSetingHeader := MessageFormat.MessageHeader_t{To: MessageFormat.ELEVATOR,ToNodeID:1, From: MessageFormat.MASTER, FromNodeID: 2, MsgType: MessageFormat.SET_LIGHT}
			light := ElevatorDriver.ButtonPlacement{Floor: floorVar%4+1, ButtonType: Elev.ButtonType(buttonTypeVar%3+1), Value: (floorVar+buttonTypeVar)%2}
			data, err := json.Marshal(light)
			if err != nil{
				fmt.Println("Errog in Marshal: ",err)
			}
			msg, err := MessageFormat.Encode_msg(lightSetingHeader,data)
			if err != nil{
				fmt.Println("Error in Encode_msg: ", err)
			}
			main_to_Elev_ch <- msg
		default:
		mutex_Ec_Ch <- true
		}
	}
}*/


