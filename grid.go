package drumbeat

// GridRes is the resolution of the grid for the pattern
type GridRes string

const (
	One4  GridRes = "1/4"
	One8  GridRes = "1/8"
	One16 GridRes = "1/16"
	One32 GridRes = "1/32"
	One64 GridRes = "1/64"
	// TODO: add triplets
)
