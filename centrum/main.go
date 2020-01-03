package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"

	"github.com/kardianos/service"
)

// Config is the runner app config structure.
type Entry struct {
	Dir  string
	Exec string
	Args []string
	Env  []string

	Stderr string
	Stdout string

	fullExec string
	cmd      *exec.Cmd
}

type Config struct {
	Name    string
	Title   string
	Desc    string
	Entries []*Entry
}

var logger service.Logger

type program struct {
	exit    chan struct{}
	service service.Service

	config *Config
}

func (p *program) Start(s service.Service) error {
	// Look for exec.
	// Verify home directory.
	var err error
	for _, entry := range p.config.Entries {
		entry.fullExec, err = exec.LookPath(entry.Exec)
		if err != nil {
			return fmt.Errorf("failed to find executable %q: %v", entry.Exec, err)
		}
	}

	go p.run()
	return nil
}

func (p *program) exec(entry *Entry) {
	entry.cmd = exec.Command(entry.fullExec, entry.Args...)
	entry.cmd.Dir = entry.Dir
	entry.cmd.Env = append(os.Environ(), entry.Env...)

	if entry.Stderr != "" {
		f, err := os.OpenFile(entry.Stderr, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
		if err != nil {
			_ = logger.Warningf("Failed to open std err %q: %v", entry.Stderr, err)
			return
		}
		defer func() {
			_ = f.Close()
		}()
		entry.cmd.Stderr = f
	}

	if entry.Stdout != "" {
		f, err := os.OpenFile(entry.Stdout, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
		if err != nil {
			_ = logger.Warningf("Failed to open std out %q: %v", entry.Stdout, err)
			return
		}
		defer func() {
			_ = f.Close()
		}()
		entry.cmd.Stdout = f
	}

	err := entry.cmd.Run()
	if err != nil {
		_ = logger.Warningf("Error running: %v", err)
	}
}

func (p *program) run() {
	_ = logger.Info("Starting ", p.config.Title)

	defer func() {
		if service.Interactive() {
			_ = p.Stop(p.service)
		} else {
			_ = p.service.Stop()
		}
	}()

	var wg sync.WaitGroup
	for _, entry := range p.config.Entries {
		wg.Add(1)
		go func(entry *Entry) {
			defer wg.Done()
			for {
				select {
				case <-p.exit:
					return
				default:
					p.exec(entry)
					time.Sleep(3 * time.Second)
				}
			}
		}(entry)
	}
	wg.Wait()
}

func (p *program) Stop(s service.Service) error {
	close(p.exit)

	_ = logger.Info("Stopping ", p.config.Title)
	for _, entry := range p.config.Entries {
		_ = entry.cmd.Process.Kill()
		//if entry.cmd.ProcessState == nil {
		//	continue
		//}
		//if entry.cmd.ProcessState.Exited() == false {
		//	_ = entry.cmd.Process.Kill()
		//}
	}

	if service.Interactive() {
		os.Exit(0)
	}
	return nil
}

func getConfigPath() (string, error) {
	fullExecPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	dir, execName := filepath.Split(fullExecPath)
	ext := filepath.Ext(execName)
	name := execName[:len(execName)-len(ext)]

	return filepath.Join(dir, name+".yaml"), nil
}

func getConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var conf Config
	err = viper.Unmarshal(&conf)
	if err != nil {
		return nil, err
	}

	//fmt.Printf("%# v", pretty.Formatter(conf))
	return &conf, nil
}

func main() {
	logLevel := flag.String("l", "error", "error level.")
	serviceFlag := flag.String("service", "", "Control the system service.")

	flag.Parse()

	configPath, err := getConfigPath()
	if err != nil {
		log.Fatal(err)
	}

	l, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.Fatal(err)
	}

	log.SetLevel(l)

	config, err := getConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	svcConfig := &service.Config{
		Name:        config.Name,
		DisplayName: config.Title,
		Description: config.Desc,
	}

	prg := &program{
		exit:   make(chan struct{}),
		config: config,
	}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	prg.service = s

	errs := make(chan error, 5)
	logger, err = s.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Print(err)
			}
		}
	}()

	if len(*serviceFlag) != 0 {
		err := service.Control(s, *serviceFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}
	err = s.Run()
	if err != nil {
		_ = logger.Error(err)
	}
}
