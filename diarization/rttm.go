package diarization

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultType        = "SPEAKER"
	defaultFileID      = "FILE1"
	defaultChannelID   = 1
	defaultSpeakerName = "speaker1"
	naString           = "<NA>"
	nsInSec            = 1000000000
)

// RTTM represents a Rich Transcription Time Marked (RTTM) file.
type RTTM struct {
	Turns []RTTMTurn `json:"turns"`
}

// NewRTTM creates a new RTTM struct.
func NewRTTM() RTTM {
	return RTTM{Turns: []RTTMTurn{}}
}

// WriteFile writes the contents in a RTTM file.
func (rttm *RTTM) WriteFile(filename string, perm os.FileMode) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, turn := range rttm.Turns {
		f.WriteString(turn.Encode() + "\n")
	}
	return f.Sync()
}

// RTTMTurn Represents a RTTM file turn
type RTTMTurn struct {
	Type                string
	FileID              string
	ChannelID           int
	TimeOnset           time.Duration
	TimeDuration        time.Duration
	OrthographyField    string
	SpeakerType         string
	SpeakerName         string
	ConfidenceScore     string
	SignalLookaheadTime string
}

func (rttm *RTTMTurn) trimVars() {
	rttm.Type = strings.TrimSpace(rttm.Type)
	rttm.FileID = strings.TrimSpace(rttm.FileID)
	rttm.OrthographyField = strings.TrimSpace(rttm.OrthographyField)
	rttm.SpeakerType = strings.TrimSpace(rttm.SpeakerType)
	rttm.SpeakerName = strings.TrimSpace(rttm.SpeakerName)
	rttm.ConfidenceScore = strings.TrimSpace(rttm.ConfidenceScore)
	rttm.SignalLookaheadTime = strings.TrimSpace(rttm.SignalLookaheadTime)
}

// Encode returns a RTTM turn string.
func (rttm *RTTMTurn) Encode() string {
	return strings.Join(rttm.encodeSlice(), " ")
}

// EncodeSlice returns RTTM encoded turn as a slice
func (rttm *RTTMTurn) encodeSlice() []string {
	rttm.trimVars()
	parts := []string{}
	// var1 = Type
	if len(rttm.Type) > 0 {
		parts = append(parts, rttm.Type)
	} else {
		parts = append(parts, defaultType)
	}
	// var2 = FileID
	if len(rttm.FileID) > 0 {
		parts = append(parts, rttm.FileID)
	} else {
		parts = append(parts, defaultFileID)
	}
	// var3 = ChannelID
	if rttm.ChannelID > 0 {
		parts = append(parts, strconv.Itoa(rttm.ChannelID))
	} else {
		parts = append(parts, strconv.Itoa(defaultChannelID))
	}
	// var4 = TimeOnset
	parts = append(parts, fmt.Sprintf("%.3f",
		float32(rttm.TimeOnset)/float32(nsInSec)))
	// var5 = TimeDuration
	parts = append(parts, fmt.Sprintf("%.3f",
		float32(rttm.TimeDuration)/float32(nsInSec)))
	// var6 = OrthographyField
	if len(rttm.OrthographyField) > 0 {
		parts = append(parts, rttm.OrthographyField)
	} else {
		parts = append(parts, naString)
	}
	// var7 = SpeakerType
	if len(rttm.SpeakerType) > 0 {
		parts = append(parts, rttm.SpeakerType)
	} else {
		parts = append(parts, naString)
	}
	// var8 = SpeakerName
	if len(rttm.SpeakerName) > 0 {
		parts = append(parts, rttm.SpeakerName)
	} else {
		parts = append(parts, defaultSpeakerName)
	}
	// var9 = ConfidenceScore
	if len(rttm.ConfidenceScore) > 0 {
		parts = append(parts, rttm.ConfidenceScore)
	} else {
		parts = append(parts, naString)
	}
	// var10 = SignalLookaheadTime
	if len(rttm.SignalLookaheadTime) > 0 {
		parts = append(parts, rttm.SignalLookaheadTime)
	} else {
		parts = append(parts, naString)
	}
	return parts
}
