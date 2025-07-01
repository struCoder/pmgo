package preparable

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/struCoder/pmgo/lib/process"
)

// ProcPreparable is a preparable with all the necessary informations to run
// a process. To actually run a process, call the Start() method.
type ProcPreparable interface {
	PrepareBin() ([]byte, error)
	Start() (process.ProcContainer, error)
	getPath() string
	Identifier() string
	getBinPath() string
	getPidPath() string
	getOutPath() string
	getErrPath() string
}

type Preparable struct {
	Name       string
	SourcePath string
	Cmd        string
	SysFolder  string
	Language   string
	KeepAlive  bool
	Args       []string
}

type BinaryPreparable struct {
	Name       string
	SourcePath string
	Cmd        string
	SysFolder  string
	Language   string
	KeepAlive  bool
	Args       []string
}

// PrepareBin will compile the Golang project from SourcePath and populate Cmd with the proper
// command for the process to be executed.
// Returns the compile command output.
func (preparable *Preparable) PrepareBin() ([]byte, error) {
	// Remove the last character '/' if present
	if len(preparable.SourcePath) > 0 && preparable.SourcePath[len(preparable.SourcePath)-1] == '/' {
		preparable.SourcePath = strings.TrimSuffix(preparable.SourcePath, "/")
	}

	binPath := preparable.getBinPath()

	cmd := ""
	cmdArgs := []string{}
	if preparable.Language == "go" {
		cmd = "go"

		// Check if SourcePath is an absolute path or relative path
		sourceCodePath := preparable.SourcePath
		if filepath.IsAbs(sourceCodePath) {
			// For absolute paths, check if it's a .go file
			if strings.HasSuffix(sourceCodePath, ".go") {
				// It's a Go source file, build it directly
				cmdArgs = []string{"build", "-o", binPath, sourceCodePath}
			} else {
				// It's a directory with Go package
				cmdArgs = []string{"build", "-o", binPath, sourceCodePath + "/."}
			}
		} else {
			// For relative paths, try modern Go modules first, fallback to GOPATH
			if _, err := os.Stat(sourceCodePath); err == nil {
				// Path exists relative to current directory
				if strings.HasSuffix(sourceCodePath, ".go") {
					cmdArgs = []string{"build", "-o", binPath, sourceCodePath}
				} else {
					cmdArgs = []string{"build", "-o", binPath, sourceCodePath + "/."}
				}
			} else {
				// Fallback to GOPATH structure, but only if GOPATH is actually set
				gopath := os.Getenv("GOPATH")
				if gopath != "" {
					gopathSrc := gopath + "/src/" + sourceCodePath
					if _, err := os.Stat(gopathSrc); err == nil {
						cmdArgs = []string{"build", "-o", binPath, gopathSrc + "/."}
					} else {
						return make([]byte, 0), fmt.Errorf("source path not found: %s (checked: %s)", sourceCodePath, gopathSrc)
					}
				} else {
					// No GOPATH set and relative path doesn't exist - try building anyway for Go modules
					if strings.HasSuffix(sourceCodePath, ".go") {
						cmdArgs = []string{"build", "-o", binPath, sourceCodePath}
					} else {
						cmdArgs = []string{"build", "-o", binPath, sourceCodePath + "/."}
					}
				}
			}
		}
	}

	preparable.Cmd = binPath
	out, err := exec.Command(cmd, cmdArgs...).Output()
	if err != nil {
		return []byte("build source fail!"), err
	}
	log.Infoln("build output: ", string(out))
	return make([]byte, 0), nil
}

// Start will execute the process based on the information presented on the preparable.
// This function should be called from inside the master to make sure
// all the watchers and process handling are done correctly.
// Returns a tuple with the process and an error in case there's any.
func (preparable *Preparable) Start() (process.ProcContainer, error) {
	proc := &process.Proc{
		Name:      preparable.Name,
		Cmd:       preparable.Cmd,
		Args:      preparable.Args,
		Path:      preparable.getPath(),
		Pidfile:   preparable.getPidPath(),
		Outfile:   preparable.getOutPath(),
		Errfile:   preparable.getErrPath(),
		KeepAlive: preparable.KeepAlive,
		Status:    &process.ProcStatus{},
	}

	err := proc.Start()
	return proc, err
}

// Identifier is a function that get proc name
func (preparable *Preparable) Identifier() string {
	return preparable.Name
}

func (preparable *Preparable) getPath() string {
	if preparable.SysFolder[len(preparable.SysFolder)-1] == '/' {
		preparable.SysFolder = strings.TrimSuffix(preparable.SysFolder, "/")
	}
	return preparable.SysFolder + "/" + preparable.Name
}

func (preparable *Preparable) getBinPath() string {
	return preparable.getPath() + "/" + preparable.Name
}

func (preparable *Preparable) getPidPath() string {
	return preparable.getBinPath() + ".pid"
}

func (preparable *Preparable) getOutPath() string {
	return preparable.getBinPath() + ".out"
}

func (preparable *Preparable) getErrPath() string {
	return preparable.getBinPath() + ".err"
}

// BinaryPreparable

// PrepareBin checks if the given binary path is a valid executable.
// Returns no bytes, but if there is an error, it will be returned.
func (preparable *BinaryPreparable) PrepareBin() ([]byte, error) {
	info, err := os.Stat(preparable.SourcePath)
	if err != nil {
		return make([]byte, 0), err
	}
	if info.Mode()&0111 == 0 {
		return make([]byte, 0), fmt.Errorf("The given source path(%s) is not executable, neither a Go source to compile", preparable.SourcePath)
	}

	err = os.MkdirAll(filepath.Dir(preparable.getOutPath()), 0755)
	if err != nil {
		return make([]byte, 0), err
	}
	preparable.Cmd = preparable.SourcePath
	return make([]byte, 0), err
}

// Start will execute the process based on the information presented on the preparable.
// This function should be called from inside the master to make sure
// all the watchers and process handling are done correctly.
// Returns a tuple with the process and an error in case there's any.
func (preparable *BinaryPreparable) Start() (process.ProcContainer, error) {
	proc := &process.Proc{
		Name:      preparable.Name,
		Cmd:       preparable.Cmd,
		Args:      preparable.Args,
		Path:      preparable.getPath(),
		Pidfile:   preparable.getPidPath(),
		Outfile:   preparable.getOutPath(),
		Errfile:   preparable.getErrPath(),
		KeepAlive: preparable.KeepAlive,
		Status:    &process.ProcStatus{},
	}

	err := proc.Start()
	return proc, err
}

// Identifier is a function that get proc name
func (preparable *BinaryPreparable) Identifier() string {
	return preparable.Name
}

func (preparable *BinaryPreparable) getPath() string {
	if preparable.SysFolder[len(preparable.SysFolder)-1] == '/' {
		preparable.SysFolder = strings.TrimSuffix(preparable.SysFolder, "/")
	}
	return preparable.SysFolder + "/" + preparable.Name
}

func (preparable *BinaryPreparable) getBinPath() string {
	return preparable.SourcePath
}

func (preparable *BinaryPreparable) getPidPath() string {
	return preparable.getPath() + "/" + preparable.Name + ".pid"
}

func (preparable *BinaryPreparable) getOutPath() string {
	return preparable.getPath() + "/" + preparable.Name + ".out"
}

func (preparable *BinaryPreparable) getErrPath() string {
	return preparable.getPath() + "/" + preparable.Name + ".err"
}
