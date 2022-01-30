package fibarohc2

import (
	"bufio"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Possible values to the used as runconfig field. Others are prohibited
const (
	TriggerAndManual string = "TRIGGER_AND_MANUAL"
	ManualOnly       string = "MANUAL_ONLY"
	Disabled         string = "DISABLED"
)

// Hc2Scene represents a LUA scene of the FibaroHC2 system
type Hc2Scene struct {
	SceneID             int    `json:"id,omitempty"`
	Name                string `json:"name,omitempty"`
	RoomID              int    `json:"roomID,omitempty"`
	RunConfig           string `json:"runConfig,omitempty"`
	MaxRunningInstances int    `json:"maxRunningInstances,omitempty"`
	Lua                 string `json:"lua,omitempty"`
	Type                string `json:"type,omitempty"`
	Autostart           bool   `json:"autostart,omitempty"`
	IsLua               bool   `json:"isLua,omitempty"`
	Visible             bool   `json:"visible,omitempty"`
}

// Hc2DebugMessage represents a debug message of the FibaroHC2 system
type Hc2DebugMessage struct {
	Timestamp int64  `json:"timestamp,omitempty"`
	Type      string `json:"type,omitempty"`
	Txt       string `json:"txt,omitempty"`
}

// NewHc2Scene creates a new Hc2Scene object properly initialized
func NewHc2Scene() (scene Hc2Scene) {
	scene = Hc2Scene{}
	scene.SceneID = -1
	scene.Visible = true
	scene.RunConfig = ManualOnly
	scene.Type = "com.fibaro.luaScene"
	scene.MaxRunningInstances = 2
	scene.Autostart = false
	return
}

var m = map[string]*regexp.Regexp{
	"gitHookComment":       regexp.MustCompile(`(?s)--\[\[ FIBARO_GIT_HOOK(.*?)--\]\]`),
	"@sceneID":             regexp.MustCompile(`@sceneID=(.*)`),
	"@name":                regexp.MustCompile(`@name=\"(.*)\"`),
	"@maxRunningInstances": regexp.MustCompile(`@maxRunningInstances=(.*)`),
	"@roomID":              regexp.MustCompile(`@roomID=(.*)`),
	"@runConfig":           regexp.MustCompile(`@runConfig=(.*)`),
	"@type":                regexp.MustCompile(`@type=\"(.*)\"`),
	"@autostart":           regexp.MustCompile(`@autostart=(.*)`),
	"@isLua":               regexp.MustCompile(`@isLua=(.*)`),
	"@visible":             regexp.MustCompile(`@visble=(.*)`),
}

// Parse parses a file input, and extracts the scene information located in comments.
// Comments are supposed to have the following format
// 	--[[ FIBARO_GIT_HOOK
// 	@sceneID=203
// 	@name="A_Trial_2"
// 	@roomID=305
// 	@runConfig=TRIGGER_AND_MANUAL
// 	@maxRunningInstances=2
// 	--]]
// If there is no header included then SceneID will be set to -1, and lua will still contain the file content
// If there is a header and the appropriate tag is not included then the field will not be touched
func (scene *Hc2Scene) Parse(input []byte) {
	// in any case we set the Lua field to the content
	scene.Lua = string(input)

	// check whether we have a header in the []byte
	i := m["gitHookComment"].FindSubmatchIndex(input)
	if i == nil {
		scene.SceneID = -1
		// don't need to try to parse the rest
		return
	}
	sceneSetInt(&scene.SceneID, "@sceneID", input)
	if sceneSetInt(&scene.RoomID, "@roomID", input) == -1 {
		// TODO: set to default
		scene.RoomID = 0
	}
	sceneSetString(&scene.Name, "@name", input)
	if sceneSetInt(&scene.MaxRunningInstances, "@maxRunningInstances", input) == -1 {
		// TODO: set to default
		scene.MaxRunningInstances = 0
	}
	sceneSetString(&scene.RunConfig, "@runConfig", input)
	sceneSetString(&scene.Type, "@type", input)
	sceneSetBool(&scene.Autostart, "@autostart", input)
	sceneSetBool(&scene.IsLua, "@isLua", input)
	sceneSetBool(&scene.Visible, "@visible", input)

}

// ParseFile parses a file provided by it's path, and extracts the scene information located in comments, i.e.
// the FIBARO_GIT_HEADER. If trim set to true, existing FIBARO_GIT_HEADERs will be removed from the Lua field.
func (scene *Hc2Scene) ParseFile(path string, trim bool) {
	_, dat := readFile(path)
	scene.Parse(dat)
	if trim {
		scene.TrimLuaHeaders()
	}
}

// TrimLuaHeaders removes all FIBARO_GIT_HOOK headers
// from the Lua field and updates the field. If multiple
// headers are included all will be removed. No other fields of the
// Hc2Scene will be modified
func (scene *Hc2Scene) TrimLuaHeaders() {
	s := []byte(scene.Lua)
	i := m["gitHookComment"].FindSubmatchIndex(s)
	a := s
	for i != nil {
		a = a[:i[0]+copy(a[i[0]:], a[i[1]:])]
		i = m["gitHookComment"].FindSubmatchIndex(a)
	}
	scene.Lua = string(a)
}

// UpdateLuaHeader deletes all headers and appends a new one
// add the end of the Lua scene with all field set to the Hc2Scene values.
func (scene *Hc2Scene) UpdateLuaHeader() {
	scene.TrimLuaHeaders()
	if string(scene.Lua[len(scene.Lua)-1:]) != "\n" {
		scene.Lua += "\n"
	}
	scene.Lua += scene.ToComment()
}

// SanityCheck performs checks on the value of some fields and returns false
// if a rule is beeing violated, true otherwise
func (scene *Hc2Scene) SanityCheck() (pass bool) {
	pass = true

	// check on runConfig
	rc := scene.RunConfig
	pass = rc == "" || ((rc == ManualOnly) || (rc == TriggerAndManual) || (rc == Disabled)) && pass
	pass = pass && (scene.Type == "" || scene.Type == "com.fibaro.luaScene")

	return
}

// ToComment translates a Hc2Scene into a LUA comment and returns it as string.m
func (scene *Hc2Scene) ToComment() string {
	var str strings.Builder

	_, _ = str.WriteString("--[[ FIBARO_GIT_HOOK - DO NOT CHANGE AS IT WILL BE DISCARDED \n")
	_, _ = str.WriteString("@sceneID=" + strconv.Itoa(scene.SceneID) + "\n")
	_, _ = str.WriteString("@name=\"" + scene.Name + "\"\n")
	_, _ = str.WriteString("@roomID=" + strconv.Itoa(scene.RoomID) + "\n")
	_, _ = str.WriteString("@autostart=" + strconv.FormatBool(scene.Autostart) + "\n")
	_, _ = str.WriteString("@runConfig=" + scene.RunConfig + "\n")
	_, _ = str.WriteString("@maxRunningInstance=" + strconv.Itoa(scene.MaxRunningInstances) + "\n")
	_, _ = str.WriteString("@type=\"" + scene.Type + "\"\n")
	_, _ = str.WriteString("@isLua=" + strconv.FormatBool(scene.IsLua) + "\n")
	_, _ = str.WriteString("--]]\n")
	return str.String()
}

func sceneSetString(v *string, k string, i []byte) string {
	submatch := m[k].FindSubmatch(i)
	result := ""
	if len(submatch) > 0 {
		result = string(submatch[1])

	}
	*v = result
	return result
}

func sceneSetInt(v *int, k string, i []byte) int {
	submatch := m[k].FindSubmatch(i)
	result := -1
	if len(submatch) > 0 {
		i, _ := strconv.Atoi(string(submatch[1]))
		result = i
	}
	*v = result
	return result
}

func sceneSetBool(v *bool, k string, i []byte) bool {
	submatch := m[k].FindSubmatch(i)
	result := *v
	if len(submatch) > 0 {
		i, _ := strconv.ParseBool(string(submatch[1]))
		result = i
	}
	*v = result
	return result
}

// ReadFile opens a file provided by it's path and returns the number of lines read and the file
// contents as byte array.
func readFile(path string) (int, []byte) {
	dat, readErr := ioutil.ReadFile(path)

	if readErr != nil {
		log.Fatal(readErr)
	}

	file, openErr := os.Open(path)
	if openErr != nil {
		log.Fatal(openErr)
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
