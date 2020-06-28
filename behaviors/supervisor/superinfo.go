package supervisor

type SuperInfo struct {
	RegDC     string
	NameSpace string
	EndPoints []string
	IsPub     bool
}

func NewSuperInfo() *SuperInfo {
	info := new(SuperInfo)
	return info
}
