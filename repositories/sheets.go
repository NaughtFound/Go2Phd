package repositories

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SheetsRepository struct {
	service       *sheets.Service
	spreadsheetID string
}

func NewSheetsRepository(clientID, clientSecret, refreshToken, spreadsheetID string) (*SheetsRepository, error) {
	ctx := context.Background()

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{sheets.SpreadsheetsReadonlyScope},
	}

	tok := &oauth2.Token{
		RefreshToken: refreshToken,
	}

	client := config.Client(ctx, tok)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Sheets client: %w", err)
	}

	return &SheetsRepository{
		service:       srv,
		spreadsheetID: spreadsheetID,
	}, nil
}

func (r *SheetsRepository) ReadRange(readRange string) ([][]any, error) {
	resp, err := r.service.Spreadsheets.Values.Get(r.spreadsheetID, readRange).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from sheet: %w", err)
	}

	return resp.Values, nil
}

func (r *SheetsRepository) ReadAllRows(sheetName string) ([][]any, error) {
	resp, err := r.service.Spreadsheets.Values.Get(r.spreadsheetID, sheetName).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from sheet: %w", err)
	}

	return resp.Values, nil
}
