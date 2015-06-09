package main

import (
    "fmt"
    //"./serial"
    "../serial_sean"
    "time" 
    //"strings"
    "os"
    "log"
    lc "libcgo"
)

var WAIT_TIME_RECV float64 =  2*1000000.0



// var RECV_LOW_BOUND int = 27  

func setLed(led string, value []byte) {
	file, err := os.OpenFile(led+"/brightness", os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	file.Write(value)
	file.Close()
}
 
func blink(led string) {
	setLed(led, []byte("1"))
	time.Sleep(1000 * time.Millisecond)
	setLed(led, []byte("0"))
}

func sleeping_func( t float64  ){
    // ttt := time.Second * 1 +  time.Duration(200)*time.Millisecond
    ttt :=   time.Duration(t)*time.Microsecond
    time.Sleep( ttt )
//    fmt.Println("wait time: ",ttt)
}

func dummy(){
    i:=1
    for true{
	i++
        sleeping_func( 100000) // 10 seconds  
    }
    
}


func main() {
    for i:=1;i<=1000;i++{
	go dummy()
    }
//    blink("/sys/class/leds/beaglebone:green:usr1")
    // c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600 , ReadTimeout: time.Microsecond*10000000}
//    c := openPort{Name: "/dev/ttyUSB0", Baud: 9600 , ReadTimeout: time.Microsecond*10000000}
    c := "/dev/ttyUSB0" 
    // s, err := serial.OpenPort(c)
    s, err := serial_sean.OpenPort(c)
    if err != nil {
       fmt.Println(err) 
    } else{
        fmt.Println("after flush")
    }

    // a linux timer 
    period := int(WAIT_TIME_RECV)   
    sstart :=  int(period/1000000) // int( current_time) / 1000000000
    nsstart := 0 // int( current_time) % 1000000000

    speriod := 0 // period / 1000000000  
    nsperiod := 0 // period % 1000000000  
    fmt.Println("sstart = ", sstart)


    fmt.Println("speriod = ", speriod)

    timer_c := lc.CreateTimerFd(1,0) 
    // lc.SetTimerFd( timer_c, 0, sstart, nsstart, speriod, nsperiod) 

    for true {

       buffer := make([]byte, 100000)
       cnt , er := s.Read(buffer)

       fmt.Printf("\nword length : %d   -> ", cnt) 
       fmt.Printf("%s", string(buffer)) 

       if er != nil {
	     fmt.Println(er) 
       }

       for k:=0; k < cnt ; k++{
            if buffer[k] == '\n' {
		send_messageAB := "123456789012345678901234567\n\n\n"
//		fmt.Printf("\nsend to A with B: %s ",  send_messageAB) 


		lc.SetTimerFd( timer_c, 0, sstart, nsstart, speriod, nsperiod) 
	        lc.ReadTimer(timer_c, 0) 
	//	sleeping_func(WAIT_TIME_RECV) 
                s.Write([]byte( send_messageAB ))
                go blink("/sys/class/leds/beaglebone:green:usr1")
		fmt.Printf("Flash@@") 
		fmt.Printf("================================================\n") 
		break 
	    } 
        }

    }
    

//    fmt.Printf("\nReceive %d bytes\n", i) 

    if err != nil {
       fmt.Println(err) 
    }

    err = s.Close()
    if err != nil {
       fmt.Println(err) 
    }
}

