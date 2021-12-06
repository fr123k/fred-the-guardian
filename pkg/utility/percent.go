package utility

import "math"

type UInt uint

func (f UInt) SeventyFivePercent() uint {
	return uint(math.Round(float64(uint(f)) * float64(0.75)))
}
