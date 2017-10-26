package godaemon

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func RunDaemon(pid string, daemon func()) {
	File, err := os.OpenFile(pid, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Println(err)
		return
	}
	info, _ := File.Stat()
	if info.Size() != 0 {
		fmt.Println("pid file is exist")
		log.Println("pid file is exist")
		return
	}
	if os.Getppid() != 1 {
		args := append([]string{os.Args[0]}, os.Args[1:]...)
		os.StartProcess(os.Args[0], args, &os.ProcAttr{Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}})
		return
	}
	File.WriteString(fmt.Sprint(os.Getpid()))
	c := make(chan os.Signal, 1)
	signal.Notify(c)
	go daemon()
	for {
		s := <-c
		fmt.Println(s)
		switch s {
		case os.Interrupt:
			fmt.Println("SIGINT")
			Exit(File)
		case os.Kill:
			fmt.Println("SIGKILL")
			Exit(File)
		case syscall.SIGTERM:
			fmt.Println("SIGTERM")
			Exit(File)
		//case syscall.SIGUSR2:
		//	fmt.Println("SIGUSR2")
		default:
			fmt.Println(s)
			Exit(File)
		}
	}

}

func Exit(F *os.File) {
	F.Close()
	os.Remove(F.Name())
	fmt.Println("bye")
}
