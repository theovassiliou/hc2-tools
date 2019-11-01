package fibarohc2

// Hc2Room represents a room in the HC2 system. Can be encoded as JSON.
type Hc2Room struct {
	RoomID    int    `json:"id"`
	Name      string `json:"name"`
	SectionID int    `json:"sectionID"`
}
