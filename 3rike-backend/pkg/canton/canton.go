// Package canton provides a stub client for the Canton JSON Ledger API.
// Replace stub methods with real DAML command submissions once contracts are deployed.
// Canton JSON Ledger API base: POST /v2/commands/submit-and-wait
// Docs: https://docs.digitalasset.com/build/3.5
package canton

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Client is the Canton JSON Ledger API client.
type Client struct {
	baseURL       string
	staticToken   string
	tokenProvider *TokenProvider
	httpClient    *http.Client
}

// New creates a Client with a static bearer token (or empty for stub mode).
func New(baseURL, token string) *Client {
	return &Client{baseURL: baseURL, staticToken: token, httpClient: &http.Client{}}
}

// NewWithTokenProvider creates a Client that auto-fetches tokens via Keycloak.
func NewWithTokenProvider(baseURL string, tp *TokenProvider) *Client {
	return &Client{baseURL: baseURL, tokenProvider: tp, httpClient: &http.Client{}}
}

func (c *Client) bearerToken() (string, error) {
	if c.tokenProvider != nil {
		return c.tokenProvider.Token()
	}
	return c.staticToken, nil
}

// TokenizeResult is returned after a tricycle is tokenized on the ledger.
type TokenizeResult struct {
	ContractID string `json:"contractId"`
}

// Tokenize submits a CreateTricycleToken command to the Canton ledger.
// Stub: returns a deterministic contract ID until real DAML contracts are deployed.
func (c *Client) Tokenize(ctx context.Context, tricycleID uint, driverParty string) (*TokenizeResult, error) {
	if c.baseURL == "" {
		// Stub mode — no Canton node configured yet.
		return &TokenizeResult{ContractID: fmt.Sprintf("stub-contract-%d", tricycleID)}, nil
	}

	payload := map[string]any{
		"commands": []map[string]any{{
			"CreateCommand": map[string]any{
				"templateId": "trike-contracts:TricycleToken:TricycleToken",
				"createArguments": map[string]any{
					"tricycleId": fmt.Sprintf("%d", tricycleID),
					"owner":      driverParty,
				},
			},
		}},
		"actAs":    []string{driverParty},
		"readAs":   []string{driverParty},
		"workflowId": fmt.Sprintf("tokenize-%d", tricycleID),
	}

	res, err := c.post(ctx, "/v2/commands/submit-and-wait", payload)
	if err != nil {
		return nil, err
	}

	var result struct {
		Transaction struct {
			Events []struct {
				Created struct {
					ContractID string `json:"contractId"`
				} `json:"created"`
			} `json:"events"`
		} `json:"transaction"`
	}
	if err := json.Unmarshal(res, &result); err != nil {
		return nil, err
	}
	if len(result.Transaction.Events) == 0 {
		return nil, fmt.Errorf("canton: no events in tokenize response")
	}
	return &TokenizeResult{ContractID: result.Transaction.Events[0].Created.ContractID}, nil
}

// FractionalizeResult is returned after fractions are created on the ledger.
type FractionalizeResult struct {
	ContractID     string `json:"contractId"`
	TotalFractions int    `json:"totalFractions"`
}

// Fractionalize submits a FractionalizeTricycle choice to the Canton ledger.
// Stub: returns deterministic values until real DAML contracts are deployed.
func (c *Client) Fractionalize(ctx context.Context, contractID string, totalFractions int, operatorParty string) (*FractionalizeResult, error) {
	if c.baseURL == "" {
		return &FractionalizeResult{ContractID: contractID, TotalFractions: totalFractions}, nil
	}

	payload := map[string]any{
		"commands": []map[string]any{{
			"ExerciseCommand": map[string]any{
				"templateId": "trike-contracts:TricycleToken:TricycleToken",
				"contractId": contractID,
				"choice":     "Fractionalize",
				"choiceArgument": map[string]any{
					"totalFractions": totalFractions,
				},
			},
		}},
		"actAs":  []string{operatorParty},
		"readAs": []string{operatorParty},
	}

	res, err := c.post(ctx, "/v2/commands/submit-and-wait", payload)
	if err != nil {
		return nil, err
	}

	var result struct {
		Transaction struct {
			Events []struct {
				Created struct {
					ContractID string `json:"contractId"`
				} `json:"created"`
			} `json:"events"`
		} `json:"transaction"`
	}
	if err := json.Unmarshal(res, &result); err != nil {
		return nil, err
	}
	if len(result.Transaction.Events) == 0 {
		return nil, fmt.Errorf("canton: no events in fractionalize response")
	}
	return &FractionalizeResult{ContractID: result.Transaction.Events[0].Created.ContractID, TotalFractions: totalFractions}, nil
}

func (c *Client) post(ctx context.Context, path string, body any) ([]byte, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if tok, err := c.bearerToken(); err != nil {
		return nil, fmt.Errorf("canton: get token: %w", err)
	} else if tok != "" {
		req.Header.Set("Authorization", "Bearer "+tok)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("canton: HTTP %d from %s", resp.StatusCode, path)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	return buf.Bytes(), err
}
