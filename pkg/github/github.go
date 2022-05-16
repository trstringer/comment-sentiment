package github

import (
	"fmt"
	"regexp"
	"strings"

	sa "github.com/trstringer/comment-sentiment/pkg/sentimentanalyzer"
)

const (
	negativeCommentSuggestion string = "... consider editing for a more positive response!"
	indicatorCommentStart     string = "<!-- ANALYSIS START -->"
	indicatorCommentEnd       string = "<!-- ANALYSIS END -->"
)

// sentimentResponse takes the Analysis and constructs the footer that
// should be added to the comment to display the analysis result.
func sentimentResponse(analysis sa.Analysis) string {
	response := fmt.Sprintf(
		"Sentiment analysis: %s %s (confidence: %.2f)",
		analysis.Sentiment,
		emojiFromSentiment(analysis.Sentiment),
		analysis.Confidence,
	)

	if analysis.Sentiment == sa.Negative {
		response = fmt.Sprintf("%s %s", response, negativeCommentSuggestion)
	}

	negativeSentences := analysis.NegativeSentences()
	if len(negativeSentences) > 0 && analysis.Sentiment != sa.Negative {
		response = fmt.Sprintf("%s Some negative sentences:", response)
		for _, negativeSentence := range negativeSentences {
			response = fmt.Sprintf("%s \"%s\"", response, negativeSentence.Text)
		}
	}

	// Once we have generated the comment modification, we need to then wrap
	// it with the indicators. This will allow future modifications of the
	// comment (e.g. on a comment update) to be able to search for the analysis
	// and replace it, etc.
	response = fmt.Sprintf("%s%s%s", indicatorCommentStart, response, indicatorCommentEnd)

	return response
}

// emojiFromSentiment converts the sentiment to a GitHub emoji string.
func emojiFromSentiment(source sa.Sentiment) string {
	var emoji string

	switch source {
	case sa.Positive:
		emoji = ":grin:"
	case sa.Negative:
		emoji = ":rage:"
	case sa.Neutral:
		emoji = ":neutral_face:"
	}

	return emoji
}

// UpdateCommentWithSentiment changes the comment text to include the analyzed
// sentiment.
func UpdateCommentWithSentiment(comment string, analysis sa.Analysis) (string, error) {
	comment, err := TrimCommentSentimentAnalysis(comment)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s\n\n%s", comment, sentimentResponse(analysis)), nil
}

// TrimCommentSentimentAnalysis removes any sentiment analysis from a comment.
func TrimCommentSentimentAnalysis(comment string) (string, error) {
	reg, err := regexp.Compile(fmt.Sprintf(
		"%s.*%s",
		indicatorCommentStart,
		indicatorCommentEnd,
	))
	if err != nil {
		return "", err
	}
	outputComment := strings.TrimRight(reg.ReplaceAllString(comment, ""), "\n")
	return outputComment, nil
}
