package Elev

import 
(
	"./ioTest"
	//"./io"
	"./channels"
	"errors"
	"fmt"
)


var lamp_channel_map 	= initLampChannelsMap()
var button_channel_map 	= initButtonChannelsMap()
const N_FLOORS = 4
const N_BUTTONS = 3
const MOTOR_SPEED = 2800

type ButtonType int
const(
	Up ButtonType = iota
	Down
	Comand
	Door
)

type MotorDir int
const(
	DirDown = iota -1
	DirStop
	DirUp

)  

func initLampChannelsMap() map[ButtonType]map[int]int{
	lamp_channel_map := map[ButtonType]map[int]int{
		Up:map[int]int{1:channels.LIGHT_UP1,2:channels.LIGHT_UP2,3:channels.LIGHT_UP3,4:channels.LIGHT_UP4},
		Down:map[int]int{1:channels.LIGHT_DOWN1,2:channels.LIGHT_DOWN2,3:channels.LIGHT_DOWN3,4:channels.LIGHT_DOWN4},
		Comand:map[int]int{1:channels.LIGHT_COMMAND1,2:channels.LIGHT_COMMAND2,3:channels.LIGHT_COMMAND3,4:channels.LIGHT_COMMAND4}}
	return lamp_channel_map
	
}

func getLampChannel(floor int, lamp ButtonType) (int, error) {
	if lamp_channel_map[lamp][floor] == 0{
		return 0, errors.New("Index out of bounds error 001")
	}
	return lamp_channel_map[lamp][floor], nil
}


func initButtonChannelsMap() map[ButtonType]map[int]int{
	button_channel_map := map[ButtonType]map[int]int{
		Up:map[int]int{1:channels.BUTTON_UP1,2:channels.BUTTON_UP2,3:channels.BUTTON_UP3,4:channels.BUTTON_UP4},
		Down:map[int]int{1:channels.BUTTON_DOWN1,2:channels.BUTTON_DOWN2,3:channels.BUTTON_DOWN3,4:channels.BUTTON_DOWN4},
		Comand:map[int]int{1:channels.BUTTON_COMMAND1,2:channels.BUTTON_COMMAND2,3:channels.BUTTON_COMMAND3,4:channels.BUTTON_COMMAND4}}
	return button_channel_map
	
}

func getButtonChannel(floor int, lamp ButtonType) (int, error) {
	if button_channel_map[lamp][floor] == 0{
		return 0, errors.New("Index out of bounds error 001")
	}
	return button_channel_map[lamp][floor], nil
}

func ElevInit() error{
	initSuccess := io.Io_init()
	if(initSuccess != 1){
		return errors.New("Unable to initialize elevator error 002")
	}
	for f := 1; f < N_FLOORS; f++ {
		var b ButtonType
		for b  = 0; b < N_BUTTONS; b++ {
			if err := ElevSetButtonLamp(b,f,0); err != nil {
				fmt.Println("f: ",f,"b: ", b)
				return err
			}
		}
	}
	ElevSetDoorOpenLamp(0)
	if err := ElevSetFloorIndicator(1); err != nil{
		return err
	}
	return nil
}

func ElevSetMotorDirection(dir MotorDir){
	if (dir == DirStop){
		io.Io_write_analog(channels.MOTOR, 0)
	} else if (dir == DirUp){
		io.Io_clear_bit(channels.MOTORDIR)
		io.Io_write_analog(channels.MOTOR, MOTOR_SPEED)
	} else if (dir == DirDown){
		io.Io_set_bit(channels.MOTORDIR)
		io.Io_write_analog(channels.MOTOR, MOTOR_SPEED)
	}
}

func ElevSetButtonLamp(button ButtonType, floor int, value int) error{
	channel, err := getLampChannel(floor, button)
	if  err != nil{
		return err
	} 
	if value >= 1{
		io.Io_set_bit(channel)
	} else {
		io.Io_clear_bit(channel)
	}
	return nil
}

func ElevSetFloorIndicator(floor int) error{
	if floor < 1 || floor > N_FLOORS {
		return errors.New("Flooor out of bounds error 003")
	}
	floor = floor - 1
	if floor & 0x02 >= 1{
        io.Io_set_bit(channels.LIGHT_FLOOR_IND1);
    } else {
        io.Io_clear_bit(channels.LIGHT_FLOOR_IND1);
    }    

    if floor & 0x01 >= 1{
        io.Io_set_bit(channels.LIGHT_FLOOR_IND2);
    } else {
        io.Io_clear_bit(channels.LIGHT_FLOOR_IND2);
    }
    return nil
}

func ElevSetDoorOpenLamp( value int) {
	if value >= 1{
		io.Io_set_bit(channels.LIGHT_DOOR_OPEN)
	} else {
		io.Io_clear_bit(channels.LIGHT_DOOR_OPEN)
	}
}

func ElevGetButtonSignal(button ButtonType, floor int) (int, error){
	channel, err := getButtonChannel(floor,button)
	if err != nil{
		return 0, err
	}
	return io.Io_read_bit(channel), nil
}

func ElevGetFloorSensorSignal() int{
	if(io.Io_read_bit(channels.SENSOR_FLOOR1) >= 1){
		return 1
	} else if (io.Io_read_bit(channels.SENSOR_FLOOR2) >= 1){
		return 2
	} else if (io.Io_read_bit(channels.SENSOR_FLOOR3) >= 1){
		return 3
	} else if (io.Io_read_bit(channels.SENSOR_FLOOR4) >= 1){
		return 4
	} else {
		return 0
	}
}





