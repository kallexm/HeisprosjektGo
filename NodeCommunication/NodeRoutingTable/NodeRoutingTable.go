package NodeRoutingTable
/*
||	File: NodeRoutingTable
||
||	Authors: 
||
||	Date: 	 Spring 2017
||	Course:  TTK4145 - Real-time Programming, NTNU
||	
||	Summary of File:
||		Contains the routing table to be used as lookup by NodeMessageRelay and which should be updated by NodeConnectionManager.
||		Defines the RoutingtableEntry_t and Routingtable_t types. and functions to interact with these.
||
*/

import
(
	"errors"
)


type RoutingEntry_t struct{
	NodeID 		uint8
	IsMaster 	bool
	IsElev 		bool
	IsNet 		bool
	IsBackup 	bool
	IsOrderDist bool
	IsExtern 	bool
	
	Receive_Ch 	<-chan 	[]byte
	Send_Ch    	chan<- 	[]byte
	Mutex_Ch	chan 	bool

}

type RoutingTable_t []RoutingEntry_t

var routingTable RoutingTable_t




func Get_reference_to_routing_table() (*RoutingTable_t) {
	return &routingTable
}




func (rt *RoutingTable_t) Add_new_routing_entries(newRoutingEntries ...RoutingEntry_t) {
	*rt = append(*rt, newRoutingEntries...)
}




func (rt *RoutingTable_t) Remove_routing_entry(nodeID uint8) (RoutingEntry_t, error){
	var removedRoutingTable RoutingEntry_t
	var err error
	err = nil
	for i := 0; i < len(*rt); i++{
		if (*rt)[i].IsExtern == true && (*rt)[i].NodeID == nodeID {
			removedRoutingTable = (*rt)[i]
			(*rt)[i] = (*rt)[len(*rt)-1]
			(*rt)[len(*rt)-1] = RoutingEntry_t{}
			(*rt) = (*rt)[:len(*rt)-1]
			return removedRoutingTable, err
		}
	}
	err = errors.New("Connection not in routing table")
	return removedRoutingTable, err
}




func (rt *RoutingEntry_t) Get_receive_Ch() <-chan []byte {
	return rt.Receive_Ch
}




func (rt *RoutingTable_t) Set_master_node(nodeID uint8) error {
	var	foundnodeIDinTable = false
	for i := 0; i < len(*rt); i++ {
		if (*rt)[i].IsMaster == true {
			(*rt)[i].IsMaster = false
		}
		if (*rt)[i].NodeID == nodeID && ( (*rt)[i].IsOrderDist == true || (*rt)[i].IsExtern == true ) {
			(*rt)[i].IsMaster  = true
			foundnodeIDinTable = true
		}
	}
	if foundnodeIDinTable == false {
		return errors.New("No entry with nodeID in routing table")
	}else{
		return nil
	}
}