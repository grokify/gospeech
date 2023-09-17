package deepgram

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/grokify/gospeech/diarization"
	"github.com/grokify/mogo/time/timeutil"
)

// Transcript represents a Deepgram transcript.
type Transcript struct {
	Metadata Metadata `json:"metadata"`
	Results  Results  `json:"results"`
}

// NewTranscriptFile attempts to read a Deepgram
// transcript file.
func NewTranscriptFile(file string) (*Transcript, error) {
	txn := Transcript{}
	bytes, err := os.ReadFile(file)
	if err != nil {
		return &txn, err
	}
	err = json.Unmarshal(bytes, &txn)
	return &txn, err
}

func (dtxn *Transcript) ChannelCount() int {
	return len(dtxn.Results.Channels)
}

func (dtxn *Transcript) MaxAltCount() int {
	maxAltCount := 0
	for _, chanRes := range dtxn.Results.Channels {
		if len(chanRes.Alternatives) > maxAltCount {
			maxAltCount = len(chanRes.Alternatives)
		}
	}
	return maxAltCount
}

// Validate checks the transcript. This initially checks
// to see if diarized word count matches transcript words
// defined by space separted strings.
func (dtxn *Transcript) Validate() error {
	for i, chanRes := range dtxn.Results.Channels {
		//fmt.Printf("NUM_CHANNEL_ALTERNATIVES: CHAN[%v]ALT[%v]\n", i, len(chanRes.Alternatives))
		for j, alt := range chanRes.Alternatives {
			txnWords := len(strings.Fields(alt.Transcript))
			spkWords := len(alt.Words)
			if txnWords != spkWords {
				return fmt.Errorf("E_MISMATCH CHAN [%v], ALT [%v] TXN_WORDS [%v] SPK_WORDS [%v]", i, j, txnWords, spkWords)
			}
			//fmt.Printf("CHAN [%v], ALT [%v] TXN_WORDS [%v] SPK_WORDS [%v]\n", i, j, txnWords, spkWords)
		}
	}
	return nil
}

type Metadata struct {
	Sha256           string    `json:"sha256"`
	Created          time.Time `json:"created"`
	Duration         float64   `json:"duration"`
	Channels         int       `json:"channels"`
	TranscriptionKey string    `json:"transcription_key"`
	RequestID        string    `json:"request_id"`
}

type Results struct {
	Channels []ChannelResult `json:"channels"`
	//Alternatives []Alternative   `json:"alternatives"`
}

type ChannelResult struct {
	Alternatives []Alternative `json:"alternatives"`
}

type Alternative struct {
	Transcript string  `json:"transcript"`
	Confidence float64 `json:"confidence"`
	Words      []Word  `json:"words"`
}

type Word struct {
	Word       string  `json:"word"`
	Start      float32 `json:"start"`
	End        float32 `json:"end"`
	Confidence float64 `json:"confidence"`
	Speaker    int     `json:"speaker"`
}

func CanonicalTranscriptFromDeepgram(dtxn *Transcript) (*diarization.Transcript, error) {
	err := dtxn.Validate()
	if err != nil {
		return nil, err
	}
	if dtxn.ChannelCount() != 1 {
		return nil, fmt.Errorf("E_CHAN_COUNT_NE_1 [%d]", dtxn.ChannelCount())
	}
	if dtxn.MaxAltCount() != 1 {
		return nil, fmt.Errorf("E_CHAN_COUNT_NE_1 [%d]", dtxn.MaxAltCount())
	}

	turns := []diarization.Turn{}
	curTurn := diarization.Turn{}

	chan0 := dtxn.Results.Channels[0]
	alt0 := chan0.Alternatives[0]
	formattedWords := strings.Fields(alt0.Transcript)

	firstWord := true
	for i, word := range alt0.Words {
		speakerName := "speaker" + strconv.Itoa(word.Speaker)
		if firstWord {
			curTurn = newTurnForWord(word, formattedWords[i])
			firstWord = false
			continue
		} else if speakerName == curTurn.SpeakerName {
			// curTurn.TimeEnd = timeutil.NewDuration(0, 0, 0, float64(word.End), 0)
			curTurn.TimeEnd = time.Duration(int64(float64(word.End) * float64(timeutil.NanosPerSecond)))
			curTurn.TimeEndRaw = strconv.FormatFloat(float64(word.End), 'E', -1, 64)
			curTurn.Text += " " + formattedWords[i]
		} else {
			turns = append(turns, curTurn)
			curTurn = newTurnForWord(word, formattedWords[i])
		}
	}
	if len(curTurn.Text) > 0 {
		turns = append(turns, curTurn)
	}
	for _, turn := range turns {
		turn.Text = strings.Join(strings.Fields(turn.Text), " ")
	}
	ctxn := diarization.NewTranscript()
	ctxn.Turns = turns
	ctxn.Inflate()
	return &ctxn, nil
}

func newTurnForWord(word Word, formattedWord string) diarization.Turn {
	return diarization.Turn{
		SpeakerName:  "speaker" + strconv.Itoa(word.Speaker),
		TimeBegin:    time.Duration(int64(float64(word.Start) * float64(timeutil.NanosPerSecond))),
		TimeBeginRaw: strconv.FormatFloat(float64(word.Start), 'E', -1, 64),
		TimeEnd:      time.Duration(int64(float64(word.End) * float64(timeutil.NanosPerSecond))),
		TimeEndRaw:   strconv.FormatFloat(float64(word.End), 'E', -1, 64),
		Text:         formattedWord,
		//TimeBegin:    timeutil.NewDuration(0, 0, 0, float64(word.Start), 0),
		//TimeEnd:      timeutil.NewDuration(0, 0, 0, float64(word.End), 0),
	}
}
