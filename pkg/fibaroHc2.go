package fibarohc2

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	resty "gopkg.in/resty.v1"
)

// SceneActionCommand are the commands to be used in an action
type SceneActionCommand int

// Start, Stop, Enable, Disable are the commands to be used when triggering an action on a scene.
const (
	Undef SceneActionCommand = iota
	Start
	Stop
	Enable
	Disable
)

func (b SceneActionCommand) String() string {
	return [...]string{"undef", "start", "stop", "enable", "disable"}[b]
}

// Set the SceneActionCommand based on the string to be parsed. Returns error if not defined.
func (b *SceneActionCommand) Set(s string) error {
	switch s {
	case "start":
		*b = Start
	case "stop":
		*b = Stop
	case "enable":
		*b = Enable
	case "disable":
		*b = Disable
	default:
		return errors.New("None of start, stop, enable, disable")
	}
	return nil
}

// FibaroHc2 represents the Fibaro HC2 system.
type FibaroHc2 struct {
	cfg FibaroConfig // contains the configuration on how to access the Fibaro system.

}

// NewFibaroHc2Config creates a new FibaroHc2 object with the configuration
// read from a file. The information in the file is JSON encoded
func NewFibaroHc2Config(file string) *FibaroHc2 {

	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		return nil
	}
	stat, _ := configFile.Stat()
	if stat.IsDir() {
		return nil
	}
	var cfg FibaroConfig
	Default(&cfg)

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&cfg)

	return &FibaroHc2{cfg}
}

// Default defines the default values for possible configuration file fields
func Default(fc *FibaroConfig) {
	fc.CreateHeader = true
}

// Config returns the configuration of the Fibaro system
func (f *FibaroHc2) Config() *FibaroConfig {
	return &f.cfg
}

// SetConfig sets the configuration of the FibaroHC2
func (f *FibaroHc2) SetConfig(fc FibaroConfig) {
	f.cfg = fc
}

// FibaroConfig represents the configuration on how to access the Fibaro HC2 system
type FibaroConfig struct {
	BaseURL      string `json:"url"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	CreateHeader bool   `json:"createHeader"`
}

func requestGet(cfg FibaroConfig, cmd string) (resp *resty.Response, err error) {
	msg := cfg.Username + ":" + cfg.Password
	encoded := "Basic " + base64.StdEncoding.EncodeToString([]byte(msg))
	log.Tracef("msg: %s\n", msg)
	log.Tracef("encoded: %s\n", encoded)

	resp, err = resty.R().
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", encoded).
		Get(cfg.BaseURL + "/api" + cmd)
	return
}

func requestPut(cfg FibaroConfig, cmd string, body []byte) (resp *resty.Response, err error) {
	msg := cfg.Username + ":" + cfg.Password
	encoded := "Basic " + base64.StdEncoding.EncodeToString([]byte(msg))

	log.Traceln(string(body))
	resp, err = resty.R().
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", encoded).
		SetBody(body).
		Put(cfg.BaseURL + "/api" + cmd)
	log.Tracef("%#v\n", resp.Status())
	return
}

func requestPost(cfg FibaroConfig, cmd string, body []byte) (resp *resty.Response, err error) {
	msg := cfg.Username + ":" + cfg.Password
	encoded := "Basic " + base64.StdEncoding.EncodeToString([]byte(msg))

	resp, err = resty.R().
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", encoded).
		SetBody(body).
		Post(cfg.BaseURL + "/api" + cmd)
	return
}

// PutOneScene upploads a scene to the FibaroHc2 system and returns the response, or any error encountered.
func (f *FibaroHc2) PutOneScene(scene Hc2Scene) (resp *resty.Response, err error) {
	// TODO We need some sanity checks here for the H2Scene
	// 1. runconfig only one of TriggerAndManual, ManualOnly,  Disabled
	// 2. luaType == "com.fibaro.luaScene"
	// What else ?

	if !scene.SanityCheck() {
		return nil, fmt.Errorf("Sanity check failed")

	}

	b, mer := json.Marshal(scene)
	if mer != nil {
		log.Errorln("Error while marshalling hc2Scene:", mer)
		return nil, mer
	}

	resp, err = requestPut(f.cfg, "/scenes/"+strconv.Itoa(scene.SceneID), b)
	if resp.StatusCode() == 404 {
		log.Fatalf("Could not upload scene \"%s\" with sceneID=%d as it does not exists in Fibaro HC2", scene.Name, scene.SceneID)
	}
	return

}

// AllScenes downloads and returns all scenes of the FibaroHC2 system.
// nil will be returned in case an error occured while downloading the
// scenes.
func (f *FibaroHc2) AllScenes() []Hc2Scene {
	resp, err := requestGet(f.cfg, "/scenes")

	var s []Hc2Scene
	json.Unmarshal(resp.Body(), &s)
	if err != nil {
		log.Errorln("Error while decoding hc2scenes: ", err)
		return nil
	}
	return s
}

// OneScene downloads and returns one scene as identified by the sceneID
func (f *FibaroHc2) OneScene(sceneID int) Hc2Scene {
	resp, err := requestGet(f.cfg, "/scenes/"+strconv.Itoa(sceneID))

	if len(resp.Body()) == 0 || err != nil {
		// scene not found
		return Hc2Scene{SceneID: -1}
	}
	var s Hc2Scene
	json.Unmarshal(resp.Body(), &s)
	return s
}

// CreateScene creates a new scene in the fibaro system with the name parameters
// set in the scene. SceneID will be updated with the new scene id being allocated.
// the header in the lua field will be updated. On sucess CreateScene will return the allocated SceneID, or -1 on error.
func (f *FibaroHc2) CreateScene(scene Hc2Scene) (newSceneID int) {
	dummyScene := Hc2Scene{
		SceneID: -1,
		Name:    scene.Name,
		Type:    "com.fibaro.luaScene",
	}

	b, merr := json.Marshal(dummyScene)
	if merr != nil {

		log.Errorln("Error while marshalling hc2Scene:", merr)
		return -1
	}
	resp, err := requestPost(f.cfg, "/scenes", b)
	if len(resp.Body()) == 0 || err != nil {
		// some kind of error
		log.Errorf("Error creating new scene with header=%v. Error:%v\n", string(b), err)
		return -1
	}
	freshScene := Hc2Scene{}
	uerr := json.Unmarshal(resp.Body(), &freshScene)
	if uerr != nil {
		log.Errorf("Error while unmarshalling received intermediate: %v", uerr)
		return -1
	}

	scene.SceneID = freshScene.SceneID
	scene.Type = freshScene.Type
	scene.UpdateLuaHeader()
	_, err = f.PutOneScene(scene)
	if err != nil {
		log.Errorf("Error while updating intermediate scene (%v): %v", scene.SceneID, err)
		return -1
	}
	return scene.SceneID
}

// OneRoom downloads and returns a room as identified by the roomID
func (f *FibaroHc2) OneRoom(roomID int) Hc2Room {
	resp, err := requestGet(f.cfg, "/rooms/"+strconv.Itoa(roomID))

	if len(resp.Body()) == 0 || err != nil {
		// room not found
		return Hc2Room{RoomID: -1}
	}

	var s Hc2Room
	json.Unmarshal(resp.Body(), &s)

	return s
}

// OneSection downloads and returns a section as identified by the sectionID
func (f *FibaroHc2) OneSection(sectionID int) Hc2Section {
	resp, err := requestGet(f.cfg, "/sections/"+strconv.Itoa(sectionID))

	if len(resp.Body()) == 0 || err != nil {
		// scene not found
		return Hc2Section{SectionID: -1}
	}

	var s Hc2Section
	json.Unmarshal(resp.Body(), &s)

	return s
}

// DebugMessages downloads and returns all debug messages for a given sceneID
func (f *FibaroHc2) DebugMessages(sceneID int) []Hc2DebugMessage {
	resp, err := requestGet(f.cfg, "/scenes/"+strconv.Itoa(sceneID)+"/debugMessages")

	if len(resp.Body()) == 0 || err != nil {
		// scene not found
		return nil
	}
	var s []Hc2DebugMessage
	json.Unmarshal(resp.Body(), &s)
	return s
}

// Action starts an actions on a given sceneID
func (f *FibaroHc2) Action(sceneID int, c SceneActionCommand) error {
	resp, err := requestPost(f.cfg, "/scenes/"+strconv.Itoa(sceneID)+"/action/"+c.String(), []byte(""))
	if err != nil {
		return err
	}
	log.Debug(resp)
	if len(resp.Body()) > 0 {
		return errors.New(resp.String())
	}
	return nil
}

func (f *FibaroHc2) settingsInfo() Hc2Info {
	resp, err := requestGet(f.cfg, "/settings/info")
	log.Debug(resp)

	if len(resp.Body()) == 0 || err != nil {
		// scene not found
		return Hc2Info{}
	}

	var s Hc2Info
	json.Unmarshal(resp.Body(), &s)

	return s
}

func (f *FibaroHc2) settingsNetwork() Hc2Network {
	resp, err := requestGet(f.cfg, "/settings/network")
	log.Debug(resp)

	if len(resp.Body()) == 0 || err != nil {
		// scene not found
		return Hc2Network{}
	}

	var s Hc2Network
	json.Unmarshal(resp.Body(), &s)

	return s
}

func (f *FibaroHc2) loginStatus() Hc2LoginStatus {
	resp, err := requestGet(f.cfg, "/loginStatus")
	log.Debug(resp)

	if len(resp.Body()) == 0 || err != nil {
		// scene not found
		return Hc2LoginStatus{}
	}

	var s Hc2LoginStatus
	json.Unmarshal(resp.Body(), &s)

	if s.Username == "" {
		h := Hc2LoginStatus{ErrorMsg: string(resp.Body())}
		return h
	}

	return s
}

// Info returns a human-readable string on information of the HC2
// TODO: Write Tests
func (f *FibaroHc2) Info(indent int) string {
	ind := strings.Repeat(" ", indent)

	i := f.settingsInfo()
	n := f.settingsNetwork()
	l := f.loginStatus()

	var str strings.Builder
	_, _ = str.WriteString("Successful connected to ...\n")
	_, _ = str.WriteString(ind + "Name         : " + i.HcName + "\n")
	_, _ = str.WriteString(ind + "Serial       : " + i.SerialNumber + "\n")
	_, _ = str.WriteString(ind + "IP           : " + n.IP + "\n")
	_, _ = str.WriteString(ind + "Version      : " + i.SoftVersion + "\n")
	_, _ = str.WriteString(ind + "ZWaveVersion : " + i.ZwaveVersion + "\n")
	_, _ = str.WriteString("\n")

	if l.ErrorMsg == "" {
		_, _ = str.WriteString("and logged in as:\n")
		_, _ = str.WriteString(ind + "User:         " + l.Username + "\n")
		_, _ = str.WriteString(ind + "Type:         " + l.Type + "\n")
	} else {
		_, _ = str.WriteString("couldn't log in with resp:\n")
		_, _ = str.WriteString(l.ErrorMsg)
	}

	return str.String()
}

// WriteInitConfigFile writes the configuration to a file. Paths required will be created if not presend.
func (f *FibaroHc2) WriteInitConfigFile(path string) (bytesWrote int) {
	os.MkdirAll(filepath.Dir(path), os.ModePerm)

	file, err := os.Create(path)
	if err != nil {
		log.Panicf("Problem creating file %s; %v\n", path, err)
	}

	defer file.Close()

	b, _ := json.MarshalIndent(f.Config(), "", " ")
	n4, err := file.Write(b)
	if err != nil {
		log.Panicf("Problem writing file %s; %v\n", path, err)
	}

	return n4

}
