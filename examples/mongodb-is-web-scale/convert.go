package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/grokify/gospeech/diarization"
	"github.com/grokify/gospeech/diarization/nvivo"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/os/osutil"
	flags "github.com/jessevdk/go-flags"
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
		txn, err := nvivo.ParseNVivoPcFile(opts.Input)
		if err != nil {
			log.Fatal(err)
		}
		//fmtutil.PrintJSON(txn)
		rttm := diarization.TranscriptToRTTM(txn)
		rttm.WriteFile(opts.Output+".rttm", 0644)
		fmt.Printf("WROTE [%v]\n", opts.Output+".rttm")
		err = osutil.WriteFileJSON(opts.Output+".json", txn, 0644, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("WROTE [%v]\n", opts.Output+".json")

		html := diarization.TranscriptWebpage(txn)
		err = ioutil.WriteFile(opts.Input+".html", []byte(html), 0644)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("WROTE [%v]\n", opts.Output+".html")
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
