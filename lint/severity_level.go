package lint

type SeverityLevel int

var Unknown SeverityLevel = 0
var Error SeverityLevel = 1
var Warning SeverityLevel = 2
var Information SeverityLevel = 3

func (sl SeverityLevel) String() string {
	switch sl {
	case Error:
		return "Error"
	case Warning:
		return "Warning"
	case Information:
		return "Information"
	default:
		return "Unknown"
	}
}
