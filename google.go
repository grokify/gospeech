package diarization

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	tu "github.com/grokify/gotilla/time/timeutil"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1p1beta1"
)

const speakerNamePrefix string = "S"

/*

https://godoc.org/google.golang.org/genproto/googleapis/cloud/speech/v1#LongRunningRecognizeResponse

*/

func LongRunningRecognizeResponseToTranscript(res *speechpb.LongRunningRecognizeResponse) (Transcript, error) {
	txn := NewTranscript()
	lastRes := res.Results[len(res.Results)-1]
	fmt.Println("NUM_ALTS[%v]\n", len(lastRes.Alternatives))
	if len(lastRes.Alternatives) != 1 {
		return txn, fmt.Errorf("E_NOT_1_LAST_ALTERNATIVE [%v]", len(lastRes.Alternatives))
	}
	alt := lastRes.Alternatives[0]
	curSpeakerTag := int32(0)
	curTurn := Turn{}
	for _, word := range alt.Words {
		if word.SpeakerTag <= 0 {
			return txn, fmt.Errorf("E_NO_SPEAKER_TAG [%v]", word.SpeakerTag)
		}
		fmt.Printf(".%v.", word.SpeakerTag)
		if curSpeakerTag == 0 {
			curSpeakerTag = word.SpeakerTag
		} else if curSpeakerTag != word.SpeakerTag {
			newTurn := Turn{
				SpeakerName: curTurn.SpeakerName,
				Text:        curTurn.Text,
				TimeBegin:   curTurn.TimeBegin,
				TimeEnd:     curTurn.TimeEnd}
			txn.Turns = append(txn.Turns, newTurn)
			curTurn = Turn{}
		}
		curSpeakerTag = word.SpeakerTag
		curTurn.Text += " " + word.Word
		curTurn.SpeakerName = speakerNamePrefix + strconv.Itoa(int(word.SpeakerTag))
		thisBeginDur := tu.DurationFromProtobuf(word.StartTime)
		thisEndDur := tu.DurationFromProtobuf(word.EndTime)
		if curTurn.TimeBegin.Nanoseconds() == 0 ||
			thisBeginDur.Nanoseconds() < curTurn.TimeBegin.Nanoseconds() {
			curTurn.TimeBegin = thisBeginDur
		}
		if curTurn.TimeEnd.Nanoseconds() == 0 ||
			thisEndDur.Nanoseconds() > curTurn.TimeEnd.Nanoseconds() {
			curTurn.TimeEnd = thisEndDur
		}
	}
	txn.Turns = append(txn.Turns, curTurn)
	txn.BuildSpeakers()
	return txn, nil
}

func ReadLongRunningRecognizeResponseFile(file string) (*speechpb.LongRunningRecognizeResponse, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return ReadLongRunningRecognizeResponse(bytes)
}

func ReadLongRunningRecognizeResponse(bytes []byte) (*speechpb.LongRunningRecognizeResponse, error) {
	res := &speechpb.LongRunningRecognizeResponse{}
	return res, json.Unmarshal(bytes, res)
}