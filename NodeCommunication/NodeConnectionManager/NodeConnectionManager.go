package NodeConnectionManager
/*
||	File: NodeConnectionManager 
||
||	Author:  Andreas Hanssen Moltumyr	
||	Partner: Martin Mostad
||	Date: 	 Spring 2017
||	Course:  TTK4145 - Real-time Programming, NTNU
||	
||	Summary of File:
||		It manages the connections a single elevator node has to all the other elevator nodes which is connected.
||			1. Initiates the network module
||			2. Establishes and stops network connections.
||      	3. Updates the routing table if a new node connects or disconnects.	
||
*/


/*[KKK]
1. Når noden starter for første gang skal den broadcast i et forsøk på å nå andre allerede eksisterende masterNode.
	a. Hvis svar skal den opptre som en slave i node nettverket og kobles opp til de andre nodene med TCP.
2. Hvis den ikke får svar skal den ta rollen som master og starte å lytte etter andre noder som prøver å koble seg på.
3. Hvis den hører en node broadcaste skal den opprette en TCP forbindelse og spawne en ny thread av type NodeSingleConnection, gi denne en channel som den lagrer i NodeRoutingTable sammen med annen informasjon om forbindelsen.
4. Sende meldinger til master om nye noder som kobler seg på eller noder som faller bort og som ikke kan motta meldinger, slik at master kan holde orden på og oppdatere køene sine riktig.
5. Oppdatere NodeRoutingTable hvis dette er nødvendig av andre grunner.
(6.) Sende en annen type broadcast av og til for å finne ut om andre master nodeNettverk eksisterer og evt. merge......
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
)

const broadCastToPort 	= "60002"
const broadCastFromPort = "60022"
const useLocalIP = false


type nodeConnectionState_t uint8
const(
	STATE_CONNECTING nodeState_t = iota
	STATE_MASTER
	STATE_SINGLE
	STATE_SLAVE
)

var nodeConnectionState 		nodeConnectionState_t
var prev_nodeConnectionState 	nodeConnectionState_t





/*
|| Thread(...) should be called as a goroutine.
||		It manages the connections a single elevator node has to all the other elevator nodes which is connected.
||		This function:
||			1. Sets up the routing table stored in file/package NodeRoutingTable with inital routing entries.
||			2. Tries to establish a TCP connection to another node which is listening on port broadCastToPort with the help of UDP broadcasting.
||			3. Starts to listen for new nodes on the network if it could not find an other node to connect to.
||			4. Updates the routing table when necessary.
*/

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
	
	
	// Start connection Sequence
	nodeConnectionState = STATE_CONNECTING
	
	for {
		if nodeConnectionState == STATE_CONNECTING {

			fmt.Println("Connecting to remote node...")
			conn, IDofNewConnectedNode, err := connect_to_other_Node(isNodeMaster, nodeID)
			if err != nil {
				checkError(err)
				fmt.Println("Failed to connect to remote node...")
				err := setMasterNodeInTable(nodeID, RoutingTable_Ch)
				checkError(err)
				fmt.Println("Node with ID", nodeID, "changed to master")
				nodeConnectionState = STATE_MASTER

			}else{
				newRoutingEntry := NodeRoutingTable.RoutingEntry_t{	NodeID: IDofNewConnectedNode	,
															IsExtern: true					,
															IsMaster: true					}
				handleNewTcpConnection(conn, newRoutingEntry, RoutingTable_Ch)
				fmt.Println("Connection to node with ID", IDofNewConnectedNode, "was successful")
				nodeConnectionState = STATE_SLAVE
			}




		}else if nodeConnectionState == STATE_MASTER {
			
			if prev_nodeConnectionState != nodeConnectionState {
				fmt.Println("Start to listen for other nodes...")
				listenerAddress, err := net.ResolveUDPAddr("udp", ":"+string(broadCastToPort))
				checkError(err)
				bCastListener, err := net.ListenUDP("udp", listenerAddress)
				fmt.Println("Address of bCastListener:", bCastListener.LocalAddr())
				checkError(err)
				defer bCastListener.Close()
			}
			
			bCastListener.SetReadDeadline(time.Now().Add(1*time.Second))
			checkError(err)
			buffer := make([]byte, 32)
			_, _, err := bCastListener.ReadFromUDP(buffer)
			nodeIDofBcaster := uint8(buffer[0])
			tcpListenerAddrOfBcaster := string(buffer[2:])
					
					
			tcpConn, err := net.DialTimeout("tcp", tcpListenerAddrOfBcaster, 1*time.Second)

			//checkError(err)
			if err == nil {
				fmt.Println("Established tcp connection:", tcpConn.LocalAddr(), "->", tcpConn.RemoteAddr())
				nodeIDmsg := make([]byte, 1)
				nodeIDmsg[0] = byte(nodeID)
				_, err := tcpConn.Write(nodeIDmsg)
				checkError(err)
				newRoutingEntry := NodeRoutingTable.RoutingEntry_t{	NodeID: nodeIDofBcaster	,
																	IsExtern: true			,
																	IsBackup: true			}
				handleNewTcpConnection(tcpConn, newRoutingEntry, RoutingTable_Ch)
				fmt.Println("Connection to node with ID", nodeIDofBcaster, "was successful")
			}




		}else if nodeConnectionState == STATE_SINGLE {
			

		}else if nodeConnectionState == STATE_SLAVE {
			select {
			case msg := <- MsgRelay_to_nodeComm_Ch:
				msgHeader, _, err := MessageFormat.Decode_msg(msg)
				checkError(err)
				if msgHeader.MsgType == MessageFormat.NODE_DISCONNECTED {
					routingTable_ptr = <- RoutingTable_Ch
					removedRoutingEntry, err := routingTable_ptr.Remove_routing_entry(msgHeader.FromNodeID)
					
					if false { //Debug print
						fmt.Println("removedRoutingEntry:", removedRoutingEntry)
						fmt.Println()
						fmt.Println("routingTable_ptr after removal:", *routingTable_ptr)
					} //Debug print
					
					RoutingTable_Ch <- routingTable_ptr
					routingTable_ptr = nil

					checkError(err)
					if err == nil {
						replyHeader := MessageFormat.MessageHeader_t{ToNodeID: msgHeader.FromNodeID				,
																 From: MessageFormat.NODE_COM				,
																 MsgType: MessageFormat.NODE_DISCONNECTED	} 
						msg, _ = MessageFormat.Encode_msg(replyHeader, "")
						removedRoutingEntry.Send_Ch <- msg
						fmt.Println("Disconnected node with ID", removedRoutingEntry.NodeID)
					}
				}else{
					// Maybe Do something
				}
			default:
				// Do Nothing
			}




		}else{

		}


	}
	NodeComm_exit_Ch <- true	
	

	for {
		//time.Sleep(1*time.Second)
		
			if isNodeMaster == true {
				//---------------------------------------------------------------net.JoinHostPort(getLocalIP(useLocalIP), broadCastToPort)
				// What to do if the node should be masterNode on the network

				// Setting up UDP listener
				fmt.Println("Start to listen for other nodes...")
				listenerAddress, err := net.ResolveUDPAddr("udp", ":"+string(broadCastToPort))
				checkError(err)
				bCastListener, err := net.ListenUDP("udp", listenerAddress)
				fmt.Println("Address of bCastListener:", bCastListener.LocalAddr())
				checkError(err)
				defer bCastListener.Close()
				
				for {
					bCastListener.SetReadDeadline(time.Now().Add(1*time.Second))
					checkError(err)
					buffer := make([]byte, 32)
					_, _, err := bCastListener.ReadFromUDP(buffer)
					nodeIDofBcaster := uint8(buffer[0])
					tcpListenerAddrOfBcaster := string(buffer[2:])
					
					
					tcpConn, err := net.DialTimeout("tcp", tcpListenerAddrOfBcaster, 1*time.Second)

					//checkError(err)
					if err == nil {
						fmt.Println("Established tcp connection:", tcpConn.LocalAddr(), "->", tcpConn.RemoteAddr())
						nodeIDmsg := make([]byte, 1)
						nodeIDmsg[0] = byte(nodeID)
						_, err := tcpConn.Write(nodeIDmsg)
						checkError(err)
						newRoutingEntry := NodeRoutingTable.RoutingEntry_t{	NodeID: nodeIDofBcaster	,
																			IsExtern: true			,
																			IsBackup: true			}
						handleNewTcpConnection(tcpConn, newRoutingEntry, RoutingTable_Ch)
						fmt.Println("Connection to node with ID", nodeIDofBcaster, "was successful")
					}
					
				
				}
				//---------------------------------------------------------------
			}
			
			if isNodeMaster == false {
				//---------------------------------------------------------------
				// what to do if the node is not masterNode
				
				//---------------------------------------------------------------
			}
			//---------------------------------------------------------------
			// Stuff that should be done no matter the value of isNodeMaster
			
			select {
			case msg := <- MsgRelay_to_nodeComm_Ch:
				msgHeader, _, err := MessageFormat.Decode_msg(msg)
				checkError(err)
				if msgHeader.MsgType == MessageFormat.NODE_DISCONNECTED {
					routingTable_ptr = <- RoutingTable_Ch
					removedRoutingEntry, err := routingTable_ptr.Remove_routing_entry(msgHeader.FromNodeID)
					
					if false { //Debug print
						fmt.Println("removedRoutingEntry:", removedRoutingEntry)
						fmt.Println()
						fmt.Println("routingTable_ptr after removal:", *routingTable_ptr)
					} //Debug print
					
					RoutingTable_Ch <- routingTable_ptr
					routingTable_ptr = nil

					checkError(err)
					if err == nil {
						replyHeader := MessageFormat.MessageHeader_t{ToNodeID: msgHeader.FromNodeID				,
																 From: MessageFormat.NODE_COM				,
																 MsgType: MessageFormat.NODE_DISCONNECTED	} 
						msg, _ = MessageFormat.Encode_msg(replyHeader, "")
						removedRoutingEntry.Send_Ch <- msg
						fmt.Println("Disconnected node with ID", removedRoutingEntry.NodeID)
					}
				}else{
					// Maybe Do something
				}
			default:
				// Do Nothing
			}
		//---------------------------------------------------------------
	}
	
	NodeComm_exit_Ch <- true
}





// Returns a valid connection if a MasterNode was found. Else it returns an error.
func connect_to_other_Node(thisNodeIsMaster bool, nodeID uint8) (net.Conn, uint8, error){

	// Setting up a UDP broadcast socket
	bCastToAddr, err	:= net.ResolveUDPAddr("udp", net.JoinHostPort("255.255.255.255"	, broadCastToPort	))
	checkError(err)
	bCastConn, err 		:= net.DialUDP("udp", nil, bCastToAddr)
	checkError(err)
	defer bCastConn.Close()
	
	// Setting up a TCP listener socket 
	listenAddr, err 	:= net.ResolveTCPAddr("tcp", net.JoinHostPort(getLocalIP(useLocalIP), "0"))
	checkError(err)
	tcpListener, err 		:= net.ListenTCP("tcp", listenAddr)
	checkError(err)
	defer tcpListener.Close()
	
	// Make broadcast message
	bCastMsg := make([]byte, 0)
	bCastMsg = append(bCastMsg, byte(nodeID))
	if thisNodeIsMaster == false {
		bCastMsg = append(bCastMsg, 0x0)
	}else{
		bCastMsg = append(bCastMsg, 0x1)
	}
	bCastMsg = append(bCastMsg, []byte(tcpListener.Addr().String())...)
	
	// Try to connect
	fmt.Println("Start connection sequence")
	i := 0
	for {
		_, err 			:= bCastConn.Write(bCastMsg)
		checkError(err)
		err 			 = tcpListener.SetDeadline(time.Now().Add(1*time.Second))
		checkError(err)
		tcpConn, err 	:= tcpListener.Accept()
		checkError(err)
		
		if err == nil {
			fmt.Println("Established tcp connection:", tcpConn.LocalAddr(), "->", tcpConn.RemoteAddr())
			buffer 		:= make([]byte, 8)
			err 	 	 = tcpConn.SetReadDeadline(time.Now().Add(1*time.Second))
			checkError(err)
			_, err 		:= tcpConn.Read(buffer)
			checkError(err)
			
			return tcpConn, uint8(buffer[0]), err
			
		}else if i >= 2 {
			return tcpConn, uint8(0), err
		}
		i++
	}
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



func handleNewTcpConnection(conn 				net.Conn								,
							newRoutingEntry 	NodeRoutingTable.RoutingEntry_t			,
							RoutingTable_Ch 	chan *NodeRoutingTable.RoutingTable_t	) {
							
	var routingTable_ptr *NodeRoutingTable.RoutingTable_t	
							
	to_newConnection_Ch		:= make(chan []byte)
	from_newConnection_Ch	:= make(chan []byte)
	newConnection_Mutex_Ch	:= make(chan bool, 1)
	newConnection_Mutex_Ch 	<- true
	// Blir channelene værende i routingTable'et selv om plassen de
	// lages går ut av scope, når det gjøres på denne måten?!?
	
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


func setMasterNodeInTable(	nodeID 			uint8, 
							RoutingTable_Ch chan *NodeRoutingTable.RoutingTable_t) error {

	var routingTable_ptr *NodeRoutingTable.RoutingTable_t

	routingTable_ptr = <- RoutingTable_Ch
	err := routingTable_ptr.Set_master_node(nodeID)
	RoutingTable_Ch <- routingTable_ptr
	routingTable_ptr = nil

	return err
}