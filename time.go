package gospeech

import (
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
)

// DurationFromProtobuf helper function to convert a `durationpb.Duration` to a `time.Duration`.
func DurationFromProtobuf(pdur *durationpb.Duration) time.Duration {
	return time.Duration(pdur.Seconds*int64(time.Second) + int64(pdur.Nanos))
}
