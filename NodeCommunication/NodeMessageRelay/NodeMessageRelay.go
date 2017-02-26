package NodeMessageRelay

/*[FFF]
1. Motta meldinger fra heis, orderDistributer og meldinger fra andre noder.
2. Sjekk De riktige feltene i meldingen.
3. Sjekk om meldingen skal lokalt, I såfall send gjennom passende channel
4. Hent connections list pekeren fra en delt kanal mellom NodeMessageRelay og og sjekk om den stemmer med noen av nodene som er koblet til gjennom TCP. Hvis den stemmer, send gjennom denne forbindelsen.
*/

import
(
	"../NodeRoutingTable"
	"../../MessageFormat"
	
	"fmt"
	"os"
)
var routingTable_ptr *NodeRoutingTable.RoutingTable_t

func NodeMessageRelay_thread (routingTable_Ch chan *NodeRoutingTable.RoutingTable_t) {
	for {
		routingTable_ptr = <- routingTable_Ch
		for i, tableEntry := range (*routingTable_ptr) {
			select {
			case receivedMsg := <- tableEntry.Receive_Ch:
				msgHeader, data, err := MessageFormat.De_Gen_Msg(receivedMsg)
				eval_error(err)
				
				fmt.Println(msgHeader, data, i)
				/* Implement routing algorithm here */
			default:
			}
		}
		routingTable_Ch <- routingTable_ptr
		routingTable_ptr = nil
	}
}


func eval_error(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}