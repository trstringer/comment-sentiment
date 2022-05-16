package github

import (
	"testing"

	sa "github.com/trstringer/comment-sentiment/pkg/sentimentanalyzer"
)

func TestUpdateCommentWithSentiment(t *testing.T) {
	testCases := []struct {
		name              string
		comment           string
		sentimentAnalysis sa.Analysis
		expected          string
	}{
		{
			name:    "positive_comment",
			comment: "this is a comment",
			sentimentAnalysis: sa.Analysis{
				Sentiment:  sa.Positive,
				Confidence: 0.9,
			},
			expected: `this is a comment

<!-- ANALYSIS START -->Sentiment analysis: Positive :grin: (confidence: 0.90)<!-- ANALYSIS END -->`,
		},
		{
			name:    "neutral_comment",
			comment: "this is a comment",
			sentimentAnalysis: sa.Analysis{
				Sentiment:  sa.Neutral,
				Confidence: 0.832,
			},
			expected: `this is a comment

<!-- ANALYSIS START -->Sentiment analysis: Neutral :neutral_face: (confidence: 0.83)<!-- ANALYSIS END -->`,
		},
		{
			name:    "negative_comment",
			comment: "this is a comment",
			sentimentAnalysis: sa.Analysis{
				Sentiment:  sa.Negative,
				Confidence: 1.0,
			},
			expected: `this is a comment

<!-- ANALYSIS START -->Sentiment analysis: Negative :rage: (confidence: 1.00) ... consider editing for a more positive response!<!-- ANALYSIS END -->`,
		},
		{
			name: "neutral_analysis_positive_comment",
			comment: `this is a comment.

this is another line.

<!-- ANALYSIS START -->Sentiment analysis: Neutral :neutral_face: (confidence: 0.83)<!-- ANALYSIS END -->`,
			sentimentAnalysis: sa.Analysis{
				Sentiment:  sa.Positive,
				Confidence: 0.9,
			},
			expected: `this is a comment.

this is another line.

<!-- ANALYSIS START -->Sentiment analysis: Positive :grin: (confidence: 0.90)<!-- ANALYSIS END -->`,
		},
		{
			name: "neutral_analysis_negative_sentence",
			comment: `this is a comment.

this is another line.`,
			sentimentAnalysis: sa.Analysis{
				Sentiment:  sa.Neutral,
				Confidence: 0.9,
				SentenceAnalyses: []sa.SentenceAnalysis{
					{
						Text:       "sentence text 1",
						Sentiment:  sa.Negative,
						Confidence: 0.9,
					},
				},
			},
			expected: `this is a comment.

this is another line.

<!-- ANALYSIS START -->Sentiment analysis: Neutral :neutral_face: (confidence: 0.90) Some negative sentences: "sentence text 1"<!-- ANALYSIS END -->`,
		},
		{
			name: "neutral_analysis_negative_sentence_non_negative_sentences",
			comment: `this is a comment.

this is another line.`,
			sentimentAnalysis: sa.Analysis{
				Sentiment:  sa.Neutral,
				Confidence: 0.9,
				SentenceAnalyses: []sa.SentenceAnalysis{
					{
						Text:       "sentence text 1",
						Sentiment:  sa.Negative,
						Confidence: 0.9,
					},
					{
						Text:       "sentence text 2",
						Sentiment:  sa.Positive,
						Confidence: 0.9,
					},
					{
						Text:       "sentence text 3",
						Sentiment:  sa.Neutral,
						Confidence: 0.9,
					},
				},
			},
			expected: `this is a comment.

this is another line.

<!-- ANALYSIS START -->Sentiment analysis: Neutral :neutral_face: (confidence: 0.90) Some negative sentences: "sentence text 1"<!-- ANALYSIS END -->`,
		},
		{
			name: "neutral_analysis_no_negative_sentences",
			comment: `this is a comment.

this is another line.`,
			sentimentAnalysis: sa.Analysis{
				Sentiment:  sa.Neutral,
				Confidence: 0.9,
				SentenceAnalyses: []sa.SentenceAnalysis{
					{
						Text:       "sentence text 1",
						Sentiment:  sa.Neutral,
						Confidence: 0.9,
					},
					{
						Text:       "sentence text 2",
						Sentiment:  sa.Positive,
						Confidence: 0.9,
					},
					{
						Text:       "sentence text 3",
						Sentiment:  sa.Neutral,
						Confidence: 0.9,
					},
				},
			},
			expected: `this is a comment.

this is another line.

<!-- ANALYSIS START -->Sentiment analysis: Neutral :neutral_face: (confidence: 0.90)<!-- ANALYSIS END -->`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual, err := UpdateCommentWithSentiment(
				testCase.comment,
				testCase.sentimentAnalysis,
			)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if actual != testCase.expected {
				t.Fatalf("Failure, expected '%s' and got '%s'", testCase.expected, actual)
			}
		})
	}
}

func TestTrimCommentSentimentAnalysis(t *testing.T) {
	testCases := []struct {
		name     string
		comment  string
		expected string
	}{
		{
			name:     "no_analysis",
			comment:  "this is a comment",
			expected: "this is a comment",
		},
		{
			name: "positive_analysis_footer",
			comment: `this is a comment

<!-- ANALYSIS START -->Sentiment analysis: Positive :grin: (confidence: 0.90)<!-- ANALYSIS END -->`,
			expected: "this is a comment",
		},
		{
			name: "negative_analysis_footer",
			comment: `this is a comment

<!-- ANALYSIS START -->Sentiment analysis: Negative :rage: (confidence: 0.90)... consider editing for a more positive response!<!-- ANALYSIS END -->`,
			expected: "this is a comment",
		},
		{
			name: "multi_paragraph_analysis_footer",
			comment: `this is a comment.

this is another comment.

<!-- ANALYSIS START -->Sentiment analysis: Positive :grin: (confidence: 0.90)<!-- ANALYSIS END -->`,
			expected: `this is a comment.

this is another comment.`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual, err := TrimCommentSentimentAnalysis(testCase.comment)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if actual != testCase.expected {
				t.Fatalf("Failure, expected '%s' and got '%s'", testCase.expected, actual)
			}
		})
	}
}
