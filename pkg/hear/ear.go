package hear

import (
	"bytes"
	"encoding/binary"
	"time"

	"github.com/gordonklaus/portaudio"
)

var audioRunning bool

func InitAudio() error {
	if !audioRunning {
		err := portaudio.Initialize()
		if err != nil {
			return err
		}
		audioRunning = true
	}
	return nil
}

func FreeAudio() {
	if audioRunning {
		portaudio.Terminate()
	}
}

const DefaultQuietTime = time.Second

type State int

const (
	Waiting State = iota
	Listening
	Asking
)

type ListenOpts struct {
	State            func(State)
	QuietDuration    time.Duration
	AlreadyListening bool
	MaxTime          time.Duration // default 默认最长7秒
}

func ListenIntoBuffer(opts ListenOpts) (*bytes.Buffer, error) {
	if opts.MaxTime <= 0 {
		opts.MaxTime = time.Second * 7
	}
	in := make([]int16, 8196)
	stream, err := portaudio.OpenDefaultStream(1, 0, 16000, len(in), in)
	if err != nil {
		return nil, err
	}

	defer stream.Close()

	err = stream.Start()
	if err != nil {
		return nil, err
	}

	var (
		buf            bytes.Buffer
		heardSomething = opts.AlreadyListening
		quiet          bool
		quietTime      = opts.QuietDuration
		quietStart     time.Time
		lastFlux       float64
	)

	vad := NewVAD(len(in))

	if quietTime == 0 {
		quietTime = DefaultQuietTime
	}

	if opts.State != nil {
		if heardSomething {
			opts.State(Listening)
		} else {
			opts.State(Waiting)
		}
	}

	start := time.Now()
reader:
	for {

		if time.Now().Sub(start) > opts.MaxTime {
			break reader
		}
		err = stream.Read()
		if err != nil {
			return nil, err
		}

		err = binary.Write(&buf, binary.LittleEndian, in)
		if err != nil {
			return nil, err
		}

		flux := vad.Flux(in)

		if lastFlux == 0 {
			lastFlux = flux
			continue
		}

		if heardSomething {
			if flux*1.75 <= lastFlux {
				if !quiet {
					quietStart = time.Now()
				} else {
					diff := time.Since(quietStart)

					if diff > quietTime {
						break reader
					}
				}

				quiet = true
			} else {
				quiet = false
				lastFlux = flux
			}
		} else {
			if flux >= lastFlux*1.75 {
				heardSomething = true
				if opts.State != nil {
					opts.State(Listening)
				}
			}

			lastFlux = flux
		}
	}

	err = stream.Stop()
	if err != nil {
		return nil, err
	}

	return &buf, nil
}
