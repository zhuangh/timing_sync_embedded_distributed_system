// modified from libcgo.c by Hao Zhuang (hao.zhuang@cs.ucsd.edu) 
// wrapper of Linux API and glibc for GO
// The author Zhou Fang, zhoufang@ucsd.edu


package libcgo

// #cgo LDFLAGS: ../ctimer/libcgo.so 
// #include <stdio.h>     
// #include <stdlib.h> 
// #include <sched.h>
// #include "libcgo.h"
import "C"     

func CreateTimerFd(clockid int, flags int) int{
	timerfd:= int( C.createTimerFd(C.int(clockid), C.int(flags)) )
	return timerfd
}

func SetTimerFd(timerfd int, flags int, v_s int, v_ns int, i_s int, i_ns int) {
	C.setTimerFd(C.int(timerfd), C.int(flags), C.int(v_s), C.int(v_ns), C.int(i_s), C.int(i_ns)) 
}

func ReadTimer(timerfd int, ifprint int) {
        C.readTimer(C.int(timerfd), C.int(ifprint))
}



