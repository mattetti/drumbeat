package drumbeat

// GridRes is the resolution of the grid for the pattern
type GridRes string

// StepsInBeat returns the number of steps to fill a beat
func (g GridRes) StepsInBeat() uint64 {
	switch g {
	case One8:
		return 2
	case One16:
		return 4
	case One32:
		return 8
	case One64:
		return 16
	}
	return 1
}

const (
	One4  GridRes = "1/4"
	One8  GridRes = "1/8"
	One16 GridRes = "1/16"
	One32 GridRes = "1/32"
	One64 GridRes = "1/64"
	// TODO: add triplets
)

// StepSize returns the size of a pattern step in ticks given its grid resolution
func (p *Pattern) StepSize() uint64 {
	switch p.Grid {
	case One4:
		return uint64(p.PPQN)
	case One8:
		return uint64(p.PPQN / 2)
	case One16:
		return uint64(p.PPQN / 4)
	case One32:
		return uint64(p.PPQN / 8)
	case One64:
		return uint64(p.PPQN / 16)
	}
	return 0
}
