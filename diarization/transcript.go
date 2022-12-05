package diarization

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/grokify/mogo/time/timeutil"
	"github.com/grokify/mogo/type/stringsutil"
)

// S2: 00:01:20.626 Sure.

const (
	rxTurnFormat string = `^([^:]+):\s+((\d+):(\d+):(\d+.\d+))\s+(.+)$`
	rxEndFormat  string = `^END:\s+((\d+):(\d+):(\d+.\d+))\s*$`
)

var rxTurn *regexp.Regexp = regexp.MustCompile(rxTurnFormat)
var rxEnd *regexp.Regexp = regexp.MustCompile(rxEndFormat)

// Transcript represents a text representation of a conversation.
type Transcript struct {
	Name          string        `json:"name"`
	Turns         []Turn        `json:"turns"`
	Speakers      SpeakerSet    `json:"speakers"`
	TotalDuration time.Duration `json:"totalDuration"`
}

// NewTranscript returns a new Transcript struct.
func NewTranscript() Transcript {
	return Transcript{
		Turns:    []Turn{},
		Speakers: SpeakerSet{SpeakersMap: map[string]Speaker{}}}
}

func (txn *Transcript) Inflate() {
	ss := SpeakerSet{SpeakersMap: map[string]Speaker{}}
	for i, turn := range txn.Turns {
		turn.Text = strings.Join(strings.Fields(turn.Text), " ")
		ss.AddTurn(turn)
		txn.Turns[i] = turn
	}
	txn.Speakers = ss
	if len(txn.Turns) > 0 {
		lastTurn := txn.Turns[len(txn.Turns)-1]
		txn.TotalDuration = lastTurn.TimeEnd
	}
}

// Turn represent what has been spoken.
type Turn struct {
	TimeBegin    time.Duration `json:"turnOnset"`
	TimeBeginRaw string        `json:"turnOnsetRaw"`
	TimeEnd      time.Duration `json:"turnEnd"`
	TimeEndRaw   string        `json:"turnEndRaw"`
	SpeakerName  string        `json:"speakerName"`
	Text         string        `json:"text"`
}

// Duration returns the duration.
func (turn *Turn) Duration() time.Duration {
	return timeutil.SubDuration(turn.TimeEnd, turn.TimeBegin)
}

// ParseTranscribeMeFile reads a TranscribeMe.com text file.
func ParseTranscribeMeFile(filename string) (*Transcript, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ParseTranscribeMe(bytes)
}

// ParseTranscribeMe reads TranscribeMe.com text file data.
func ParseTranscribeMe(bytes []byte) (*Transcript, error) {
	txn := &Transcript{
		Turns:    []Turn{},
		Speakers: SpeakerSet{SpeakersMap: map[string]Speaker{}}}
	speakersMap := map[string]int{}
	lines := strings.Split(string(bytes), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if m := rxTurn.FindStringSubmatch(line); len(m) > 0 {
			p := Turn{
				SpeakerName:  strings.TrimSpace(m[1]),
				TimeBeginRaw: strings.TrimSpace(m[2]),
				Text:         strings.TrimSpace(m[6])}
			dur, err := time.ParseDuration(
				fmt.Sprintf("%sh%sm%ss", m[3], m[4], m[5]))
			if err != nil {
				return nil, err
			}
			p.TimeBegin = dur
			txn.Turns = append(txn.Turns, p)
			speakersMap[p.SpeakerName] = 1
		} else if m := rxEnd.FindStringSubmatch(line); len(m) > 0 {
			dur, err := time.ParseDuration(
				fmt.Sprintf("%sh%sm%ss", m[2], m[3], m[4]))
			if err != nil {
				return nil, err
			}
			txn.TotalDuration = dur
		}
	} /*
		for s := range speakersMap {
			txn.SpeakerNames = append(txn.SpeakerNames, s)
		}
	sort.Strings(txn.SpeakerNames)*/
	return txn, nil
}

// WriteJSON writes out the parsed transcript as a pretty
// printed file.
func (txn *Transcript) WriteJSON(filename string, perm os.FileMode) error {
	data, err := json.MarshalIndent(txn, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, perm)
}

type SpeakerSet struct {
	ByAppearance []string
	SpeakersMap  map[string]Speaker
}

// AddTurn adds a turn to the speaker information
func (ss *SpeakerSet) AddTurn(turn Turn) {
	speakerName := strings.TrimSpace(turn.SpeakerName)
	speaker, ok := ss.SpeakersMap[speakerName]
	if !ok {
		speaker = Speaker{}
		ss.ByAppearance = append(ss.ByAppearance, speakerName)
	}
	speaker.Name = speakerName
	speaker.Turns += 1
	speaker.TotalDuration = timeutil.SumDurations(speaker.TotalDuration, turn.Duration())
	ss.SpeakersMap[speakerName] = speaker
}

// SpeakerNameIndex returns the position where the speaker appears in the transcript.
func (ss *SpeakerSet) SpeakerNameIndex(speakerName string) int {
	return stringsutil.SliceIndex(ss.ByAppearance, speakerName, false, nil)
}

// Speaker represents a speaker including numbers of turns spoken and total duration spoken.
type Speaker struct {
	Name          string
	Turns         int32
	TotalDuration time.Duration
}
