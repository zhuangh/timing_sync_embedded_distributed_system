package serial_sean

import (
	"fmt"
	"os"
//	"sync"
	"syscall"
	"unsafe"
    )




func OpenPort(name string) (f *os.File, err error) {

    fmt.Println("use Sean's serial_sean.go") 
    f, err = os.OpenFile(name, syscall.O_RDWR|syscall.O_NOCTTY, 0666)
	if err != nil {
	    return nil, err
	}

    defer func() {
	if err != nil && f != nil {
	    f.Close()
	}
    }()

    fd := f.Fd()

    //	Set serial port 'name' to 115200/8/N/1 in RAW mode (i.e. no pre-process of received data
    // and pay special attention to Cc field, this tells the serial port to not return until at
    // at least 22 bytes have been read. This is a tunable parameter they may help in Lab 3
    t := syscall.Termios{
            Iflag:  syscall.IGNPAR,
	    Cflag:  syscall.CS8 | syscall.CREAD | syscall.CLOCAL | syscall.B9600,
	    Cc:     [32]uint8{syscall.VMIN: 3},
	    Ispeed: syscall.B9600,
	    Ospeed: syscall.B9600,
    }


    // Syscall to apply these parameters
    _, _, errno := syscall.Syscall6(
				syscall.SYS_IOCTL,
				uintptr(fd),
				uintptr(syscall.TCSETS),
				uintptr(unsafe.Pointer(&t)),
				0,
				0,
				0,
			       )

    if errno != 0 {
	return nil, errno
    }

    return f, nil
}


