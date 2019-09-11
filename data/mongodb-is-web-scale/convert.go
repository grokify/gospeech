package main

import (
	"fmt"
	"log"

	"github.com/grokify/go-diarization"
	"github.com/grokify/gotilla/fmt/fmtutil"
)

func main() {
	file := "episode_1_mongo_db_is_web_scale_b2F-DItXtZs_transcribeme.txt"
	txn, err := diarization.ParseTranscribeMeFile(file)
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.PrintJSON(txn)
	rttm := diarization.TranscriptToRTTM(txn)
	rttm.WriteFile(file+".rttm", 0644)
	fmt.Println("DONE")
}
