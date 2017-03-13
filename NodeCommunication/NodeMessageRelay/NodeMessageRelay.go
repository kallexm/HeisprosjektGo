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
	//"os"
	//"time"
)




var routingTable_ptr *NodeRoutingTable.RoutingTable_t

func Thread (routingTable_Ch chan *NodeRoutingTable.RoutingTable_t) {
	fmt.Println("Starting messageRelay")
	for {
		//time.Sleep(time.Millisecond*50)
		routingTable_ptr = <- routingTable_Ch

		//fmt.Printf("%+v", *routingTable_ptr)
		//fmt.Println()

		for i, tableEntry := range (*routingTable_ptr) {
			select {
			case tableEntry.Mutex_Ch <- true:
				continueFor := true
				for continueFor {
					select {
					case receivedMsg := <- tableEntry.Receive_Ch:
						msgHeader, data, err := MessageFormat.Decode_msg(receivedMsg)

						if false {
							eval_error(err)
							fmt.Println(tableEntry)
							fmt.Println("Message in Relay:", msgHeader, data, i)
						}
						
						for j, searchTableEntry := range (*routingTable_ptr) {
							if false {
							eval_error(err)
							fmt.Println(tableEntry)
							fmt.Println("Message in Relay:", msgHeader, data, j)
							}
							//-----------------------------------------------------------------
							// Routing Entries
							if msgHeader.To == MessageFormat.MASTER && searchTableEntry.IsMaster == true {
								if msgHeader.From == MessageFormat.ELEVATOR && msgHeader.FromNodeID == 0 {
									msgHeader.FromNodeID = tableEntry.NodeID
									sendMsg, _ := MessageFormat.Encode_msg(msgHeader, data)
									searchTableEntry.Send_Ch <- sendMsg
									break
								}else{
									searchTableEntry.Send_Ch <- receivedMsg
									break
								}
								
							}else if msgHeader.To == MessageFormat.ELEVATOR && (searchTableEntry.IsElev == true || searchTableEntry.IsExtern == true) && msgHeader.ToNodeID == searchTableEntry.NodeID {
								searchTableEntry.Send_Ch <- receivedMsg
								break
								
							}else if msgHeader.To == MessageFormat.BACKUP && (searchTableEntry.IsBackup == true || searchTableEntry.IsExtern == true) {
								searchTableEntry.Send_Ch <- receivedMsg
								break
								
							}else if msgHeader.To == MessageFormat.NODE_COM && searchTableEntry.IsNet == true {
								searchTableEntry.Send_Ch <- receivedMsg
								break
								
							}else if msgHeader.To == MessageFormat.ORDER_DIST && searchTableEntry.IsOrderDist == true {
								searchTableEntry.Send_Ch <- receivedMsg
								break
								
							}else if j == len(*routingTable_ptr){
								// Drop package if it doesn't match any of the above filters/masks
								fmt.Println("Package with header", msgHeader, "dropped because of no matching filters!")
							}
							//-----------------------------------------------------------------
						}
					default:
						// Do nothing
					}
					select {
					case <- tableEntry.Mutex_Ch:
						continueFor = false
					default:
						// Do nothing
					}
				}
			default:
				// Do nothing
			}
			
			
		}
		routingTable_Ch <- routingTable_ptr
		routingTable_ptr = nil
	}
}


func eval_error(err error) {
	if err != nil {
		fmt.Println(err)
		//os.Exit(0)
	}
}

