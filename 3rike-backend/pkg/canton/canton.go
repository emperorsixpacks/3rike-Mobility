// Package canton provides a client for the Canton JSON Ledger API.
// Canton JSON Ledger API base: POST /v2/commands/submit-and-wait-for-transaction-tree
// Docs: https://docs.digitalasset.com/build/3.5
package canton

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// TrikeContractsPkgID is the package hash for trike-contracts-1.0.0.
// DAML Int and Decimal fields must be serialized as strings per the JSON Ledger API spec.
const TrikeContractsPkgID = "fb36f0cfcfe6c4b2a24f458f5ba06bfc697fa0584b13f44ae3d3568a294d4c19"

// PartyID constructs a Canton party ID from a Keycloak subject claim and the
// participant node fingerprint (the fixed suffix after :: for all users on the same node).
// e.g. PartyID("70af8ee8-...", "1220195a...") → "70af8ee8-...::1220195a..."
func PartyID(keycloakSub, participantFingerprint string) string {
	if keycloakSub == "" || participantFingerprint == "" {
		return ""
	}
	return keycloakSub + "::" + participantFingerprint
}

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
func (c *Client) Tokenize(ctx context.Context, tricycleID uint, driverParty string) (*TokenizeResult, error) {
	if c.baseURL == "" {
		return &TokenizeResult{ContractID: fmt.Sprintf("stub-contract-%d", tricycleID)}, nil
	}

	payload := map[string]any{
		"commandId":  fmt.Sprintf("tokenize-%d-%d", tricycleID, time.Now().UnixNano()),
		"workflowId": fmt.Sprintf("tokenize-%d", tricycleID),
		"actAs":      []string{driverParty},
		"readAs":     []string{driverParty},
		"commands": []map[string]any{{
			"CreateCommand": map[string]any{
				"templateId": TrikeContractsPkgID + ":TricycleToken:TricycleToken",
				"createArguments": map[string]any{
					"operator":       driverParty,
					"driver":         driverParty,
					"tricycleId":     fmt.Sprintf("%d", tricycleID),
					"make":           "Unknown",
					"model":          "Unknown",
					"isEV":           false,
					"priceUSD":       "0.0",
					"totalFractions": "0",
					"fractionalized": false,
					"weeksRemaining": "70",
				},
			},
		}},
	}

	res, err := c.post(ctx, "/v2/commands/submit-and-wait-for-transaction-tree", payload)
	if err != nil {
		return nil, err
	}
	contractID, err := extractContractID(res)
	if err != nil {
		return nil, err
	}
	return &TokenizeResult{ContractID: contractID}, nil
}

// FractionalizeResult is returned after fractions are created on the ledger.
type FractionalizeResult struct {
	ContractID     string `json:"contractId"`
	TotalFractions int    `json:"totalFractions"`
}

// Fractionalize exercises the Fractionalize choice on a TricycleToken contract.
func (c *Client) Fractionalize(ctx context.Context, contractID string, totalFractions int, operatorParty string) (*FractionalizeResult, error) {
	if c.baseURL == "" {
		return &FractionalizeResult{ContractID: contractID, TotalFractions: totalFractions}, nil
	}

	payload := map[string]any{
		"commandId": fmt.Sprintf("fractionalize-%s-%d", contractID[:8], time.Now().UnixNano()),
		"actAs":     []string{operatorParty},
		"readAs":    []string{operatorParty},
		"commands": []map[string]any{{
			"ExerciseCommand": map[string]any{
				"templateId": TrikeContractsPkgID + ":TricycleToken:TricycleToken",
				"contractId": contractID,
				"choice":     "Fractionalize",
				"choiceArgument": map[string]any{
					"investors":    []string{operatorParty},
					"unitsEach":    fmt.Sprintf("%d", totalFractions),
					"pricePerUnit": "1.0",
				},
			},
		}},
	}

	res, err := c.post(ctx, "/v2/commands/submit-and-wait-for-transaction-tree", payload)
	if err != nil {
		return nil, err
	}
	newContractID, err := extractContractID(res)
	if err != nil {
		return nil, err
	}
	return &FractionalizeResult{ContractID: newContractID, TotalFractions: totalFractions}, nil
}

// WalletBalance holds a user's CC balance from the Canton wallet.
type WalletBalance struct {
	Round               int    `json:"round"`
	EffectiveUnlockedQty string `json:"effective_unlocked_qty"`
	EffectiveLockedQty   string `json:"effective_locked_qty"`
	TotalHoldingFees     string `json:"total_holding_fees"`
}

// GetWalletBalance fetches the CC balance for the authenticated user from the validator wallet API.
func (c *Client) GetWalletBalance(ctx context.Context, validatorURL string) (*WalletBalance, error) {
	tok, err := c.bearerToken()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, validatorURL+"/api/validator/v0/wallet/balance", nil)
	if err != nil {
		return nil, err
	}
	if tok != "" {
		req.Header.Set("Authorization", "Bearer "+tok)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("canton: wallet balance HTTP %d", resp.StatusCode)
	}
	var bal WalletBalance
	return &bal, json.NewDecoder(resp.Body).Decode(&bal)
}
func extractContractID(data []byte) (string, error) {
	var result struct {
		TransactionTree struct {
			EventsById map[string]struct {
				CreatedTreeEvent struct {
					Value struct {
						ContractID string `json:"contractId"`
					} `json:"value"`
				} `json:"CreatedTreeEvent"`
			} `json:"eventsById"`
		} `json:"transactionTree"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return "", fmt.Errorf("canton: parse response: %w", err)
	}
	for _, ev := range result.TransactionTree.EventsById {
		if id := ev.CreatedTreeEvent.Value.ContractID; id != "" {
			return id, nil
		}
	}
	return "", fmt.Errorf("canton: no contractId in response")
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
		var buf bytes.Buffer
		_, _ = buf.ReadFrom(resp.Body)
		return nil, fmt.Errorf("canton: HTTP %d from %s: %s", resp.StatusCode, path, buf.String())
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	return buf.Bytes(), err
}
