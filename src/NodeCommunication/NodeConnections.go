package NodeConnections 

import
{
	"fmt"
}

//[FFF]Denne verdien bør hentes fra en config fil og må være forskjellig for alle nodene.
const NodeID = 1

type singleConnection struct{
	ip string
	TCPport string
	isMasterNode bool
	nodeIdentity uint8
}

//[FFF]Blir vanskelig å impementere en del av funksjonene med maps datatypen
connectionTable = make(map[uint8]singleConnection)


//[FFF]Trengs denne funksjonen for å initsialisere egen connection?
func init(){
	connectionTable[1] = singleConnection{ip; GET_IP(), TCPport: GET_TCP_PORT(), isMasterNode: false}
}


func Get_master_node_IP_and_port() {
	
}

func Get_elevator_node_IP_and_port(n int){

}

func Add_new_elevator_node(ip string, TCPport string, nodeIdentity uint8){

} 

func Set_master_node(n int){
	
}




