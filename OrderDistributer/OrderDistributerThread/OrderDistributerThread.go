package OrderDistributerThread



import
(
	"../../MessageFormat"
	"../OrderQueue"
	"../OrderEvaluator"
	"../../ElevatorControl/ElevatorStructs"
	"encoding/json"
	"fmt"
	//"time"
)


type orderDistributerState_t uint8
const(
	STATE_MASTER orderDistributerState_t = iota
	STATE_SLAVE
)

type backUpStruct struct{
	elevators map[Id_t]Elev
	disabeledElevators map[Id_t]Elev
	orders []Order
	orderIdNr int
}

var orderDistributerState 		orderDistributerState_t
var prev_orderDistributerState 	orderDistributerState_t

func Thread(from_NodeComm_Ch 			<-chan 	[]byte	,
			to_NodeComm_Ch 				chan<- 	[]byte	,
			OrderDist_NodeComm_Mutex_Ch chan 	bool	,
			OrderDist_exit_Ch 			chan<- 	bool	,
			nodeID						uint8			) {

	orderDistributerState 		= STATE_SLAVE
	prev_orderDistributerState 	= STATE_SLAVE

	// Code to generate local elevator struct
	// Code to setup queue
	// Initialization code


	for {
		// ==========[ Begin STATE_MASTER ]===========
		if orderDistributerState == STATE_MASTER {

			// ------[ Entry Action ]-------
			if prev_orderDistributerState != orderDistributerState {
				fmt.Println("OD: [STATE_MASTER]")

				prev_orderDistributerState = orderDistributerState
			}

			// ------[ When in state, do ]-------
			select {
			case resciveMsg := <- from_NodeComm_Ch:
				resciveMsgHeader, data, err := MessageFormat.Decode_msg(resciveMsg)

				if resciveMsgHeader.From == MessageFormat.ELEVATOR {
					<- OrderDist_NodeComm_Mutex_Ch
				}

				if false { 					//Dummy if
					fmt.Println(data, err) 	//Dummy print
				}							//Dummy if

				switch resciveMsgHeader.MsgType {
				case MessageFormat.ORDER_FINISHED_BY_ELEVATOR:
					fmt.Println("ORDER_FINISHED_BY_ELEVATOR:", data)
					OrderQueue.OrderCompleet(int(resciveMsgHeader.FromNodeID))
					sendBackupToSlave(to_NodeComm_Ch)
					redistributeOrders(to_NodeComm_Ch)
					//

				case MessageFormat.NEW_ELEVATOR_REQUEST:
					fmt.Println("NEW_ELEVATOR_REQUEST:", data)
					decodeNewElevatorRequest(resciveMsgHeader, data)
					sendBackupToSlave(to_NodeComm_Ch)

					// Implement

				case MessageFormat.ELEVATOR_STATUS_DATA:
					fmt.Println("ELEVATOR_STATUS_DATA:", data)
					// Implement

				case MessageFormat.NODE_CONNECTED:
					fmt.Println("NODE_CONNECTED:", uint8(data[0]))	
					// Implement			
					// See if one has got an deactivated elevator struct that matches
					// the id in data (uint8/byte):
					// If yes: activate struct
					// If no:  generate a new elevator struct for that id, if struct not in
					// 		   activated elevator structs. Ignore if in activated elevator structs.

				case MessageFormat.NODE_DISCONNECTED:
					fmt.Println("NODE_DISCONNECTED:", uint8(data[0]))
					// Implement
					// See if one has got an activated elevator struct that matches
					// the id in data (unit8/byte):
					// If yes: deactivate struct
					// If no:  ignore

				case MessageFormat.CHANGE_TO_MASTER:
					//fmt.Println("CHANGE_TO_MASTER")
					// Do nothing

				case MessageFormat.CHANGE_TO_SLAVE:
					fmt.Println("CHANGE_TO_SLAVE")
					orderDistributerState = STATE_SLAVE
				}

				if resciveMsgHeader.From == MessageFormat.ELEVATOR {
					OrderDist_NodeComm_Mutex_Ch <- true
				}

			}


			// ------[ Exit Action ]-------
			if orderDistributerState != STATE_MASTER {


			}
		// ==========[ End STATE_MASTER ]===========




		// ==========[ Begin STATE_SLAVE ]===========
		}else if orderDistributerState == STATE_SLAVE {

			// ------[ Entry Action ]-------
			if prev_orderDistributerState != orderDistributerState {
				fmt.Println("OD: [STATE_SLAVE]")

				prev_orderDistributerState = orderDistributerState
			}

			// ------[ When in state, do ]-------
			select {
			case resciveMsg := <- from_NodeComm_Ch:
				resciveMsgHeader, data, err := MessageFormat.Decode_msg(resciveMsg)

				if false { 					//Dummy if
					fmt.Println(data, err) 	//Dummy print
				}							//Dummy if

				switch  resciveMsgHeader.MsgType {
				case MessageFormat.BACKUP_DATA_TRANSFER:
					fmt.Println("BACKUP_DATA_TRANSFER")
					// Implement

				case MessageFormat.CHANGE_TO_MASTER:
					fmt.Println("CHANGE_TO_MASTER")
					orderDistributerState = STATE_MASTER

				case MessageFormat.CHANGE_TO_SLAVE:
					//fmt.Println("CHANGE_TO_SLAVE")
					// Do nothing
				}

			}
			// ------[ Exit Action ]-------
			if orderDistributerState != STATE_SLAVE {
				
			}
		// ==========[ End STATE_SLAVE ]===========
		}
	}

}

func generateMsg(msgType MessageFormat.MsgType_t, toNodeId int, to MessageFormat.Address_t, inputStruct interface{}) []byte{
	msgHeader := MessageFormat.MessageHeader_t{To: to,ToNodeID: uint8(toNodeId), From: MessageFormat.MASTER,MsgType:msgType}
	data, err := json.Marshal(inputStruct)
	if err != nil{
		fmt.Println("Error in GenerateMsg")
	}
	msg, err := MessageFormat.Encode_msg(msgHeader, data)
	if err != nil{
		fmt.Println("Error in GenerateMsg Msg Encoding")
	}
	return msg
	
}

func redistributeOrders(to_NodeComm_Ch chan<- 	[]byte){
	elevators := OrderQueue.GetElevators()
	ordersToBeAsigned := OrderEvaluator.CalculateOrderAsignment(OrderQueue.GetOrders(), OrderQueue.GetElevators())
	for id,_ := range ordersToBeAsigned{
		elevators[id].ChangeCurentOrder(ordersToBeAsigned[id])
		newOrder := ElevatorStructs.Order{Floor: (*ordersToBeAsigned[id]).Floor,OrderDir: ElevatorStructs.Dir((*ordersToBeAsigned[id]).OrderType)}
		to_NodeComm_Ch <- generateMsg(MessageFormat.NEW_ORDER_TO_ELEVATOR, int(id), MessageFormat.ELEVATOR,newOrder)
	}
}

func sendBackupToSlave(to_NodeComm_Ch chan<- 	[]byte){
	backUp := backUpStruct{}
	backUp.orders := OrderQueue.GetOrders()
	backUp.elevatos := OrderQueue.GetElevators()
	backUp.disabeledElevators := OrderQueue.GetDisabeledElevators()
	backUp.orderIdNr := OrderQueue.GetORderIdNr()
	to_NodeComm_Ch <- generateMsg(MessageFormat.BACKUP_DATA_TRANSFER,0,MessageFormat.BACKUP, backUp)
}

func decodeNewElevatorRequest(msgHeader MessageFormat.MessageHeader_t, data [byte]){
	var newOrder ElevatorStructs.Order
	fi err:= json.Unmarshal(data, &newOrder); err != nil{
		fmt.Println("Error in decodeNewElevatorRequest: ", err)
	}
	OrderQueue.AddOrder(OrderQueue.Order{Floor: newOrder.Floor, OrderType: OrderQueue.OrderType_t(newOrder.OrderDir),})
	
}

func setLights(floor int, value int, orderType OrderQueue.OrderType_t, id OrderQueue.Id_t, to_NodeComm_Ch chan<- 	[]byte){
	if ( orderType == OrderQueue.Comand){
		to_NodeComm_Ch <- generateMsg(MessageFormat.SET_LIGHT, int(id), MessageFormat.ELEVATOR, ElevatorStructs.ButtonPlacement{Floor: floor, ButtonType:  })
	}
	
}


/*
func Thread(from_NodeComm_Ch 			<-chan 	[]byte	,
			to_NodeComm_Ch 				chan<- 	[]byte	,
			OrderDist_NodeComm_Mutex_Ch chan 	bool	,
			OrderDist_exit_Ch 			chan<- 	bool	) {
	
	for {
		select {
		case msg := <- from_NodeComm_Ch:
			receivedMsgHeader, data, err := MessageFormat.Decode_msg(msg)
			CheckError(err)
			fmt.Println("Message received:", string(data), receivedMsgHeader)
			
		default:
			time.Sleep(100*time.Millisecond)
		}
	}
	
	OrderDist_exit_Ch <- true
}


func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}*/