/*
Master package is the main package that keeps everything running as it should be. It's responsible for starting, stopping and deleting processes. It also will keep an eye on the Watcher in case a process dies so it can restart it again.

RemoteMaster is responsible for exporting the main pmgo operations as HTTP requests. If you want to start a Remote Server, run:

- remoteServer := master.StartRemoteMasterServer(dsn, configFile)

It will start a remote master and return the instance.

To make remote requests, use the Remote Client by instantiating using:

- remoteClient, err := master.StartRemoteClient(dsn, timeout)

It will start the remote client and return the instance so you can use to initiate requests, such as:

- remoteClient.StartGoBin(sourcePath, name, keepAlive, args)
*/

package master

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sync"

	"time"

	"github.com/struCoder/pmgo/lib/preparable"
	"github.com/struCoder/pmgo/lib/process"
	"github.com/struCoder/pmgo/lib/utils"
	"github.com/struCoder/pmgo/lib/watcher"

	log "github.com/sirupsen/logrus"
)

// Master is the main module that keeps everything in place and execute
// the necessary actions to keep the process running as they should be.
type Master struct {
	sync.Mutex

	SysFolder string           // SysFolder is the main pmgo folder where the necessary config files will be stored.
	PidFile   string           // PidFille is the pmgo pid file path.
	OutFile   string           // OutFile is the pmgo output log file path.
	ErrFile   string           // ErrFile is the pmgo err log file path.
	Watcher   *watcher.Watcher // Watcher is a watcher instance.

	Procs      map[string]process.ProcContainer // Procs is a map containing all procs started on pmgo.
	configFile string                           // custom config filename
}

// DecodableMaster is a struct that the config toml file will decode to.
// It is needed because toml decoder doesn't decode to interfaces, so the
// Procs map can't be decoded as long as we use the ProcContainer interface
type DecodableMaster struct {
	SysFolder string
	PidFile   string
	OutFile   string
	ErrFile   string

	Watcher *watcher.Watcher

	Procs map[string]*process.Proc
}

// InitMaster will start a master instance with configFile.
// It returns a Master instance.
func InitMaster(configFile string) *Master {
	watcher := watcher.InitWatcher()
	decodableMaster := &DecodableMaster{}
	decodableMaster.Procs = make(map[string]*process.Proc)

	err := utils.SafeReadTomlFile(configFile, decodableMaster)
	if err != nil {
		panic(err)
	}

	procs := make(map[string]process.ProcContainer)
	for k, v := range decodableMaster.Procs {
		procs[k] = v
	}
	// We need this hack because toml decoder doesn't decode to interfaces
	master := &Master{
		SysFolder:  decodableMaster.SysFolder,
		PidFile:    decodableMaster.PidFile,
		OutFile:    decodableMaster.OutFile,
		ErrFile:    decodableMaster.ErrFile,
		Watcher:    decodableMaster.Watcher,
		Procs:      procs,
		configFile: configFile,
	}

	if master.SysFolder == "" {
		os.MkdirAll(path.Dir(configFile), 0777)
		master.SysFolder = path.Dir(configFile) + "/"
	}
	master.Watcher = watcher
	err = master.Revive()
	if err != nil {
		log.Error(err)
	} else {
		log.Infof("All procs revived...")
	}
	go master.WatchProcs()
	// go master.SaveProcsLoop()
	go master.UpdateStatus()
	return master
}

// ProcInfo will return process detail info with map
func (master *Master) ProcInfo(proName string) map[string]string {
	proc := master.Procs[proName]
	procDetailInfo := make(map[string]string)
	if proc != nil {
		procStatus := proc.GetStatus()
		procDetailInfo["pid"] = fmt.Sprintf("%d", proc.GetPid())
		procDetailInfo["outFile"] = proc.GetOutFile()
		procDetailInfo["pidFile"] = proc.GetPidFile()
		procDetailInfo["errorFile"] = proc.GetErrFile()
		procDetailInfo["path"] = proc.GetPath()
		procDetailInfo["name"] = proc.GetName()
		procDetailInfo["uptime"] = procStatus.Uptime
		procDetailInfo["status"] = procStatus.Status
		procDetailInfo["restart"] = fmt.Sprintf("%d", procStatus.Restarts)
	}

	return procDetailInfo
}

// WatchProcs will keep the procs running forever.
func (master *Master) WatchProcs() {
	for proc := range master.Watcher.RestartProc() {
		if !proc.ShouldKeepAlive() {
			master.Lock()
			master.updateStatus(proc)
			master.Unlock()
			log.Infof("Proc %s does not have keep alive set. Will not be restarted.", proc.Identifier())
			continue
		}
		log.Infof("Restarting proc %s.", proc.Identifier())
		if proc.IsAlive() {
			log.Warnf("Proc %s was supposed to be dead, but it is alive.", proc.Identifier())
		}
		master.Lock()
		err := master.restart(proc)
		master.Unlock()
		if err != nil {
			log.Warnf("Could not restart process %s due to %s.", proc.Identifier(), err)
		}
	}
}

// Prepare will compile the source code into a binary and return a preparable
// ready to be executed.
func (master *Master) Prepare(sourcePath string, name string, language string, keepAlive bool, args []string, cwd string, fromBin bool) (preparable.ProcPreparable, []byte, error) {
	var procPreparable preparable.ProcPreparable
	if fromBin {
		procPreparable = &preparable.BinaryPreparable{
			Name:       name,
			SourcePath: sourcePath,
			SysFolder:  master.SysFolder,
			Language:   language,
			KeepAlive:  keepAlive,
			Args:       args,
			Cwd:        cwd,
		}
	} else {
		procPreparable = &preparable.Preparable{
			Name:       name,
			SourcePath: sourcePath,
			SysFolder:  master.SysFolder,
			Language:   language,
			KeepAlive:  keepAlive,
			Args:       args,
			Cwd:        cwd,
		}
	}

	output, err := procPreparable.PrepareBin()
	return procPreparable, output, err
}

// RunPreparable will run procPreparable and add it to the watch list in case everything goes well.
func (master *Master) RunPreparable(procPreparable preparable.ProcPreparable) error {
	master.Lock()
	defer master.Unlock()
	if _, ok := master.Procs[procPreparable.Identifier()]; ok {
		log.Warnf("Proc %s already exist.", procPreparable.Identifier())
		return errors.New("Trying to start a process that already exist.")
	}
	proc, err := procPreparable.Start()
	if err != nil {
		return err
	}
	master.Procs[proc.Identifier()] = proc
	master.saveProcsWrapper()
	master.Watcher.AddProcWatcher(proc)
	proc.SetStatus("running")
	return nil
}

// ListProcs will return a list of all procs.
func (master *Master) ListProcs() []process.ProcContainer {
	procsList := []process.ProcContainer{}
	for _, v := range master.Procs {
		procsList = append(procsList, v)
	}
	return procsList
}

// RestartProcess will restart a process.
func (master *Master) RestartProcess(name string) error {
	if proc, ok := master.Procs[name]; ok {
		master.Lock()
		err := master.restart(proc)
		master.Unlock()
		return err
	}
	return errors.New("unknow process.")

}

// StartProcess will a start a process.
func (master *Master) StartProcess(name string) error {
	master.Lock()
	defer master.Unlock()
	if proc, ok := master.Procs[name]; ok {
		return master.start(proc)
	}
	return errors.New("Unknown process.")
}

// StopProcess will stop a process with the given name.
func (master *Master) StopProcess(name string) error {
	master.Lock()
	defer master.Unlock()
	if proc, ok := master.Procs[name]; ok {
		return master.stop(proc)
	}
	return errors.New("Unknown process.")
}

// DeleteProcess will delete a process and all its files and childs forever.
func (master *Master) DeleteProcess(name string) error {
	master.Lock()
	defer master.Unlock()
	log.Infof("Trying to delete proc %s", name)
	if proc, ok := master.Procs[name]; ok {
		err := master.stop(proc)
		if err != nil {
			return err
		}
		delete(master.Procs, name)
		err = master.delete(proc)
		if err != nil {
			return err
		}
		log.Infof("Successfully deleted proc %s", name)
	}
	return nil
}

// Revive will revive all procs listed on ListProcs. This should ONLY be called
// during Master startup.
func (master *Master) Revive() error {
	master.Lock()
	defer master.Unlock()
	procs := master.ListProcs()
	log.Info("Reviving all processes")
	for id := range procs {
		proc := procs[id]
		if !proc.ShouldKeepAlive() {
			log.Infof("Proc %s does not have KeepAlive set. Will not revive it.", proc.Identifier())
			continue
		}
		log.Infof("Reviving proc %s", proc.Identifier())
		err := master.start(proc)
		if err != nil {
			return fmt.Errorf("Failed to revive proc %s due to %s", proc.Identifier(), err)
		}
	}
	return nil
}

// NOT thread safe method. Lock should be acquire before calling it.
func (master *Master) start(proc process.ProcContainer) error {
	if !proc.IsAlive() {
		err := proc.Start()
		if err != nil {
			return err
		}
		master.Watcher.AddProcWatcher(proc)
		proc.SetStatus("running")
		proc.SetUptime()
		master.saveProcsWrapper()
	}
	return nil
}

func (master *Master) delete(proc process.ProcContainer) error {
	return proc.Delete()
}

// NOT thread safe method. Lock should be acquire before calling it.
func (master *Master) stop(proc process.ProcContainer) error {
	if proc.IsAlive() {
		waitStop := master.Watcher.StopWatcher(proc.Identifier())
		err := proc.GracefullyStop()
		if err != nil {
			return err
		}
		if waitStop != nil {
			<-waitStop
			proc.NotifyStopped()
			proc.SetStatus("stopped")
			proc.SetUptime()
		}
		log.Infof("Proc %s successfully stopped.", proc.Identifier())
		master.saveProcsWrapper()
	}
	return nil
}

// UpdateStatus will update a process status every 30s.
func (master *Master) UpdateStatus() {
	for {
		master.Lock()
		procs := master.ListProcs()
		for id := range procs {
			proc := procs[id]
			master.updateStatus(proc)
		}
		master.Unlock()
		time.Sleep(30 * time.Second)
	}
}

func (master *Master) updateStatus(proc process.ProcContainer) {
	if proc.IsAlive() {
		proc.SetStatus("running")
	} else {
		proc.NotifyStopped()
		proc.SetStatus("stopped")
	}
}

// NOT thread safe method. Lock should be acquire before calling it.
func (master *Master) restart(proc process.ProcContainer) error {
	// restat count +1
	proc.AddRestart()
	err := master.stop(proc)
	if err != nil {
		return err
	}
	err = master.start(proc)
	master.saveProcsWrapper()
	return err
}

// SaveProcsLoop will loop forever to save the list of procs onto the proc file.
// func (master *Master) SaveProcsLoop() {
// 	for {
// 		log.Infof("Saving list of procs.")
// 		master.Lock()
// 		master.saveProcsWrapper()
// 		master.Unlock()
// 		time.Sleep(5 * time.Minute)
// 	}
// }

// Stop will stop pmgo and save all of its running procs.
func (master *Master) Stop() error {
	log.Info("Stopping pmgo...")
	procs := master.ListProcs()
	for id := range procs {
		proc := procs[id]
		log.Infof("Stopping proc %s", proc.Identifier())
		master.stop(proc)
	}
	log.Info("Saving and returning list of procs.")
	return master.saveProcsWrapper()
}

// SaveProcs will save a list of procs onto a file inside configPath.
// Returns an error in case there's any.
func (master *Master) SaveProcs() error {
	master.Lock()
	defer master.Unlock()
	return master.saveProcsWrapper()
}

// NOT Thread Safe. Lock should be acquired before calling it.
func (master *Master) saveProcsWrapper() error {
	configPath := master.getConfigPath()
	return utils.SafeWriteTomlFile(master, configPath)
}

func (master *Master) getConfigPath() string {
	return master.configFile
}

// IsExistProc current proc whether exist
func (master *Master) IsExistProc(procName string) (bool, error) {
	if proc, ok := master.Procs[procName]; ok {
		if !proc.IsAlive() {
			err := master.RestartProcess(procName)
			if err != nil {
				return false, err
			}
		}
		return true, nil
	}
	return false, nil
}
