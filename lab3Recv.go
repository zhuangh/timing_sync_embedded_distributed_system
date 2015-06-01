package main

import (
    "fmt"
    "./serial"
    "time" 
    //"strings"
    "os"
    "log"
)

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
//    blink("/sys/class/leds/beaglebone:green:usr1")
    c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600} // , ReadTimeout: time.Microsecond*2000000}
    s, err := serial.OpenPort(c)
    if err != nil {
       fmt.Println(err) 
    } else{
        fmt.Println("after flush")
    }
    i := 1
    for true {
       buffer := make([]byte, 1280)
       cnt , er := s.Read(buffer)
       fmt.Printf("%s", string(buffer)) 
       if er != nil {
	     fmt.Println(er) 
       }
       i += cnt 

       for k:=0; k < cnt ; k++{
            if buffer[k] == '\n' {
		send_messageAB := "ack\n\n\n"
		fmt.Printf("\nsend to A with B: %s ",  send_messageAB) 
		 _, err = s.Write([]byte( send_messageAB ))
		 blink("/sys/class/leds/beaglebone:green:usr1")
                 fmt.Printf("\nFlash@@\n") 
                 fmt.Printf("\n================================================\n") 
//		 fmt.Println("reset buffer")
		 break 
	    }
        }
    }
    

    fmt.Printf("\nReceive %d bytes\n", i) 

    if err != nil {
       fmt.Println(err) 
    }

    err = s.Close()
    if err != nil {
       fmt.Println(err) 
    }
}

/*
 		 tnow := time.Now() 
		 sst := tnow.Format(time.StampMicro)
                 fmt.Println("not rock back ", buf)

		 // B
		 tnow_msg := fmt.Sprintf("%02d:%s:%02d:%02d:%02d:%02d:%06d",
                        tnow.Day(), strings.ToUpper(sst[0:3]) , tnow.Year(),
                        tnow.Hour(), tnow.Minute(), tnow.Second(), tnow.Nanosecond()/1000 ) 

		fmt.Printf("\nsend buf and tnow: %s",tnow_msg) 
		send_messageAB := string(buf)+":"+string(tnow_msg)+"\n"
		fmt.Printf("\nsend to A with B: %s ",  send_messageAB) 
		 _, err = s.Write([]byte( send_messageAB ))
*/
