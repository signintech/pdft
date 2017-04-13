package core

//TmplJSON tmpl.json
type TmplJSON struct {
	Pdf    string      `json:"pdf"`
	Fonts  []FontJSON  `json:"fonts"`
	Fields []FieldJSON `json:"fields"`
}

//FontJSON font
type FontJSON struct {
	Font string `json:"font"`
	File string `json:"file"`
}

//FieldJSON field
type FieldJSON struct {
	Key    *string  `json:"key"`
	Font   *string  `json:"font"`
	Size   *int     `json:"size"`
	Page   *int     `json:"page"`
	X      *float64 `json:"x"`
	Y      *float64 `json:"y"`
	W      *float64 `json:"w"`
	H      *float64 `json:"h"`
	Align  *string  `json:"align"`  //LEFT , CENTER , RIGHT
	VAlign *string  `json:"valign"` //TOP,MIDDLE,BOTTOM
}
