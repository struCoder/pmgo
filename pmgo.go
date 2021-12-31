/*
PMGO is a lightweight process manager written in Golang for Golang applications. It helps you keep all of your applications alive forever, if you want to. You can also reload, start, stop, delete and query status on the fly.

PMGO also provide a way to start a process by compiling a Golang project source code.

The main PMGO module is the Master module, it's the glue that keep everything running as it should be.

If you need to use the remote version of PMGO, take a look at RemoteMaster on Master package.

To use the remote version of PMGO, use:

- remoteServer := master.StartRemoteMasterServer(dsn, configFile)

It will start a remote master and return the instance.

To make remote requests, use the Remote Client by instantiating using:

- remoteClient, err := master.StartRemoteClient(dsn, timeout)

It will start the remote client and return the instance so you can use to initiate requests, such as:

- remoteClient.StartGoBin(sourcePath, name, keepAlive, args)
*/
package main

import (
	"sync"

	"github.com/struCoder/pmgo/lib/cli"
	"github.com/struCoder/pmgo/lib/master"
	"github.com/struCoder/pmgo/lib/utils"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/sevlyar/go-daemon"

	"os"
	"os/signal"
	"path"
	"syscall"

	"fmt"

	"time"

	log "github.com/sirupsen/logrus"
)

var (
	app             = kingpin.New("pmgo", "Aguia Process Manager.")
	serveConfigFile = app.Flag("config-file", "Config file location").String()
	dns             = app.Flag("dns", "TCP Dns host.").Default(":9876").String()
	timeout         = 30 * time.Second

	serveStop = app.Command("kill", "Kill daemon pmgo.")

	serve = app.Command("serve", "Create pmgo daemon.")

	resurrect = app.Command("resurrect", "Resurrect all previously save processes.")

	start           = app.Command("start", "start and daemonize an app.")
	startSourcePath = start.Arg("start go file", "go file.").Required().String()
	startName       = start.Arg("name", "Process name.").Required().String()
	binFile         = start.Arg("binary", "compiled golang file").Bool()
	startKeepAlive  = true
	startCwd        = start.Flag("cwd", "Change the working directory").String()
	startArgs       = start.Flag("args", "External args.").Strings()

	restart     = app.Command("restart", "Restart a process.")
	restartName = restart.Arg("name", "Process name.").Required().String()

	stop     = app.Command("stop", "Stop a process.")
	stopName = stop.Arg("name", "Process name.").Required().String()

	delete     = app.Command("delete", "Delete a process.")
	deleteName = delete.Arg("name", "Process name.").Required().String()

	save = app.Command("save", "Save a list of processes onto a file.")

	status = app.Command("list", "Get pmgo list.")

	version        = app.Command("version", "get version")
	currentVersion = "0.6.1"

	info     = app.Command("info", "Describe importance parameters of a process id")
	infoName = info.Arg("name", "process name").Required().String()
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case serveStop.FullCommand():
		checkRemoteMasterServer()
		cli := cli.InitCli(*dns, timeout)
		cli.DeleteAllProcess()
		stopRemoteMasterServer()
	case serve.FullCommand():
		log.Warn("Server will auto start and this command will be delete")
		startRemoteMasterServer()
	case resurrect.FullCommand():
		fmt.Println("This feature will not support. sorry")
	case start.FullCommand():
		checkRemoteMasterServer()
		cli := cli.InitCli(*dns, timeout)
		cli.StartGoBin(*startSourcePath, *startName, startKeepAlive, *startArgs, *startCwd, *binFile)
		cli.Status()
	case restart.FullCommand():
		checkRemoteMasterServer()
		cli := cli.InitCli(*dns, timeout)
		cli.RestartProcess(*restartName)
		cli.Status()
	case stop.FullCommand():
		checkRemoteMasterServer()
		cli := cli.InitCli(*dns, timeout)
		cli.StopProcess(*stopName)
		cli.Status()
	case delete.FullCommand():
		checkRemoteMasterServer()
		cli := cli.InitCli(*dns, timeout)
		cli.DeleteProcess(*deleteName)
	case save.FullCommand():
		cli := cli.InitCli(*dns, timeout)
		cli.Save()
	case status.FullCommand():
		checkRemoteMasterServer()
		cli := cli.InitCli(*dns, timeout)
		cli.Status()
	case version.FullCommand():
		fmt.Println(currentVersion)
	case info.FullCommand():
		checkRemoteMasterServer()
		cli := cli.InitCli(*dns, timeout)
		cli.ProcInfo(*infoName)
	}
}

func isDaemonRunning(ctx *daemon.Context) (bool, *os.Process, error) {
	d, err := ctx.Search()

	if err != nil {
		return false, d, err
	}

	if err := d.Signal(syscall.Signal(0)); err != nil {
		return false, d, err
	}

	return true, d, nil
}

var ctx *daemon.Context

func getCtx() *daemon.Context {
	if ctx != nil {
		return ctx
	}
	if *serveConfigFile == "" {
		configFile := os.Getenv("PMGO_CONFIG_FILE")
		if configFile != "" {
			*serveConfigFile = configFile
		} else {
			folderPath := os.Getenv("HOME")
			*serveConfigFile = folderPath + "/.pmgo/config.toml"
		}
	}

	dir := path.Dir(*serveConfigFile)
	os.MkdirAll(dir, 0755)

	decodableMaster := master.DecodableMaster{}
	err := utils.SafeReadTomlFile(*serveConfigFile, &decodableMaster)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	if decodableMaster.SysFolder != "" {
		dir = decodableMaster.SysFolder
	}

	if decodableMaster.PidFile == "" {
		decodableMaster.PidFile = path.Join(dir, "main.pid")
	} else {
		os.MkdirAll(path.Dir(decodableMaster.PidFile), 0755)
	}

	if decodableMaster.OutFile == "" {
		decodableMaster.OutFile = path.Join(dir, "main.log")
	} else {
		os.MkdirAll(path.Dir(decodableMaster.OutFile), 0755)
	}

	ctx = &daemon.Context{
		PidFileName: decodableMaster.PidFile,
		PidFilePerm: 0644,
		LogFileName: decodableMaster.OutFile,
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
	}
	return ctx
}

// if RemoteMasterServer not running, just run
func checkRemoteMasterServer() {
	ctx := getCtx()
	if ok, _, _ := isDaemonRunning(ctx); !ok {
		startRemoteMasterServer()
	}
}

var waitedForSignal os.Signal

func waitForChildSignal(wg *sync.WaitGroup) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGUSR1, syscall.SIGUSR2)
	wg.Add(1)
	go func() {
		waitedForSignal = <-signalChan
		wg.Done()
	}()
}

func kill(pid int, signal os.Signal) error {
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	defer p.Release()
	return p.Signal(signal)
}

func startRemoteMasterServer() {
	var wg sync.WaitGroup
	waitForChildSignal(&wg)
	ctx := getCtx()
	if ok, _, _ := isDaemonRunning(ctx); ok {
		log.Info("pmgo daemon is already running.")
		return
	}

	d, err := ctx.Reborn()
	if err != nil {
		log.Fatalf("Failed to reborn daemon due to %+v.", err)
	}

	if d != nil {
		wg.Wait()
		if waitedForSignal == syscall.SIGUSR1 {
			log.Info("daemon started")
			return
		}
	} else {
		kill(os.Getpid(), syscall.SIGUSR1)
		wg.Wait()
		defer ctx.Release()
	}

	log.Info("Starting remote master server...")
	remoteMaster := master.StartRemoteMasterServer(*dns, *serveConfigFile)

	// send signal to parent's process to kill goroutine
	kill(os.Getppid(), syscall.SIGUSR1)
	sigsKill := make(chan os.Signal, 1)
	signal.Notify(sigsKill,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	<-sigsKill
	log.Info("Received signal to stop...")
	err = remoteMaster.Stop()
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func stopRemoteMasterServer() {
	log.Info("pmgo stopping...")
	ctx := getCtx()
	if ok, p, _ := isDaemonRunning(ctx); ok {
		if err := p.Signal(syscall.Signal(syscall.SIGQUIT)); err != nil {
			log.Fatalf("Failed to kill daemon %v", err)
		}
	} else {
		ctx.Release()
		log.Info("instance is not running.")
	}
	log.Info("pmgo daemon terminated")
}
