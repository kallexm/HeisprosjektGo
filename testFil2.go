package main

import
(
	"fmt"
	"time"
	"net"
)



func main(){
	fmt.Println(getLocalIP(true))
	
}


func getNanoSecTime() int64 {
	return (time.Now().UnixNano() - (time.Now().UnixNano()/100000)*100000)
}




func getLocalIP(useLocalIP bool) string {
	if useLocalIP {
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
	}else{
		return "127.0.0.1"
	}
}