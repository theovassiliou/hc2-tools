package main

import "fmt"

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
