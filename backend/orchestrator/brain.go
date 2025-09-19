package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"
)

// TaskProfile represents the analyzed characteristics of an AI task
type TaskProfile struct {
	Type           string    `json:"type"`            // e.g., "text-generation", "classification", "translation"
	Complexity     float64   `json:"complexity"`      // 0-1 score
	TokenEstimate  int       `json:"token_estimate"`
	Priority       int       `json:"priority"`        // 1-5 scale
	MaxBudget      float64   `json:"max_budget"`
	RequiredCAPs   []string  `json:"required_caps"`   // Required capabilities
	TimeConstraint time.Duration `json:"time_constraint"`
}

// ProviderScore represents the evaluated score of a provider for a specific task
type ProviderScore struct {
	Provider     *Provider
	TotalScore   float64
	CostScore    float64
	SpeedScore   float64
	QualityScore float64
	ReasoningLog []string
}

// AIBrain handles intelligent routing decisions
type AIBrain struct {
	providers     []*Provider
	metrics       *MetricsCollector
	cache         *Cache
	loadBalancer  *LoadBalancer
}

func NewAIBrain() *AIBrain {
	return &AIBrain{
		metrics:      NewMetricsCollector(),
		cache:        NewCache(),
		loadBalancer: NewLoadBalancer(),
	}
}

// AnalyzeTask examines the prompt and requirements to create a task profile
func (b *AIBrain) AnalyzeTask(ctx context.Context, prompt string, options map[string]interface{}) (*TaskProfile, error) {
	// Use a simple classifier first (can be enhanced with ML later)
	profile := &TaskProfile{
		Type:       detectTaskType(prompt),
		Complexity: calculateComplexity(prompt),
		TokenEstimate: estimateTokenCount(prompt),
		Priority:    getPriorityFromOptions(options),
		MaxBudget:   getBudgetFromOptions(options),
	}

	// Analyze required capabilities
	profile.RequiredCAPs = analyzeRequiredCapabilities(prompt, options)
	
	return profile, nil
}

// SelectProvider chooses the optimal provider for the task
func (b *AIBrain) SelectProvider(ctx context.Context, profile *TaskProfile) (*RoutingResult, error) {
	scores := make([]*ProviderScore, 0)

	// Score each provider
	for _, provider := range b.providers {
		if !providerMeetsRequirements(provider, profile) {
			continue
		}

		score := b.scoreProvider(ctx, provider, profile)
		scores = append(scores, score)
	}

	// Sort by total score
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].TotalScore > scores[j].TotalScore
	})

	if len(scores) == 0 {
		return nil, fmt.Errorf("no suitable provider found for the task")
	}

	// Get top 3 providers as fallbacks
	fallbacks := make([]Provider, 0)
	for i := 1; i < min(len(scores), 3); i++ {
		fallbacks = append(fallbacks, *scores[i].Provider)
	}

	return &RoutingResult{
		Provider:   scores[0].Provider.Name,
		Model:      scores[0].Provider.Model,
		Confidence: scores[0].TotalScore,
		Reasoning:  strings.Join(scores[0].ReasoningLog, "; "),
		Fallbacks:  fallbacks,
	}, nil
}

// scoreProvider evaluates a provider for a specific task
func (b *AIBrain) scoreProvider(ctx context.Context, provider *Provider, profile *TaskProfile) *ProviderScore {
	score := &ProviderScore{
		Provider:     provider,
		ReasoningLog: make([]string, 0),
	}

	// Get real-time metrics
	metrics := b.metrics.GetProviderMetrics(provider.Name)
	
	// Cost Score (0-1)
	estimatedCost := float64(profile.TokenEstimate) * provider.CostPer1K / 1000
	costScore := 1.0 - (estimatedCost / profile.MaxBudget)
	score.CostScore = costScore
	score.ReasoningLog = append(score.ReasoningLog, 
		fmt.Sprintf("Cost efficiency: %.2f (estimated cost: $%.4f)", costScore, estimatedCost))

	// Speed Score (0-1)
	speedScore := 1.0 - (float64(provider.AvgLatency) / float64(profile.TimeConstraint.Milliseconds()))
	score.SpeedScore = speedScore
	score.ReasoningLog = append(score.ReasoningLog, 
		fmt.Sprintf("Speed score: %.2f (avg latency: %dms)", speedScore, provider.AvgLatency))

	// Quality Score (0-1) based on historical metrics
	qualityScore := calculateQualityScore(metrics, profile)
	score.QualityScore = qualityScore
	score.ReasoningLog = append(score.ReasoningLog, 
		fmt.Sprintf("Quality score: %.2f (success rate: %.2f%%)", qualityScore, metrics.SuccessRate*100))

	// Calculate weighted total score
	score.TotalScore = (costScore * 0.4) + (speedScore * 0.3) + (qualityScore * 0.3)

	return score
}

// Helper functions
func detectTaskType(prompt string) string {
	// Simple keyword-based detection (can be enhanced with ML)
	switch {
	case strings.Contains(strings.ToLower(prompt), "classify"):
		return "classification"
	case strings.Contains(strings.ToLower(prompt), "translate"):
		return "translation"
	case strings.Contains(strings.ToLower(prompt), "summarize"):
		return "summarization"
	default:
		return "text-generation"
	}
}

func calculateComplexity(prompt string) float64 {
    // Basic complexity calculation (can be enhanced)
    wordCount := len(strings.Fields(prompt))
    if wordCount <= 0 { return 0 }
    ratio := float64(wordCount) / 1000.0
    if ratio > 1.0 { ratio = 1.0 }
    return ratio
}

func estimateTokenCount(prompt string) int {
	// Rough estimate (can be enhanced)
	return len(strings.Fields(prompt)) * 1.3
}

func providerMeetsRequirements(provider *Provider, profile *TaskProfile) bool {
	// Check if provider supports all required capabilities
	for _, cap := range profile.RequiredCAPs {
		found := false
		for _, providerCap := range provider.Capabilities {
			if cap == providerCap {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func calculateQualityScore(metrics *ProviderMetrics, profile *TaskProfile) float64 {
	// Combine various quality metrics
	successWeight := 0.5
	errorWeight := 0.3
	timeoutWeight := 0.2

	successScore := metrics.SuccessRate
	errorScore := 1.0 - metrics.ErrorRate
	timeoutScore := 1.0 - metrics.TimeoutRate

	return (successScore * successWeight) + 
	       (errorScore * errorWeight) + 
	       (timeoutScore * timeoutWeight)
}
