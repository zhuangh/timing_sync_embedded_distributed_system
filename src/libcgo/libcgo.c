// libcgo.c
// wrapper of Linux API and glibc for GO
// zhoufang@ucsd.edu

#include "libcgo.h"

#define _GNU_SOURCE 
#include <pthread.h>
#include <stdio.h>     
#include <stdlib.h> 
#include <sched.h>
#include <sys/timerfd.h>
#include <sys/syscall.h>
#include <stdint.h>

int getTid() {
  return syscall(SYS_gettid);
}

// change os scheduler
int setSched(int tid, int policy, int priority) {
  struct sched_param param;
  param.sched_priority = priority;
  return sched_setscheduler(tid, policy, &param);
}

// timer (file descriptor)
int createTimerFd(int clockid, int flags) {
  int timerfd = timerfd_create(CLOCK_MONOTONIC, 0); 
  if(timerfd == -1){
    printf("Create timer fails\n");
  }
  return timerfd;
}

void setTimerFd(int timerfd, int flags, int v_s, int v_ns, int i_s, int i_ns) {
  struct itimerspec new_value;  
  new_value.it_value.tv_sec = v_s;    // initial time
  new_value.it_value.tv_nsec = v_ns;
  new_value.it_interval.tv_sec = i_s; //interval time
  new_value.it_interval.tv_nsec = i_ns;

  if(timerfd_settime(timerfd, 0, &new_value, NULL) == -1) 
    printf("timerfd_settime"); 
}

void readTimer(int timerfd, int ifprint) {
  uint64_t exp; 
  int len; 
  len = read(timerfd, &exp, sizeof(uint64_t)); 
  if(ifprint){
    printf("Timer fd: %d  len: %d  timeout %llu\n", timerfd, len, (unsigned long long) exp);
  }
}




void bindCPU(int coreid) {
  // bind to a processor
  cpu_set_t my_set;
  CPU_ZERO(&my_set);
  CPU_SET(coreid, &my_set);

  // bind process to processor
  if (sched_setaffinity(0, sizeof(my_set), &my_set) < 0) {
    printf("Error: sched_setaffinity\n");
  }
  
}

void checkCPU(){
  cpu_set_t my_set;
  int j;

  if (sched_getaffinity(0, sizeof(my_set), &my_set) <0 ) {
    printf("Error: sched_setaffinity\n");
  } else {
    printf(" CPU: ");
    for (j = 0; j < CPU_SETSIZE; j++)
      if (CPU_ISSET(j, &my_set)) printf(" %d", j);
    printf("\n");
  }

}
