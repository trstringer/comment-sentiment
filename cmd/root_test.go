package cmd

import (
	"testing"

	sa "github.com/trstringer/comment-sentiment/pkg/sentimentanalyzer"
	fake "github.com/trstringer/comment-sentiment/pkg/sentimentanalyzer/azure"
	// gh "github.com/trstringer/comment-sentiment/pkg/github"
)

func TestResponseComment(t *testing.T) {
	testCases := []struct {
		name              string
		comment           string
		expectedSentiment sa.Sentiment
	}{
		{
			name:              "positive_comment",
			comment:           "i love this",
			expectedSentiment: sa.Positive,
		},
	}

	fakeSentimentService := &fake.FakeSentimentService{}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			analysis, err := fakeSentimentService.AnalyzeSentiment(testCase.comment)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if analysis.Sentiment != testCase.expectedSentiment {
				t.Fatalf("Unexpected sentiment. Expected %s got %s", testCase.expectedSentiment, analysis.Sentiment)
			}
		})
	}
}
