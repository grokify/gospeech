// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Command captionasync sends audio data to the Google Speech API
// and prints its transcript.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"

	//"github.com/abc/def"
	speech "cloud.google.com/go/speech/apiv1p1beta1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1p1beta1"
)

const usage = `Usage: captionasync <audiofile>

Audio file must be a 16-bit signed little-endian encoded
with a sample rate of 16000.

The path to the audio file may be a GCS URI (gs://...).

Set credentials with GOOGLE_APPLICATION_CREDENTIALS

https://cloud.google.com/docs/authentication/production
https://cloud.google.com/speech-to-text/docs/multiple-voices
https://godoc.org/google.golang.org/genproto/googleapis/cloud/speech/v1p1beta1#RecognitionConfig
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}

	var sendFunc func(io.Writer, *speech.Client, string) error

	path := os.Args[1]
	if strings.Contains(path, "://") {
		sendFunc = sendGCS
	} else {
		sendFunc = send
	}

	ctx := context.Background()
	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if err := sendFunc(os.Stdout, client, os.Args[1]); err != nil {
		log.Fatal(err)
	}
}

// [START speech_transcribe_async]

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
			Encoding:        speechpb.RecognitionConfig_MP3,
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

// [END speech_transcribe_async]

// [START speech_transcribe_async_gcs]

func sendGCS(w io.Writer, client *speech.Client, gcsURI string) error {
	ctx := context.Background()

	// Send the contents of the audio file with the encoding and
	// and sample rate information to be transcripted.
	req := &speechpb.LongRunningRecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:                   speechpb.RecognitionConfig_MP3,
			SampleRateHertz:            44100,
			LanguageCode:               "en-US",
			EnableSpeakerDiarization:   true,
			DiarizationSpeakerCount:    2,
			EnableAutomaticPunctuation: true,
			DiarizationConfig: &speechpb.SpeakerDiarizationConfig{
				EnableSpeakerDiarization: true,
				MinSpeakerCount:          2,
				MaxSpeakerCount:          2,
			},
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Uri{Uri: gcsURI},
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

	writeOutput := true
	if writeOutput {
		pathParts := strings.Split(gcsURI, "/")
		outFile := pathParts[len(pathParts)-1] + "_transcript_gcs-standard.json"
		bytes, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(outFile, bytes, 0644)
		if err != nil {
			return err
		}
	}
	// Print the results.
	// Response is: LongRunningRecognizeResponse
	// https://godoc.org/google.golang.org/genproto/googleapis/cloud/speech/v1p1beta1#LongRunningRecognizeResponse
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			fmt.Fprintf(w, "\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
		}
	}

	return nil
}

// TypeName returns the name of a struct.
// stackoverflow-answerId:1908967
func TypeName(myvar interface{}) (res string) {
	t := reflect.TypeOf(myvar)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
		res += "*"
	}
	return res + t.Name()
}

// [END speech_transcribe_async_gcs]
