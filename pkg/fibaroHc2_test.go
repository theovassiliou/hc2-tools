package fibarohc2

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
)

const ConfigFileName string = "../configs/configTest.json"

func TestFibaroHc2_OneScene(t *testing.T) {
	cfg := NewFibaroHc2Config(ConfigFileName).Config()
	httpmock.ActivateNonDefault(cfg.client.GetClient())
	defer httpmock.DeactivateAndReset()

	type args struct {
		sceneID      int
		fixtureName  string
		responseCode int
		method       string
		route        string
	}
	tests := []struct {
		name string
		args args
		want Hc2Scene
	}{
		{
			"get one short scene",
			args{
				55,
				"../test/scene55.json",
				200,
				"GET",
				"http://192.10.66.55/api/scenes/55",
			},
			Hc2Scene{
				SceneID:             55,
				Name:                "PressST_VD_Button",
				RoomID:              0,
				RunConfig:           ManualOnly,
				MaxRunningInstances: 2,
				Type:                "com.fibaro.luaScene",
				Autostart:           false,
				IsLua:               true,
			},
		},
		{
			"get one long scene",
			args{
				146,
				"../test/scene146.json",
				200,
				"GET",
				"http://192.10.66.55/api/scenes/146",
			},
			Hc2Scene{
				SceneID:             146,
				Name:                "VSLSchlafzimmer",
				RoomID:              5,
				RunConfig:           TriggerAndManual,
				MaxRunningInstances: 2,
				Type:                "com.fibaro.luaScene",
				Autostart:           false,
				IsLua:               true,
			},
		},
		{
			"get one non existing scene",
			args{
				147,
				"",
				404,
				"GET",
				"http://192.10.66.55/api/scenes/147",
			},
			Hc2Scene{
				SceneID: -1,
			},
		},
	}
	for _, tt := range tests {
		fixture, _ := ioutil.ReadFile(tt.args.fixtureName)
		responder := httpmock.NewBytesResponder(tt.args.responseCode, fixture)
		httpmock.RegisterResponder(tt.args.method, tt.args.route, responder)

		t.Run(tt.name, func(t *testing.T) {
			f := &FibaroHc2{
				cfg: *cfg,
			}
			got := f.OneScene(tt.args.sceneID)
			AssertHc2SceneWOLuaEqual(t, got, tt.want)
		})
	}
}

func testFibaroHc2OneSceneParametrized(t *testing.T) {
	configuration := NewFibaroHc2Config(ConfigFileName).Config()
	httpmock.ActivateNonDefault(configuration.client.GetClient())
	defer httpmock.DeactivateAndReset()

	fixture, _ := ioutil.ReadFile("../test/scene55.json")
	responder := httpmock.NewBytesResponder(200, fixture)
	httpmock.RegisterResponder("GET", "http://192.10.66.55/api/scenes/55", responder)

	f := &FibaroHc2{
		cfg: *configuration,
	}

	got := f.OneScene(55)
	log.Println(got)
	t.Fail()
}

func testFibaroHc2OneSceneParametrized2(t *testing.T) {
	configuration := NewFibaroHc2Config(ConfigFileName).Config()
	httpmock.ActivateNonDefault(configuration.client.GetClient())
	defer httpmock.DeactivateAndReset()

	fixture, _ := ioutil.ReadFile("../test/scene105NoHeader.json")
	responder := httpmock.NewBytesResponder(200, fixture)
	httpmock.RegisterResponder("GET", "http://192.10.66.55/api/scenes/105", responder)

	f := &FibaroHc2{
		cfg: *configuration,
	}

	got := f.OneScene(105)
	log.Println(got)
	t.Fail()
}

func TestFibaroHc2_AllScenes(t *testing.T) {
	cfg := NewFibaroHc2Config(ConfigFileName).Config()
	httpmock.ActivateNonDefault(cfg.client.GetClient())
	defer httpmock.DeactivateAndReset()

	fixture, _ := ioutil.ReadFile("../test/scenes.json")
	responder := httpmock.NewBytesResponder(200, fixture)
	httpmock.RegisterResponder("GET", "http://192.10.66.55/api/scenes", responder)

	f := &FibaroHc2{
		cfg: *cfg,
	}
	got := f.AllScenes()
	AssertEqual(t, len(got), 93)
	AssertHc2SceneWOLuaEqual(t, got[1], Hc2Scene{
		SceneID:             19,
		Name:                "VSLFlur",
		RoomID:              10,
		RunConfig:           TriggerAndManual,
		MaxRunningInstances: 2,
		Type:                "com.fibaro.luaScene",
		Autostart:           false,
		IsLua:               true,
	})
}

func TestFibaroHc2_OneRoom(t *testing.T) {
	cfg := NewFibaroHc2Config(ConfigFileName).Config()
	httpmock.ActivateNonDefault(cfg.client.GetClient())
	defer httpmock.DeactivateAndReset()

	type args struct {
		roomID       int
		fixtureName  string
		responseCode int
		method       string
		route        string
	}
	tests := []struct {
		name string
		args args
		want Hc2Room
	}{
		{
			"get one room",
			args{
				5,
				"../test/room5.json",
				200,
				"GET",
				"http://192.10.66.55/api/rooms/5",
			},
			Hc2Room{
				Name:      "Schlafzimmer",
				RoomID:    5,
				SectionID: 4,
			},
		},
		{
			"get nonexisting room",
			args{
				6,
				"",
				404,
				"GET",
				"http://192.10.66.55/api/rooms/6",
			},
			Hc2Room{
				RoomID: -1,
			},
		},
	}
	for _, tt := range tests {

		fixture, _ := ioutil.ReadFile(tt.args.fixtureName)
		responder := httpmock.NewBytesResponder(tt.args.responseCode, fixture)
		httpmock.RegisterResponder(tt.args.method, tt.args.route, responder)

		t.Run(tt.name, func(t *testing.T) {
			f := &FibaroHc2{
				cfg: *cfg,
			}
			if got := f.OneRoom(tt.args.roomID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FibaroHc2.OneRoom() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFibaroHc2_OneSection(t *testing.T) {
	cfg := NewFibaroHc2Config(ConfigFileName).Config()
	httpmock.ActivateNonDefault(cfg.client.GetClient())
	defer httpmock.DeactivateAndReset()

	type args struct {
		sectionID    int
		fixtureName  string
		responseCode int
		method       string
		route        string
	}
	tests := []struct {
		name string
		args args
		want Hc2Section
	}{
		{
			"get one section",
			args{
				4,
				"../test/section4.json",
				200,
				"GET",
				"http://192.10.66.55/api/sections/4",
			},
			Hc2Section{
				Name:      "NickyTheo",
				SectionID: 4,
			},
		},
		{
			"get nonexisting section",
			args{
				5,
				"",
				404,
				"GET",
				"http://192.10.66.55/api/sections/5",
			},
			Hc2Section{
				SectionID: -1,
			},
		},
	}
	for _, tt := range tests {

		fixture, _ := ioutil.ReadFile(tt.args.fixtureName)
		responder := httpmock.NewBytesResponder(tt.args.responseCode, fixture)
		httpmock.RegisterResponder(tt.args.method, tt.args.route, responder)

		t.Run(tt.name, func(t *testing.T) {
			f := &FibaroHc2{
				cfg: *cfg,
			}
			if got := f.OneSection(tt.args.sectionID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FibaroHc2.OneSection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFibaroHc2_PutOneScene(t *testing.T) {
	cfg := NewFibaroHc2Config(ConfigFileName).Config()
	httpmock.ActivateNonDefault(cfg.client.GetClient())
	defer httpmock.DeactivateAndReset()

	type args struct {
		sceneFile string

		fixtureName  string
		responseCode int
		method       string
		route        string
	}
	tests := []struct {
		name     string
		args     args
		wantResp string
		wantErr  bool
	}{
		{
			"put one scene (assuming existing)",
			args{
				"../test/shortHeader.lua",
				"",
				200,
				http.MethodPut,
				"http://192.10.66.55/api/scenes/203",
			},
			"",
			false,
		},
		{
			"put one scene (wrong runConfig)",
			args{
				"../test/shortHeader2.lua",
				"../test/scene205_404.json",
				404,
				http.MethodPut,
				"http://192.10.66.55/api/scenes/205",
			},
			`{"type":"ERROR","reason":"Invalid runConfig parameter","message":"Invalid runConfig parameter"}`,
			true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FibaroHc2{
				cfg: *cfg,
			}

			fixture, _ := ioutil.ReadFile(tt.args.fixtureName)
			responder := httpmock.NewBytesResponder(tt.args.responseCode, fixture)
			httpmock.RegisterResponder(tt.args.method, tt.args.route, responder)
			var scene Hc2Scene
			scene.ParseFile(tt.args.sceneFile, false)
			gotResp, err := f.PutOneScene(scene)
			if (err != nil) && tt.wantErr {
				return
			} else if (err != nil) && !tt.wantErr {
				t.Errorf("FibaroHHc2.PutOneScene() produced error = %v", err)
				return
			}
			var i, j interface{}
			if !reflect.DeepEqual(json.Unmarshal([]byte(gotResp.String()), &i), json.Unmarshal([]byte(tt.wantResp), &j)) {
				t.Errorf("FibaroHc2.PutOneScene() = %v, want %v", gotResp.String(), tt.wantResp)
			}
		})
	}
}

func TestFibaroHc2_SetConfig(t *testing.T) {
	cfg := NewFibaroHc2Config(ConfigFileName).Config()
	f := FibaroHc2{}
	f.SetConfig(*cfg)
	if got := f.Config(); !reflect.DeepEqual(got, cfg) {
		t.Errorf("FibaroHc2.Config() = %v, want %v", got, cfg)
	}
}

func TestOneSceneAndUpdate(t *testing.T) {
	cfg := NewFibaroHc2Config(ConfigFileName).Config()
	httpmock.ActivateNonDefault(cfg.client.GetClient())
	defer httpmock.DeactivateAndReset()

	fixture, _ := ioutil.ReadFile("../test/scene55InvHeader.json")
	responder := httpmock.NewBytesResponder(200, fixture)
	httpmock.RegisterResponder(http.MethodGet, "http://192.10.66.55/api/scenes/55", responder)

	f := &FibaroHc2{
		cfg: *cfg,
	}
	gotScene := f.OneScene(55)
	secondScene := Hc2Scene{}
	secondScene.Parse([]byte(gotScene.Lua))
	if gotScene.SceneID == secondScene.SceneID {
		t.Errorf("scenes are supposed to have different sceneIDs = %v, want %v", gotScene, secondScene)
	}
	gotScene.UpdateLuaHeader()
	secondScene.Parse([]byte(gotScene.Lua))
	if gotScene.SceneID != secondScene.SceneID {
		t.Errorf("scenes are supposed to have same sceneIDs = %v, want %v", gotScene, secondScene)
	}
}

func TestCreateNewFile(t *testing.T) {
	cfg := NewFibaroHc2Config(ConfigFileName).Config()
	httpmock.ActivateNonDefault(cfg.client.GetClient())
	defer httpmock.DeactivateAndReset()

	fixture, _ := ioutil.ReadFile("../test/postSceneResponse.json")
	responderPost := httpmock.NewBytesResponder(200, fixture)
	responderPut := httpmock.NewBytesResponder(200, []byte(""))
	httpmock.RegisterResponder(http.MethodPost, "/api/scenes", responderPost)
	httpmock.RegisterResponder(http.MethodPut, "/api/scenes/206", responderPut)

	fixtureResponse, _ := ioutil.ReadFile("../test/scene206.json")
	responderGet := httpmock.NewBytesResponder(200, fixtureResponse)
	httpmock.RegisterResponder(http.MethodGet, "/api/scenes/206", responderGet)

	f := &FibaroHc2{
		cfg: *cfg,
	}

	// We have a file
	sceneFile, _ := ioutil.ReadFile("../test/shortHeaderNew.lua")

	// we check for header by
	// parsing the file and checking the sceneID to -1
	loadedScene := Hc2Scene{
		Visible: true,
	}
	loadedScene.Parse(sceneFile)

	// We check for sceneID == -1
	if loadedScene.SceneID != -1 {
		t.Errorf("Where expected no header, i.e. SceneID == -1, but was %v", loadedScene.SceneID)
	}

	if loadedScene.Name == "" {
		// TODO: Move "unnamed" to some kind of constant repo
		loadedScene.Name = "unnamed"
	}

	// as an effect when parsing we have all supported fields filled
	// including the luaField

	// now we can CreateScene()
	sceneIDCreated := f.CreateScene(loadedScene)
	if sceneIDCreated == -1 {
		t.Errorf("Scene creation failed, loaded scene was  %v", loadedScene)
	}

	sceneCreated := f.OneScene(sceneIDCreated)
	want := Hc2Scene{IsLua: true, Name: loadedScene.Name, SceneID: sceneIDCreated, Type: "com.fibaro.luaScene", Visible: true, RoomID: 305, RunConfig: TriggerAndManual, MaxRunningInstances: 4}
	AssertHc2SceneWOLuaEqual(t, sceneCreated, want)

}

// TODO: Add test cases for action and test.
