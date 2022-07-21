package hear

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHear(t *testing.T) {
	b, err := Hear2Wav(time.Second*10, time.Second*2)
	assert.Nil(t, err)

	err = ioutil.WriteFile("./result.wav", b, os.ModePerm)
	assert.Nil(t, err)
}
