package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

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
	cmdName = "hc2DownloadScene"
)

var conf = config{}

type config struct {
	LogLevel log.Level `help:"Log level, one of panic, fatal, error, warn or warning, info, debug, trace"`
	CfgFile  string    `help:"The config file to use"`
	Init     bool      `help:"Create a default config file as defined by cfg-file, if set. If not set ~/.hc2-tools/config.json will be created."`
	Test     bool      `help:"Just print information about the contacted HC2 system"`

	User     string `opts:"group=HC2" help:"Username for HC2 authentication"`
	Password string `opts:"group=HC2" help:"Password for HC2 authentication"`
	URL      string `opts:"group=HC2" help:"URL of the Fibaro HC2 system, in the form http://..."`

	CreateHeader bool   `opts:"group=Scene" help:"If set create the FIBARO_GIT_HEADER if none present"`
	SceneID      int    `opts:"group=Scene" help:"The sceneId that shall be used. If none given, all scenes will be downloaded."`
	Dir          string `opts:"group=Scene" help:"Where to search for the included libraries"`
}

func main() {
	workingHomeDir, _ := homedir.Dir()

	conf = config{
		CfgFile:      workingHomeDir + "/" + hc2.Hc2DefaultConfigFile,
		CreateHeader: true,
		SceneID:      -1,
		LogLevel:     log.InfoLevel,
		Dir:          "./download",
	}

	//parse config
	opts.New(&conf).
		Repo(hc2.RepoName).
		Version(hc2.FormatFullVersion(cmdName, version, branch, commit)).
		Parse()

	log.SetLevel(conf.LogLevel)

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

	if conf.SceneID == -1 {
		allScenes := f.AllScenes()
		log.Infof("Processing %d scenes\n", len(allScenes))

		var bytesWrote int
		var filesCreated int

		for i, aScene := range allScenes {
			amountOfBytes := writeFile(f, conf.Dir, f.OneScene(aScene.SceneID))
			log.Debugf("%d: Wrote %d:%s \n", i, aScene.SceneID, aScene.Name)
			bytesWrote += amountOfBytes
			filesCreated++
		}

		log.Infof("retrieved %d scenes\n", len(allScenes))
		log.Infof("created %d files\n", filesCreated)
		log.Infof("wrote %d bytes\n", bytesWrote)

	} else {
		s := f.OneScene(conf.SceneID)
		if s.SceneID == -1 {
			log.Fatalf("scene with id %d does not exists\n", conf.SceneID)
		}
		bytesWrote := writeFile(f, conf.Dir, s)
		log.Infof("retrieved scene %d", conf.SceneID)
		log.Infof("wrote %d bytes\n", bytesWrote)
		log.Infof("created file: %s\n", s.Name)

	}
}

func writeFile(fib *hc2.FibaroHc2, baseDir string, scene hc2.Hc2Scene) (bytesWrote int) {

	room := fib.OneRoom(scene.RoomID)
	section := fib.OneSection(room.SectionID)
	path := filepath.Join(baseDir, section.Name, room.Name)
	os.MkdirAll(path, os.ModePerm)
	file := filepath.Join(path, scene.Name+".lua")

	i := 0
	// check whether it exists and create a unique if yes
	e, err := os.Open(file)
	defer e.Close()

	for err == nil {
		log.Debugf("   File %s exists. Creating new fileName", file)
		file = filepath.Join(path, scene.Name+"_"+strconv.Itoa(i)+".lua")
		e, err = os.Open(file)
		defer e.Close()
		i++
	}

	f, err := os.Create(file)
	defer f.Close()

	if err != nil {
		log.Printf("Problem creating file %s; %v", file, err)
		return 0
	}

	check(err)
	w := bufio.NewWriter(f)
	n4, err := w.WriteString(scene.Lua + "\n")
	check(err)

	var isPresent hc2.Hc2Scene
	isPresent.Parse([]byte(scene.Lua))
	if isPresent.SceneID == -1 {
		log.Infoln("No LuaSpec in file: " + file)
		log.Infoln("Adding ... ")

		n5, err := w.WriteString(scene.ToComment())
		check(err)
		n4 += n5
	}
	w.Flush()
	return n4

}

func check(e error) {
	if e != nil {
		log.Panic(e)
	}
}
