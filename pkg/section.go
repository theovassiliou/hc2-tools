package fibarohc2

// Hc2Section represents a section in the HC2 system. Can be encoded as JSON.
type Hc2Section struct {
	SectionID int    `json:"id"`
	Name      string `json:"name"`
}
