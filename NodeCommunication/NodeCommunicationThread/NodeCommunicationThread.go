package NodeCommunicationThread

import
(
	"fmt"
)



func NC_thr(OD_to_NC_Ch <-chan []byte, NC_to_OD_Ch chan<- []byte, NC_exit_Ch chan<- bool) {
	str := "Hello, this is a long string that is going to be changed into a byte array, pased through a channel, converted to string again and printed out\n"
	byteArray := []byte(str)
	NC_to_OD_Ch <- byteArray
	
	byteArray = <-OD_to_NC_Ch
	str = string(byteArray)

	fmt.Println(str)
	fmt.Println(byteArray)
	
	NC_exit_Ch <- true
}