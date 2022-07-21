package pcm2wav

import (
	_ "embed"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed test.pcm
var b []byte

func TestPcm2Wav(t *testing.T) {
	resultWav, err := Pcm2Wav(b, 1, 16000, 16)
	if err != nil {
		assert.Error(t, err)
	}

	err = ioutil.WriteFile("./result.wav", resultWav, 0666)
	assert.Nil(t, err)
}
