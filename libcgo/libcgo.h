// modified from libcgo.c by Hao Zhuang (hao.zhuang@cs.ucsd.edu) 
// wrapper of Linux API and glibc for GO
// The author Zhou Fang, zhoufang@ucsd.edu


int createTimer(int clockid, int flags);
void setTimer(int timerfd, int flags, int v_s, int v_ns, int i_s, int i_ns);
void readTimer(int timerfd, int ifprint);



