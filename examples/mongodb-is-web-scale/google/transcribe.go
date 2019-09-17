package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	option "google.golang.org/api/option"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

/*

Sample: https://github.com/GoogleCloudPlatform/golang-samples/blob/master/speech/captionasync/captionasync.go

https://cloud.google.com/speech-to-text/docs/async-recognize#speech-async-recognize-gcs-go

https://cloud.google.com/speech-to-text/docs/async-recognize


$ go build captionasync.go
# command-line-arguments
./captionasync.go:79:21: undefined: "google.golang.org/genproto/googleapis/cloud/speech/v1".RecognitionConfig_MP3

*/

var opts struct {
	CredentialsFile string `short:"c" long:"credentials" description:"Path to credentials file" required:"true"`
}

func main() {
	err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	file := "episode_1_mongo_db_is_web_scale_b2F-DItXtZs.mp3"

	cfg := option.WithCredentialsFile(opts.CredentialsFile)

	send(os.Stdout, file)

	fmt.Println("DONE")
}

func send(w io.Writer, client *speech.Client, filename string) error {
	ctx := context.Background()
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// Send the contents of the audio file with the encoding and
	// and sample rate information to be transcripted.
	req := &speechpb.LongRunningRecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        "MP3",
			SampleRateHertz: 44100,
			LanguageCode:    "en-US",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	}

	op, err := client.LongRunningRecognize(ctx, req)
	if err != nil {
		return err
	}
	resp, err := op.Wait(ctx)
	if err != nil {
		return err
	}

	// Print the results.
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			fmt.Fprintf(w, "\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
		}
	}
	return nil
}
