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
	// Simple penalty system, could be more complex with LOC normalization
	finalScore := math.Max(0, baseScore-deduction)

	return ScoreBreakdown{
		OverallScore: finalScore,
		StaticScore:  finalScore,
		AIConfidence: 0.85, // Mock confidence
	}
}
