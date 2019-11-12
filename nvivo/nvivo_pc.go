package nvivo

import (
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/grokify/go-diarization"
	tu "github.com/grokify/gotilla/time/timeutil"
)

const (
	nVivoLinePCFormat string = `^((\d{2}):(\d{2}):(\d{2}),(\d{3}))\s+-\s+((\d{2}):(\d{2}):(\d{2}),(\d{3}))\s+(.+)\t([^\t]+)$`
	nVivoLineVars     int    = 13
)

var rxNVivoLinePC *regexp.Regexp = regexp.MustCompile(nVivoLinePCFormat)

// ParseNVivoPcFile parses a NVivo PC file. This file has
// begin and end times for each turn.
func ParseNVivoPcFile(file string) (*diarization.Transcript, error) {
	tr := &diarization.Transcript{
		Turns:    []diarization.Turn{},
		Speakers: diarization.SpeakerSet{SpeakersMap: map[string]diarization.Speaker{}}}
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return tr, err
	}
	lines := strings.Split(string(bytes), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		turn, err := ParseNVivoTurnLine(line)
		if err != nil {
			return tr, err
		}
		tr.Turns = append(tr.Turns, turn)
		tr.Speakers.AddTurn(turn)
	}

	return tr, nil
}

func ParseNVivoTurnLine(turnString string) (diarization.Turn, error) {
	turn := diarization.Turn{}
	m := rxNVivoLinePC.FindStringSubmatch(turnString)
	if len(m) != nVivoLineVars {
		return turn, nil
	}

	turn.TimeBeginRaw = m[1]
	durBegin, err := tu.ParseDurationInfoStrings(m[2], m[3], m[4], m[5], "", "")
	if err != nil {
		return turn, err
	}
	turn.TimeBegin = durBegin.Duration()
	turn.TimeEndRaw = m[6]
	durEnd, err := tu.ParseDurationInfoStrings(m[7], m[8], m[9], m[10], "", "")
	if err != nil {
		return turn, err
	}
	turn.TimeEnd = durEnd.Duration()
	turn.Text = strings.TrimSpace(m[11])
	turn.SpeakerName = strings.TrimSpace(m[12])
	return turn, nil
}