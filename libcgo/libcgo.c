// modified from libcgo.c by Hao Zhuang (hao.zhuang@cs.ucsd.edu) 
// wrapper of Linux API and glibc for GO
// The author Zhou Fang, zhoufang@ucsd.edu

#include "libcgo.h"

#define _GNU_SOURCE 
#include <pthread.h>
#include <stdio.h>     
#include <stdlib.h> 
#include <sys/timerfd.h>
#include <sys/syscall.h>
#include <stdint.h>

// timer (file descriptor)
int createTimer(int clockid, int flags) {
    int timerfd = timerfd_create(CLOCK_MONOTONIC, 0); 
    if(timerfd == -1){
	printf("Create timer fails\n");
    }
    return timerfd;
}


void setTimer(int timerfd, int flags, int v_s, int v_ns, int i_s, int i_ns){
    struct itimerspec new_value;  
    new_value.it_value.tv_sec = v_s;    // initial time
    new_value.it_value.tv_nsec = v_ns;
    new_value.it_interval.tv_sec = i_s; //interval time
    new_value.it_interval.tv_nsec = i_ns;

    if(timerfd_settime(timerfd, 0, &new_value, NULL) == -1) {
	printf("timerfd_settime"); 
    }
}

void readTimer(int timerfd, int ifprint) {
    uint64_t exp; 
    int len; 
    len = read(timerfd, &exp, sizeof(uint64_t)); 
    if(ifprint){
	printf("Timer fd: %d  len: %d  timeout %llu\n", timerfd, len, (unsigned long long) exp);
    }
}


