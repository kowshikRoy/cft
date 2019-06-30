package cmd

type Standings struct {
	Status  string `json:"status"`
	Comment string `json:"comment"`
	Result  Result `json:"result"`
}
type Contest struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	Type                string `json:"type"`
	Phase               string `json:"phase"`
	Frozen              bool   `json:"frozen"`
	DurationSeconds     int    `json:"durationSeconds"`
	StartTimeSeconds    int    `json:"startTimeSeconds"`
	RelativeTimeSeconds int    `json:"relativeTimeSeconds"`
}
type Problems struct {
	ContestID int      `json:"contestId"`
	Index     string   `json:"index"`
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	Points    float64  `json:"points"`
	Rating    int      `json:"rating"`
	Tags      []string `json:"tags"`
}
type Members struct {
	Handle string `json:"handle"`
}

type ProblemResults struct {
	Points                    float64 `json:"points"`
	RejectedAttemptCount      int     `json:"rejectedAttemptCount"`
	Type                      string  `json:"type"`
	BestSubmissionTimeSeconds int     `json:"bestSubmissionTimeSeconds,omitempty"`
}

type Party struct {
	ContestID        int       `json:"contestId"`
	Members          []Members `json:"members"`
	ParticipantType  string    `json:"participantType"`
	TeamID           int       `json:"teamId"`
	TeamName         string    `json:"teamName"`
	Ghost            bool      `json:"ghost"`
	Room             int       `json:"room"`
	StartTimeSeconds int       `json:"startTimeSeconds"`
}
type Rows struct {
	Party                 Party            `json:"party,omitempty"`
	Rank                  int              `json:"rank"`
	Points                float64          `json:"points"`
	Penalty               int              `json:"penalty"`
	SuccessfulHackCount   int              `json:"successfulHackCount"`
	UnsuccessfulHackCount int              `json:"unsuccessfulHackCount"`
	ProblemResults        []ProblemResults `json:"problemResults"`
}
type Row struct {
}
type Result struct {
	Contest  Contest    `json:"contest"`
	Problems []Problems `json:"problems"`
	Rows     []Rows     `json:"rows"`
}
