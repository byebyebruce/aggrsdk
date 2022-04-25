package aliyun

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

var (
	cfg      Aliyun
	endpoint = os.Getenv("oss_ep")
	bucket   = os.Getenv("oss_bucket")
	app      = os.Getenv("aliyun_app")
)

func init() {
	cfg.AK = os.Getenv("aliyun_ak")
	cfg.AS = os.Getenv("aliyun_as")
}

func TestAASR(t *testing.T) {
	aasr := NewAASR(cfg, app)
	ret, err := aasr.Do("https://mp4audio.oss-cn-beijing.aliyuncs.com/a.mp4.wav")
	fmt.Println(ret, err)
}

func TestTranslate(t *testing.T) {
	trans := NewTranslator(cfg)
	ret, err := trans.Do("sue is a graphic designer, oh, that's interesting. you know, our restaurant needs a new menu design. oh, really, you do that sort of thing, yes, i've designed logos and menus for several restaurants in town. let's get together and talk over some ideas, great, so here's my card. ", "en", "zh")
	fmt.Println(ret, err)
}

func TestOSS(t *testing.T) {
	o := NewOSS(cfg, endpoint, bucket)
	b := &bytes.Buffer{}
	b.WriteString("haha")
	ret, err := o.Upload("test/test.txt", b)
	fmt.Println(ret, err)
}
