
package main

import (
    "fmt"
    //"./serial"
    "./serial_sean"
    "time" 
    "strconv"
    "strings"
    "os"
    "log"
)



var max_delay float64 = -1  
var min_delay float64 = 1e15 
var delay_sum float64 = 0
var delay_avg float64 = 0
//var delay_cnt float64 = 0 
var delay_cnt int = 0 

var timeout_val = time.Second * 10  
var WAIT_TIME float64 =  5*1000000.0
var WAIT_TIME_RECV float64 =  2*1000000.0

var XBEE_RESOLUTION float64 = 90000.0 // 90 milli second 
var TIME_OUT float64 = 10*1000000.0 
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

func timestamp()( send_msg string) {
    t := time.Now()
    sst := (t.Format(time.StampMicro))

    send_msg = fmt.Sprintf("%02d:%s:%02d:%02d:%02d:%02d:%06d\n",
              t.Day(), strings.ToUpper(sst[0:3]) , t.Year(),
               t.Hour(), t.Minute(), t.Second(), t.Nanosecond()/1000 )
//    send_msg += send_msg 
//    send_msg += send_msg 
    return send_msg 
}


func parse_timestamp( buf string)(tt float64) {
        h , _ := strconv.ParseFloat((buf[12:14]), 64)
	m , _ := strconv.ParseFloat((buf[15:17]),64)
	s , _ := strconv.ParseFloat((buf[18:20]),64)
	ms , _ := strconv.ParseFloat((buf[21:27]),64)
//	fmt.Println(h,m,s,ms) 
	tt =  ((((60*h+m))*60 ) + s ) * 1000000   + ms 
	return tt 
}

func sleeping_func( t float64  ){
    // ttt := time.Second * 1 +  time.Duration(200)*time.Millisecond
    ttt :=   time.Duration(t)*time.Microsecond
    time.Sleep( ttt )
//    fmt.Println("wait time: ",ttt)
}



func time_diff_now( buf string ) (tdiff float64) {
    tnow := time.Now() 
    tnow_ns := tnow.Nanosecond() /1000  
    tnow_flt := float64( (((60* float64(tnow.Hour()) + float64(tnow.Minute()) )*60)+ float64(tnow.Second()) )*1000000 + float64(tnow_ns) ) // A'  
    tt := parse_timestamp(buf) 
    tdiff = tnow_flt -tt
    return tdiff 
}


func SendAndFlash(delay float64) {

    //c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600, ReadTimeout: timeout_val}
//    s, err := serial.OpenPort(c)
    c := "/dev/ttyUSB0" 
    s, err := serial_sean.OpenPort(c)

    for true{
       send_msg := timestamp() 
       s.Write([]byte(send_msg))
       fmt.Printf("%s\n", send_msg)
       sleeping_func( delay) 
       go blink("/sys/class/leds/beaglebone:green:usr1")
       fmt.Println("@@FLASH ========")
       sleeping_func( WAIT_TIME ) 
    }

    err = s.Close()
    if err != nil {
       fmt.Println(err) 
    }
 
}



func sync( round int ) (  float64 ) {


    // c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600,  ReadTimeout: timeout_val}
    // s, err := serial.OpenPort(c)
    c := "/dev/ttyUSB0" 
    s, err := serial_sean.OpenPort(c)

    if err != nil {
       fmt.Println(err) 
    } else {
        fmt.Println("after flush")
    }

    d:=0.0
    counter := 0 
    delay_cnt:=0

    fmt.Printf("Start system's sync\n")


    for ii :=1 ; ii<= round;ii++ {

	fmt.Println("\n---- Testing @iter = ", ii)
	
	counter++
	send_msg := timestamp() 
        fmt.Printf("%s\n", send_msg)

        s.Write([]byte(send_msg))

	if err != nil {
            fmt.Println(err)
	}

	Tag := true
	for Tag {
	    Tag = false

            buffer := make([]byte, 100000)
            cnt , err := s.Read(buffer)
	    cur_pt := 0 
	    if err != nil {
		fmt.Println(err) 
	    }



	    fmt.Println("cnt: ", cnt) 

/*
	    if cnt >= RECV_LOW_BOUND {
		    tdiff := time_diff_now( send_msg ) -  WAIT_TIME_RECV 

		    if tdiff > TIME_OUT || tdiff<0 {
			Tag = false
			break 
		    }else{ 
			delay_sum += tdiff
			delay_cnt += 1
			if max_delay < tdiff{ 
			    max_delay = tdiff 
			}
			if min_delay > tdiff{
			    min_delay = tdiff 
			}
		       Tag = false
		       break
		    }
       	    } 


*/

	    for k:=0 ; k < cnt; k++{
	       if  buffer[k] == '\n'{
		    tdiff := time_diff_now( send_msg ) -  WAIT_TIME_RECV 
		    fmt.Println(" currnt delay = ", tdiff) 
		    fmt.Println(" delay cnt = ", delay_cnt) 

		    if tdiff > TIME_OUT || tdiff < XBEE_RESOLUTION {
			Tag = false
			break 
		    }else{ 
			delay_sum += tdiff
			delay_cnt += 1

			fmt.Println(" current delay sum = ", delay_sum) 
			fmt.Println(" current delay cnt = ", delay_cnt) 
	
			if max_delay < tdiff{ 
			    max_delay = tdiff 
			}
			if min_delay > tdiff{
			    min_delay = tdiff 
			}

		       Tag = false
		       break
		    }
               }else{
		    fmt.Printf("%c",buffer[k]) 
		    cur_pt++
		    Tag = true 
               }
	 } // for k
	fmt.Println();
       } // for Tag

    
    sleeping_func(WAIT_TIME ) 

    } // sync times 

    d = delay_sum/(float64(delay_cnt))

    err = s.Close()
    if err != nil {
       fmt.Println(err) 
    }


    return d
}

func main(){
   round := 5 
   delay_avg := sync(round) 
//    delay_avg := 300000.0
    fmt.Println("RTT+overhead: ", delay_avg ) 
    SendAndFlash(delay_avg + WAIT_TIME_RECV) 
}




