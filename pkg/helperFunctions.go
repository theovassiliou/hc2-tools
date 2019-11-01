package fibarohc2

import (
	"reflect"
	"testing"
)

// AssertEqual checks if values are equal
func AssertEqual(t *testing.T, a interface{}, b interface{}) {

	if a == b {
		return
	}
	// debug.PrintStack()
	t.Errorf("Received %v (type %v), expected %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}

// AssertHc2SceneEqual checks if two Hc2Scenes are equal, provided that the wanted
// Lua field contains a filename instead of a the actual lua script
func AssertHc2SceneEqual(t *testing.T, got Hc2Scene, want Hc2Scene) {
	AssertEqual(t, got.SceneID, want.SceneID)
	AssertEqual(t, got.Name, want.Name)
	AssertEqual(t, got.RoomID, want.RoomID)
	AssertEqual(t, got.RunConfig, want.RunConfig)
	AssertEqual(t, got.MaxRunningInstances, want.MaxRunningInstances)
	_, luaFileContents := readFile(want.Lua)
	AssertEqual(t, got.Lua, string(luaFileContents))
	AssertEqual(t, got.Type, want.Type)
	AssertEqual(t, got.Autostart, want.Autostart)
	AssertEqual(t, got.IsLua, want.IsLua)
}

// AssertHc2SceneWOLuaEqual checks if two Hc2Scenes are equal, without looking into the lua field
func AssertHc2SceneWOLuaEqual(t *testing.T, got Hc2Scene, want Hc2Scene) {
	AssertEqual(t, got.SceneID, want.SceneID)
	AssertEqual(t, got.Name, want.Name)
	AssertEqual(t, got.RoomID, want.RoomID)
	AssertEqual(t, got.RunConfig, want.RunConfig)
	AssertEqual(t, got.MaxRunningInstances, want.MaxRunningInstances)
	AssertEqual(t, got.Type, want.Type)
	AssertEqual(t, got.Autostart, want.Autostart)
	AssertEqual(t, got.IsLua, want.IsLua)
}
