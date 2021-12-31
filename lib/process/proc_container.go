package process

import (
	"errors"
	"os"
	"strconv"
	"syscall"

	"github.com/struCoder/pmgo/lib/utils"
)

// ProcContainer is a interface that about proc
type ProcContainer interface {
	Start() error
	ForceStop() error
	GracefullyStop() error
	Restart() error
	Delete() error
	IsAlive() bool
	Identifier() string
	ShouldKeepAlive() bool
	AddRestart()
	NotifyStopped()
	SetStatus(status string)
	SetUptime()
	SetSysInfo()
	GetPid() int
	GetStatus() *ProcStatus
	Watch() (*os.ProcessState, error)
	release()
	GetOutFile() string
	GetPidFile() string
	GetPath() string
	GetErrFile() string
	GetName() string
}

// Proc is a os.Process wrapper with Status and more info that will be used on Master to maintain
// the process health.
type Proc struct {
	Name      string
	Cmd       string
	Cwd       string
	Args      []string
	Path      string
	Pidfile   string
	Outfile   string
	Errfile   string
	KeepAlive bool
	Pid       int
	Status    *ProcStatus
	process   *os.Process
}

// Start will execute the command Cmd that should run the process. It will also create an out, err and pidfile
// in case they do not exist yet.
// Returns an error in case there's any.
func (proc *Proc) Start() error {
	outFile, err := utils.GetFile(proc.Outfile)
	if err != nil {
		return err
	}
	errFile, err := utils.GetFile(proc.Errfile)
	if err != nil {
		return err
	}
	if proc.Cwd == "" {
		proc.Cwd, _ = os.Getwd()
	}
	procAtr := &os.ProcAttr{
		Dir: proc.Cwd,
		Env: os.Environ(),
		Files: []*os.File{
			os.Stdin,
			outFile,
			errFile,
		},
	}
	args := append([]string{proc.Name}, proc.Args...)
	process, err := os.StartProcess(proc.Cmd, args, procAtr)
	if err != nil {
		return err
	}
	proc.process = process
	proc.Pid = proc.process.Pid
	err = utils.WriteFile(proc.Pidfile, []byte(strconv.Itoa(proc.process.Pid)))
	if err != nil {
		return err
	}
	proc.Status.InitUptime()
	proc.Status.SetStatus("started")
	return nil
}

// ForceStop will forcefully send a SIGKILL signal to process killing it instantly.
// Returns an error in case there's any.
func (proc *Proc) ForceStop() error {
	if proc.process != nil {
		err := proc.process.Signal(syscall.SIGKILL)
		proc.Status.SetStatus("stopped")
		proc.release()
		return err
	}
	return errors.New("Process does not exist.")
}

// GracefullyStop will send a SIGTERM signal asking the process to terminate.
// The process may choose to die gracefully or ignore this signal completely. In that case
// the process will keep running unless you call ForceStop()
// Returns an error in case there's any.
func (proc *Proc) GracefullyStop() error {
	if proc.process != nil {
		err := proc.process.Signal(syscall.SIGTERM)
		proc.Status.SetStatus("asked to stop")
		return err
	}
	return errors.New("Process does not exist.")
}

// Restart will try to gracefully stop the process and then Start it again.
// Returns an error in case there's any.
func (proc *Proc) Restart() error {
	if proc.IsAlive() {
		err := proc.GracefullyStop()
		if err != nil {
			return err
		}
	}
	return proc.Start()
}

// Delete will delete everything created by this process, including the out, err and pid file.
// Returns an error in case there's any.
func (proc *Proc) Delete() error {
	proc.release()
	err := utils.DeleteFile(proc.Outfile)
	if err != nil {
		return err
	}
	err = utils.DeleteFile(proc.Errfile)
	if err != nil {
		return err
	}
	return os.RemoveAll(proc.Path)
}

// IsAlive will check if the process is alive or not.
// Returns true if the process is alive or false otherwise.
func (proc *Proc) IsAlive() bool {
	p, err := os.FindProcess(proc.Pid)
	if err != nil {
		return false
	}
	return p.Signal(syscall.Signal(0)) == nil
}

// Watch will stop execution and wait until the process change its state. Usually changing state, means that the process died.
// Returns a tuple with the new process state and an error in case there's any.
func (proc *Proc) Watch() (*os.ProcessState, error) {
	return proc.process.Wait()
}

// Will release the process and remove its PID file
func (proc *Proc) release() {
	if proc.process != nil {
		proc.process.Release()
	}
	utils.DeleteFile(proc.Pidfile)
}

// NotifyStopped that process was stopped so we can set its PID to -1
func (proc *Proc) NotifyStopped() {
	proc.Pid = -1
}

// AddRestart is add one restart to proc status
func (proc *Proc) AddRestart() {
	proc.Status.AddRestart()
}

// GetPid will return proc current PID
func (proc *Proc) GetPid() int {
	return proc.Pid
}

// GetOutFile will return proc out file
func (proc *Proc) GetOutFile() string {
	return proc.Outfile
}

// GetErrFile will return proc error file
func (proc *Proc) GetErrFile() string {
	return proc.Errfile
}

// GetPidFile will return proc pid file
func (proc *Proc) GetPidFile() string {
	return proc.Pidfile
}

// GetPath will return proc path
func (proc *Proc) GetPath() string {
	return proc.Path
}

// GetStatus will return proc current status
func (proc *Proc) GetStatus() *ProcStatus {
	if !proc.IsAlive() {
		proc.ResetUpTime()
	} else {
		// update uptime
		proc.SetUptime()
	}
	// update cpu and memory
	proc.SetSysInfo()

	return proc.Status
}

// SetStatus will set proc status
func (proc *Proc) SetStatus(status string) {
	proc.Status.SetStatus(status)
}

// SetUptime will set Uptime
func (proc *Proc) SetUptime() {
	proc.Status.SetUptime()
}

// ResetUpTime will set Uptime
func (proc *Proc) ResetUpTime() {
	proc.Status.ResetUptime()
}

// SetSysInfo will get current proc cpu and memory usage
func (proc *Proc) SetSysInfo() {
	proc.Status.SetSysInfo(proc.process.Pid)
}

// Identifier is that will be used by watcher to keep track of its processes
func (proc *Proc) Identifier() string {
	return proc.Name
}

// ShouldKeepAlive will returns true if the process should be kept alive or not
func (proc *Proc) ShouldKeepAlive() bool {
	return proc.KeepAlive
}

// GetName will return current proc name
func (proc *Proc) GetName() string {
	return proc.Name
}
