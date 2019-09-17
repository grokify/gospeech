package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/grokify/go-diarization"
	"github.com/grokify/gotilla/fmt/fmtutil"
	iom "github.com/grokify/gotilla/io/ioutilmore"
	tu "github.com/grokify/gotilla/time/timeutil"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1p1beta1"
)

const speakerNamePrefix string = "S"

/*

https://godoc.org/google.golang.org/genproto/googleapis/cloud/speech/v1#LongRunningRecognizeResponse

*/

func main() {
	file := "../mongodb-is-web-scale/web-scale_b2F-DItXtZs.mp3_tr_google_standard.json"

	//file = "/Users/john.wang/jwdev/JGo/gopath/src/github.com/grokify/golang-samples/speech/captionasync/episode_1_mongo_db_is_web_scale_b2F-DItXtZs.mp3_transcript_gcs-standard.json"

	res, err := ReadLongRunningRecognizeResponseFile(file)
	if err != nil {
		log.Fatal(err)
	}
	// fmtutil.PrintJSON(res)
	fmt.Printf("NUM_RESULTS [%v]\n", len(res.Results))

	for i, resp := range res.Results {
		for _, alt := range resp.Alternatives {
			// Printf("[%v] %v\n", i, alt.Transcript)
			for _, wrd := range alt.Words {
				//fmtutil.PrintJSON(wrd)
				if wrd.SpeakerTag > 0 {
					fmt.Printf("%v ", i)
				}
			}
		}
	}

	lastRes := res.Results[len(res.Results)-1]
	fmt.Println("NUM_ALTS[%v]\n", len(lastRes.Alternatives))
	if len(lastRes.Alternatives) != 1 {
		panic("A")
	}
	alt := lastRes.Alternatives[0]
	//turns := []diarization.Turn{}
	txn := diarization.NewTranscript()
	curSpeakerTag := int32(0)
	curTurn := diarization.Turn{}
	for _, word := range alt.Words {
		if word.SpeakerTag <= 0 {
			panic("A")
		}
		fmt.Printf(".%v.", word.SpeakerTag)
		if curSpeakerTag == 0 {
			curSpeakerTag = word.SpeakerTag
		} else if curSpeakerTag != word.SpeakerTag {
			newTurn := diarization.Turn{
				SpeakerName: curTurn.SpeakerName,
				Text:        curTurn.Text,
				TimeBegin:   curTurn.TimeBegin,
				TimeEnd:     curTurn.TimeEnd}
			txn.Turns = append(txn.Turns, newTurn)
			curTurn = diarization.Turn{}
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
	fmtutil.PrintJSON(txn)

	err = iom.WriteFileJSON(file+".json", txn, 0644, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	rttm := diarization.TranscriptToRTTM(&txn)
	err = rttm.WriteFile(file+".rttm", 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DONE")
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
