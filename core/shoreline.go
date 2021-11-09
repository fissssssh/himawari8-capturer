package core

type Shoreline uint8

const (
	Ignore Shoreline = iota
	Red
	Green
	Yellow
)

func (c Shoreline) String() string {
	switch c {
	case Ignore:
		return ""
	case Red:
		return "ff0000"
	case Green:
		return "00ff00"
	case Yellow:
		return "ffff00"
	default:
		return ""
	}
}
