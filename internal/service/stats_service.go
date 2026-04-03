package service

import (
	"math"
)

type StatsResult struct {
	CR_A        float64 `json:"cr_a"`
	CR_B        float64 `json:"cr_b"`
	Uplift      float64 `json:"uplift"`
	ZScore      float64 `json:"z_score"`
	PValue      float64 `json:"p_value"`
	Significant bool    `json:"significant"`
}

func CalculateStats(imprA, convA, imprB, convB int) StatsResult {
	if imprA == 0 || imprB == 0 {
		return StatsResult{
			CR_A:        0,
			CR_B:        0,
			Uplift:      0,
			ZScore:      0,
			PValue:      1,
			Significant: false,
		}
	}

	p1 := float64(convA) / float64(imprA)
	p2 := float64(convB) / float64(imprB)

	pPool := float64(convA+convB) / float64(imprA+imprB)

	se := math.Sqrt(pPool * (1 - pPool) * (1/float64(imprA) + 1/float64(imprB)))

	if se == 0 {
		return StatsResult{}
	}

	z := (p2 - p1) / se

	pValue := 2 * (1 - normalCDF(math.Abs(z)))

	uplift := 0.0
	if p1 > 0 {
		uplift = (p2 - p1) / p1
	}

	return StatsResult{
		CR_A:        p1,
		CR_B:        p2,
		Uplift:      uplift,
		ZScore:      z,
		PValue:      pValue,
		Significant: pValue < 0.05,
	}
}
func normalCDF(x float64) float64 {
	return 0.5 * (1 + math.Erf(x/math.Sqrt2))
}