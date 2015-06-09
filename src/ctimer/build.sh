
#change GOPATH, GOBIN to your path
export GOPATH=$GOPATH:/home/debian/final/:/home/zhuang/Dropbox/15sp_cse237b/final/ 
#export GOROOT=$GOROOT:/home/debian/final/:/home/zhuang/Dropbox/15sp_cse237b/final/ 
export PATH=$PATH:$GOPATH
export GOBIN=$GOBIN:/home/debian/final/bin
export GOMAXPROCS=10

# dynamic link
gcc -shared -fPIC ../libcgo/libcgo.c -o libcgo.so
#gcc -shared -fPIC ./libcgo/libcgo.c -o ./libcgo/libcgo.so

# static link
# gcc ../libcgo/libcgo.c -c
# ar rv libcgo.a libcgo.o > /dev/null 2>&1

go build ../libcgo/libcgo.go
go build send_ctimer.go
go build recv_ctimer.go
go build dist_send_ctimer.go
go build dist_recv_ctimer.go
#go install main

