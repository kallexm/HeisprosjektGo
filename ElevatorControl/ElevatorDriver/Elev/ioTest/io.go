package io 

import "fmt"

func Io_set_bit(channel int){
	fmt.Println("set bit: ", channel)
}

func Io_clear_bit(channel int){
	fmt.Println("clear bit: ", channel)
}

func Io_read_bit(channel int) int{
	fmt.Println("read bit: ", channel)
	return 1
}

func Io_read_analog(channel int) int{
	fmt.Println("read annalog: ", channel)
	return 1
}

func Io_write_analog(channel int, value int) {
	fmt.Println("write analgo: ",channel)
}

func Io_init()int{
	fmt.Println("init")
	return 1
}