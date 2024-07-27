package render

// FieldInfos FieldInfo slice
type FieldInfos []FieldInfo

func (f FieldInfos) toMap() map[string]FieldInfo {
	m := make(map[string]FieldInfo)
	for _, i := range f {
		m[i.Key] = i
	}
	return m
}

// FieldInfo field position
type FieldInfo struct {
	Key        string
	PageNum    int
	X, Y       float64
	W, H       float64
	Align      int
	Font       string
	Size       int
	IsWrapText bool
}
