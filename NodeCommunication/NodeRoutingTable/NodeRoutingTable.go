package NodeRoutingTable
/*
||	File: NodeRoutingTable
||
||	Author:  Andreas Hanssen Moltumyr	
||	Partner: Martin Mostad
||	Date: 	 Spring 2017
||	Course:  TTK4145 - Real-time Programming, NTNU
||	
||	Summary of File:
||		Contains the routing table to be used as lookup by NodeMessageRelay and which should be updated by NodeConnectionManager.
||		Defines the RoutingtableEntry_t and Routingtable_t types. and functions to interact with these.
||
*/

/*[KKK]
1. Holde en connectionlist som benyttes av NodeMessageRelay og oppdateres av NodeConnectionManager
2. connectionlisten bør være redundant lagret.
3. En peker til dette connectionlisten ligger i en channel som er delt mellom NodeMessageRelay og NodeConnectionManager. Bare en av disse kan holde pekeren om gangen. Når en av nodene er ferding med å bruke den skal den sendes tilbake til channelen.
*/

import
(
	//"fmt"
)


type RoutingEntry_t struct{
	NodeID uint8
	IsMaster bool
	IsElev bool
	IsNet bool
	IsBackup bool
	IsExtern bool
	
	Receive_Ch <-chan []byte
	Send_Ch    chan<- []byte
}

type RoutingTable_t []RoutingEntry_t


var routingTable RoutingTable_t


func Get_ptr_to_routing_table() (*RoutingTable_t) {
	return &routingTable
}

func (rt *RoutingTable_t) Add_new_routing_entries(newRoutingEntries ...RoutingEntry_t) {
	*rt = append(*rt, newRoutingEntries...)
}

func (rt *RoutingTable_t) Remove_routing_entry(nodeID uint8) {
	for i := 0; i < len(*rt); i++{
		if (*rt)[i].IsExtern == true && (*rt)[i].NodeID == nodeID {
			(*rt)[i] = (*rt)[len(*rt)-1]
			(*rt)[len(*rt)-1] = RoutingEntry_t{}
			(*rt) = (*rt)[:len(*rt)-1]
			break
		}
	}
}

func (rt *RoutingEntry_t) Get_receive_Ch() <-chan []byte {
	return rt.Receive_Ch
}

/*func (rt *RoutingEntry_t) Get_send_Ch() chan<- []byte {
	return rt.send_Ch
}*/

/*func (rt *RoutingTable_t) Contains_entry_with(nodeID uint8) {
//connectionTable = make([]singleConnection, 0)
}*/

/*
func Get_master_node_IP_and_port() (string, string, error){
	var err error
	for i := 0; i < len(connectionTable); i++ {
		if connectionTable[i].isMasterNode == true {
			return connectionTable[i].ip, connectionTable[i].TCPport, err
		}
	}
}


func Set_master_node(n int){
	for i := 0; i < len(connectionTable); i++ {
		if connectionTable[i].isMasterNode == true {
			connectionTable[i].isMasterNode = false
		}
	}
	connectionTable[n].isMasterNode = true
	return err
}*/