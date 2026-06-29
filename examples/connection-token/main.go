// Connection token flow example — mirrors the Python UserClient workflow.
//
// Demonstrates the full asset connection lifecycle:
//  1. Get self asset accounts (available accounts for current user on an asset)
//  2. Create a connection token
//  3. Get connection token auth info (secret)
//  4. Get client launch URL
//  5. Create a super connection token (for SSO)
//  6. Get SSO login URL
//
// Environment variables:
//
//	JUMPSERVER_URL        — base URL
//	JUMPSERVER_KEY_ID     — access key ID
//	JUMPSERVER_SECRET_ID  — access key secret
//	JUMPSERVER_ASSET_ID   — asset ID to connect to (required)
//	JUMPSERVER_ACCOUNT_ID — account ID to use (required)
//
// Run:
//
//	go run ./examples/connection-token
package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"

	jumpserver "github.com/jumpserver-south/jumpserver-sdk-go"
	"github.com/jumpserver-south/jumpserver-sdk-go/model"
)

func main() {
	baseURL := os.Getenv("JUMPSERVER_URL")
	keyID := os.Getenv("JUMPSERVER_KEY_ID")
	secretID := os.Getenv("JUMPSERVER_SECRET_ID")
	assetID := "62571e56-406a-4caa-82a1-265b80642f7e"
	accountID := "0d1cee31-dbb8-42a1-b592-352e751c6426"

	if baseURL == "" || keyID == "" || secretID == "" {
		log.Fatal("JUMPSERVER_URL, JUMPSERVER_KEY_ID, JUMPSERVER_SECRET_ID required")
	}
	if assetID == "" || accountID == "" {
		log.Fatal("JUMPSERVER_ASSET_ID and JUMPSERVER_ACCOUNT_ID required")
	}

	client := jumpserver.NewClient(
		jumpserver.WithBaseURL(baseURL),
		jumpserver.WithAccessKeyAuth(keyID, secretID),
	)
	ctx := context.Background()

	// ============================================================
	// 1. Get self asset accounts (equivalent to Python get_self_asset_accounts)
	// ============================================================
	fmt.Println("=== 1. Get Self Asset Accounts ===")
	accounts, _, err := client.Permissions.GetSelfAssetAccounts(ctx, assetID)
	if err != nil {
		log.Printf("GetSelfAssetAccounts: %v (may require user-level permission)", err)
	} else {
		prettyPrint("accounts", accounts)
	}

	// ============================================================
	// 2. Create connection token (equivalent to Python create_token)
	// ============================================================
	fmt.Println("\n=== 2. Create Connection Token ===")
	connReq := &model.ConnectionTokenRequest{
		Asset:         assetID,
		Account:       accountID,
		Protocol:      "ssh",
		ConnectMethod: "web_cli",
		InputUsername: "",
		InputSecret:   "",
		ConnectOptions: map[string]any{
			"appletConnectMethod":    "client",
			"charset":                "default",
			"is_backspace_as_ctrl_h": false,
			"keyboard_layout":        "en-us-qwerty",
			"rdp_client_option":      []string{"full_screen", "drives_redirect"},
			"rdp_color_quality":      "32",
			"rdp_resolution":         "auto",
			"rdp_smart_size":         "0",
			"resolution":             "auto",
		},
	}
	token, _, err := client.Auth.CreateConnectionToken(ctx, connReq)
	if err != nil {
		log.Fatalf("CreateConnectionToken: %v", err)
	}
	fmt.Printf("Token ID: %s\n", token.ID)
	fmt.Printf("Token Value: %s\n", token.Value)

	// ============================================================
	// 3. Get connection token auth info (equivalent to Python
	//    get_connect_token_auth_info via super-connection-token/secret/)
	// ============================================================
	fmt.Println("\n=== 3. Get Connection Token Auth Info (Secret) ===")
	secret, _, err := client.Auth.GetSuperConnectionTokenSecret(ctx, token.ID)
	if err != nil {
		log.Fatalf("GetSuperConnectionTokenSecret: %v", err)
	}
	prettyPrint("auth info (JSON)", secret)

	// Base64 encoded version (for remote app parameters)
	secretJSON, _ := json.Marshal(secret)
	fmt.Printf("\nAuth info (Base64):\n%s\n", base64.StdEncoding.EncodeToString(secretJSON))

	// ============================================================
	// 4. Get client URL (equivalent to Python get_client_url)
	// ============================================================
	fmt.Println("\n=== 4. Get Client URL ===")
	clientURL, _, err := client.Auth.GetClientURL(ctx, token.ID)
	if err != nil {
		log.Printf("GetClientURL: %v", err)
	} else {
		fmt.Printf("Client URL: %s\n", clientURL)
		fmt.Println("(Paste into browser to launch local client)")
	}

	// ============================================================
	// 5. Create super connection token (equivalent to Python create_super_token)
	// ============================================================
	fmt.Println("\n=== 5. Create Super Connection Token ===")
	profile, _, err := client.Users.Profile(ctx)
	if err != nil {
		log.Fatalf("Users.Profile: %v", err)
	}

	superReq := &model.ConnectionTokenRequest{
		User:          profile.ID,
		Asset:         assetID,
		Account:       accountID,
		Protocol:      "ssh",
		ConnectMethod: "web_cli",
		ConnectOptions: map[string]any{
			"appletConnectMethod":    "client",
			"charset":                "default",
			"is_backspace_as_ctrl_h": false,
			"keyboard_layout":        "en-us-qwerty",
			"rdp_client_option":      []string{"full_screen", "drives_redirect"},
			"rdp_color_quality":      "32",
			"rdp_resolution":         "auto",
			"rdp_smart_size":         "0",
			"resolution":             "auto",
		},
	}
	superToken, _, err := client.Auth.CreateSuperConnectionToken(ctx, superReq)
	if err != nil {
		log.Fatalf("CreateSuperConnectionToken: %v", err)
	}
	fmt.Printf("Super Token ID: %s\n", superToken.ID)

	// ============================================================
	// 6. Get SSO login URL (equivalent to Python get_login_url)
	// ============================================================
	fmt.Println("\n=== 6. Get SSO Login URL ===")
	ssoReq := &model.SSOLoginRequest{
		Username: profile.Username,
		Next:     fmt.Sprintf("/lion/connect/?token=%s", superToken.ID),
	}
	ssoResult, _, err := client.Auth.SSOLoginURL(ctx, ssoReq)
	if err != nil {
		log.Printf("SSOLoginURL: %v (SSO may not be enabled)", err)
	} else {
		prettyPrint("SSO login result", ssoResult)
		if url, ok := ssoResult["login_url"].(string); ok {
			fmt.Printf("Login URL: %s\n", url)
		}
	}

	// ============================================================
	// Bonus: SSO login to Luna web terminal
	// ============================================================
	fmt.Println("\n=== Bonus: SSO Login to Luna ===")
	ssoReq2 := &model.SSOLoginRequest{
		Username: profile.Username,
		Next:     "/luna/",
	}
	ssoResult2, _, err := client.Auth.SSOLoginURL(ctx, ssoReq2)
	if err != nil {
		log.Printf("SSOLoginURL (luna): %v", err)
	} else {
		prettyPrint("Luna SSO result", ssoResult2)
	}

	fmt.Println("\n✓ Connection token flow complete")
}

func prettyPrint(label string, v any) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("%s: %v\n", label, v)
		return
	}
	fmt.Printf("%s:\n%s\n", label, string(data))
}
