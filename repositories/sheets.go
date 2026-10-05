package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"naughtfound.github.io/go2phd/models"
)

type SheetsRepository struct {
	service       *sheets.Service
	spreadsheetID string
}

func NewSheetsRepository(credentialsPath string, spreadsheetID string) (*SheetsRepository, error) {
	ctx := context.Background()

	b, err := os.ReadFile(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read credentials file: %w", err)
	}

	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret: %w", err)
	}

	client := getClient(config)
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Sheets client: %w", err)
	}

	return &SheetsRepository{
		service:       srv,
		spreadsheetID: spreadsheetID,
	}, nil
}

func NewSheetsRepositoryFromRefreshToken(clientID, clientSecret, refreshToken, spreadsheetID string) (*SheetsRepository, error) {
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

func (r *SheetsRepository) GetAllPositions(sheetName string) ([]models.Position, error) {
	rows, err := r.ReadAllRows(sheetName)
	if err != nil {
		return nil, err
	}

	if len(rows) <= 1 {
		return []models.Position{}, nil
	}

	var positions []models.Position

	for i, row := range rows[1:] {
		pos := parseRowToPosition(i+1, row)
		positions = append(positions, pos)
	}

	return positions, nil
}

func parseRowToPosition(id int, row []any) models.Position {
	pos := models.Position{
		ID: id,
	}

	getCol := func(idx int) string {
		if idx < len(row) && row[idx] != nil {
			return strings.TrimSpace(fmt.Sprintf("%v", row[idx]))
		}
		return ""
	}

	// Column Mapping (assuming sheet column order):
	// Col 0: Title
	// Col 1: Priority
	// Col 2: Tags (comma-separated: "AI, Go, Microservices")
	// Col 3: Status
	// Col 4: University Name
	// Col 5: Due date (e.g. "2026-11-15" or "2026-11-15T00:00:00Z")
	// Col 6: RollingBased ("true"/"false" or "TRUE"/"FALSE")
	// Col 7: Links (comma-separated or single URL)

	pos.Title = getCol(0)
	pos.Priority = models.Priority(getCol(1))

	rawTags := getCol(2)
	if rawTags != "" {
		tagList := strings.SplitSeq(rawTags, ",")
		for t := range tagList {
			if trimmed := strings.TrimSpace(t); trimmed != "" {
				pos.Tags = append(pos.Tags, trimmed)
			}
		}
	}

	pos.Status = models.Status(getCol(3))

	uniName := getCol(4)
	if uniName != "" {
		pos.University = models.University{Name: uniName}
	}

	rawDueDate := getCol(5)
	if rawDueDate != "" {
		if parsedDate, err := parseDate(rawDueDate); err == nil {
			pos.DueDate = &parsedDate
		}
	}

	rawRolling := strings.ToLower(getCol(6))
	if b, err := strconv.ParseBool(rawRolling); err == nil {
		pos.RollingBased = b
	} else if rawRolling == "yes" || rawRolling == "y" {
		pos.RollingBased = true
	}

	rawLinks := getCol(7)
	if rawLinks != "" {
		linkList := strings.SplitSeq(rawLinks, ",")
		for l := range linkList {
			if trimmed := strings.TrimSpace(l); trimmed != "" {
				pos.Links = append(pos.Links, trimmed)
			}
		}
	}

	return pos
}

func parseDate(dateStr string) (time.Time, error) {
	formats := []string{
		"2006-01-02",
		"2006-01-02T15:04:05Z07:00",
		"01/02/2006",
		"02/01/2006",
		"2006/01/02",
	}

	for _, fmtStr := range formats {
		if t, err := time.Parse(fmtStr, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unsupported date format: %s", dateStr)
}

func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
