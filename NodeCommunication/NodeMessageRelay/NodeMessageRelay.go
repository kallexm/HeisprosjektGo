package NodeMessageRelay
/*
||	File: NodeMessageRelay
||
||	Author:  Andreas Hanssen Moltumyr	
||	Partner: Martin Mostad
||	Date: 	 Spring 2017
||	Course:  TTK4145 - Real-time Programming, NTNU
||	
||	Summary of File:
||		Contains a goroutine function which gets messages from channels and
||		forwards them based on lookup in the routing table in NodeRoutingTable
||		which is maintained by NodeConnectionManager.
||
*/

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
				msgHeader, data, err := MessageFormat.Decode_msg(receivedMsg)
				eval_error(err)
				
				fmt.Println("Message in Relay:", msgHeader, data, i)
				//-----------------------------------------------------------------
				// Implement routing algorithm here
				for i, searchTableEntry := range (*routingTable_ptr) {
					if msgHeader.To == MessageFormat.MASTER && searchTableEntry.IsMaster == true {
						searchTableEntry.Send_Ch <- receivedMsg
						break
						
					}else if msgHeader.To == MessageFormat.ELEVATOR && (searchTableEntry.IsElev == true || searchTableEntry.IsExtern == true) && msgHeader.ToNodeID == searchTableEntry.NodeID {
						searchTableEntry.Send_Ch <- receivedMsg
						break
						
					}else if msgHeader.To == MessageFormat.BACKUP && (searchTableEntry.IsBackup == true || searchTableEntry.IsExtern == true) && msgHeader.ToNodeID == searchTableEntry.NodeID {
						searchTableEntry.Send_Ch <- receivedMsg
						break
						
					}else if msgHeader.To == MessageFormat.NODE_COM && searchTableEntry.IsNet == true {
						searchTableEntry.Send_Ch <- receivedMsg
						break
						
					}else if msgHeader.To == MessageFormat.ORDER_DIST && searchTableEntry.IsOrderDist == true {
						searchTableEntry.Send_Ch <- receivedMsg
						break
						
					}else if i == len(*routingTable_ptr){
						// Drop package if it doesn't match any of the above filters/masks
						fmt.Println("Package with header", msgHeader, "dropped because of no matching filters!")
					}
				}
				
				//-----------------------------------------------------------------
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

