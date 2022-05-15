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
	Sentiment        Sentiment
	Confidence       float32
	SentenceAnalyses []SentenceAnalysis
}

// SentenceAnalysis represents individual sentence analysis.
type SentenceAnalysis struct {
	Sentiment  Sentiment
	Confidence float32
	Text       string
}

type sentimentAnalyzer interface {
	analyzeSentiment(string) (Analysis, error)
}

// GetSentiment calls a sentimentService and retrieves the corresponding
// analysis.
func GetSentiment(svc sentimentAnalyzer, comment string) (Analysis, error) {
	return svc.analyzeSentiment(comment)
}

// NegativeSentences returns any negative sentences.
func (a Analysis) NegativeSentences() []SentenceAnalysis {
	negativeSentences := []SentenceAnalysis{}
	for _, sentenceAnalysis := range a.SentenceAnalyses {
		if sentenceAnalysis.Sentiment == Negative {
			negativeSentences = append(negativeSentences, sentenceAnalysis)
		}
	}

	return negativeSentences
}
