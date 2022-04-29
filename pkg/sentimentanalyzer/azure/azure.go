package azure

import (
	sa "github.com/trstringer/comment-sentiment/pkg/sentimentanalyzer"
)

// SentimentService represents the cognitive services language resource.
type SentimentService struct {
	key string
}

// NewSentimentService generates a new Azure sentiment service.
func NewSentimentService(key string) *SentimentService {
	return &SentimentService{
		key: key,
	}
}

// AnalyzeSentiment makes a call to cognitive services to analyze the
// sentiment.
func (a SentimentService) AnalyzeSentiment() (*sa.Analysis, error) {
	return nil, nil
}
