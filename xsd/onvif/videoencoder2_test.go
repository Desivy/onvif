package onvif

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/IOTechSystems/onvif/xsd"
)

// Media2 returns GovLength and Profile as ATTRIBUTES of the configuration
// element. encoding/xml never matches an attribute against an element tag,
// so a caller reading these as elements would always see them as unset.
func TestVideoEncoder2Configuration_UnmarshalsAttributes(t *testing.T) {
	const body = `<VideoEncoder token="VEC_1" GovLength="40" Profile="Main">
		<Name>enc1</Name>
		<Encoding>H264</Encoding>
		<Resolution><Width>640</Width><Height>480</Height></Resolution>
		<RateControl ConstantBitRate="false">
			<FrameRateLimit>20.0</FrameRateLimit>
			<BitrateLimit>512</BitrateLimit>
		</RateControl>
	</VideoEncoder>`

	var got VideoEncoder2Configuration
	if err := xml.Unmarshal([]byte(body), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.GovLength == nil {
		t.Fatal("GovLength is nil — attribute not parsed")
	}
	if int(*got.GovLength) != 40 {
		t.Errorf("GovLength = %d, want 40", int(*got.GovLength))
	}
	if got.Profile == nil || string(*got.Profile) != "Main" {
		t.Errorf("Profile = %v, want Main", got.Profile)
	}
	if got.RateControl == nil || got.RateControl.ConstantBitRate == nil {
		t.Fatal("RateControl.ConstantBitRate is nil — attribute not parsed")
	}
}

// The request side must emit them as attributes too; a conformant device
// ignores or faults on element forms.
func TestVideoEncoder2ConfigurationRequest_MarshalsAttributes(t *testing.T) {
	gov := xsd.Int(40)
	prof := xsd.String("Main")
	cbr := xsd.Boolean(false)
	req := VideoEncoder2ConfigurationRequest{
		GovLength:   &gov,
		Profile:     &prof,
		RateControl: &VideoRateControl2Request{ConstantBitRate: &cbr},
	}
	out, err := xml.Marshal(&req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(out)
	if !strings.Contains(s, `GovLength="40"`) {
		t.Errorf("GovLength not emitted as attribute: %s", s)
	}
	if !strings.Contains(s, `Profile="Main"`) {
		t.Errorf("Profile not emitted as attribute: %s", s)
	}
	if !strings.Contains(s, `ConstantBitRate="false"`) {
		t.Errorf("ConstantBitRate not emitted as attribute: %s", s)
	}
}

// Devices advertise the framerates they support as a space-separated
// ATTRIBUTE list, not as the FrameRateRange the fork models today. A caller
// snapping to the nearest supported value needs the discrete list; a
// min/max range cannot express it.
func TestVideoEncoder2ConfigurationOptions_FrameRatesSupportedAttribute(t *testing.T) {
	const body = `<Options FrameRatesSupported="30 25 20 15 12 10">
		<Encoding>H264</Encoding>
		<ResolutionsAvailable><Width>640</Width><Height>480</Height></ResolutionsAvailable>
	</Options>`

	var got VideoEncoder2ConfigurationOptions
	if err := xml.Unmarshal([]byte(body), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if string(got.FrameRatesSupported) != "30 25 20 15 12 10" {
		t.Errorf("FrameRatesSupported = %q — attribute not parsed", string(got.FrameRatesSupported))
	}
}

// GovLengthRange and ProfilesSupported are also space-separated ATTRIBUTE
// lists on a Profile-T device ("1 300", "Main Main10"), not the Min/Max
// element pair (GovLengthRange) or repeated element ([]xsd.String,
// ProfilesSupported) the fork modeled before; ConstantBitRateSupported is a
// plain boolean attribute the fork modeled as an untagged element.
func TestVideoEncoder2ConfigurationOptions_RemainingAttributes(t *testing.T) {
	const body = `<Options GovLengthRange="1 300" FrameRatesSupported="30 25 20 15 12 10" ProfilesSupported="Main Main10" ConstantBitRateSupported="true">
		<Encoding>H264</Encoding>
		<ResolutionsAvailable><Width>640</Width><Height>480</Height></ResolutionsAvailable>
	</Options>`

	var got VideoEncoder2ConfigurationOptions
	if err := xml.Unmarshal([]byte(body), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if string(got.GovLengthRange) != "1 300" {
		t.Errorf("GovLengthRange = %q — attribute not parsed", string(got.GovLengthRange))
	}
	if string(got.ProfilesSupported) != "Main Main10" {
		t.Errorf("ProfilesSupported = %q — attribute not parsed", string(got.ProfilesSupported))
	}
	if got.ConstantBitRateSupported == nil {
		t.Fatal("ConstantBitRateSupported is nil — attribute not parsed")
	}
	if bool(*got.ConstantBitRateSupported) != true {
		t.Errorf("ConstantBitRateSupported = %v, want true", *got.ConstantBitRateSupported)
	}
}
