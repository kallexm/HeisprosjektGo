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
||		OrderDistributerThread.Thread(...) had to states
|| 			STATE_MASTER: 	It acts as an order distributer for all connected elevators.
||						  	1. It gets new orders over the the network from all nodes.
||						  	2. It Calculate a cost function and assigns the orders to the
||						  		 best positioned elevator.
||						  	3. It sends backup versions of the queue to the other
||							 	orderDistributerThreads which are in STATE_SLAVE.
||			
||			STATE_SLAVE:	1. It stores backup versions of the queue and other state data
||								sent by the only OrderDistributer in STATE_MASTER.
||							2. It is ready to take over the job as OrderDistributer
||								in STATE_MASTER if the node disconnects from the master node.
||								
||
*/


import
(
	"../OrderQueue"
	"../OrderEvaluator"
	"../../MessageFormat"
	"../../ElevatorControl/ElevatorStructs"
	
	"encoding/json"
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


	for {
		// ==========[ Begin STATE_MASTER ]===========
		if orderDistributerState == STATE_MASTER {

			// ------[ Entry Action ]-------
			if prev_orderDistributerState != orderDistributerState {

				prev_orderDistributerState = orderDistributerState
			}

			// ------[ When in state, do ]-------
			select {
			case resciveMsg := <- from_NodeComm_Ch:
				resciveMsgHeader, data, _ := MessageFormat.Decode_msg(resciveMsg)

				<- OrderDist_NodeComm_Mutex_Ch

				switch resciveMsgHeader.MsgType {
				case MessageFormat.ORDER_FINISHED_BY_ELEVATOR:
					order := OrderQueue.GetElevatorCurentOrder(int(resciveMsgHeader.FromNodeID))
					OrderQueue.OrderComplete(int(resciveMsgHeader.FromNodeID))
					redistributeOrders(to_NodeComm_Ch)
					sendBackupToSlave(to_NodeComm_Ch)
					setLights(
						order.Floor,
						0,
						order.OrderType,
						OrderQueue.Id_t(resciveMsgHeader.FromNodeID),
						to_NodeComm_Ch)

				case MessageFormat.NEW_ELEVATOR_REQUEST:
					newOrder := decodeNewElevatorRequest(resciveMsgHeader, data)
					OrderQueue.AddOrder(
						OrderQueue.Order{
							Floor: 				newOrder.Floor, 
							OrderType: 			OrderQueue.OrderType_t(newOrder.OrderDir),
							DesignatedElevator: OrderQueue.Id_t(resciveMsgHeader.FromNodeID),
							Cost: 				map[OrderQueue.Id_t]int{}					}		)
					redistributeOrders(to_NodeComm_Ch)
					sendBackupToSlave(to_NodeComm_Ch)
					setLights(
						newOrder.Floor,
						1,OrderQueue.OrderType_t(newOrder.OrderDir),
						OrderQueue.Id_t(resciveMsgHeader.FromNodeID),
						to_NodeComm_Ch									)
					newMsg, _ := MessageFormat.Encode_msg(
						MessageFormat.MessageHeader_t{
							To: 		MessageFormat.ELEVATOR, 
							ToNodeID: 	resciveMsgHeader.FromNodeID, 
							From:		MessageFormat.MASTER, 
							MsgType: 	MessageFormat.NEW_ELEVATOR_REQUEST_ACCEPTED}, 
						[]byte{}							)
					to_NodeComm_Ch <- newMsg

				case MessageFormat.ELEVATOR_STATUS_DATA:
					newStatus := decodeNewElevatorStatusData(resciveMsgHeader, data)
					OrderQueue.ChangeElevatorPosition(int(resciveMsgHeader.FromNodeID),newStatus)

				case MessageFormat.NODE_CONNECTED:
					OrderQueue.AddElevator(int(data[0]))
					iterateOrderListAndSetLights(to_NodeComm_Ch)

				case MessageFormat.NODE_DISCONNECTED:
					OrderQueue.RemoveElevator(int(data[0]))

				case MessageFormat.CHANGE_TO_MASTER:
					// Do nothing

				case MessageFormat.CHANGE_TO_SLAVE:
					orderDistributerState = STATE_SLAVE

				case MessageFormat.MERGE_ORDERS_REQUEST:
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
					backupData := decodeBackupDataTransfer(resciveMsgHeader, data)
					OrderQueue.BackupWrite(
						backupData.elevators 			,
						backupData.disabeledElevators 	,
						backupData.orders 				,
						backupData.orderIdNr 			)

				case MessageFormat.CHANGE_TO_MASTER:
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
	msgHeader := MessageFormat.MessageHeader_t{
		To: 		to 						,
		ToNodeID: 	uint8(toNodeId) 		,
		From: 		MessageFormat.MASTER 	,
		MsgType:	msgType 				}

	data, _ := json.Marshal(inputStruct)
	msg, _ 	:= MessageFormat.Encode_msg(msgHeader, data)
	return msg
}



func redistributeOrders(to_NodeComm_Ch chan<-[]byte){
	elevators := OrderQueue.GetElevators()
	ordersToBeAsigned := OrderEvaluator.CalculateOrderAssignment(OrderQueue.GetOrders(), OrderQueue.GetElevators())
	for id,_ := range ordersToBeAsigned {
		elevators[id] = elevators[id].ChangeCurentOrder(ordersToBeAsigned[id].OrderId)
		newOrder := ElevatorStructs.Order{
			Floor: 		(*ordersToBeAsigned[id]).Floor,
			OrderDir: 	ElevatorStructs.Dir((*ordersToBeAsigned[id]).OrderType)}
		to_NodeComm_Ch <- generateMsg(MessageFormat.NEW_ORDER_TO_ELEVATOR, int(id), MessageFormat.ELEVATOR,newOrder)
	}
}



func sendBackupToSlave(to_NodeComm_Ch chan<- []byte){
	backUp 						:= backUpStruct{}
	backUp.orders 				= OrderQueue.GetOrders()
	backUp.elevators 			= OrderQueue.GetElevators()
	backUp.disabeledElevators 	= OrderQueue.GetDisabeledElevators()
	backUp.orderIdNr 			= OrderQueue.GetOrderIdNr()
	to_NodeComm_Ch <- generateMsg(MessageFormat.BACKUP_DATA_TRANSFER,0,MessageFormat.BACKUP, backUp)
}



func sendMergeOrdersRequest(to_NodeComm_Ch chan<- []byte, nodeID uint8) {
	dataToMerge 					:= backUpStruct{}
	dataToMerge.orders 				= OrderQueue.GetOrders()
	dataToMerge.elevators 			= OrderQueue.GetElevators()
	dataToMerge.disabeledElevators 	= OrderQueue.GetDisabeledElevators()
	dataToMerge.orderIdNr 			= 0
	data, _ := json.Marshal(dataToMerge)

	msgHeader := MessageFormat.MessageHeader_t{	
		To: 			MessageFormat.MASTER 				,
		ToNodeID: 		uint8(0) 							,
		From: 			MessageFormat.BACKUP 				,
		FromNodeID: 	nodeID 								,
		MsgType: 		MessageFormat.MERGE_ORDERS_REQUEST	}
	mergeMsg, _ := MessageFormat.Encode_msg(msgHeader, data)
	
	to_NodeComm_Ch <- mergeMsg
}
	


func decodeNewElevatorRequest(msgHeader MessageFormat.MessageHeader_t, data []byte) ElevatorStructs.Order{
	var newOrder ElevatorStructs.Order
	json.Unmarshal(data, &newOrder)
	return newOrder
}



func decodeBackupDataTransfer(msgHeader MessageFormat.MessageHeader_t, data []byte) backUpStruct {
	var backupData backUpStruct
	json.Unmarshal(data, &backupData)
	return backupData
}



func decodeMergeOrdersRequest(msgHeader MessageFormat.MessageHeader_t, data []byte)	backUpStruct {
	var slaveOrders backUpStruct
	json.Unmarshal(data, &slaveOrders)
	return slaveOrders
}



func decodeNewElevatorStatusData(msgHeader MessageFormat.MessageHeader_t, data []byte) OrderQueue.Position{
	var newPosition ElevatorStructs.Position
	json.Unmarshal(data, &newPosition)
	return OrderQueue.Position{Floor: newPosition.Floor, Dir: OrderQueue.Dir_t(newPosition.Dir)}
}



func setLights(floor int, value int, orderType OrderQueue.OrderType_t, id OrderQueue.Id_t, to_NodeComm_Ch chan<- 	[]byte){
	if orderType == OrderQueue.Command {
		setLightStruct := ElevatorStructs.ButtonPlacement{
			Floor:		 	floor 					,
			ButtonType:	 	ElevatorStructs.Command	,
			Value: 			value 					}

		setLightMsg := generateMsg(
			MessageFormat.SET_LIGHT	,
			int(id)					, 
			MessageFormat.ELEVATOR  ,
			setLightStruct			)
		to_NodeComm_Ch <- setLightMsg
		return
	}
	elevators := OrderQueue.GetElevators()
	for id, _ := range elevators{
		setLightStruct := ElevatorStructs.ButtonPlacement{
			Floor:		 	floor 									,
			ButtonType:	 	ElevatorStructs.ButtonType(orderType)	,
			Value: 			value 									}

		setLightMsg := generateMsg(
			MessageFormat.SET_LIGHT ,
			int(id) 				,
			MessageFormat.ELEVATOR 	,
			setLightStruct			)
		to_NodeComm_Ch <- setLightMsg
	}
}



func iterateOrderListAndSetLights(to_NodeComm_Ch chan<- []byte) {
	orders := OrderQueue.GetOrders()
	for _, singleOrder := range orders {
		setLights(
			singleOrder.Floor 				,
			1								, 
			singleOrder.OrderType 			, 
			singleOrder.DesignatedElevator	, 
			to_NodeComm_Ch 					)
	}
}







