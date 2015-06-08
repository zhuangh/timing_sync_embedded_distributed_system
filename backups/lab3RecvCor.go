
package main

import (
    "fmt"
    "./serial"
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
    //    time.Sleep(1000 * time.Millisecond)
}


func main() {
    // MAX_BIT:=10000000

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
    buf := make([] byte, 1280)  
    cur_pt := 0 

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

       ttt = time.Second * 4 +  time.Duration(d*1000)*time.Millisecond
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
    
     Tag := true

     for Tag {
	Tag = false
       for k:=0 ; k < cnt; k++{
           if buffer[k] == '\n'{
                tnow := time.Now() 
                tnow_ns := tnow.Nanosecond() /1000  

	        tnow_flt := float64( (((60* float64(tnow.Hour()) + float64(tnow.Minute()) )*60)+ float64(tnow.Second()) )*1000000 + float64(tnow_ns) ) // A'  

                fmt.Println(tnow.Hour() ) 
                fmt.Println(tnow.Minute() ) 
                fmt.Println(tnow.Second() ) 
		fmt.Println(float64(tnow_ns)) 
                fmt.Println(tnow_ns) 

//		fmt.Println( (((60*tnow.Hour() + tnow.Minute())*60) +tnow.Second() )*1000000 )  

		h , _ := strconv.ParseFloat(string(buf[12:14]), 64)
		m , _ := strconv.ParseFloat(string(buf[15:17]),64)
		s , _ := strconv.ParseFloat(string(buf[18:20]),64)
		ms , _ := strconv.ParseFloat(string(buf[21:27]),64)

		hb , _ := strconv.ParseFloat(string(buf[28:30]), 64)
		mb , _ := strconv.ParseFloat(string(buf[31:33]),64)
		sb , _ := strconv.ParseFloat(string(buf[34:36]),64)
		msb , _ := strconv.ParseFloat(string(buf[37:43]),64)

                fmt.Println(h,m,s,ms) 

                tt :=  ((((60*h+m))*60 ) + s ) * 1000000   + ms 

                ttb :=  ((((60*hb+mb))*60 ) + sb ) * 1000000   + msb 
		fmt.Printf("\n from B %d\n", ttb)


                tdiff := tnow_flt -tt

		delay_sum += tdiff
		delay_cnt += 1

                if max_delay < tdiff{ 
                   max_delay = tdiff 
                }
                if min_delay > tdiff{
                   min_delay = tdiff 
                }

                fmt.Println("roger!\n", string(buf[0:27])) 
  		sst := tnow.Format(time.StampMicro)  

		//A'  
                tnow_msg := fmt.Sprintf("%02d:%s:%02d:%02d:%02d:%02d:%06d\r\n", 
                        tnow.Day(), strings.ToUpper(sst[0:3]) , tnow.Year(),
                        tnow.Hour(), tnow.Minute(), tnow.Second(), tnow.Nanosecond()/1000 ) 


                fmt.Println( tnow_msg, tdiff)

                fmt.Println( tnow_flt ," ", tt , "= " ,tdiff ) 
                fmt.Println( "The xbee delay = " ,tdiff ) 
                cur_pt = 0 
                fmt.Println(" ------------- Clear buffer --------------- ")

		Tag = false
		break

           }else{
		buf[cur_pt] =  (buffer[k]) 
	        cur_pt++
		Tag = true 
		 //  fmt.Println("current buffer ", buffer)
           }
       }
     }
   }

    fmt.Printf("\nReceive %d bytes\n", i) 

    delay_avg = delay_sum / float64(delay_cnt)/2 
    fmt.Printf("\n averay delay = %f, max_delay = %f , min_delay = %f  @ delay_cnt = %d \n", delay_avg, max_delay/2, min_delay/2, delay_cnt ) 

    if err != nil {
       fmt.Println(err) 
    }

    err = s.Close()
    if err != nil {
       fmt.Println(err) 
    }
}




