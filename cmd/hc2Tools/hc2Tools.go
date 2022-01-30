package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/jpillora/opts"
	log "github.com/sirupsen/logrus"
	hc2 "github.com/theovassiliou/hc2-tools/pkg"
)

//set this via ldflags (see https://stackoverflow.com/q/11354518)
// version is the current version number as tagged via git tag 1.0.0 -m 'A message'
var (
	version = hc2.Version
	commit  string
	branch  string
	cmdName = "hc2Tools"
)

type config struct {
	Username string `opts:"group=HC2" help:"Username for HC2 authentication"`
	Password string `opts:"group=HC2" help:"Password for HC2 authentication"`
	URL      string `opts:"group=HC2" help:"URL of the Fibaro HC2 system, in the form http://aa.bb.cc.dd"`

	LogLevel log.Level `help:"Log level, one of panic, fatal, error, warn or warning, info, debug, trace"`
}

const shortUsage = "Adding, modifying and deleting items from a ShopShop list"

var conf config

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
		Version(hc2.FormatFullVersion(cmdName, version, branch, commit)).
		AddCommand(
			opts.New(&devices{}).
				Summary(devicesUsage)).
		AddCommand(
			opts.New(&showGlobalVar{}).
				Summary(showGlobalVarUsage)).
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
		Username: conf.Username,
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

func getGlobalVariables(globalVars []string) []hc2.Hc2GlobalVar {
	var allVars []hc2.Hc2GlobalVar
	var selectedVars []hc2.Hc2GlobalVar
	allVars = f.AllGlobalVariables()

	if globalVars == nil {
		return allVars
	}
	m := make(map[string]hc2.Hc2GlobalVar)
	for _, gv := range allVars {
		m[gv.Name] = gv
	}

	for _, name := range globalVars {
		h, ok := m[name]
		if ok {
			h = f.GlobalVariable(name)
			selectedVars = append(selectedVars, h)
		}

	}

	return selectedVars
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
