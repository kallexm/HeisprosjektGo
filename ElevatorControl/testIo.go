package main

import 
(
	"./io"
	"time"
)


func main() {
	io.Io_init()
	for{
		io.Io_clear_bit(0x300+13)
		time.Sleep(time.Second*1)
		io.Io_set_bit(0x300+13)
		time.Sleep(time.Second*1)
	}

}