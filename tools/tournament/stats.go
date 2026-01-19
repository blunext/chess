package main

import (
	"math"
)

// EloDiff calculates Elo difference and 95% confidence interval
func EloDiff(wins, draws, losses int) (eloDiff, eloError float64) {
	total := float64(wins + draws + losses)
	if total == 0 {
		return 0, 0
	}

	// Score from engine1's perspective
	score := (float64(wins) + 0.5*float64(draws)) / total

	// Elo difference
	if score <= 0 || score >= 1 {
		if score >= 1 {
			return 800, 0 // Perfect score, cap at +800
		}
		return -800, 0 // Zero score, cap at -800
	}
	eloDiff = -400 * math.Log10(1/score-1)

	// Standard error using normal approximation
	// Variance of score = p(1-p)/n
	variance := score * (1 - score) / total
	stdErr := math.Sqrt(variance)

	// Convert to Elo error (derivative of Elo formula)
	// d(Elo)/d(score) = 400 / (ln(10) * score * (1-score))
	if score > 0.01 && score < 0.99 {
		dElo := 400 / (math.Ln10 * score * (1 - score))
		eloError = 1.96 * stdErr * dElo // 95% CI
	} else {
		eloError = 200 // Large uncertainty at extreme scores
	}

	return eloDiff, eloError
}

// LOS calculates Likelihood Of Superiority (probability that engine1 is stronger)
func LOS(wins, draws, losses int) float64 {
	if wins+losses == 0 {
		return 0.5
	}

	// Using normal approximation to binomial
	n := float64(wins + losses)
	p := float64(wins) / n

	// Z-score for p > 0.5
	z := (p - 0.5) * math.Sqrt(n) / 0.5

	// CDF of standard normal
	return 0.5 * (1 + erf(z/math.Sqrt2))
}

// SPRT performs Sequential Probability Ratio Test
// Returns LLR (Log-Likelihood Ratio) and conclusion if test stopped
// elo0, elo1 are the hypothesis bounds (e.g., -5, 0 for "not weaker" test)
func SPRT(wins, draws, losses int, elo0, elo1 float64) (llr float64, conclusion string) {
	total := float64(wins + draws + losses)
	if total < 10 {
		return 0, ""
	}

	w := float64(wins) / total
	d := float64(draws) / total
	l := float64(losses) / total

	// Convert Elo to expected scores
	p0 := 1 / (1 + math.Pow(10, -elo0/400))
	p1 := 1 / (1 + math.Pow(10, -elo1/400))

	// 3-point model for draws
	// Assume draw probability stays roughly constant
	w0 := p0 - d/2
	l0 := 1 - p0 - d/2
	w1 := p1 - d/2
	l1 := 1 - p1 - d/2

	// LLR calculation
	if w0 <= 0 || w1 <= 0 || l0 <= 0 || l1 <= 0 {
		return 0, ""
	}

	llr = total * (w*math.Log(w1/w0) + l*math.Log(l1/l0))

	// SPRT bounds (alpha=0.05, beta=0.05)
	alpha := 0.05
	beta := 0.05
	lowerBound := math.Log(beta / (1 - alpha))
	upperBound := math.Log((1 - beta) / alpha)

	if llr >= upperBound {
		return llr, "H0 REJECTED - Engine1 is not weaker"
	}
	if llr <= lowerBound {
		return llr, "H1 REJECTED - Engine1 may be weaker"
	}

	return llr, "" // Continue testing
}

// erf is the error function (Gauss error function)
func erf(x float64) float64 {
	// Approximation using Horner's method
	// Abramowitz and Stegun approximation
	a1 := 0.254829592
	a2 := -0.284496736
	a3 := 1.421413741
	a4 := -1.453152027
	a5 := 1.061405429
	p := 0.3275911

	sign := 1.0
	if x < 0 {
		sign = -1
		x = -x
	}

	t := 1.0 / (1.0 + p*x)
	y := 1.0 - (((((a5*t+a4)*t)+a3)*t+a2)*t+a1)*t*math.Exp(-x*x)

	return sign * y
}
