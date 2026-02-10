package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/rajanmehta/ai-go-code-review/analyzer"
	"github.com/rajanmehta/ai-go-code-review/config"
	"github.com/rajanmehta/ai-go-code-review/llm"
	"github.com/rajanmehta/ai-go-code-review/review"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.LoadConfig()

	if len(os.Args) < 2 {
		logger.Error("Usage: goreview review <filename_or_directory>")
		os.Exit(1)
	}

	command := os.Args[1]
	if command != "review" {
		logger.Error("Unknown command", "command", command)
		os.Exit(1)
	}

	target := os.Args[2]
	logger.Info("Starting analysis", "target", target, "model", cfg.OllamaModel)

	// 1. Static Analysis
	issues, err := analyzer.RunStaticAnalysis(target)
	if err != nil {
		logger.Error("Error running static analysis", "error", err)
		os.Exit(1)
	}

	// 2. LLM Review
	llmClient := llm.NewOllamaClient(cfg.OllamaModel)

	// Read file content for LLM
	var contentBuilder strings.Builder
	var totalBytes int
	const maxBytes = 10000 // Limit context sent to LLM

	info, err := os.Stat(target)
	if err != nil {
		logger.Error("Error checking target", "error", err)
		os.Exit(1)
	}

	if info.IsDir() {
		logger.Info("Scanning directory for LLM review", "path", target)
		err := filepath.Walk(target, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
				if totalBytes >= maxBytes {
					return nil // Skip if max bytes reached
				}

				fileContent, err := os.ReadFile(path)
				if err != nil {
					logger.Warn("Failed to read file", "path", path, "error", err)
					return nil
				}

				if totalBytes+len(fileContent) > maxBytes {
					// Truncate if necessary or just stop adding files
					// For now, let's just stop adding to avoid partial files
					logger.Warn("Context limit reached, stopping file collection", "limit", maxBytes)
					return filepath.SkipDir
				}

				contentBuilder.WriteString(fmt.Sprintf("\n--- File: %s ---\n", path))
				contentBuilder.Write(fileContent)
				totalBytes += len(fileContent)
			}
			return nil
		})
		if err != nil {
			logger.Error("Error walking directory", "error", err)
		}
	} else {
		fileContent, err := os.ReadFile(target)
		if err != nil {
			logger.Error("Error reading file", "error", err)
			os.Exit(1)
		}
		contentBuilder.WriteString(fmt.Sprintf("\n--- File: %s ---\n", target))
		contentBuilder.Write(fileContent)
	}

	if contentBuilder.Len() > 0 {
		review, err := llmClient.ReviewCode(contentBuilder.String())
		if err != nil {
			logger.Error("LLM Review Error", "error", err)
		} else {
			logger.Info("AI Review Generated")
			fmt.Println("\n--- AI Review ---")
			fmt.Println(review)
			fmt.Println("-----------------")
		}
	} else {
		logger.Warn("No content found to review")
	}

	loc := strings.Count(contentBuilder.String(), "\n")

	// 3. Generate Report
	review.GenerateReport(issues, loc)
}
