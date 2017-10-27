package godaemon

import (
	"fmt"
	levelLogger "github.com/cxt90730/LevelLogger-go"
	"os"
	"os/signal"
	"syscall"
)

var dLogger *levelLogger.LevelLogger

func RunDaemon(pidFile string, daemon func(), logger *levelLogger.LevelLogger) error {
	dLogger = logger
	File, err := os.OpenFile(pidFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		printLog(levelLogger.LogError, err)
		return err
	}
	info, _ := File.Stat()
	if info.Size() != 0 {
		printLog(levelLogger.LogError, "pid file is exist")
		return err
	}
	if os.Getppid() != 1 {
		args := append([]string{os.Args[0]}, os.Args[1:]...)
		os.StartProcess(os.Args[0], args, &os.ProcAttr{Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}})
		return err
	}
	File.WriteString(fmt.Sprint(os.Getpid()))
	c := make(chan os.Signal)
	signal.Notify(c)
	printLog(levelLogger.LogInfo, "Daemon is running...")
	go daemon()
	for {
		s := <-c
		fmt.Println(s)
		switch s {
		case os.Interrupt:
			printLog(levelLogger.LogInfo, "RECV SIGINT")
			Exit(File)
		case os.Kill:
			printLog(levelLogger.LogInfo, "RECV SIGKILL")
			Exit(File)
		case syscall.SIGTERM:
			printLog(levelLogger.LogInfo, "RECV SIGTERM")
			Exit(File)
		//case syscall.SIGUSR2:
		//	fmt.Println("SIGUSR2")
		default:
			printLog(levelLogger.LogInfo, s)
			Exit(File)
		}
	}
}

func printLog(level int, v ...interface{}) {
	if dLogger != nil {
		dLogger.PrintLevelLog(level, v...)
	}
}

func Exit(F *os.File) {
	F.Close()
	os.Remove(F.Name())
	printLog(levelLogger.LogInfo, "Daemon exit!")
}
