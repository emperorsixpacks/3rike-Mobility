// Package canton_test contains integration tests for the Canton JSON Ledger API.
//
// These tests run against the HackCanton devnet and require OIDC credentials.
// They are skipped automatically when the env vars are not set.
//
// Required env vars:
//
//	OIDC_USERNAME   your hackathon email
//	OIDC_PASSWORD   your hackathon password
//	CANTON_PARTY    your Canton party ID (e.g. "Alice::abc123...")
//
// Optional:
//
//	CANTON_URL      defaults to the devnet JSON Ledger API
//	OIDC_TOKEN_URL  defaults to the hackathon Keycloak URL
//	OIDC_CLIENT_ID  defaults to the hackathon client ID
package canton_test

import (
	"context"
	"os"
	"testing"

	"github.com/3rike12/3rike-backend/pkg/canton"
)

const (
	defaultCantonURL  = "https://ledger-api-json.participant.hackcanton-01.devnet.naas.noders.services"
	defaultOIDCURL    = "https://keycloak.naas.noders.services/realms/noders-appsfactory/protocol/openid-connect/token"
	defaultClientID   = "web-app-ui-hackcanton-01-devnet"
)

func integrationClient(t *testing.T) (*canton.Client, string) {
	t.Helper()

	username := os.Getenv("OIDC_USERNAME")
	password := os.Getenv("OIDC_PASSWORD")
	party := os.Getenv("CANTON_PARTY")

	if username == "" || password == "" || party == "" {
		t.Skip("skipping integration test: set OIDC_USERNAME, OIDC_PASSWORD, CANTON_PARTY")
	}

	tokenURL := getEnvOr("OIDC_TOKEN_URL", defaultOIDCURL)
	clientID := getEnvOr("OIDC_CLIENT_ID", defaultClientID)
	cantonURL := getEnvOr("CANTON_URL", defaultCantonURL)

	tp := canton.NewTokenProvider(tokenURL, clientID, username, password)
	return canton.NewWithTokenProvider(cantonURL, tp), party
}

func getEnvOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// TestTokenProvider_FetchesToken verifies that the provider can obtain a token from Keycloak.
func TestTokenProvider_FetchesToken(t *testing.T) {
	username := os.Getenv("OIDC_USERNAME")
	password := os.Getenv("OIDC_PASSWORD")
	if username == "" || password == "" {
		t.Skip("skipping: set OIDC_USERNAME and OIDC_PASSWORD")
	}

	tokenURL := getEnvOr("OIDC_TOKEN_URL", defaultOIDCURL)
	clientID := getEnvOr("OIDC_CLIENT_ID", defaultClientID)

	tp := canton.NewTokenProvider(tokenURL, clientID, username, password)
	tok, err := tp.Token()
	if err != nil {
		t.Fatalf("Token() error: %v", err)
	}
	if tok == "" {
		t.Fatal("Token() returned empty string")
	}
	t.Logf("token obtained (len=%d)", len(tok))

	// Second call should return cached token.
	tok2, err := tp.Token()
	if err != nil {
		t.Fatalf("Token() second call error: %v", err)
	}
	if tok != tok2 {
		t.Error("expected cached token on second call")
	}
}

// TestTokenize_CreatesContractOnLedger submits a CreateTricycleToken command and
// verifies a contract ID is returned.
func TestTokenize_CreatesContractOnLedger(t *testing.T) {
	client, party := integrationClient(t)

	res, err := client.Tokenize(context.Background(), 1, party)
	if err != nil {
		t.Fatalf("Tokenize() error: %v", err)
	}
	if res.ContractID == "" {
		t.Fatal("Tokenize() returned empty contractId")
	}
	t.Logf("contractId: %s", res.ContractID)
}

// TestFractionalize_ExercisesChoiceOnLedger tokenizes a tricycle then fractionalizes it,
// verifying the full backend→contract round-trip.
func TestFractionalize_ExercisesChoiceOnLedger(t *testing.T) {
	client, party := integrationClient(t)
	ctx := context.Background()

	tok, err := client.Tokenize(ctx, 2, party)
	if err != nil {
		t.Fatalf("Tokenize() error: %v", err)
	}
	t.Logf("tokenized contractId: %s", tok.ContractID)

	frac, err := client.Fractionalize(ctx, tok.ContractID, 10, party)
	if err != nil {
		t.Fatalf("Fractionalize() error: %v", err)
	}
	if frac.ContractID == "" {
		t.Fatal("Fractionalize() returned empty contractId")
	}
	if frac.TotalFractions != 10 {
		t.Errorf("expected 10 fractions, got %d", frac.TotalFractions)
	}
	t.Logf("fractionalized contractId: %s, fractions: %d", frac.ContractID, frac.TotalFractions)
}
