package gospeech

import (
	"time"

	"github.com/grokify/mogo/time/timeutil"
	"google.golang.org/protobuf/types/known/durationpb"
)

// DurationFromProtobuf helper function to convert a `durationpb.Duration` to a `time.Duration`.
func DurationFromProtobuf(pdur *durationpb.Duration) time.Duration {
	return time.Duration(pdur.Seconds*timeutil.NanosPerSecond + int64(pdur.Nanos))
}
