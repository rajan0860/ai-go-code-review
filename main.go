package main

import (
	"fmt"
	"os"

	"github.com/rajanmehta/ai-go-code-review/analyzer"
	"github.com/rajanmehta/ai-go-code-review/llm"
	"github.com/rajanmehta/ai-go-code-review/review"
)

func main() {
	fmt.Println("GoReview: AI-Powered Code Review")
	if len(os.Args) < 2 {
		fmt.Println("Usage: goreview review <filename_or_directory>")
		os.Exit(1)
	}

	command := os.Args[1]
	if command != "review" {
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}

	target := os.Args[2]
	fmt.Printf("Analyzing %s...\n", target)

	// 1. Static Analysis
	issues, err := analyzer.RunStaticAnalysis(target)
	if err != nil {
		fmt.Printf("Error running static analysis: %v\n", err)
		os.Exit(1)
	}

	// 2. LLM Review
	// Default to local Ollama instance with qwen2.5-coder
	model := "qwen2.5-coder:latest"
	if envModel := os.Getenv("OLLAMA_MODEL"); envModel != "" {
		model = envModel
	}

	llmClient := llm.NewOllamaClient(model)
	// OpenAI fallback (optional, commented out for now)
	// apiKey := os.Getenv("OPENAI_API_KEY")
	// if apiKey != "" {
	// 	llmClient = llm.NewOpenAIClient(apiKey)
	// }

	// Read file content for LLM
	content, err := os.ReadFile(target)
	if err == nil {
		review, err := llmClient.ReviewCode(string(content))
		if err != nil {
			fmt.Printf("LLM Review Error: %v\n", err)
		} else {
			fmt.Println("\n--- AI Review ---")
			fmt.Println(review)
			fmt.Println("-----------------")
		}
	} else {
		// If target is a directory, avoiding reading all files for now to save tokens
		fmt.Println("Skipping LLM review for directory target (TODO: implement directory summary)")
	}

	// 3. Generate Report
	review.GenerateReport(issues)
}
