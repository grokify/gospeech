package diarization

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
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
	Turns         []Turn        `json:"turns"`
	SpeakerNames  []string      `json:"speakerNames"`
	TotalDuration time.Duration `json:"totalDuration"`
}

// Turn represent what has been spoken.
type Turn struct {
	TurnOnset    time.Duration `json:"turnOnset"`
	TurnOnsetRaw string        `json:"turnOnsetRaw"`
	SpeakerName  string        `json:"speakerName"`
	Text         string        `json:"text"`
}

// ParseTranscribeMeFile reads a TranscribeMe.com text file.
func ParseTranscribeMeFile(filename string) (*Transcript, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ParseTranscribeMe(bytes)
}

// ParseTranscribeMe reads TranscribeMe.com text file data.
func ParseTranscribeMe(bytes []byte) (*Transcript, error) {
	txn := &Transcript{
		Turns:        []Turn{},
		SpeakerNames: []string{}}
	speakersMap := map[string]int{}
	lines := strings.Split(string(bytes), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		m := rxTurn.FindStringSubmatch(line)
		if len(m) > 0 {
			p := Turn{
				SpeakerName:  strings.TrimSpace(m[1]),
				TurnOnsetRaw: strings.TrimSpace(m[2]),
				Text:         strings.TrimSpace(m[6])}
			durStr := fmt.Sprintf("%vh%vm%vs", m[3], m[4], m[5])
			dur, err := time.ParseDuration(durStr)
			if err != nil {
				panic("Y")
				return nil, err
			}
			p.TurnOnset = dur
			txn.Turns = append(txn.Turns, p)
			speakersMap[p.SpeakerName] = 1
			continue
		}
		m2 := rxEnd.FindStringSubmatch(line)
		if len(m2) > 0 {
			durStr := fmt.Sprintf("%vh%vm%vs", m2[2], m2[3], m2[4])
			dur, err := time.ParseDuration(durStr)
			if err != nil {
				panic("Z")
				return nil, err
			}
			txn.TotalDuration = dur
		}
	}
	for s := range speakersMap {
		txn.SpeakerNames = append(txn.SpeakerNames, s)
	}
	sort.Strings(txn.SpeakerNames)
	return txn, nil
}

// WriteJSON writes out the parsed transcript as a pretty
// printed file.
func (txn *Transcript) WriteJSON(filename string, perm os.FileMode) error {
	data, err := json.MarshalIndent(txn, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, perm)
}
