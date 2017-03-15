package OrderDistributerThread
/*
||	File: OrderDistributerThread.go
||
||	Authors: 
||
||	Date: 	 Spring 2017
||	Course:  TTK4145 - Real-time Programming, NTNU
||	
||	Summary of File:
||		OrderDistributerThread.Thread(...) receives requests from all elevators 
||
*/


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
	elevators 			map[OrderQueue.Id_t]OrderQueue.Elev
	disabeledElevators 	map[OrderQueue.Id_t]OrderQueue.Elev
	orders 				[]OrderQueue.Order
	orderIdNr 			int
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
	OrderQueue.AddElevator(int(nodeID))

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
				resciveMsgHeader, data, _ := MessageFormat.Decode_msg(resciveMsg)

				<- OrderDist_NodeComm_Mutex_Ch

				switch resciveMsgHeader.MsgType {
				case MessageFormat.ORDER_FINISHED_BY_ELEVATOR:
					fmt.Println("ORDER_FINISHED_BY_ELEVATOR:", data)
					order := OrderQueue.GetElevatorCurentOrder(int(resciveMsgHeader.FromNodeID))
					fmt.Println("order GetElevatorCurentOrder: ", order)
					OrderQueue.OrderCompleet(int(resciveMsgHeader.FromNodeID))
					redistributeOrders(to_NodeComm_Ch)
					sendBackupToSlave(to_NodeComm_Ch)
					//fmt.Println("order.Floor: ", order.Floor, "order.OrderType: ",order.OrderType )
					setLights(order.Floor,0,order.OrderType,OrderQueue.Id_t(resciveMsgHeader.FromNodeID),to_NodeComm_Ch)

					//

				case MessageFormat.NEW_ELEVATOR_REQUEST:
					//DU må redistrubuerer ordere
					fmt.Println("NEW_ELEVATOR_REQUEST:", data)
					newOrder := decodeNewElevatorRequest(resciveMsgHeader, data)
					OrderQueue.AddOrder(OrderQueue.Order{Floor: newOrder.Floor, OrderType: OrderQueue.OrderType_t(newOrder.OrderDir),DesignatedElevator: OrderQueue.Id_t(resciveMsgHeader.FromNodeID) ,Cost: map[OrderQueue.Id_t]int{}})
					redistributeOrders(to_NodeComm_Ch)
					sendBackupToSlave(to_NodeComm_Ch)
					setLights(newOrder.Floor,1,OrderQueue.OrderType_t(newOrder.OrderDir),OrderQueue.Id_t(resciveMsgHeader.FromNodeID),to_NodeComm_Ch)
					newMsg, err := MessageFormat.Encode_msg(MessageFormat.MessageHeader_t{To: MessageFormat.ELEVATOR, ToNodeID: resciveMsgHeader.FromNodeID, From:MessageFormat.MASTER, MsgType: MessageFormat.NEW_ELEVATOR_REQUEST_ACCEPTED}, []byte{})
					if err != nil{
						fmt.Println("Noe gikk galt i lagingen av ACCEPTED meldingen")
					}
					fmt.Println("ACCEPTED send")
					elevators := OrderQueue.GetElevators()
					fmt.Println("Elevators :", elevators)
					to_NodeComm_Ch <- newMsg

				case MessageFormat.ELEVATOR_STATUS_DATA:
					fmt.Println("ELEVATOR_STATUS_DATA:", data)
					newStatus := decodeNewElevatorStatusData(resciveMsgHeader, data)
					OrderQueue.ChangeElevatorPosition(int(resciveMsgHeader.FromNodeID),newStatus)
					elevators :=  OrderQueue.GetElevators()
					fmt.Println("Elevators: ", elevators)

				case MessageFormat.NODE_CONNECTED:
					fmt.Println("NODE_CONNECTED:", uint8(data[0]))
					OrderQueue.AddElevator(int(data[0]))
					iterateOrderListAndSetLights(to_NodeComm_Ch)

				case MessageFormat.NODE_DISCONNECTED:
					fmt.Println("NODE_DISCONNECTED:", uint8(data[0]))
					OrderQueue.RemoveElevator(int(data[0]))

				case MessageFormat.CHANGE_TO_MASTER:
					// Do nothing

				case MessageFormat.CHANGE_TO_SLAVE:
					fmt.Println("CHANGE_TO_SLAVE")
					orderDistributerState = STATE_SLAVE

				case MessageFormat.MERGE_ORDERS_REQUEST:
					fmt.Println("MERGE_ORDERS_REQUEST")
					slaveOrders := decodeMergeOrdersRequest(resciveMsgHeader, data)
					OrderQueue.MergeOrderFromSlave(slaveOrders.elevators, slaveOrders.disabeledElevators, slaveOrders.orders)
					redistributeOrders(to_NodeComm_Ch)
					iterateOrderListAndSetLights(to_NodeComm_Ch)
				}

				OrderDist_NodeComm_Mutex_Ch <- true

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
				<- OrderDist_NodeComm_Mutex_Ch
				sendMergeOrdersRequest(to_NodeComm_Ch, nodeID)
				OrderDist_NodeComm_Mutex_Ch <- true 

				prev_orderDistributerState = orderDistributerState
			}

			// ------[ When in state, do ]-------
			select {
			case resciveMsg := <- from_NodeComm_Ch:
				resciveMsgHeader, data, _ := MessageFormat.Decode_msg(resciveMsg)

				switch  resciveMsgHeader.MsgType {
				case MessageFormat.BACKUP_DATA_TRANSFER:
					fmt.Println("BACKUP_DATA_TRANSFER")
					backupData := decodeBackupDataTransfer(resciveMsgHeader, data)
					OrderQueue.BackupWrite(backupData.elevators, backupData.disabeledElevators, backupData.orders, backupData.orderIdNr)

				case MessageFormat.CHANGE_TO_MASTER:
					fmt.Println("CHANGE_TO_MASTER")
					orderDistributerState = STATE_MASTER

				case MessageFormat.CHANGE_TO_SLAVE:
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


func redistributeOrders(to_NodeComm_Ch chan<-[]byte){
	elevators := OrderQueue.GetElevators()
	ordersToBeAsigned := OrderEvaluator.CalculateOrderAsignment(OrderQueue.GetOrders(), OrderQueue.GetElevators())
	fmt.Println("Orders to be asigned : ", ordersToBeAsigned)
	for id,_ := range ordersToBeAsigned{
		fmt.Println("Orders to be asigned : ", ordersToBeAsigned[id])
		elevators[id] = elevators[id].ChangeCurentOrder(ordersToBeAsigned[id].OrderId)
		newOrder := ElevatorStructs.Order{Floor: (*ordersToBeAsigned[id]).Floor,OrderDir: ElevatorStructs.Dir((*ordersToBeAsigned[id]).OrderType)}
		to_NodeComm_Ch <- generateMsg(MessageFormat.NEW_ORDER_TO_ELEVATOR, int(id), MessageFormat.ELEVATOR,newOrder)
	}
}


func sendBackupToSlave(to_NodeComm_Ch chan<- 	[]byte){
	backUp := backUpStruct{}
	backUp.orders 				= OrderQueue.GetOrders()
	backUp.elevators 			= OrderQueue.GetElevators()
	backUp.disabeledElevators 	= OrderQueue.GetDisabeledElevators()
	backUp.orderIdNr 			= OrderQueue.GetOrderIdNr()
	fmt.Println("Backup to slave send")
	to_NodeComm_Ch <- generateMsg(MessageFormat.BACKUP_DATA_TRANSFER,0,MessageFormat.BACKUP, backUp)
}


func sendMergeOrdersRequest(to_NodeComm_Ch chan<- []byte, nodeID uint8) {
	dataToMerge := backUpStruct{}
	dataToMerge.orders 				= OrderQueue.GetOrders()
	dataToMerge.elevators 			= OrderQueue.GetElevators()
	dataToMerge.disabeledElevators 	= OrderQueue.GetDisabeledElevators()
	dataToMerge.orderIdNr 			= 0
	data, err := json.Marshal(dataToMerge)
	if err != nil{
		fmt.Println("Error in sendMergeOrdersRequest")
	}

	msgHeader := MessageFormat.MessageHeader_t{	
		To: 			MessageFormat.MASTER 				,
		ToNodeID: 		uint8(0) 							,
		From: 			MessageFormat.BACKUP 				,
		FromNodeID: 	nodeID 								,
		MsgType: 		MessageFormat.MERGE_ORDERS_REQUEST	}
	mergeMsg, err := MessageFormat.Encode_msg(msgHeader, data)
	if err != nil{
		fmt.Println("Error in sendMergeOrdersRequest Msg Encoding")
	}
	to_NodeComm_Ch <- mergeMsg
}
	

func decodeNewElevatorRequest(msgHeader MessageFormat.MessageHeader_t, data []byte) ElevatorStructs.Order{
	var newOrder ElevatorStructs.Order
	if err:= json.Unmarshal(data, &newOrder); err != nil {
		fmt.Println("Error in decodeNewElevatorRequest: ", err)
	}
	return newOrder
}


func decodeBackupDataTransfer(msgHeader MessageFormat.MessageHeader_t, data []byte) backUpStruct {
	var backupData backUpStruct
	if err := json.Unmarshal(data, &backupData); err != nil {
		fmt.Println("Error in decodeBackupDataTransfer: ", err)
	}
	return backupData
}


func decodeMergeOrdersRequest(msgHeader MessageFormat.MessageHeader_t, data []byte)	backUpStruct {
	var slaveOrders backUpStruct
	if err := json.Unmarshal(data, &slaveOrders); err != nil {
		fmt.Println("Error in decodeMergeOrdersRequest: ", err)
	}
	return slaveOrders
}


func setLights(floor int, value int, orderType OrderQueue.OrderType_t, id OrderQueue.Id_t, to_NodeComm_Ch chan<- 	[]byte){
	if ( orderType == OrderQueue.Comand){
		to_NodeComm_Ch <- generateMsg(MessageFormat.SET_LIGHT, int(id), MessageFormat.ELEVATOR, ElevatorStructs.ButtonPlacement{Floor: floor, ButtonType: ElevatorStructs.Comand, Value: value})
		return
	}
	elevators := OrderQueue.GetElevators()
	for id, _ := range elevators{
		to_NodeComm_Ch <- generateMsg(MessageFormat.SET_LIGHT, int(id), MessageFormat.ELEVATOR, ElevatorStructs.ButtonPlacement{Floor: floor, ButtonType: ElevatorStructs.ButtonType(orderType), Value: value})
	}
	
}


func iterateOrderListAndSetLights(to_NodeComm_Ch chan<- []byte) {
	orders := OrderQueue.GetOrders()
	for _, singleOrder := range orders {
		setLights(singleOrder.Floor, 1, singleOrder.OrderType, singleOrder.DesignatedElevator, to_NodeComm_Ch)
	}
}



func decodeNewElevatorStatusData(msgHeader MessageFormat.MessageHeader_t, data []byte) OrderQueue.Position{
	var newPosition ElevatorStructs.Position
	if err := json.Unmarshal(data, &newPosition); err != nil{
		fmt.Println("Error in decodeNewElevatorRequest: ", err)
	}
	return OrderQueue.Position{Floor: newPosition.Floor, Dir: OrderQueue.Dir_t(newPosition.Dir)}
}



