package diarization

import (
	"github.com/grokify/gotilla/time/timeutil"
)

// TranscriptToRTTM converts a transcript struct to a RTTM
// struct.
func TranscriptToRTTM(txn *Transcript) RTTM {
	rttm := NewRTTM()
	l := len(txn.Turns)
	for i, turn := range txn.Turns {
		rttmTurn := RTTMTurn{
			SpeakerName: turn.SpeakerName,
			TurnOnset:   turn.TurnOnset}
		if i == l-1 {
			rttmTurn.TurnDuration = timeutil.SubDuration(txn.TotalDuration, turn.TurnOnset)
		} else {
			turnNext := txn.Turns[i+1]
			rttmTurn.TurnDuration = timeutil.SubDuration(turnNext.TurnOnset, turn.TurnOnset)
		}
		rttm.Turns = append(rttm.Turns, rttmTurn)
	}
	return rttm
}
