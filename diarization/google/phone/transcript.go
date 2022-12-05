package phone

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/grokify/go-transcribe/diarization"
)

type Transcript struct {
	Results []Result `json:"results"`
}

func NewTranscriptFile(file string) (*diarization.Transcript, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	gtxn := Transcript{}
	err = json.Unmarshal(bytes, &gtxn)
	if err != nil {
		return nil, err
	}

	return NewTranscriptGoogleResponse(gtxn)
}

func NewTranscriptGoogleResponse(gtxn Transcript) (*diarization.Transcript, error) {
	//ctxn := diarization.NewTranscript()
	if len(gtxn.Results) == 0 {
		return nil, fmt.Errorf("No Results [%v]", len(gtxn.Results))
	}
	lastResult := gtxn.Results[len(gtxn.Results)-1]
	altCount := len(lastResult.Alternatives)
	if altCount != 1 {
		panic("Multiple Alternatives")
	}
	alt := lastResult.Alternatives[0]
	turns := []diarization.Turn{}
	curTurn := diarization.Turn{}
	prevSpeakerTag := int64(-1)
	for i, word := range alt.Words {
		if i == 0 {
			curTurn = newTurnForWord(word)
		} else if word.SpeakerTag != prevSpeakerTag {
			turns = append(turns, curTurn)
			curTurn = newTurnForWord(word)
		} else {
			durationEnd, err := time.ParseDuration(word.EndTime)
			if err != nil {
				panic(err)
			}
			curTurn.TimeEnd = durationEnd
			curTurn.TimeEndRaw = word.EndTime
			curTurn.Text += " " + word.Word
		}
		prevSpeakerTag = word.SpeakerTag
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

type Result struct {
	Alternatives []Alternative `json:"alternatives"`
}

type Alternative struct {
	Transcript string  `json:"transcript"`
	Confidence float64 `json:"confidence"`
	Words      []Word  `json:"words"`
}

type Word struct {
	StartTime  string  `json:"startTime"`
	EndTime    string  `json:"endTime"`
	Word       string  `json:"word"`
	Confidence float64 `json:confidence`
	SpeakerTag int64   `json:"speakerTag"`
}

func newTurnForWord(word Word) diarization.Turn {
	durationBegin, err := time.ParseDuration(word.StartTime)
	if err != nil {
		panic(err)
	}
	durationEnd, err := time.ParseDuration(word.EndTime)
	if err != nil {
		panic(err)
	}
	return diarization.Turn{
		SpeakerName:  "speaker" + strconv.Itoa(int(word.SpeakerTag)),
		TimeBegin:    durationBegin,
		TimeBeginRaw: word.StartTime,
		TimeEnd:      durationEnd,
		TimeEndRaw:   word.EndTime,
		Text:         word.Word}
}
