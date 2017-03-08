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

const broadCastToPort = "60002"
const broadCastFromPort = "60022"
var isNodeMaster bool


/*
|| NodeConnectionManager_thread(...) should be called as a goroutine.
||		It manages the connections a single elevator node has to all the other elevator nodes which is connected.
||		This function:
||			1. Sets up the routing table stored in file/package NodeRoutingTable with inital routing entries.
||			2. Tries to establish a TCP connection to another node which is listening on port broadCastToPort with the help of UDP broadcasting.
||			3. Starts to listen for new nodes on the network if it could not find an other node to connect to.
||			4. Updates the routing table when necessary.
*/
func NodeConnectionManager_thread(from_OrderDist_Ch <-chan []byte, to_OrderDist_Ch chan<- []byte,
								  from_ElevCtrl_Ch  <-chan []byte, to_ElevCtrl_Ch  chan<- []byte,
								  NodeComm_exit_Ch  chan<- bool  , nodeID uint8                 ) {
	
	//Setting up initial starting
	nodeComm_to_MsgRelay_Ch := make(chan []byte)
	MsgRelay_to_nodeComm_Ch := make(chan []byte)
	
	RoutingTable_Ch := make(chan *NodeRoutingTable.RoutingTable_t, 1)
	routingTable_ptr := NodeRoutingTable.Get_ptr_to_routing_table()
	routingTable_ptr.Add_new_routing_entries(NodeRoutingTable.RoutingEntry_t{NodeID: nodeID, 			   Receive_Ch: from_OrderDist_Ch,       Send_Ch: to_OrderDist_Ch},
											 NodeRoutingTable.RoutingEntry_t{NodeID: nodeID, IsElev: true, Receive_Ch: from_ElevCtrl_Ch,		Send_Ch: to_ElevCtrl_Ch},
											 NodeRoutingTable.RoutingEntry_t{NodeID: nodeID, IsNet:  true, Receive_Ch: nodeComm_to_MsgRelay_Ch, Send_Ch: MsgRelay_to_nodeComm_Ch})
	
	
	
	RoutingTable_Ch <- routingTable_ptr
	routingTable_ptr = nil
	go NodeMessageRelay.NodeMessageRelay_thread(RoutingTable_Ch)
	
	
	// Start connection Sequence
	isNodeMaster = false
	
	conn, IDofNewConnectedNode, err := connect_to_other_Node(isNodeMaster, nodeID)
	if err != nil {
		CheckError(err)
		isNodeMaster = true
	}else{
		to_newConnection_Ch := make(chan []byte)
		from_newConnection_Ch := make(chan []byte)
		
		go NodeSingleConnection.HandleConnection(conn, IDofNewConnectedNode, to_newConnection_Ch, from_newConnection_Ch)
		
		newRoutingEntry := NodeRoutingTable.RoutingEntry_t{NodeID: IDofNewConnectedNode,
														   IsExtern: true,
														   IsMaster: true,
														   Receive_Ch: from_newConnection_Ch,
														   Send_Ch: to_newConnection_Ch}
		routingTable_ptr = <- RoutingTable_Ch
		routingTable_ptr.Add_new_routing_entries(newRoutingEntry)
		RoutingTable_Ch <- routingTable_ptr
		routingTable_ptr = nil
	}
	
	for {
		if isNodeMaster == true {
			
		}
		
		if isNodeMaster == false {
			
		}
	}

	
	NodeComm_exit_Ch <- true
}





// Returns a valid connection if a MasterNode was found. Else it returns an error.
func connect_to_other_Node(thisNodeIsMaster bool, nodeID uint8) (net.Conn, uint8, error){

	// Setting up a UDP broadcast socket
	bCastToAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort("255.255.255.255", broadCastToPort))
	CheckError(err)
	bCastFromAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(GetLocalIP(), broadCastFromPort))
	CheckError(err)
	bCastConn, err := net.DialUDP("udp", bCastFromAddr, bCastToAddr)
	CheckError(err)
	defer bCastConn.Close()
	
	// Setting up a TCP listener socket
	listenForMasterAddr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(GetLocalIP(), "0"))
	CheckError(err)
	listenForMaster, err := net.ListenTCP("tcp", listenForMasterAddr)
	CheckError(err)
	defer listenForMaster.Close()
	
	// Make broadcast message
	bCastMsg := make([]byte, 0)
	bCastMsg = append(bCastMsg, byte(nodeID))
	if thisNodeIsMaster == false {
		bCastMsg = append(bCastMsg, 0x0)
	}else{
		bCastMsg = append(bCastMsg, 0x1)
	}
	bCastMsg = append(bCastMsg, []byte(listenForMaster.Addr().String())...)
	
	// [KKK] Debug prints
	//fmt.Println(bCastMsg)
	//fmt.Println(string(bCastMsg))
	
	// Try to connect
	i := 0
	for {
		_, err := bCastConn.Write(bCastMsg)
		CheckError(err)
		err = listenForMaster.SetDeadline(time.Now().Add(1*time.Second))
		CheckError(err)
		tcpConnFromMaster, err := listenForMaster.Accept()
		//CheckError(err)
		
		if err == nil {
			buffer := make([]byte, 256)
			err = tcpConnFromMaster.SetReadDeadline(time.Now().Add(1*time.Second))
			CheckError(err)
			_, err := tcpConnFromMaster.Read(buffer)
			CheckError(err)
			msgHeader, _, err := MessageFormat.Decode_msg(buffer)
			CheckError(err)	
			
			return tcpConnFromMaster, msgHeader.FromNodeID, err
			
		}else if i >= 2 {
			return tcpConnFromMaster, uint8(0), err
		}
		i++
	}
}



func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}



func GetLocalIP() string {
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