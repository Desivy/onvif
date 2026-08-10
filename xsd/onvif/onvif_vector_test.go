package onvif

import (
	"encoding/xml"
	"strings"
	"testing"
)

func TestVector2D_ZeroAxesAreSerialized(t *testing.T) {
	out, err := xml.Marshal(struct {
		XMLName xml.Name `xml:"Velocity"`
		PanTilt Vector2D `xml:"PanTilt"`
	}{PanTilt: Vector2D{X: 0, Y: -0.5}})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got := string(out)
	if !strings.Contains(got, `x="0"`) {
		t.Errorf("zero X attribute was dropped; ONVIF declares x as use=\"required\"\ngot: %s", got)
	}
	if !strings.Contains(got, `y="-0.5"`) {
		t.Errorf("Y attribute missing\ngot: %s", got)
	}
}

func TestVector1D_ZeroAxisIsSerialized(t *testing.T) {
	out, err := xml.Marshal(struct {
		XMLName xml.Name `xml:"Velocity"`
		Zoom    Vector1D `xml:"Zoom"`
	}{Zoom: Vector1D{X: 0}})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(string(out), `x="0"`) {
		t.Errorf("zero X attribute was dropped\ngot: %s", string(out))
	}
}
