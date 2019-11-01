package fibarohc2

// Hc2Network represents the network information of the HC2 system. Can be encoded as JSON.
type Hc2Network struct {
	Dhcp                bool   `json:"dhcp"`
	IP                  string `json:"ip"`
	Mask                string `json:"mask"`
	Gateway             string `json:"gateway"`
	DNS                 string `json:"dns"`
	RemoteAccess        bool   `json:"remoteAccess"`
	RemoteAccessSupport int    `json:"remoteAccessSupport"`
}

// Hc2Info represents the general information of the HC2 system. Can be encoded as JSON.
type Hc2Info struct {
	SerialNumber string `json:"serialNumber,omitempty"`
	HcName       string `json:"hcName,omitempty"`
	Mac          string `json:"mac,omitempty"`
	SoftVersion  string `json:"softVersion,omitempty"`
	Beta         bool   `json:"beta,omitempty"`
	ZwaveVersion string `json:"zwaveVersion,omitempty"`
	ZwaveRegion  string `json:"zwaveRegion,omitempty"`
	ServerStatus int    `json:"serverStatus,omitempty"`
	// "defaultLanguage": "",
	// "sunsetHour": "",
	// "sunriseHour": "",
	// "hotelMode": "bool",
	// "updateStableAvailable": "bool",
	// "temperatureUnit": "",
	// "updateBetaAvailable": "bool",
	// "batteryLowNotification": "bool",
	// "smsManagement": "bool"
}

// Hc2LoginStatus represents information on the logged user. Can be encoded as JSON.
type Hc2LoginStatus struct {
	Status   bool   `json:"status,omitempty"`
	UserID   int    `json:"userID,omitempty"`
	Username string `json:"username,omitempty"`
	Type     string `json:"type,omitempty"`
	ErrorMsg string
}
