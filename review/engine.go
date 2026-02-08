package review

import (
	"fmt"

	"github.com/rajanmehta/ai-go-code-review/analyzer"
)

func GenerateReport(issues []analyzer.Issue) {
	fmt.Println("Review Report:")

	score := CalculateQualityScore(issues, 0)
	fmt.Printf("Quality Score: %.1f/10\n", score.OverallScore)

	if len(issues) == 0 {
		fmt.Println("No issues found.")
		return
	}

	for _, issue := range issues {
		fmt.Printf("[%s] %s:%d: %s\n", issue.Severity, issue.File, issue.Line, issue.Message)
	}
}
