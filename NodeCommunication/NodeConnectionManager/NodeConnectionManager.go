package NodeConnectionManager
/*
||	File: NodeConnectionManager 
||
||	Authors:  
||
||	Date: 	 Spring 2017
||	Course:  TTK4145 - Real-time Programming, NTNU
||	
||	Summary of File:
||		It manages the connections a single elevator node has to all the other elevator nodes which is connected.
||			1. Initiates the network module
||			2. Establishes and ends network connections.
||      	3. Updates the routing table if a new node connects or disconnects.	
||
*/

import
(
	"../NodeSingleConnection"
	"../NodeMessageRelay"
	"../NodeRoutingTable"
	"../../MessageFormat"
	
	"fmt"
	"net"
	"time"
	"math/rand"
)


const BROADCAST_TO_PORT					= "60002"
const USE_LOCAL_IP						= false

const ALLOWED_CONNECTION_RETRIES		= 3
const TCP_CONNECTING_TIMEOUT_OFFSET 	= 500
const TCP_CONNECTING_TIMEOUT_VARIANCE	= 500

const MASTER_TCP_DIAL_TIMEOUT_TIME		= 200


type nodeConnectionState_t uint8
const(
	STATE_CONNECTING nodeConnectionState_t = iota
	STATE_MASTER
	STATE_SINGLE
	STATE_SLAVE
)

var nodeConnectionState 		nodeConnectionState_t
var prev_nodeConnectionState 	nodeConnectionState_t





func Thread(from_OrderDist_Ch 			<-chan 	[]byte	,
			to_OrderDist_Ch 			chan<- 	[]byte	,
			OrderDist_NodeComm_Mutex_Ch chan 	bool	,
		  	from_ElevCtrl_Ch  			<-chan 	[]byte	,
		  	to_ElevCtrl_Ch  			chan<- 	[]byte	,
		  	ElevCtrl_NodeComm_Mutex_Ch	chan 	bool	,
		  	NodeComm_exit_Ch  			chan<- 	bool 	,
		  	nodeID 						uint8           ) {
	
	// Initialization of channels and routingTable
	nodeComm_to_MsgRelay_Ch		:= make(chan []byte)
	MsgRelay_to_nodeComm_Ch 	:= make(chan []byte)
	nodeComm_MsgRelay_Mutex_Ch	:= make(chan bool)
	
	RoutingTable_Ch		:= make(chan *NodeRoutingTable.RoutingTable_t, 1)
	routingTable_ptr	:= NodeRoutingTable.Get_reference_to_routing_table()
	routingTable_ptr.Add_new_routing_entries(NodeRoutingTable.RoutingEntry_t{	NodeID:			nodeID 						,
																			 	IsOrderDist:	true 						,
																			 	Receive_Ch:		from_OrderDist_Ch			,
																			 	Send_Ch: 		to_OrderDist_Ch				,
																			 	Mutex_Ch:		OrderDist_NodeComm_Mutex_Ch	}	,

											 NodeRoutingTable.RoutingEntry_t{	NodeID:			nodeID 						,
											 								 	IsElev:			true 						,
											 								 	Receive_Ch:		from_ElevCtrl_Ch			,
											 								 	Send_Ch:		to_ElevCtrl_Ch				,
											 								 	Mutex_Ch:		ElevCtrl_NodeComm_Mutex_Ch	}	,

											 NodeRoutingTable.RoutingEntry_t{	NodeID:			nodeID 						,
											 									IsNet:			true 						,
											 									Receive_Ch:		nodeComm_to_MsgRelay_Ch		,
											 									Send_Ch: 		MsgRelay_to_nodeComm_Ch		,
											 									Mutex_Ch:		nodeComm_MsgRelay_Mutex_Ch	}	)
	RoutingTable_Ch <- routingTable_ptr
	routingTable_ptr = nil

	go NodeMessageRelay.Thread(RoutingTable_Ch)
	

	rand.Seed(int64(nodeID)*int64(time.Now().Second()))

	
	nodeConnectionState 		= STATE_CONNECTING
	prev_nodeConnectionState 	= STATE_CONNECTING
	
	var bCastListener 			*net.UDPConn
	var numberOfConnectedSlaves uint8

	var nodeDisconnectedList 	[]uint8
	var nodeConnectedList 		[]uint8

	for {
		// ==========[ Begin STATE_CONNECTING ]===========
		if nodeConnectionState == STATE_CONNECTING {

			// ------[ Entry Action ]-------
			if prev_nodeConnectionState != nodeConnectionState {
				printFromNET("Begin STATE_CONNECTING", nodeID, nodeConnectionState)
				prev_nodeConnectionState = nodeConnectionState
			}

			// ------[ When in state, do ]-------
			if getLocalIP(USE_LOCAL_IP) == "" {
				nodeConnectionState = STATE_SINGLE
			}else{
				printFromNET("Connecting to remote node...", nodeID, nodeConnectionState)
				conn, IDofNewConnectedNode, err := connect_to_other_Node(nodeID)

				if err != nil {
					printFromNET("Failed to connect to remote node...", nodeID, nodeConnectionState)
					nodeConnectionState = STATE_MASTER

				}else{
					newRoutingEntry := NodeRoutingTable.RoutingEntry_t{	NodeID: IDofNewConnectedNode	,
																		IsExtern: true					,
																		IsMaster: true					}
					handleNewTcpConnection(conn, newRoutingEntry, RoutingTable_Ch)
					err := setMasterNodeInTable(IDofNewConnectedNode, RoutingTable_Ch)
					checkError(err)
					printFromNET("Connected to remote node"+string(IDofNewConnectedNode), nodeID, nodeConnectionState)
					nodeConnectionState = STATE_SLAVE
				}
			}
			
			// ------[ Exit Action ]-------
			if nodeConnectionState != STATE_CONNECTING {
				if nodeConnectionState == STATE_MASTER || nodeConnectionState == STATE_SINGLE {
					select {
					case <- nodeComm_MsgRelay_Mutex_Ch:
						msgHeader := MessageFormat.MessageHeader_t{	To: 		MessageFormat.ORDER_DIST		,
																	ToNodeID: 	nodeID 							,
																	From:		MessageFormat.NODE_COM			,
																	FromNodeID: nodeID 							,
																	MsgType:	MessageFormat.CHANGE_TO_MASTER	}
						changeToMasterMsg, err := MessageFormat.Encode_msg(msgHeader, []byte(""))
						checkError(err)

						tryToSendTimer := time.NewTimer(500*time.Millisecond)
						select {
						case nodeComm_to_MsgRelay_Ch <- changeToMasterMsg:
						case <- tryToSendTimer.C:
						}
						nodeComm_MsgRelay_Mutex_Ch <- true
					}
				}else if nodeConnectionState == STATE_SLAVE {
					select {
					case <- nodeComm_MsgRelay_Mutex_Ch:
						msgHeader := MessageFormat.MessageHeader_t{	To: 		MessageFormat.ORDER_DIST 		,
																	ToNodeID: 	nodeID 							,
																	From:		MessageFormat.NODE_COM			,
																	FromNodeID: nodeID 							,
																	MsgType:	MessageFormat.CHANGE_TO_SLAVE	}
						changeToSlaveMsg, err := MessageFormat.Encode_msg(msgHeader, []byte(""))
						checkError(err)

						tryToSendTimer := time.NewTimer(500*time.Millisecond)
						select {
						case nodeComm_to_MsgRelay_Ch <- changeToSlaveMsg:
						case <- tryToSendTimer.C:
						}
						nodeComm_MsgRelay_Mutex_Ch <- true
					}
				}
			}
		// ==========[ End STATE_CONNECTING ]===========




		// ==========[ Begin STATE_MASTER ]===========
		}else if nodeConnectionState == STATE_MASTER {
			
			// ------[ Entry Action ]-------
			if prev_nodeConnectionState != nodeConnectionState {
				printFromNET("Begin STATE_MASTER", nodeID, nodeConnectionState)
				err := setMasterNodeInTable(nodeID, RoutingTable_Ch)

				printFromNET("Starts to listen for remote nodes...", nodeID, nodeConnectionState)
				listenerAddress, err := net.ResolveUDPAddr("udp", ":"+string(BROADCAST_TO_PORT))
				checkError(err)
				bCastListener, err = net.ListenUDP("udp", listenerAddress)
				checkError(err)
				if err != nil {
					printFromNET("UDP port"+BROADCAST_TO_PORT+"already in use...", nodeID, nodeConnectionState)
					nodeConnectionState = STATE_CONNECTING
					continue
				}
				prev_nodeConnectionState = nodeConnectionState
			}
			
			// ------[ When in state, do ]-------
			bCastListener.SetReadDeadline(time.Now().Add(1*time.Second))
			buffer := make([]byte, 32)
			_, _, err := bCastListener.ReadFromUDP(buffer)
			if err == nil {
				nodeIDofBcaster := uint8(buffer[0])
				tcpListenerAddrOfBcaster := string(buffer[1:])
				tcpConn, err := net.DialTimeout("tcp", tcpListenerAddrOfBcaster, MASTER_TCP_DIAL_TIMEOUT_TIME*time.Millisecond)
				checkError(err)
				if err == nil {
					printMsg := "Established tcp connection:"+tcpConn.LocalAddr().String()+"->"+tcpConn.RemoteAddr().String()
					printFromNET(printMsg, nodeID, nodeConnectionState)
					nodeIDmsg := make([]byte, 1)
					nodeIDmsg[0] = byte(nodeID)
					_, err := tcpConn.Write(nodeIDmsg)
					checkError(err)
					if err == nil {
						newRoutingEntry := NodeRoutingTable.RoutingEntry_t{	NodeID: nodeIDofBcaster	,
																			IsExtern: true			,
																			IsBackup: true			}
						handleNewTcpConnection(tcpConn, newRoutingEntry, RoutingTable_Ch)
						printMsg := "Connection to node with ID"+string(nodeIDofBcaster)+"was successful"
						printFromNET(printMsg, nodeID, nodeConnectionState)
						nodeConnectedList = append(nodeConnectedList, nodeIDofBcaster)
						numberOfConnectedSlaves++
					}else{
						tcpConn.Close()
					}
				}
			}


			select {
			case msg := <- MsgRelay_to_nodeComm_Ch:
				msgHeader, _, err := MessageFormat.Decode_msg(msg)
				checkError(err)
				if msgHeader.MsgType == MessageFormat.NODE_DISCONNECTED {
					routingTable_ptr = <- RoutingTable_Ch
					removedRoutingEntry, err := routingTable_ptr.Remove_routing_entry(msgHeader.FromNodeID)
					RoutingTable_Ch <- routingTable_ptr
					routingTable_ptr = nil

					checkError(err)
					if err == nil {
						replyHeader := MessageFormat.MessageHeader_t{ToNodeID: msgHeader.FromNodeID				,
																	 From: MessageFormat.NODE_COM				,
																	 MsgType: MessageFormat.NODE_DISCONNECTED	} 
						msg, _ = MessageFormat.Encode_msg(replyHeader, []byte(""))
						removedRoutingEntry.Send_Ch <- msg
						printFromNET("Disconnected node with ID"+string(removedRoutingEntry.NodeID), nodeID, nodeConnectionState)
						nodeDisconnectedList = append(nodeDisconnectedList, msgHeader.FromNodeID)
						numberOfConnectedSlaves--
					}
				}else{

				}
			case <- nodeComm_MsgRelay_Mutex_Ch:
				if len(nodeConnectedList) > 0 {
					connectedNodeID := []byte{byte(nodeConnectedList[len(nodeConnectedList)-1])}
					msgHeader := MessageFormat.MessageHeader_t{	To: 		MessageFormat.MASTER 			,
																ToNodeID: 	nodeID 							,
																From:		MessageFormat.NODE_COM			,
																FromNodeID: nodeID 							,
																MsgType:	MessageFormat.NODE_CONNECTED	}
					newNodeConnectedMsg, err := MessageFormat.Encode_msg(msgHeader, connectedNodeID)
					checkError(err)

					tryToSendTimer := time.NewTimer(500*time.Millisecond)
					select {
					case nodeComm_to_MsgRelay_Ch <- newNodeConnectedMsg:
						nodeConnectedList = nodeConnectedList[:len(nodeConnectedList)-1]
					case <- tryToSendTimer.C:
					}
				}
				if len(nodeDisconnectedList) > 0 {
					disconnectedNodeID := []byte{byte(nodeDisconnectedList[len(nodeDisconnectedList)-1])}
					msgHeader := MessageFormat.MessageHeader_t{	To: 		MessageFormat.MASTER 	,
																ToNodeID: 	nodeID 					,
																From:		MessageFormat.NODE_COM	,
																FromNodeID: nodeID 					,
																MsgType:	MessageFormat.NODE_DISCONNECTED 			}
					nodeDisconnectedMsg, err := MessageFormat.Encode_msg(msgHeader, disconnectedNodeID)
					checkError(err)
					
					tryToSendTimer := time.NewTimer(500*time.Millisecond)
					select {
					case nodeComm_to_MsgRelay_Ch <- nodeDisconnectedMsg:
						nodeDisconnectedList = nodeDisconnectedList[:len(nodeDisconnectedList)-1]
					case <- tryToSendTimer.C:
					}
				}
				nodeComm_MsgRelay_Mutex_Ch <- true
			}

			if numberOfConnectedSlaves == 0 && getLocalIP(USE_LOCAL_IP) == "" {
				nodeConnectionState = STATE_SINGLE
			}

			
			// ------[ Exit Action ]-------
			if nodeConnectionState != STATE_MASTER {
				bCastListener.Close()
			}
		// ==========[ End STATE_MASTER ]===========





		// ==========[ Begin STATE_SINGLE ]===========
		}else if nodeConnectionState == STATE_SINGLE {

			// ------[ Entry Action ]-------
			if prev_nodeConnectionState != nodeConnectionState {
				printFromNET("Begin STATE_SINGLE", nodeID, nodeConnectionState)
				err := setMasterNodeInTable(nodeID, RoutingTable_Ch)
				checkError(err)
				prev_nodeConnectionState = nodeConnectionState
			}

			// ------[ When in state, do ]-------
			if getLocalIP(USE_LOCAL_IP) != "" {
				nodeConnectionState = STATE_CONNECTING
			}

			// ------[ Exit Action ]-------
			if nodeConnectionState != STATE_SINGLE {
				
			}
		// ==========[ End STATE_SINGLE ]===========





		// ==========[ Begin STATE_SLAVE ]===========
		}else if nodeConnectionState == STATE_SLAVE {

			// ------[ Entry Action ]-------
			if prev_nodeConnectionState != nodeConnectionState {
				printFromNET("Begin STATE_SLAVE", nodeID, nodeConnectionState)
				prev_nodeConnectionState = nodeConnectionState
			}

			// ------[ When in state, do ]-------
			select {
			case msg := <- MsgRelay_to_nodeComm_Ch:
				msgHeader, _, err := MessageFormat.Decode_msg(msg)
				checkError(err)
				if msgHeader.MsgType == MessageFormat.NODE_DISCONNECTED {
					routingTable_ptr = <- RoutingTable_Ch
					removedRoutingEntry, err := routingTable_ptr.Remove_routing_entry(msgHeader.FromNodeID)
					RoutingTable_Ch <- routingTable_ptr
					routingTable_ptr = nil

					checkError(err)
					if err == nil {
						replyHeader := MessageFormat.MessageHeader_t{ToNodeID: msgHeader.FromNodeID				,
																 From: MessageFormat.NODE_COM				,
																 MsgType: MessageFormat.NODE_DISCONNECTED	} 
						msg, _ = MessageFormat.Encode_msg(replyHeader, []byte(""))
						removedRoutingEntry.Send_Ch <- msg
						printFromNET("Disconnected node with ID"+string(removedRoutingEntry.NodeID), nodeID, nodeConnectionState)
						nodeConnectionState = STATE_CONNECTING
					}
				}
			default:
			}

			// ------[ Exit Action ]-------
			if nodeConnectionState != STATE_SLAVE {
				
			}
		// ==========[ End STATE_SLAVE ]===========
		}
	}
	NodeComm_exit_Ch <- true	
}	





// Returns a valid connection if a MASTER node was found. Else it returns an error.
func connect_to_other_Node(nodeID uint8) (net.Conn, uint8, error){

	// Setting up a UDP broadcast socket
	bCastToAddrStr	 	:= net.JoinHostPort("255.255.255.255", BROADCAST_TO_PORT)
	bCastToAddr, _		:= net.ResolveUDPAddr("udp", bCastToAddrStr)
	bCastConn, err 		:= net.DialUDP("udp", nil, bCastToAddr)
	checkError(err)
	
	// Setting up a TCP listener socket
	listenAddrStr		:= net.JoinHostPort(getLocalIP(USE_LOCAL_IP), "0")
	listenAddr, _	 	:= net.ResolveTCPAddr("tcp", listenAddrStr)
	tcpListener, err 	:= net.ListenTCP("tcp", listenAddr)
	checkError(err)
	defer tcpListener.Close()
	
	// Make broadcast message
	bCastMsg := []byte{byte(nodeID)}
	bCastMsg = append(bCastMsg, []byte(tcpListener.Addr().String())...)
	
	// Generate random SetDeadline wait time
	randomNumber := time.Duration(rand.Int63n(TCP_CONNECTING_TIMEOUT_VARIANCE) + TCP_CONNECTING_TIMEOUT_OFFSET)

	// Try to connect with other node over TCP
	numOfRetries := 0
	for {
		_, err 			:= bCastConn.Write(bCastMsg)
		checkError(err)
		err 			 = tcpListener.SetDeadline(time.Now().Add(randomNumber*time.Millisecond))
		checkError(err)
		tcpConn, err 	:= tcpListener.Accept()
		checkError(err)
		
		if err == nil {
			printMsg := "Established tcp connection:"+tcpConn.LocalAddr().String()+"->"+tcpConn.RemoteAddr().String()
			printFromNET(printMsg, nodeID, nodeConnectionState)

			buffer 		:= make([]byte, 8)
			err 	 	 = tcpConn.SetReadDeadline(time.Now().Add(1*time.Second))
			checkError(err)
			_, err 		:= tcpConn.Read(buffer)
			checkError(err)
			
			return tcpConn, uint8(buffer[0]), err
			
		}else if numOfRetries >= ALLOWED_CONNECTION_RETRIES {
			return tcpConn, uint8(0), err
		}
		numOfRetries++
	}
}




// Add new connection to routing table and spawns a thread to manage the connection
func handleNewTcpConnection(conn 				net.Conn								,
							newRoutingEntry 	NodeRoutingTable.RoutingEntry_t			,
							RoutingTable_Ch 	chan *NodeRoutingTable.RoutingTable_t	) {
							
	var routingTable_ptr *NodeRoutingTable.RoutingTable_t	
							
	to_newConnection_Ch		:= make(chan []byte)
	from_newConnection_Ch	:= make(chan []byte)
	newConnection_Mutex_Ch	:= make(chan bool, 1)
	newConnection_Mutex_Ch 	<- true
	
	newRoutingEntry.Receive_Ch 	= from_newConnection_Ch
	newRoutingEntry.Send_Ch 	= to_newConnection_Ch
	newRoutingEntry.Mutex_Ch	= newConnection_Mutex_Ch
	
	IDofNewConnectedNode := newRoutingEntry.NodeID
	
	go NodeSingleConnection.HandleConnection(	conn 					,
												IDofNewConnectedNode	,
												to_newConnection_Ch		,
												from_newConnection_Ch	,
												newConnection_Mutex_Ch 	)
	
	routingTable_ptr = <- RoutingTable_Ch
	routingTable_ptr.Add_new_routing_entries(newRoutingEntry)
	RoutingTable_Ch <- routingTable_ptr
	routingTable_ptr = nil
}




func checkError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}




func getLocalIP(useLocalIP bool) string {
	if useLocalIP {
		return "127.0.0.1"
	}else{
		address, err := net.InterfaceAddrs()
		if err != nil {
			return ""
		}
		for _, addr := range address {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String()
				}
			}
		}
		return ""
	}
}




func setMasterNodeInTable(nodeID uint8, RoutingTable_Ch chan *NodeRoutingTable.RoutingTable_t) error {
	var routingTable_ptr *NodeRoutingTable.RoutingTable_t

	routingTable_ptr = <- RoutingTable_Ch
	err := routingTable_ptr.Set_master_node(nodeID)
	RoutingTable_Ch <- routingTable_ptr
	routingTable_ptr = nil

	return err
}




func printFromNET(str string, nodeID uint8, currentState nodeConnectionState_t) {
	var currentStateStr string
	if currentState == STATE_CONNECTING {
		currentStateStr = "CONNECTING"
	}else if currentState == STATE_MASTER {
		currentStateStr = "MASTER"
	}else if currentState == STATE_SINGLE {
		currentStateStr = "SINGLE"
	}else if currentState == STATE_SLAVE {
		currentStateStr = "SLAVE"
	}else {
		currentStateStr = "n/a"
	}

	fmt.Println("[ NET |", nodeID, "|", currentStateStr, "]:", str)
}