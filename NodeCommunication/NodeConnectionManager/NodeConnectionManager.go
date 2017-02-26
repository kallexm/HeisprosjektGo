package NodeConnectionManager
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
	//"../NodeSingleConnection"
	"../NodeMessageRelay"
	"../NodeRoutingTable"
	"../../MessageFormat"
	
	"fmt"
	"net"
	"time"
	
)

const broadCastToPort = "60002"
const broadCastFromPort = "60022"

func NodeConnectionManager_thread(from_OrderDist_Ch <-chan []byte, to_OrderDist_Ch chan<- []byte,
								  from_ElevCtrl_Ch  <-chan []byte, to_ElevCtrl_Ch  chan<- []byte,
								  NodeComm_exit_Ch  chan<- bool  , nodeID uint8                 ) {
								  
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
	
	
	bCastToAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort("255.255.255.255", broadCastToPort))
	CheckError(err)
	bCastFromAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(GetLocalIP(), broadCastFromPort))
	CheckError(err)
	bCastConn, err := net.DialUDP("udp", bCastFromAddr, bCastToAddr)
	CheckError(err)
	defer bCastConn.Close()
	
	listenForMasterAddr, err := net.ResolveTCPAddr("tcp", ":0")
	CheckError(err)
	listenForMaster, err := net.ListenTCP("tcp", listenForMasterAddr)
	CheckError(err)
	defer listenForMaster.Close()
	
	fmt.Println("Searching for other nodes...")
	bCastMsg := make([]byte, 0)
	bCastMsg = append(bCastMsg, byte(nodeID), byte(MessageFormat.MASTER))
	bCastMsg = append(bCastMsg, []byte(listenForMaster.Addr().String())...)
	
	i := 0
	for {
		_, err := bCastConn.Write(bCastMsg)
		CheckError(err)
		err = listenForMaster.SetDeadline(time.Now().Add(2*time.Second))
		CheckError(err)
		tcpConnFromMaster, err := listenForMaster.Accept()
		CheckError(err)
		if err == nil {
			fmt.Println("Connection succeeded...")
			fmt.Println(tcpConnFromMaster)
		}else{
			i++
			fmt.Println(err)
			if i >= 3 {
				fmt.Println("Failed on third attempt...\nStarting as masterNode...")
				break
			}
		}
	}
	
	
	
	
	
	NodeComm_exit_Ch <- true
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