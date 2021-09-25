package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

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

type devices struct {
	DeviceIds []int `type:"arg" name:"deviceId" help:"device to retrieve. All if no deviceIDs given."`
	All       bool  `help:"show also invisble devices"`
}

const devicesUsage = "Lists devices, all if no deviceID given"

func (cmd *devices) Run() {
	var allDevices = getDevices(cmd.DeviceIds)

	var i = 0
	for _, device := range allDevices {
		if cmd.DeviceIds == nil && (device.Visible || cmd.All) {
			i++
			fmt.Printf("%d %s: %s with ID: %d\n", i, device.Name, device.Type, device.ID)
		} else if selected(cmd.DeviceIds, device.ID) {
			fmt.Printf("%#v\n\n", device)
		}
	}

}

type showRemoteController struct {
	DeviceIds []int `type:"arg" name:"deviceId" help:"Display button features. All if no deviceIDs given."`
	All       bool  `help:"show also invisble devices"`
}

const showRemoteControllerUsage = "List button features, all if no deviceID given"

func (cmd *showRemoteController) Run() {
	var allDevices = getDevices(cmd.DeviceIds)

	tmpl := parsedTemplate("printbuttonfeatures", "templates/printButtonFeatures.template")

	var i = 0

	for _, device := range allDevices {

		if device.Implements("zwaveCentralScene") && (device.Visible || cmd.All) {
			i++
			if cmd.DeviceIds == nil {
				fmt.Printf("%d %s: %s with ID: %d \n", i, device.Name, device.Type, device.ID)
			} else if selected(cmd.DeviceIds, device.ID) {
				var s []hc2.Key
				json.Unmarshal([]byte(device.Properties.CentralSceneSupport.(string)), &s)
				fmt.Printf("\n%d %s: %s with ID: %d", i, device.Name, device.Type, device.ID)
				err := tmpl.Execute(os.Stdout, s)
				if err != nil {
					log.Panic(err)
				}
			}
		}
	}

}

type createSceneActivationScript struct {
	DeviceIds []int `type:"arg" name:"deviceId" help:"DeviceId for which to create the lua-scipt."`
	Vd        bool  `help:"create also the corresponding VD script"`
}

const createSceneActivationScriptUsage = "Create a template lua script for a SceneActivation device"

func (cmd *createSceneActivationScript) Run() {
	var allDevices = getDevices(cmd.DeviceIds)

	pCSHTemplate := parsedTemplate("printCentralSceneHandler", "templates/printCentralSceneHandler.lua.template")

	for _, device := range allDevices {
		err := pCSHTemplate.Execute(os.Stdout, device)
		if err != nil {
			log.Panic(err)
		}
	}

}

type showSceneActivation struct {
	DeviceIds []int `type:"arg" name:"deviceId" help:"Display scene activation module. All if no deviceIDs given."`
	All       bool  `help:"show also invisble devices"`
}

const showSceneActivationrUsage = "List button features, all if no deviceID given"

func (cmd *showSceneActivation) Run() {
	var allDevices = getDevices(cmd.DeviceIds)

	var i = 0

	for _, device := range allDevices {
		if device.Implements("zwaveSceneActivation") && (device.Visible || cmd.All) {
			i++
			if cmd.DeviceIds == nil {
				fmt.Printf("%d %s: %s with ID: %d\n", i, device.Name, device.Type, device.ID)
			} else if selected(cmd.DeviceIds, device.ID) {
				fmt.Printf("%#v\n\n", device)
			}
		}
	}
}

type showHues struct {
	DeviceIds []int `type:"arg" name:"deviceId" help:"device to retrieve"`
	All       bool  `help:"show also invisble devices"`
	VslStyle  bool  `type:"flag"`
}

const showHuesUsage = "Print current HUE values"

func (cmd *showHues) Run() {
	var allDevices = getDevices(cmd.DeviceIds)

	pVslTemplate := parsedTemplate("printhuevaluesVSL", "templates/printHueValuesVSL.template")
	pTemplate := parsedTemplate("printhuevalues", "templates/printHueValues.template")

	var i = 0
	for _, device := range allDevices {
		if device.Type == "com.fibaro.philipsHueLight" && ((device.Visible) || cmd.All) {
			i++
			if cmd.DeviceIds == nil || selected(cmd.DeviceIds, device.ID) {
				fmt.Printf("%d ", i)
				var err error
				if cmd.VslStyle {
					err = pVslTemplate.Execute(os.Stdout, device)
				} else {
					err = pTemplate.Execute(os.Stdout, device)
				}
				if err != nil {
					log.Panic(err)
				}
			}
		}
	}

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
		Version(formatFullVersion("hc2Tools", version, branch, commit)).
		AddCommand(
			opts.New(&devices{}).
				Summary(devicesUsage)).
		AddCommand(
			opts.New(&showHues{}).
				Summary(showHuesUsage)).
		AddCommand(
			opts.New(&showRemoteController{}).
				Summary(showRemoteControllerUsage)).
		AddCommand(
			opts.New(&showSceneActivation{}).
				Summary(showSceneActivationrUsage)).
		AddCommand(
			opts.New(&createSceneActivationScript{}).
				Summary(createSceneActivationScriptUsage)).
		Parse()
	log.SetLevel(conf.LogLevel)

	hcConfig := hc2.FibaroConfig{
		Username: conf.User,
		Password: conf.Password,
		BaseURL:  conf.URL,
	}
	hc2.Default(&hcConfig)

	f = &hc2.FibaroHc2{}
	f.SetConfig(hcConfig)
	log.Traceln(f)

	if cmd.IsRunnable() {
		cmd.Run()
	}
}
func getDevices(deviceIDs []int) []hc2.Hc2Device {
	var allDevices []hc2.Hc2Device
	if deviceIDs == nil {
		allDevices = f.AllDevices()
	} else {
		for _, id := range deviceIDs {
			allDevices = append(allDevices, *f.OneDevice(id))
		}
	}
	return allDevices
}

func selected(deviceIds []int, deviceID int) bool {
	for _, device := range deviceIds {
		if deviceID == device {
			return true
		}
	}
	return false
}

func parsedTemplate(name, filename string) *template.Template {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Panic(err)
	}
	pVslTemplate, err := template.New(name).Parse(string(file))
	if err != nil {
		log.Panic(err)
	}
	return pVslTemplate
}

func formatFullVersion(cmdName, version, branch, commit string) string {
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
