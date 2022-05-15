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
	json.Unmarshal(body, sentimentAnalysis)
	lol := string(body)
	_ = lol

	if len(sentimentAnalysis.Documents) == 0 {
		return nil, fmt.Errorf("unexpectedly no analysis returned")
	}
	resultSentiment := sentimentAnalysis.Documents[0].sentiment()
	resultConfidence := sentimentAnalysis.Documents[0].ConfidenceScores.confidence(resultSentiment)
	resultAnalysis := sa.Analysis{
		Sentiment:  resultSentiment,
		Confidence: resultConfidence,
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

func (t textAnalyticsResponseDocument) sentiment() sa.Sentiment {
	var sentiment sa.Sentiment
	switch t.Sentiment {
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
