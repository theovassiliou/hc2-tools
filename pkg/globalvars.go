package fibarohc2

type Hc2GlobalVar struct {
	Name       string   `json:"name"`
	Value      string   `json:"value"`
	ReadOnly   bool     `json:"readOnly"`
	IsEnum     bool     `json:"isEnum"`
	EnumValues []string `json:"enumValues"`
	Created    int      `json:"created"`
	Modified   int      `json:"modified"`
}
