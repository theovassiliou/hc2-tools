package main

import (
	"fmt"
	"strings"

	"github.com/jpillora/opts"
	log "github.com/sirupsen/logrus"
	hc2 "github.com/theovassiliou/hc2-tools/pkg"
)

//set this via ldflags (see https://stackoverflow.com/q/11354518)
const pVersion = ".1"

// version is the current version number as tagged via git tag 1.0.0 -m 'A message'
var (
	version = "1.0" + pVersion + "-src"
	commit  string
	branch  string
)

type config struct {
	User     string `opts:"group=HC2" help:"Username for HC2 authentication"`
	Password string `opts:"group=HC2" help:"Password for HC2 authentication"`
	URL      string `opts:"group=HC2" help:"URL of the Fibaro HC2 system, in the form http://aa.bb.cc.dd"`

	LogLevel log.Level `help:"Log level, one of panic, fatal, error, warn or warning, info, debug, trace"`
}

const shortUsage = "Adding, modifying and deleting items from a ShopShop list"

var conf config

type allDevices struct {
	ItemDescription []string `type:"arg" name:"description" help:"item to add"`
	Quantity        string   `type:"flag"`
}

const addUsage = "Lists all devices"

func (cmd *allDevices) Run() {

	fmt.Println(f.AllDevices())
}

// FormatFullVersion formats for a cmdName the version number based on version, branch and commit
func FormatFullVersion(cmdName, version, branch, commit string) string {
	var parts = []string{cmdName}

	if version != "" {
		parts = append(parts, version)
	} else {
		parts = append(parts, "unknown")
	}

	if branch != "" || commit != "" {
		if branch == "" {
			branch = "unknown"
		}
		if commit == "" {
			commit = "unknown"
		}
		git := fmt.Sprintf("(git: %s %s)", branch, commit)
		parts = append(parts, git)
	}

	return strings.Join(parts, " ")
}

var f *hc2.FibaroHc2

func main() {

	conf = config{
		LogLevel: log.DebugLevel,
	}

	//parse config
	cmd := opts.
		New(&conf).
		Summary(shortUsage).
		PkgRepo().
		UserConfigPath().
		Version(FormatFullVersion("hc2Tools", version, branch, commit)).
		AddCommand(
			opts.New(&allDevices{}).
				Summary(addUsage)).
		Parse()
	log.SetLevel(conf.LogLevel)

	hcConfig := hc2.FibaroConfig{
		Username: conf.User,
		Password: conf.Password,
		BaseURL:  conf.URL,
	}

	f = &hc2.FibaroHc2{}
	f.SetConfig(hcConfig)
	log.Traceln(f)

	if cmd.IsRunnable() {
		cmd.Run()
	}
}
