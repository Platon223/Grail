package auth

import (
	"fmt"
	"os"
	"context"
	"os/exec"
	"runtime"
	"encoding/json"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gmailapi "google.golang.org/api/gmail/v1"
	//"google.golang.org/api/option"
)


func configDir() string {
	home, _ := os.UserHomeDir()
	dir := home + "/.config/grail"
	os.Mkdir(dir, 0700)
	return dir
}


func credentialsPath() string { return configDir() + "/credentials.json" }
func tokenPath() string { return configDir() + "/token.json" } 

func HasCredentials() bool {
	_, err := os.Stat(credentialsPath())
	return err == nil
}

func HasToken() bool {
	_, err := os.Stat(tokenPath())
	return err == nil
}


func SaveCredentials(srcPath string) error {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("couldn't read credentials file: %w", err)
	}
	return os.WriteFile(credentialsPath(), data, 0600)
}


func LoadConfig() (*oauth2.Config, error) {
    data, err := os.ReadFile(credentialsPath())
    if err != nil {
        return nil, fmt.Errorf("credentials not found, run: grail auth setup")
    }
    config, err := google.ConfigFromJSON(data, gmailapi.GmailReadonlyScope)
    if err != nil {
        return nil, err
    }
    return config, nil
}


func GetClient(config *oauth2.Config) (*oauth2.Token, error) {
	if HasToken() {
		token, err := loadToken()
		if err == nil && token.Valid() {
			return token, nil
		}

		if err == nil {
			tokenSource := config.TokenSource(context.Background(), token)
			newToken, err := tokenSource.Token()
			if err == nil {
				saveToken(newToken)
				return newToken, nil
			}
		}
	}

	return tokenFromBrowser(config)
}

func tokenFromBrowser(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Println("Opening browser for Gmail authorization...")
	openBrowser(authURL)
	fmt.Printf("\nIf browser didn't open, visit this URL:\n%s\n\n", authURL)
	fmt.Print("Paste the authorization code here: ")

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		return nil, fmt.Errorf("failed to read code: %w", err)
	}

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	saveToken(token)
	fmt.Println("\n✓ Authorization successful")
	return token, nil
}

func loadToken() (*oauth2.Token, error) {
	f, err := os.Open(tokenPath())
	if err != nil {
		return nil, err
	}
	defer f.Close()

	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	return token, err
}

func saveToken(token *oauth2.Token) {
	f, err := os.OpenFile(tokenPath(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func openBrowser(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	default:
		cmd = "xdg-open"
		args = []string{url}
	}

	exec.Command(cmd, args...).Start()
}
