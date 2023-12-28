package interactor

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/c18t/nippo-cli/internal/adapter/presenter"
	"github.com/c18t/nippo-cli/internal/core"
	"github.com/c18t/nippo-cli/internal/usecase/port"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

type initSaveDriveTokenInteractor struct {
	presenter presenter.InitSaveDriveTokenPresenter
}

func NewInitSaveDriveTokenInteractor(presenter presenter.InitSaveDriveTokenPresenter) port.InitSaveDriveTokenUsecase {
	return &initSaveDriveTokenInteractor{presenter}
}

func (u *initSaveDriveTokenInteractor) Handle(input *port.InitSaveDriveTokenUsecaseInputData) {
	output := &port.InitSaveDriveTokenUsecaseOutputData{}

	dataDir := core.Cfg.GetDataDir()
	b, err := os.ReadFile(path.Join(dataDir, "credentials.json"))
	if err != nil {
		u.presenter.Suspend(fmt.Errorf("unable to read client secret file: %v", err))
		return
	}

	config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope, drive.DriveReadonlyScope)
	if err != nil {
		u.presenter.Suspend(fmt.Errorf("unable to parse client secret file to config: %v", err))
		return
	}

	tok := getTokenFromWeb(config)
	saveToken(path.Join(dataDir, "token.json"), tok)

	u.presenter.Complete(output)
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Printf("Unable to cache oauth token: %v\n", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		fmt.Printf("Unable to read authorization code %v\n", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		fmt.Printf("Unable to retrieve token from web %v\n", err)
	}
	return tok
}
