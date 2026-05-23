package inbox

import (
	"fmt"
	"context"
    gmailapi "google.golang.org/api/gmail/v1"
    "google.golang.org/api/option"
	"github.com/Platon223/Grail/internal/domain/auth"
)

func GetMainInboxEmails() (*gmailapi.Service, error) {	

	if !auth.HasCredentials() {
		err := fmt.Errorf("Not set up yet. Run: grail auth setup")
        return nil, err
    }

    config, err := auth.LoadConfig()
    if err != nil {
        return nil, err
    }

    token, err := auth.GetClient(config)
    if err != nil {
        return nil, err
    }

    ctx := context.Background()
    tokenSource := config.TokenSource(ctx, token)
    svc, err := gmailapi.NewService(ctx, option.WithTokenSource(tokenSource))
    if err != nil {
        return nil, err
    }

	return svc, nil
}
