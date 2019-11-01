package fibarohc2

import (
	"encoding/json"
	"testing"
)

func TestSceneStringTable(t *testing.T) {

	tables := []struct {
		input Hc2Scene
		want  string
	}{
		{
			Hc2Scene{
				Name: "A Script Name",
			},
			"--[[ FIBARO_GIT_HOOK - DO NOT CHANGE AS IT WILL BE DISCARDED \n" +
				"@sceneID=0\n" +
				"@name=\"A Script Name\"\n" +
				"@roomID=0\n" +
				"@autostart=false\n" +
				"@runConfig=\n" +
				"@maxRunningInstance=0\n" +
				"@type=\"\"\n" +
				"@isLua=false\n" +
				"--]]\n"},
		{
			Hc2Scene{
				SceneID: 22,
			},
			"--[[ FIBARO_GIT_HOOK - DO NOT CHANGE AS IT WILL BE DISCARDED \n" +
				"@sceneID=22\n" +
				"@name=\"\"\n" +
				"@roomID=0\n" +
				"@autostart=false\n" +
				"@runConfig=\n" +
				"@maxRunningInstance=0\n" +
				"@type=\"\"\n" +
				"@isLua=false\n" +
				"--]]\n"},
	}

	for _, tc := range tables {
		// fmt.Println(tc.input)
		// fmt.Print("")
		AssertEqual(t, tc.input.ToComment(), tc.want)
	}

}

func TestSceneMarshal(t *testing.T) {
	tables := []struct {
		input Hc2Scene
		want  string
	}{
		{
			Hc2Scene{
				Name: "A Script Name",
			},
			`{"name":"A Script Name"}`},
		{
			Hc2Scene{
				SceneID: 22,
			},
			`{"id":22}`},
	}

	for _, tc := range tables {
		s, _ := json.Marshal(tc.input)
		AssertEqual(t, string(s), tc.want)
	}
}

func TestSceneUnmarshal(t *testing.T) {
	data := []byte(`
		{
			"id":44,
			"name":"A Script Name",
			"roomID":0,
			"autostart":false,
			"runConfig":"TRIGGER_AND_RUN",
			"maxRunningInstances":0
		}
		`)
	var aScene Hc2Scene
	json.Unmarshal(data, &aScene)
	AssertEqual(t, aScene.ToComment(),
		`--[[ FIBARO_GIT_HOOK - DO NOT CHANGE AS IT WILL BE DISCARDED 
@sceneID=44
@name="A Script Name"
@roomID=0
@autostart=false
@runConfig=TRIGGER_AND_RUN
@maxRunningInstance=0
@type=""
@isLua=false
--]]
`)
}

func TestHc2Scene_ToComment(t *testing.T) {
	type fields struct {
		SceneID             int
		Name                string
		RoomID              int
		RunConfig           string
		MaxRunningInstances int
		Lua                 string
		Type                string
		Autostart           bool
		IsLua               bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Valid Empty",
			fields{},
			`--[[ FIBARO_GIT_HOOK - DO NOT CHANGE AS IT WILL BE DISCARDED 
@sceneID=0
@name=""
@roomID=0
@autostart=false
@runConfig=
@maxRunningInstance=0
@type=""
@isLua=false
--]]
`},
		{"Valid FullStructAllTrue",
			fields{
				SceneID:             22,
				Name:                "Marvel",
				RoomID:              5,
				RunConfig:           ManualOnly,
				MaxRunningInstances: 2,
				Lua:                 "LuaScript",
				Type:                "com.script",
				Autostart:           true,
				IsLua:               true,
			},
			`--[[ FIBARO_GIT_HOOK - DO NOT CHANGE AS IT WILL BE DISCARDED 
@sceneID=22
@name="Marvel"
@roomID=5
@autostart=true
@runConfig=MANUAL_ONLY
@maxRunningInstance=2
@type="com.script"
@isLua=true
--]]
`},
		{"Valid FullStructAllFalse",
			fields{
				SceneID:             22,
				Name:                "Marvel",
				RoomID:              5,
				RunConfig:           ManualOnly,
				MaxRunningInstances: 2,
				Lua:                 "LuaScript",
				Type:                "com.script",
				Autostart:           false,
				IsLua:               false,
			},
			`--[[ FIBARO_GIT_HOOK - DO NOT CHANGE AS IT WILL BE DISCARDED 
@sceneID=22
@name="Marvel"
@roomID=5
@autostart=false
@runConfig=MANUAL_ONLY
@maxRunningInstance=2
@type="com.script"
@isLua=false
--]]
`},
		{"Valid FullStructAllNegate",
			fields{
				SceneID:             -1,
				Name:                "Marvel",
				RoomID:              -1,
				RunConfig:           ManualOnly,
				MaxRunningInstances: -1,
				Lua:                 "LuaScript",
				Type:                "com.script",
				Autostart:           true,
				IsLua:               true,
			},
			`--[[ FIBARO_GIT_HOOK - DO NOT CHANGE AS IT WILL BE DISCARDED 
@sceneID=-1
@name="Marvel"
@roomID=-1
@autostart=true
@runConfig=MANUAL_ONLY
@maxRunningInstance=-1
@type="com.script"
@isLua=true
--]]
`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scene := &Hc2Scene{
				SceneID:             tt.fields.SceneID,
				Name:                tt.fields.Name,
				RoomID:              tt.fields.RoomID,
				RunConfig:           tt.fields.RunConfig,
				MaxRunningInstances: tt.fields.MaxRunningInstances,
				Lua:                 tt.fields.Lua,
				Type:                tt.fields.Type,
				Autostart:           tt.fields.Autostart,
				IsLua:               tt.fields.IsLua,
			}
			if got := scene.ToComment(); got != tt.want {
				t.Errorf("Hc2Scene.ToComment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHc2Scene_ParseFile(t *testing.T) {

	type args struct {
		path string
	}
	tests := []struct {
		name   string
		fields Hc2Scene
		args   args
	}{
		{
			"Reading shortHeader.lua",
			Hc2Scene{
				SceneID:   203,
				RoomID:    305,
				Name:      "A_Trial_2",
				Lua:       "../test/shortHeader.lua",
				RunConfig: "TRIGGER_AND_MANUAL",
			},
			args{"../test/shortHeader.lua"},
		},
		{
			"Reading shortHeader2.lua",
			Hc2Scene{
				SceneID:   205,
				RoomID:    305,
				Name:      "A_Trial_2",
				Lua:       "../test/shortHeader2.lua",
				RunConfig: "MANUAL",
			},
			args{"../test/shortHeader2.lua"},
		},
		{
			"Reading shortHeader3.lua",
			Hc2Scene{
				SceneID:   205,
				RoomID:    305,
				Name:      "A_Trial_2",
				Lua:       "../test/shortHeader3.lua",
				RunConfig: "MANUAL_ONLY",
			},
			args{"../test/shortHeader3.lua"},
		},
		{
			"Reading SZAllLightsOff.lua, no header",
			Hc2Scene{
				SceneID: -1,
				Lua:     "../test/SZAllLightsOff.lua",
			},
			args{"../test/SZAllLightsOff.lua"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scene := &Hc2Scene{}
			scene.ParseFile(tt.args.path, false)
			AssertHc2SceneEqual(t, *scene, Hc2Scene(tt.fields))
		})
	}
}

func Test_sceneSetInt(t *testing.T) {
	var i int
	type args struct {
		v *int
		k string
		i []byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"int 1",
			args{
				&i,
				"@sceneID",
				[]byte(`@sceneID=22`),
			},
			22,
		},
		{
			"int 2",
			args{
				&i,
				"@sceneID",
				[]byte(`@sceneID=0`),
			},
			0,
		},
		{
			"int 3",
			args{
				&i,
				"@sceneID",
				[]byte(`@sceneID=-1`),
			},
			-1,
		},
		{
			"int inv",
			args{
				&i,
				"@sceneID",
				[]byte(`@sceneID=hallo`),
			},
			0,
		},
		{
			"int inv2",
			args{
				&i,
				"@sceneID",
				[]byte(`@sceneD=22`),
			},
			-1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sceneSetInt(tt.args.v, tt.args.k, tt.args.i); got != tt.want {
				t.Errorf("sceneSetInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sceneSetBool(t *testing.T) {
	var b bool
	type args struct {
		v *bool
		k string
		i []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"bool 1",
			args{
				&b,
				"@autostart",
				[]byte(`@autostart=true`),
			},
			true,
		},
		{
			"bool 2",
			args{
				&b,
				"@autostart",
				[]byte(`@autostart=false`),
			},
			false,
		},
		{
			"bool inv",
			args{
				&b,
				"@autostart",
				[]byte(`@autostart=xx`),
			},
			false,
		},
		{
			"bool inv2",
			args{
				&b,
				"@autostart",
				[]byte(`@autos=true`),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sceneSetBool(tt.args.v, tt.args.k, tt.args.i); got != tt.want {
				t.Errorf("sceneSetBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sceneSetString(t *testing.T) {
	var s string
	type args struct {
		v *string
		k string
		i []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"string 1",
			args{
				&s,
				"@name",
				[]byte(`@name=""`),
			},
			"",
		},
		{
			"string 2",
			args{
				&s,
				"@name",
				[]byte(`@name="hello"`),
			},
			"hello",
		},
		{
			"string 3",
			args{
				&s,
				"@name",
				[]byte(`@name="-1"`),
			},
			"-1",
		},
		{
			"string 4",
			args{
				&s,
				"@name",
				[]byte(`@name="This is a test"`),
			},
			"This is a test",
		},
		{
			"string inv1",
			args{
				&s,
				"@name",
				[]byte(`@name="`),
			},
			"",
		},
		{
			"string inv2",
			args{
				&s,
				"@name",
				[]byte(`@name=22`),
			},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sceneSetString(tt.args.v, tt.args.k, tt.args.i); got != tt.want {
				t.Errorf("sceneSetString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"reading shortHeader.lua",
			args{
				"../test/shortHeader.lua",
			},
			16,
		},
		// TODO: Add test cases for non existing paths.
		// The problem with this is that the application terminates.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := readFile(tt.args.path)
			if got != tt.want {
				t.Errorf("ReadFile() got = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestLoadConfigInvalid(t *testing.T) {
	f := NewFibaroHc2Config("../configs/unknwon.json.example")
	if f != nil {
		t.Fail()
	}
}
func TestLoadConfig(t *testing.T) {
	f := NewFibaroHc2Config("../configs/config.json.example")
	config := f.Config()
	AssertEqual(t, config.Username, "admin")
	AssertEqual(t, config.Password, "theAdminPassword")
	AssertEqual(t, config.BaseURL, "http://hc.2.ip.address")
}

func TestParseMultiHeaderScene(t *testing.T) {
	_, s := readFile("../test/doubleHeader.lua")
	t.Run("Parsing double headers", func(t *testing.T) {
		scene1 := NewHc2Scene()
		scene1.Parse(s)
		AssertHc2SceneWOLuaEqual(t, scene1, Hc2Scene{
			Name:      "A_Trial_2",
			SceneID:   203,
			RoomID:    306,
			RunConfig: TriggerAndManual,
		})
	})

	i := m["gitHookComment"].FindIndex(s)
	if i == nil {
		t.Errorf("m[\"gitHookComment\"] returned  %v want non-nil", i)
	}
	t.Run("Parsing modified header", func(t *testing.T) {

		if i != nil {
			// This modifies the original slice (s)
			a := s[:i[0]+copy(s[i[0]:], s[i[1]:])]

			scene2 := NewHc2Scene()
			scene2.Parse(a)
			AssertHc2SceneWOLuaEqual(t, scene2, Hc2Scene{
				Name:      "A_Trial_4",
				SceneID:   204,
				RoomID:    305,
				RunConfig: TriggerAndManual,
			})
		} else {
			t.FailNow()
		}
	})
}

func TestParseTrimLuaHeader(t *testing.T) {
	scene := NewHc2Scene()
	scene.ParseFile("../test/doubleHeader.lua", false)
	scene.TrimLuaHeaders()
	scene2 := NewHc2Scene()
	scene2.Parse([]byte(scene.Lua))
	AssertHc2SceneWOLuaEqual(t, scene2, NewHc2Scene())

	scene = NewHc2Scene()
	scene.ParseFile("../test/doubleHeader.lua", true)
	AssertEqual(t, scene.Lua, scene2.Lua)
}

func TestHc2Scene_SanityCheck(t *testing.T) {
	type fields struct {
		SceneID             int
		Name                string
		RoomID              int
		RunConfig           string
		MaxRunningInstances int
		Lua                 string
		Type                string
		Autostart           bool
		IsLua               bool
		Visible             bool
	}
	tests := []struct {
		name     string
		fields   fields
		wantPass bool
	}{
		{
			"Empty Scene",
			fields{},
			true,
		},
		{
			"Only SceneID",
			fields{SceneID: 2},
			true,
		},
		{
			"Wrong RunningConfig",
			fields{RunConfig: "HALLO"},
			false,
		},
		{
			"Wrong type",
			fields{Type: "com."},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scene := &Hc2Scene{
				SceneID:             tt.fields.SceneID,
				Name:                tt.fields.Name,
				RoomID:              tt.fields.RoomID,
				RunConfig:           tt.fields.RunConfig,
				MaxRunningInstances: tt.fields.MaxRunningInstances,
				Lua:                 tt.fields.Lua,
				Type:                tt.fields.Type,
				Autostart:           tt.fields.Autostart,
				IsLua:               tt.fields.IsLua,
				Visible:             tt.fields.Visible,
			}
			if gotPass := scene.SanityCheck(); gotPass != tt.wantPass {
				t.Errorf("Hc2Scene.SanityCheck() = %v, want %v", gotPass, tt.wantPass)
			}
		})
	}
}
