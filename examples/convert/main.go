package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/grokify/gospeech/diarization"
	"github.com/grokify/gospeech/diarization/deepgram"
	"github.com/grokify/gospeech/diarization/google"
	"github.com/grokify/gospeech/diarization/google/phone"
	"github.com/grokify/gospeech/diarization/nvivo"
	"github.com/grokify/mogo/fmt/fmtutil"
	flags "github.com/jessevdk/go-flags"
)

type Options struct {
	InputFile string `short:"i" long:"input" description:"Input file" required:"true"`
	InputType string `short:"t" long:"type" description:"Input file type" required:"true"`
}

func main() {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	switch opts.InputType {
	case "deepgram":
		processDeepgram(opts)
	case "google":
		processGoogle(opts)
	case "nvivopc":
		processNvivopc(opts)
	default:
		log.Fatal(fmt.Sprintf("Error processing Input Type [%v]", opts.InputType))
	}
	fmt.Println("DONE")
}

func processGoogle(opts Options) {
	if 1 == 0 {
		txn, err := google.NewTranscriptFileLongRunningResponse(opts.InputFile)
		if err != nil {
			log.Fatal(err)
		}
		html := diarization.TranscriptWebpage(txn)
		err = ioutil.WriteFile("_transcript.html", []byte(html), 0644)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("WROTE [%v]\n", "_transcript.html")
	}
	txn, err := phone.NewTranscriptFile(opts.InputFile)
	if err != nil {
		log.Fatal(err)
	}
	html := diarization.TranscriptWebpage(txn)
	err = ioutil.WriteFile("_transcript.html", []byte(html), 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("WROTE [%v]\n", "_transcript.html")
}

func processNvivopc(opts Options) {
	txn, err := nvivo.ParseNVivoPcFile(opts.InputFile)
	if err != nil {
		log.Fatal(err)
	}
	html := diarization.TranscriptWebpage(txn)
	err = ioutil.WriteFile("_transcript.html", []byte(html), 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("WROTE [%v]\n", "_transcript.html")
}

func processDeepgram(opts Options) {
	dtxn, err := deepgram.NewTranscriptFile(opts.InputFile)
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.PrintJSON(dtxn)

	fmt.Printf("NUM_CHANNEL_RESULTS: [%v]\n", len(dtxn.Results.Channels))

	err = dtxn.Validate()
	if err != nil {
		log.Fatal(err)
	}

	ctxn, err := deepgram.CanonicalTranscriptFromDeepgram(dtxn)
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.PrintJSON(ctxn)

	html := diarization.TranscriptWebpage(ctxn)
	err = ioutil.WriteFile("_transcript.html", []byte(html), 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("WROTE [%v]\n", "_transcript.html")

	/*
		for i, chanRes := range dtxn.Results.Channels {
			fmt.Printf("NUM_CHANNEL_ALTERNATIVES: CHAN[%v]ALT[%v]\n", i, len(chanRes.Alternatives))
			for j, alt := range chanRes.Alternatives {
				txnWords := len(strings.Fields(alt.Transcript))
				spkWords := len(alt.Words)
				fmt.Printf("CHAN [%v], ALT [%v] TXN_WORDS [%v] SPK_WORDS [%v]\n", i, j, txnWords, spkWords)
			}
		}
	*/

}
