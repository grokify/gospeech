package main

import (
	"fmt"
	"log"

	"github.com/grokify/go-transcribe/diarization"
	"github.com/grokify/go-transcribe/diarization/google"
	"github.com/grokify/mogo/fmt/fmtutil"
	iom "github.com/grokify/mogo/io/ioutilmore"
)

/*

https://godoc.org/google.golang.org/genproto/googleapis/cloud/speech/v1#LongRunningRecognizeResponse

*/

func main() {
	file := "../mongodb-is-web-scale/web-scale_b2F-DItXtZs.mp3_tr_google_standard.json"

	res, err := google.ReadLongRunningRecognizeResponseFile(file)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("NUM_RESULTS [%v]\n", len(res.Results))

	checkSpeakerTags := false
	if checkSpeakerTags {
		for i, resp := range res.Results {
			for _, alt := range resp.Alternatives {
				for _, wrd := range alt.Words {
					if wrd.SpeakerTag > 0 {
						fmt.Printf("%v ", i)
					}
				}
			}
		}
	}

	txn, err := google.LongRunningRecognizeResponseToTranscript(res)
	if err != nil {
		log.Fatal(err)
	}
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
