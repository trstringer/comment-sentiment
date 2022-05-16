package azure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	sa "github.com/trstringer/comment-sentiment/pkg/sentimentanalyzer"
)

// SentimentService represents the cognitive services language resource.
type SentimentService struct {
	endpoint string
	key      string
}

type confidenceScores struct {
	Positive float32 `json:"positive"`
	Negative float32 `json:"negative"`
	Neutral  float32 `json:"neutral"`
}

type textAnalyticsResponseDocument struct {
	ID               string           `json:"id"`
	Sentiment        string           `json:"sentiment"`
	ConfidenceScores confidenceScores `json:"confidenceScores"`
	Sentences        []sentence       `json:"sentences"`
}

type sentence struct {
	Sentiment        string           `json:"sentiment"`
	ConfidenceScores confidenceScores `json:"confidenceScores"`
	Text             string           `json:"text"`
}

type textAnalyticsResponse struct {
	Documents []textAnalyticsResponseDocument `json:"documents"`
}

// NewSentimentService generates a new Azure sentiment service.
func NewSentimentService(endpoint, key string) *SentimentService {
	return &SentimentService{
		endpoint: endpoint,
		key:      key,
	}
}

// AnalyzeSentiment makes a call to cognitive services to analyze the
// sentiment.
func (a SentimentService) AnalyzeSentiment(text string) (*sa.Analysis, error) {
	textMarshalled, err := formatDocument(text)
	if err != nil {
		return nil, fmt.Errorf("error creating format document: %w", err)
	}

	textAnalyticsURL := fmt.Sprintf("%s/text/analytics/v3.2-preview.1/sentiment", a.endpoint)

	req, err := http.NewRequest(
		http.MethodPost,
		textAnalyticsURL,
		bytes.NewBuffer(textMarshalled),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating new request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Ocp-Apim-Subscription-Key", a.key)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("unexpected status from language service: HTTP %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	sentimentAnalysis := &textAnalyticsResponse{}
	if err := json.Unmarshal(body, sentimentAnalysis); err != nil {
		return nil, fmt.Errorf("error unmarshalling text analysis")
	}

	if len(sentimentAnalysis.Documents) == 0 {
		return nil, fmt.Errorf("unexpectedly no analysis returned")
	}
	resultSentiment := sentimentFromString(sentimentAnalysis.Documents[0].Sentiment)
	resultConfidence := sentimentAnalysis.Documents[0].ConfidenceScores.confidence(resultSentiment)
	sentenceAnalyses := []sa.SentenceAnalysis{}
	for _, sentence := range sentimentAnalysis.Documents[0].Sentences {
		sentenceSentiment := sentimentFromString(sentence.Sentiment)
		sentenceSentimentConfidence := sentence.ConfidenceScores.confidence(sentenceSentiment)
		sentenceAnalyses = append(
			sentenceAnalyses,
			sa.SentenceAnalysis{
				Text:       sentence.Text,
				Confidence: sentenceSentimentConfidence,
				Sentiment:  sentenceSentiment,
			},
		)
	}
	resultAnalysis := sa.Analysis{
		Sentiment:        resultSentiment,
		Confidence:       resultConfidence,
		SentenceAnalyses: sentenceAnalyses,
	}

	return &resultAnalysis, nil
}

func formatDocument(text string) ([]byte, error) {
	documentStructure := map[string][]map[string]string{
		"documents": {
			{
				"id":   "1",
				"text": text,
			},
		},
	}
	return json.Marshal(documentStructure)
}

func sentimentFromString(rawSentiment string) sa.Sentiment {
	var sentiment sa.Sentiment
	switch rawSentiment {
	case "positive":
		sentiment = sa.Positive
	case "negative":
		sentiment = sa.Negative
	case "neutral":
		sentiment = sa.Neutral
	}
	return sentiment
}

func (c confidenceScores) confidence(sentiment sa.Sentiment) float32 {
	switch sentiment {
	case sa.Positive:
		return c.Positive
	case sa.Negative:
		return c.Negative
	default:
		return c.Neutral
	}
}
