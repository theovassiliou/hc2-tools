package main

import (
	"fmt"

	"github.com/amimof/huego"
)

const ip = "192.168.178.49"
const user = "J9drcGiExAsE6UcTuAnLcAz7h3rctWpGZezT1E16"
const hueRoomId = "4"

func main() {
	bridge := huego.New(ip, user)
	groups, _ := bridge.GetGroups()

	fmt.Printf("Found %d groups\n", len(groups))
	for _, group := range groups {
		fmt.Printf("group.Name=%v, id=%v\n", group.Name, group.ID)
	}

	scenes, _ := bridge.GetScenes()

	fmt.Printf("--- code snipped for hue scene selection ---\n")

	fmt.Printf("local hueRoomId = %v;\n", hueRoomId)
	fmt.Printf("local ip = \"%v\";\n", ip)
	fmt.Printf("local user = \"%v\";\n", user)

	fmt.Println("local sceneArray = {")
	for _, aScene := range scenes {
		if aScene.Group == hueRoomId {
			fmt.Printf("    		{name = \"%v\", scene = \"%v\" }, \n", aScene.Name, aScene.ID)
		}
	}
	fmt.Println("}")
}
