package main

import (
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	hc2 "github.com/theovassiliou/hc2-tools/pkg"
)

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
