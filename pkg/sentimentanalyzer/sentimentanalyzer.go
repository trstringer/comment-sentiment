package sentimentanalyzer

// Sentiment is the type representation of a sentiment.
type Sentiment int

//go:generate stringer -type=Sentiment
const (
	Positive Sentiment = iota
	Negative
	Neutral
)

// Analysis represents the sentiment analysis.
type Analysis struct {
	Sentiment  Sentiment
	Confidence float32
}

type sentimentAnalyzer interface {
	analyzeSentiment(string) (Analysis, error)
}

// GetSentiment calls a sentimentService and retrieves the corresponding
// analysis.
func GetSentiment(svc sentimentAnalyzer, comment string) (Analysis, error) {
	return svc.analyzeSentiment(comment)
}
