package speller

import (
	"encoding/json"
	"fmt"
	"github.com/Khasmag06/kode-notes/internal/models"
	"net/http"
	"net/url"
)

type service struct {
	client *http.Client
	apiURL string
}

func New() *service {
	return &service{
		client: &http.Client{},
		apiURL: "https://speller.yandex.net/services/spellservice.json/checkText"}
}

func (s *service) CheckText(text string) ([]models.SpellError, error) {
	textEncoded := url.QueryEscape(text)
	urlParams := url.Values{}
	urlParams.Add("text", textEncoded)
	reqURL := fmt.Sprintf("%s?%s", s.apiURL, urlParams.Encode())

	resp, err := s.client.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	defer resp.Body.Close()

	var spellErrors []models.SpellError
	err = json.NewDecoder(resp.Body).Decode(&spellErrors)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return spellErrors, nil
}
