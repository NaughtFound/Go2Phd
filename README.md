# Go2PhD - Academic Application Tracker & Analytics API

**Go2PhD** is a Go-powered stateless REST backend designed to track and analyze PhD and academic position applications. By integrating directly with Google Sheets as a dynamic database via Google OAuth 2.0, Go2PhD enables users to manage application pipelines, due dates, priorities, and generate real-time analytics without managing a traditional SQL/NoSQL database.

---

## Features

* 🔐 **Stateless Google OAuth 2.0 Auth**: Authenticate seamlessly using Google OAuth consent and refresh tokens passed via HTTP headers.
* 📊 **Real-Time Analytics & Stats**: Aggregated count mappings for application statuses, priorities, tags, and rolling applications.
* 🌍 **Dynamic Grouping**: Group position metrics dynamically by `country`, `university`, `priority`, `status`, or `tags`.
* 🟩 **Google Sheets Database**: Reads directly from Google Sheets using `FORMATTED_VALUE` rendering for clean text parsing.
* ⚡ **Lightweight & Fast**: Built with Go and the Gin Web Framework.

---

## Prerequisites

* [Go 1.20+](https://golang.org/doc/install) installed on your system.
* A **Google Cloud Platform (GCP)** account.
* A **Google Sheet** configured to store your application data.

---

## Step 1: Set Up Google Cloud Platform (GCP) Credentials

1. Go to the [Google Cloud Console](https://console.cloud.google.com/).
2. Create a new project named **Go2PhD Tracker**.
3. Enable required APIs:
* Go to **APIs & Services > Library**.
* Search for **Google Sheets API** and click **Enable**.


4. Configure OAuth Consent Screen:
* Go to **APIs & Services > OAuth consent screen**.
* Select **External** user type and fill in app contact details.
* Under **Test users**, add your personal Google email address.


5. Create Credentials:
* Go to **APIs & Services > Credentials** > **Create Credentials** > **OAuth client ID**.
* Select **Web application**.
* Set **Authorized redirect URIs** to `http://localhost:8080/auth/callback` (or your frontend callback URL).
* Copy your **Client ID** and **Client Secret**.



---

## Step 2: Prepare Your Google Sheet

Create a Google Sheet and copy its **Spreadsheet ID** from the URL bar.

Set up the top header row in your sheet with the following structure:

| Title | Priority | Tags | Status | University Name | Country | Due Date | Rolling Based | Links |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| PhD in Computer Vision | High | Machine Learning, CV | Applied | Technical University of Munich | Germany | 5/1/2026 | false | [https://example.com](https://example.com) |

### Allowed Enum Values

* **Priority**: `High`, `Medium`, `Low`
* **Status**: `Not started`, `In progress`, `Applied`, `Emailed`, `Interview`, `Rejected`, `Shortlisted`, `Accepted`, `Canceled`, `No Answer`

---

## Step 3: Installation & Configuration

Clone the repository and install dependencies:

```bash
git clone https://github.com/naughtfound/go2phd.git
cd go2phd
go mod download

```

Create a `.env` file in the project root:

```env
HOST=localhost
PORT=8080
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_REDIRECT_URI=http://localhost:8080/auth/callback

```

Run the application:

```bash
go run main.go

```

The server will start at `http://localhost:8080`.

---

## API Reference & Usage

Since the backend is stateless, API endpoints rely on custom HTTP headers for runtime authentication and sheet configuration.

### Required Request Headers

| Header Name | Required | Description |
| --- | --- | --- |
| `X-Refresh-Token` | **Yes** | Google OAuth refresh token obtained during login |
| `X-Spreadsheet-ID` | **Yes** | Google Sheet ID containing your application data |

---

### 1. Get All Positions

Retrieves all positions parsed from the specified Google Sheet tab.

* **Endpoint**: `GET /positions`
* **Query Parameters**:
* `sheet` *(optional, default: `Sheet1`)*: Worksheet tab name.



**Curl Example**:

```bash
curl -X GET "http://localhost:8080/positions?sheet=Sheet1" \
  -H "X-Refresh-Token: YOUR_REFRESH_TOKEN" \
  -H "X-Spreadsheet-ID: YOUR_SPREADSHEET_ID"

```

**Response Example**:

```json
[
  {
    "id": 1,
    "title": "PhD in Computer Vision",
    "priority": "High",
    "tags": ["Machine Learning", "CV"],
    "status": "Applied",
    "university": {
      "name": "Technical University of Munich",
      "country": "Germany"
    },
    "due_date": "2026-05-01",
    "rolling_based": false,
    "links": ["https://example.com"]
  }
]

```

---

### 2. Get Statistics & Analytics

Calculates aggregate count statistics (overall and optional grouped breakdowns).

* **Endpoint**: `GET /positions/stats`
* **Query Parameters**:
* `sheet` *(optional, default: `Sheet1`)*: Worksheet tab name.
* `group_by` *(optional)*: Dimension to group by (`country`, `university`, `status`, `priority`, or `tags`).



**Curl Example (Grouped by Country)**:

```bash
curl -X GET "http://localhost:8080/positions/stats?group_by=country" \
  -H "X-Refresh-Token: YOUR_REFRESH_TOKEN" \
  -H "X-Spreadsheet-ID: YOUR_SPREADSHEET_ID"

```

**Response Example**:

```json
{
  "overall": {
    "total_positions": 5,
    "status_breakdown": {
      "Not started": 1,
      "In progress": 1,
      "Applied": 1,
      "Emailed": 0,
      "Interview": 0,
      "Rejected": 1,
      "Shortlisted": 0,
      "Accepted": 1,
      "Canceled": 0,
      "No Answer": 0
    },
    "priority_breakdown": {
      "High": 3,
      "Medium": 2,
      "Low": 0
    },
    "tags_breakdown": {
      "Machine Learning": 2,
      "CV": 1
    },
    "rolling_based_count": 2
  },
  "grouped": [
    {
      "group_key": "Germany",
      "stats": {
        "total_positions": 2,
        "status_breakdown": {
          "Not started": 0,
          "In progress": 0,
          "Applied": 1,
          "Emailed": 0,
          "Interview": 0,
          "Rejected": 0,
          "Shortlisted": 0,
          "Accepted": 1,
          "Canceled": 0,
          "No Answer": 0
        },
        "priority_breakdown": {
          "High": 2,
          "Medium": 0,
          "Low": 0
        },
        "tags_breakdown": {
          "Machine Learning": 1
        },
        "rolling_based_count": 1
      }
    }
  ]
}

```

---

## Project Architecture

```
.
├── go.mod
├── go.sum
├── main.go
├── LICENSE
├── README.md
├── handlers/
│   ├── auth.go            # Authentication & OAuth HTTP handlers
│   ├── position.go        # Positions & statistics endpoints
│   └── university.go      # University listing handlers
├── models/
│   ├── position.go        # Application models & priority/status enums
│   ├── stats.go           # Overall & Grouped statistics DTOs
│   └── university.go      # University structure definitions
├── repositories/
│   ├── sheets.go          # Google Sheets API client & row parsers
│   └── university.go      # University data access logic
└── services/
    ├── auth.go            # Google OAuth token exchange & service flow
    ├── position.go        # Application position retrieval logic
    ├── stats.go           # Dynamic analytics aggregation & grouping
    └── university.go      # University business logic

```

---

## License

This project is licensed under the [MIT License](https://www.google.com/search?q=LICENSE).