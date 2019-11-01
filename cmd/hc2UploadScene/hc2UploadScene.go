/*
hc2UploadScene provides functionality to upload a lua scene to a Fibaro HC2 system.

	Usage: hc2UploadScene [options] <lua-script>

	<lua-script> the file to be uploaded

	Options:
	--log-level, -l    Log level, one of panic, fatal, error, warn or warning, info, debug, trace
						(default info)
	--cfg-file, -c     The config file to use (default /Users/the/go/src/github.com/theovassiliou/hc2-tools/configs/config.json)
	--create-header    Create the FIBARO_GIT_HEADER if set
	--dont-upload, -d  Don't upload the file but print only
	--version, -v      display version
	--help, -h         display help

	Scene options:
	--scene-id, -s     The sceneId that shall be used. If none given, create a new scene and implies
						createHeader if header is missing (default -1)
	--room-id, -r      The roomId that shall be used. Implies createHeader if header is missing
						(default -1)
	--scene-name       The scene name that shall be used. If none given and no header in file, than
						take filename without file extenion and implies createHeader if header is
						missing

	Require Expand options:
	--dont-expand      Don't expand the require statements
	--expand-path, -e  Where to search for the included libraries

	Version:
		0.0.1-src

	Read more:
		github.com/theovassiliou/hc2-tools
*/
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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
	cmdName = "hc2UploadScene"
)

var conf = config{}

type config struct {
	LuaScript string    `type:"arg" help:"<lua-script> the file to be uploaded"`
	LogLevel  log.Level `help:"Log level, one of panic, fatal, error, warn or warning, info, debug, trace"`
	CfgFile   string    `help:"The config file to use"`
	Init      bool      `help:"Create a default config file as defined by cfg-file, if set. If not set ~/.hc2-tools/config.json will be created."`
	Test      bool      `help:"Just print information about the contacted HC2 system"`

	CreateHeader bool `help:"Create the FIBARO_GIT_HEADER if set"`

	User     string `opts:"group=HC2" help:"Username for HC2 authentication"`
	Password string `opts:"group=HC2" help:"Password for HC2 authentication"`
	URL      string `opts:"group=HC2" help:"URL of the Fibaro HC2 system, in the form http://..."`

	DontUpload bool `help:"Don't upload the file but print only"`

	SceneID   int    `opts:"group=Scene" help:"The sceneId that shall be used. If none given, create a new scene and implies createHeader if header is missing"`
	RoomID    int    `opts:"group=Scene" help:"The roomId that shall be used. Implies createHeader if header is missing"`
	SceneName string `opts:"group=Scene" help:"The scene name that shall be used. If none given and no header in file, than take filename without file extenion and implies createHeader if header is missing"`

	DontExpand bool   `opts:"group=Require Expand" help:"Don't expand the require statements"`
	ExpandPath string `opts:"group=Require Expand" help:"Where to search for the included libraries"`
}

// TODO: Find a reasonable way to parametrize the libary2Ignore feature
var m = map[string]*regexp.Regexp{
	"uncommented":   regexp.MustCompile(`(?m)^\s*require\(('|")(.*)('|")\);?`),
	"commented":     regexp.MustCompile(`(?m)^-+\s*require\(('|")(.*)('|")\)`),
	"block-library": regexp.MustCompile(``),
	"ignoreExpand":  regexp.MustCompile(`library2Ignore`),
}

func main() {
	workingHomeDir, _ := homedir.Dir()

	conf = config{
		CfgFile:      workingHomeDir + "/.hc2-tools/config.json",
		CreateHeader: false,
		DontUpload:   false,
		DontExpand:   false,
		SceneID:      -1,
		RoomID:       -1,
		SceneName:    "",
		LogLevel:     log.InfoLevel,
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
	// Assumptions:
	// CommandLine Parameters overrule file content, file content overrules config-defauls
	/* Plan:
	- Fill scene with config-details
	- Read file
	- Parse header if any
	- Overwrite with parameters, if any
	- Update header
	- Upload file based on header information
	*/
	var shallUpdateHeader = false
	hc2Scene := hc2.NewHc2Scene()

	hc2Scene.ParseFile(conf.LuaScript, false) // will exit if no such file

	if hc2Scene.SceneID == -1 {
		// There was no header

		// base filename without suffix
		name := strings.TrimSuffix(filepath.Base(conf.LuaScript), filepath.Ext(filepath.Base(conf.LuaScript)))
		hc2Scene.Name = name
		if filepath.Ext(filepath.Base(conf.LuaScript)) == ".lua" {
			hc2Scene.IsLua = true
		}
		if conf.SceneName != "" {
			hc2Scene.Name = conf.SceneName
			shallUpdateHeader = true
		}

		if conf.CreateHeader {
			shallUpdateHeader = true
		}
	} else {
		// There was a header, update values with commandline parameters

		if conf.SceneName != "" {
			hc2Scene.Name = conf.SceneName
			shallUpdateHeader = true
		} else if hc2Scene.Name == "" {
			name := strings.TrimSuffix(filepath.Base(conf.LuaScript), filepath.Ext(filepath.Base(conf.LuaScript)))
			hc2Scene.Name = name
			shallUpdateHeader = true
		}

	}

	if conf.RoomID != -1 {
		hc2Scene.RoomID = conf.RoomID
		shallUpdateHeader = true
	}

	if conf.SceneID != -1 {
		hc2Scene.SceneID = conf.SceneID
		shallUpdateHeader = true
	}

	if shallUpdateHeader {
		hc2Scene.UpdateLuaHeader()
	}

	var sb strings.Builder

	if !conf.DontExpand {
		for m["uncommented"].MatchString(hc2Scene.Lua) {

			// we have at least one uncommented require statement
			scanner := bufio.NewScanner(strings.NewReader(hc2Scene.Lua))
			for scanner.Scan() {
				if len(scanner.Text()) == 0 {
				} else {
					if m["uncommented"].MatchString(scanner.Text()) {
						sb.WriteString((m["uncommented"].ReplaceAllStringFunc(scanner.Text(), replaceFunc)))
					} else {
						sb.WriteString(scanner.Text())
					}
				}
				sb.WriteString("\n")
			}
			check(scanner.Err())
			hc2Scene.Lua = sb.String()
			sb.Reset()
		}
	}

	if !conf.DontUpload {
		if hc2Scene.SceneID == -1 {
			// we have to create a new scene in fibaro
			f.CreateScene(hc2Scene)
		} else {
			// assume that there exists a scene with this sceneId
			f.PutOneScene(hc2Scene)
		}
		os.Exit(0)
	}

	if hc2Scene.SceneID == -1 {
		log.Debugln("Creating a new scene")
	} else {
		log.Debugf("Updating scene %d\n", hc2Scene.SceneID)
	}
	log.Println(hc2Scene)
}

func slComment(s string) string {
	return "--^ " + s
}

func expandFile(reqStat, reqPar string) string {
	// get the path prefix
	pathPrefix := conf.ExpandPath

	// read the file (from pathPrefix)
	nL, f := readFile(filepath.Join(pathPrefix, reqPar+".lua"))

	// hopefully we found the file
	if nL <= 0 {
		return slComment(reqStat + " <-- FILE NOT FOUND")
	}

	var sb strings.Builder
	// building the expanded require statement
	sb.WriteString(slComment(reqStat))
	sb.WriteString(`
-- LIBRARY BEGIN -------------------------
-- DO NOT MODIFY THE CODE
`)
	sb.Write(f)
	sb.WriteString("\n-- LIBRARY END -------------------------\n")

	return sb.String()
}

func replaceFunc(s string) string {

	// in case we find the ignoreExpand key we are just
	// commenting the require statement
	if m["ignoreExpand"].MatchString(s) {
		return slComment(s)
	}

	i := m["uncommented"].FindStringSubmatch(s)

	return expandFile(s, i[2])
}

// ReadFile opens a file provided by it's path and returns the number of lines read and the file
// contents as byte array.
func readFile(path string) (int, []byte) {
	dat, readErr := ioutil.ReadFile(path)

	if readErr != nil {
		return 0, []byte{}
		// log.Fatal(readErr)
	}

	file, openErr := os.Open(path)
	if openErr != nil {
		return 0, []byte{}
		// log.Fatal(openErr)
	}
	defer file.Close()

	var noOfLines int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		noOfLines++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return noOfLines, dat
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
