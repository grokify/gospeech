package main

type IbmTranscriptionResults struct {
	Results       []Result       `json:"results"`
	ResultIndex   int            `json:"resultIndex"`
	SpeakerLabels []SpeakerLabel `json:"speakerLabels"`
}

type Result struct {
	FinalResults bool          `json:"finalResults"`
	Alternatives []Alternative `json:"alternatives"`
}

type Alternative struct {
	Transcript     string      `json:"transcript"`
	Confidence     float32     `json:"confidence"`
	Timestamps     []Timestamp `json:"timestamps"`
	WordConfidence []WordConfidence
}

type Timestamp struct {
	EndTime   float32 `json:"endTime"`
	StartTime float32 `json:"startTime"`
	Word      string  `json:"word"`
}

type WordConfidence struct {
	Confidence float32 `json:"confidence"`
	Word       string  `json:"word"`
}

// SpeakerLabel isan object for speaker. From and to indicate
// transcript timestamp in seconds
type SpeakerLabel struct {
	From         float32 `json:"from"`
	To           float32 `json:"to"`
	Speaker      int     `json:"speaker"`
	Confidence   float32 `json:"confidence"`
	FinalResults bool    `json:"finalResults"`
}
