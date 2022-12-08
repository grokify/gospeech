# Go Speech

[![Build Status][build-status-svg]][build-status-link]
[![Go Report Card][goreport-svg]][goreport-link]
[![Docs][docs-godoc-svg]][docs-godoc-link]
[![License][license-svg]][license-link]

Tools to test diarization for speech-to-text voice recognition systems.

Initially, it is designed to convert transcripts from TranscribeMe.com to [Rich Transcription Time Marked (RTTM) files](https://github.com/nryant/dscore#rttm). You can read more from the following diarization evaluation tool:

* https://github.com/nryant/dscore

Install the following Python pre-requisites before running dscore:

```
$ pip install tabulate intervaltree numpy scipy
```

## Usage

See the following example:

[`mongodb-is-web-scale`](data/mongodb-is-web-scale)

 [build-status-svg]: https://github.com/grokify/goauth/workflows/test/badge.svg
 [build-status-link]: https://github.com/grokify/gospeech/actions/workflows/test.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/gospeech
 [goreport-link]: https://goreportcard.com/report/github.com/grokify/gospeech
 [docs-godoc-svg]: https://img.shields.io/badge/docs-godoc-blue.svg
 [docs-godoc-link]: https://godoc.org/github.com/grokify/gospeech
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-link]: https://github.com/grokify/gospeech/blob/master/LICENSE
