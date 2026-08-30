package control

type DecisionFuzzy struct{}

type Strategy string

const (
	OFF      = "OFF"
	ELECTRIC = "ELECTRIC"
	HYBRID   = "HYBRID"
	FROM_CH  = "FROM_CH"

	HYBRID_DELTA_T    = 10.0
	HYBRID_MAX_DHW_T  = 45.0
	EFFICIENT_DELTA_T = 8.0

	MIN_PELLET_LEVEL = 25.0
)

type fuzzyStrategy struct {
	Electric     float64
	Hybrid       float64
	FromHotWater float64
}

func (DecisionFuzzy) DecideStrategy(
	dhwTemp float64,
	chTemp float64,
	pelletLevelPercent float64,
) Strategy {
	// Hard limits.
	if dhwTemp >= 70.0 {
		return ELECTRIC
	}

	// if pelletLevelPercent <= 25.0 {
	// 	return ELECTRIC
	// }

	deltaT := chTemp - dhwTemp

	s := fuzzyStrategyScores(dhwTemp, deltaT, pelletLevelPercent)

	switch {
	case s.FromHotWater >= s.Hybrid && s.FromHotWater >= s.Electric:
		return FROM_CH

	case s.Hybrid >= s.Electric:
		return HYBRID

	default:
		return ELECTRIC
	}
}

func fuzzyStrategyScores(
	dhwTemp float64,
	deltaT float64,
	pelletLevelPercent float64,
) fuzzyStrategy {

	// How much does DHW need heating?
	cold := leftShoulder(dhwTemp, 30.0, 42.0)
	warm := triangle(dhwTemp, 35.0, 45.0, 55.0)
	hot := rightShoulder(dhwTemp, 48.0, 60.0)

	// How useful is the heat available from CO?
	lowDelta := leftShoulder(deltaT, 0.0, 15.0)
	highDelta := rightShoulder(deltaT, 5.0, 20.0)

	// How much pellet is left
	lowPelletLevel := leftShoulder(pelletLevelPercent, 20.0, 50.0)
	mediumPelletLevel := triangle(pelletLevelPercent, 20.0, 50.0, 80.0)
	highPelletLevel := rightShoulder(pelletLevelPercent, 65.0, 85.0)

	return fuzzyStrategy{
		Electric: max(
			lowDelta,
			lowPelletLevel,
		),

		Hybrid: min(
			max(cold, warm),
			highDelta,
			max(mediumPelletLevel, highPelletLevel),
		),

		FromHotWater: min(
			highDelta,
			highPelletLevel,
			max(cold, warm, hot),
		),
	}
}

func triangle(x, left, center, right float64) float64 {
	if x <= left || x >= right {
		return 0.0
	}

	if x == center {
		return 1.0
	}

	if x < center {
		return (x - left) / (center - left)
	}

	return (right - x) / (right - center)
}

func leftShoulder(x, full, zero float64) float64 {
	if x <= full {
		return 1.0
	}

	if x >= zero {
		return 0.0
	}

	return (zero - x) / (zero - full)
}

func rightShoulder(x, zero, full float64) float64 {
	if x <= zero {
		return 0.0
	}

	if x >= full {
		return 1.0
	}

	return (x - zero) / (full - zero)
}

func min(values ...float64) float64 {
	result := values[0]

	for _, value := range values[1:] {
		if value < result {
			result = value
		}
	}

	return result
}

func max(values ...float64) float64 {
	result := values[0]

	for _, value := range values[1:] {
		if value > result {
			result = value
		}
	}

	return result
}
