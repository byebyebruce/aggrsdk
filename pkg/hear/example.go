package hear

import (
	"time"

	"github.com/byebyebruce/aggrsdk/pkg/pcm2wav"
	"github.com/gordonklaus/portaudio"
)

// ref https://github.dev/evanphx/hear

// Hear2Pcm 听
// 生成pcm buffer
// numInputChannels 1
// umOutputChannels 0
// sampleRate 16000
func Hear2Pcm(maxDuration, quietDurationToStop time.Duration) ([]byte, error) {
	err := portaudio.Initialize()
	if err != nil {
		return nil, err
	}
	defer portaudio.Terminate()

	opts := ListenOpts{
		QuietDuration:    quietDurationToStop,
		AlreadyListening: true,
		MaxTime:          maxDuration,
	}

	bf, err := ListenIntoBuffer(opts)
	if err != nil {
		return nil, err
	}
	return bf.Bytes(), nil

}

// Hear2Wav 听生成wav buffer
func Hear2Wav(maxDuration, quietDurationToStop time.Duration) ([]byte, error) {
	b, err := Hear2Pcm(maxDuration, quietDurationToStop)
	if err != nil {
		return nil, err
	}
	return pcm2wav.Pcm2Wav(b, 1, 16000, 16)
}
