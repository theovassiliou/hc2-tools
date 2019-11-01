package main

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"

	"github.com/jpillora/opts"

	hc2 "github.com/theovassiliou/hc2-tools/pkg"
)

//set this via ldflags (see https://stackoverflow.com/q/11354518)
var (
	version = hc2.Version
	commit  string
	branch  string
	cmdName = "hc2SceneInteract"
)

var conf = config{}

type config struct {
	Action   hc2.SceneActionCommand `help:"Triggers a scene action for sceneID. One of start, stop, enable, disable."`
	GetDebug bool                   `help:"Retrieve after starting the action the debug messages, while respecting the tail-flag. Ignored with action enable or disable"`

	Init bool `help:"Create a default config file as defined by cfg-file, if set. If not set ~/.hc2-tools/config.json will be created."`
	Test bool `help:"Just print information about the contacted HC2 system"`

	User     string `opts:"group=HC2" help:"Username for HC2 authentication"`
	Password string `opts:"group=HC2" help:"Password for HC2 authentication"`
	URL      string `opts:"group=HC2" help:"URL of the Fibaro HC2 system, in the form http://..."`

	SceneID  int       `opts:"group=Generic command" help:"The sceneId that shall be used."`
	File     string    `opts:"group=Generic command" help:"sceneID is taken from <lua-script-file> with fibaro header. sceneID flag is ignored."`
	Tail     bool      `opts:"group=Generic command" help:"The -t option causes get-debug to not stop when all debug messages are read, but rather to wait for additional data to be appended to the input."`
	CfgFile  string    `help:"The config file to use"`
	LogLevel log.Level `help:"Log level, one of panic, fatal, error, warn or warning, info, debug, trace"`
}

var m = map[string]*regexp.Regexp{
	"span-color-open":  regexp.MustCompile(`<span style=\"color:([a-z]+);\">`),
	"span-color-close": regexp.MustCompile(`</span>`),
}

func main() {
	workingHomeDir, _ := homedir.Dir()

	conf = config{
		CfgFile:  workingHomeDir + "/" + hc2.Hc2DefaultConfigFile,
		LogLevel: log.InfoLevel,
	}

	//parse config
	opts.New(&conf).
		Repo(hc2.RepoName).
		Version(hc2.FormatFullVersion(cmdName, version, branch, commit)).
		Parse()

	log.SetLevel(conf.LogLevel)

	hc2Scene := hc2.NewHc2Scene()

	var f *hc2.FibaroHc2
	if !conf.Init {
		f = hc2.NewFibaroHc2Config(conf.CfgFile)
	}

	if f == nil {
		if conf.User == "" && conf.Password == "" && conf.URL == "" {
			log.Fatalf("Could not read config file (%s) and no parameters given.\n"+
				" Consider using --init to create a config file\n", conf.CfgFile)
		} else if conf.User != "" && conf.Password != "" && conf.URL != "" {
			var confg hc2.FibaroConfig
			hc2.Default(&confg)
			f = &hc2.FibaroHc2{}
			f.SetConfig(confg)
		} else {
			log.Fatalf("Not all login parameters provided. Aborting.")

		}
	}

	if conf.User != "" {
		cfg := f.Config()
		log.Tracef("Configured user %s\n", conf.User)
		cfg.Username = conf.User
	}

	if conf.Password != "" {
		cfg := f.Config()
		cfg.Password = conf.Password
	}

	if conf.URL != "" {
		cfg := f.Config()
		cfg.BaseURL = conf.URL
	}

	if conf.Init {
		var i int
		var filePath string
		if conf.CfgFile != "" {
			filePath = conf.CfgFile
		} else {
			filePath = workingHomeDir + "/" + hc2.Hc2DefaultConfigFile
		}
		i = f.WriteInitConfigFile(filePath)
		log.Debugf("Wrote to file %s %d bytes", filePath, i)
	}

	if conf.Test {
		fmt.Println(f.Info(2))
		os.Exit(0)
	}

	if conf.Action == hc2.Undef && !conf.GetDebug {
		log.Errorln("Nothing to be done. Aborting.")
	}

	if conf.SceneID != 0 {
		hc2Scene.SceneID = conf.SceneID
	}
	if hc2Scene.SceneID == -1 {
		if conf.File == "" {
			log.Fatalln("No SceneID and no file given. Aborting.")
		}
		hc2Scene.ParseFile(conf.File, false) // will exit if no such file
		if hc2Scene.SceneID == -1 {
			log.Fatalf("No SceneID included in file %s. Aborting", conf.File)
		}
	}
	runAction(f, hc2Scene.SceneID)
	runGetMessage(f, hc2Scene.SceneID)
}

func runAction(f *hc2.FibaroHc2, sceneID int) error {
	if conf.Action == hc2.Undef {
		return nil
	}

	err := f.Action(sceneID, conf.Action)
	if err != nil {
		log.Fatalln("Error " + err.Error() + " while trying to trigger action: " + conf.Action.String() + ". Aborting.")
	}

	if conf.GetDebug {
		// Let's give the HC2 a little bit time to generate new messages.
		time.Sleep(1 * time.Second)
	}
	return nil
}

func runGetMessage(f *hc2.FibaroHc2, sceneID int) error {

	if !conf.GetDebug {
		return nil
	}

	var lenPrinted int
	var firstTimestamp int64

	for {
		dm := f.DebugMessages(sceneID)
		if len(dm) > 0 {
			if dm[0].Timestamp > firstTimestamp {
				firstTimestamp = dm[0].Timestamp
				lenPrinted = 0
			}
		}
		for ; lenPrinted < len(dm); lenPrinted++ {
			value := dm[lenPrinted]
			t := time.Unix(value.Timestamp, 0)
			var s string
			s = m["span-color-open"].ReplaceAllString(value.Txt, "")
			s = m["span-color-close"].ReplaceAllString(s, "")
			fmt.Printf("[%s] %s: %s\n", value.Type, t.Format("15:04:05"), s)
		}
		if !conf.Tail {
			break
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}
