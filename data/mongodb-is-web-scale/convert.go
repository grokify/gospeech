package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/grokify/go-diarization"
	"github.com/grokify/gotilla/fmt/fmtutil"
	iom "github.com/grokify/gotilla/io/ioutilmore"
	"github.com/jessevdk/go-flags"
)

var opts struct {
	Input string `short:"i" long:"input" description:"transcription file" required:"true"`
	Type  string `short:"t" long:"type" description:"file type" required:"true"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	switch strings.ToLower(strings.TrimSpace(opts.Type)) {
	case "nvivopc":
		txn, err := diarization.ParseNVivoPcFile(opts.Input)
		if err != nil {
			log.Fatal(err)
		}
		//fmtutil.PrintJSON(txn)
		rttm := diarization.TranscriptToRTTM(txn)
		rttm.WriteFile(opts.Input+".rttm", 0644)
		fmt.Printf("WROTE [%v]\n", opts.Input+".rttm")
		err = iom.WriteFileJSON(opts.Input+".json", txn, 0644, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("WROTE [%v]\n", opts.Input+".json")
	case "transcribeme":
		txn, err := diarization.ParseTranscribeMeFile(opts.Input)
		if err != nil {
			log.Fatal(err)
		}
		fmtutil.PrintJSON(txn)
		rttm := diarization.TranscriptToRTTM(txn)
		rttm.WriteFile(opts.Input+".rttm", 0644)
	}

	fmt.Println("DONE")
}
