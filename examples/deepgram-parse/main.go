package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/grokify/go-transcribe/diarization"
	"github.com/grokify/go-transcribe/diarization/deepgram"
	"github.com/grokify/mogo/fmt/fmtutil"
	flags "github.com/jessevdk/go-flags"
)

type Options struct {
	Input string `short:"i" long:"input" description:"Input file" required:"true"`
}

func main() {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	dtxn, err := deepgram.NewTranscriptFile(opts.Input)
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
	fmt.Println("DONE")
}
