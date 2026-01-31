package cmn

// generate:reset
type ResetableStruct3 struct {
	i     int
	str   string
	strP  *string
	s     []int
	m     map[string]string
	child *ResetableStruct3
}
