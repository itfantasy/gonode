package supervisor

type SuperInfo struct {
	Url      string
	RegComp  string
	LogComp  string
	MoniComp string
}

func NewSuperInfo() *SuperInfo {
	info := new(SuperInfo)
	return info
}
