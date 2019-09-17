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
	Type   string `short:"t" long:"type" description:"file type" required:"true"`
	Input  string `short:"i" long:"input" description:"transcription file" required:"true"`
	Output string `short:"o" long:"output" description:"transcription file" required:"false"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}
	opts.Output = strings.TrimSpace(opts.Output)
	if len(opts.Output) == 0 {
		opts.Output = opts.Input
	}

	switch strings.ToLower(strings.TrimSpace(opts.Type)) {
	case "nvivopc":
		txn, err := diarization.ParseNVivoPcFile(opts.Input)
		if err != nil {
			log.Fatal(err)
		}
		//fmtutil.PrintJSON(txn)
		rttm := diarization.TranscriptToRTTM(txn)
		rttm.WriteFile(opts.Output+".rttm", 0644)
		fmt.Printf("WROTE [%v]\n", opts.Output+".rttm")
		err = iom.WriteFileJSON(opts.Output+".json", txn, 0644, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("WROTE [%v]\n", opts.Output+".json")
	case "transcribeme":
		panic("use NVivo PC instead of TranscribeMe native")
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
