package test

import "testing"
import "github.com/lmendes86/Kalbi/sdp"

func TestSDPParser(t *testing.T) {
	byteMsg := []byte(msg)
	x := sdp.Parse(byteMsg)
	t.Log(string(x.Time.String()))

}
