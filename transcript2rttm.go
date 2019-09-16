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
			TimeOnset:   turn.TimeBegin}
		if turn.TimeEnd.Nanoseconds() > turn.TimeBegin.Nanoseconds() {
			rttmTurn.TimeDuration = turn.Duration()
		} else {
			if i == l-1 {
				rttmTurn.TimeDuration = timeutil.SubDuration(txn.TotalDuration, turn.TimeBegin)
			} else {
				turnNext := txn.Turns[i+1]
				rttmTurn.TimeDuration = timeutil.SubDuration(turnNext.TimeBegin, turn.TimeBegin)
			}
		}
		rttm.Turns = append(rttm.Turns, rttmTurn)
	}
	return rttm
}
