// libcgo.go
// wrapper of Linux API and glibc for GO
// zhoufang@ucsd.edu

package libcgo

// #cgo LDFLAGS: ../ctimer/libcgo.so 
// #include <stdio.h>     
// #include <stdlib.h> 
// #include <sched.h>
// #include "libcgo.h"
import "C"     

func GetTid() int {
        return int(C.getTid())
}

func GetSched(tid int) int {
        tid2 := C.__pid_t(tid)
	sched := int(C.sched_getscheduler(tid2))
	return sched
}

func SetSched(tid int, policy int, priority int) int {
	return int( C.setSched(C.int(tid), C.int(policy), C.int(priority)) )
}

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

// bind current OS thread to a CPU core
// it is useful when comparing RT thread with normal thread
// because the comparison is evident when threads are on the same CPU
func BindCPU(coreid int) {
	C.bindCPU(C.int(coreid))
}

func CheckCPU() {
	C.checkCPU()
}
