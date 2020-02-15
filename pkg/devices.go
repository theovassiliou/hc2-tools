package fibarohc2

// Hc2Device represents a device in the HC2 system. Can be encoded as JSON.
type Hc2Device struct {
	Enabled    bool     `json:"enabled"`
	ID         int      `json:"id"`
	Interfaces []string `json:"interfaces"`
	Name       string   `json:"name"`
	ParentID   int      `json:"parentId"`
	RoomID     int      `json:"roomID"`
	Type       string   `json:"type"`
	Visible    bool     `json:"visible"`
}
