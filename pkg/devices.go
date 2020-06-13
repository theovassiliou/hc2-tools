package fibarohc2

import "strconv"

// Hc2Device represents a device in the HC2 system. Can be encoded as JSON.
type Hc2Device struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	ID         int      `json:"id"`
	RoomID     int      `json:"roomID"`
	Interfaces []string `json:"interfaces"`
	ParentID   int      `json:"parentId"`
	Enabled    bool     `json:"enabled"`
	Visible    bool     `json:"visible"`
	Properties struct {
		Dead                interface{} `json:"dead"`
		Energy              interface{} `json:"energy"`
		Power               interface{} `json:"power"`
		Sat                 interface{} `json:"sat"`
		Bri                 interface{} `json:"bri"`
		Ct                  interface{} `json:"ct"`
		Hue                 interface{} `json:"hue"`
		On                  interface{} `json:"on"`
		Value               interface{} `json:"value"`
		CentralSceneSupport interface{} `json:"centralSceneSupport"`
	} `json:"properties"`
}
type Key struct {
	KeyAttribute []string `json:"keyAttributes"`
	KeyId        int      `json:"keyId"`
}

func (d Hc2Device) Implements(name string) bool {
	for _, iN := range d.Interfaces {
		if iN == name {
			return true
		}
	}
	return false
}

func (d Hc2Device) GetValue(value string) int {
	switch value {
	case "bri":
		v, _ := strconv.Atoi(d.Properties.Bri.(string))
		return v
	case "sat":
		v, _ := strconv.Atoi(d.Properties.Sat.(string))
		return v
	case "ct":
		v, _ := strconv.Atoi(d.Properties.Ct.(string))
		return v
	case "hue":
		v, _ := strconv.Atoi(d.Properties.Hue.(string))
		return v
	case "power":
		v, _ := strconv.Atoi(d.Properties.Power.(string))
		return v
	case "energy":
		v, _ := strconv.Atoi(d.Properties.Energy.(string))
		return v
	case "value":
		v, _ := strconv.Atoi(d.Properties.Value.(string))
		return v
	}
	return -1
}

// Don't forget to cast the interface type in case you are using it
