package NodeConnections 

import
(
	"fmt"
)


const NodeID = 1

type singleConnection struct{
	ip string
	TCPport string
	isMasterNode bool
	nodeIdentity uint8
}


connectionTable = make([]singleConnection, 0)


//[FFF]Trengs denne funksjonen for å initsialisere egen connection?
func init(){
	connectionTable[1] = singleConnection{ip; GET_IP(), TCPport: GET_TCP_PORT(), isMasterNode: false}
}


func Get_master_node_IP_and_port() (string, string, error){
	var err error
	for i := 0; i < len(connectionTable); i++ {
		if connectionTable[i].isMasterNode == true {
			return connectionTable[i].ip, connectionTable[i].TCPport, err
		}
	}
}

func Get_elevator_node_IP_and_port(n int){
	return connectionTable[n].ip, connectionTable[n].TCPport, err
}

func Add_new_elevator_node(ip string, TCPport string, nodeIdentity uint8){
	singleConnection 
	connectionTable = append(connectionTable, singleConnection{ip: ip, TCPport: TCPport, nodeIdentity: nodeIndentity})
	return err
} 

func Set_master_node(n int){
	for i := 0; i < len(connectionTable); i++ {
		if connectionTable[i].isMasterNode == true {
			connectionTable[i].isMasterNode = false
		}
	}
	connectionTable[n].isMasterNode = true
	return err
}