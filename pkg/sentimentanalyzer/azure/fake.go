package azure

import (
	"strings"

	sa "github.com/trstringer/comment-sentiment/pkg/sentimentanalyzer"
)

// FakeSentimentService is a fake sentiment service.
type FakeSentimentService struct{}

// AnalyzeSentiment returns a fake analysis.
func (f FakeSentimentService) AnalyzeSentiment(text string) (*sa.Analysis, error) {
	var confidence float32 = 0.9
	var resultSentiment sa.Sentiment

	if strings.Contains(text, "love") {
		resultSentiment = sa.Positive
	} else if strings.Contains(text, "hate") {
		resultSentiment = sa.Negative
	} else {
		resultSentiment = sa.Neutral
	}

	return &sa.Analysis{
		Sentiment:  resultSentiment,
		Confidence: confidence,
	}, nil
}
