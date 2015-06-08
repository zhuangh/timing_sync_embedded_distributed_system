// libcgo.h
// wrapper of Linux API and glibc for GO
// zhoufang@ucsd.edu

int getTid();
int setSched(int tid, int policy, int priority);

int createTimerFd(int clockid, int flags);
void setTimerFd(int timerfd, int flags, int v_s, int v_ns, int i_s, int i_ns);
void readTimer(int timerfd, int ifprint);

void bindCPU(int coreid);
void checkCPU();
