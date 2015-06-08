
package main

import (
    "fmt"
    "./serial"
    "time" 
//    "strconv"
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


func main() {

    c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600,  ReadTimeout: time.Second*3}
    s, err := serial.OpenPort(c)

    if err != nil {
       fmt.Println(err) 
    }  
 
    if err != nil {
        fmt.Println(err) 
    } else{
        fmt.Println("after flush")
    }

    i := 1
   // buf := make([] byte, 1280)  
   // cur_pt := 0 

    d:=0.0
    counter := 0 

    //ttt:=time.Second * (4+1) // 1 for the led
    ttt:=time.Second * (4+1) // 1 for the led
    t_set := 5

    fmt.Printf("Start system's characeterization\n",d);
    for true {
	
	if delay_cnt == 10 { 
	     d = (delay_sum / float64(delay_cnt)/1000000 - float64(t_set)) 
             fmt.Printf("average delay is %f", d);
             fmt.Printf("\nFinish characeterization\n");
	}

	counter++

        t := time.Now()
        sst := (t.Format(time.StampMicro))

        send_msg := fmt.Sprintf("%02d:%s:%02d:%02d:%02d:%02d:%06d\n",
              t.Day(), strings.ToUpper(sst[0:3]) , t.Year(),
              t.Hour(), t.Minute(), t.Second(), t.Nanosecond()/1000 )

        fmt.Printf("%s\n", send_msg)

        _, err = s.Write([]byte(send_msg))
	if err != nil {
		fmt.Println(err)
	}

       // ttt = time.Second * 4 +  time.Duration(d*1000)*time.Millisecond
       ttt = time.Second * 1 +  time.Duration(200)*time.Millisecond
       fmt.Println("wait time: ",ttt)
       time.Sleep( ttt )
       fmt.Println("Flash")
       blink("/sys/class/leds/beaglebone:green:usr1")

       buffer := make([]byte, 1280)
       cnt , er := s.Read(buffer)

       fmt.Printf("%s", string(buffer)) 
       if er != nil {
	     fmt.Println(er) 
       }

       i += cnt 
       fmt.Printf("cnt = ",cnt)
   }

    if err != nil {
       fmt.Println(err) 
    }

    err = s.Close()
    if err != nil {
       fmt.Println(err) 
    }
}




