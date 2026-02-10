package review

import (
	"math"

	"github.com/rajanmehta/ai-go-code-review/analyzer"
)

type ScoreBreakdown struct {
	OverallScore float64
	StaticScore  float64
	AIConfidence float64
}

func CalculateQualityScore(issues []analyzer.Issue, linesOfCode int) ScoreBreakdown {
	deduction := 0.0

	for _, issue := range issues {
		switch issue.Severity {
		case "error":
			deduction += 2.0
		case "warning":
			deduction += 0.5
		case "info":
			deduction += 0.1
		default:
			deduction += 0.2
		}
	}

	baseScore := 10.0

	// Normalize deduction based on Lines of Code (LOC)
	// Larger files are expected to have more issues, so we dampen the penalty.
	// Factor = 1.0 for small files, decreases as LOC increases.
	// Using a log-based damping factor: factor = 1 / log10(max(LOC, 10))

	dampingFactor := 1.0
	if linesOfCode > 10 {
		dampingFactor = 1.0 / math.Log10(float64(linesOfCode))
	}

	adjustedDeduction := deduction * dampingFactor
	finalScore := math.Max(0, baseScore-adjustedDeduction)

	return ScoreBreakdown{
		OverallScore: finalScore,
		StaticScore:  finalScore,
		AIConfidence: 0.85, // Mock confidence
	}
}
