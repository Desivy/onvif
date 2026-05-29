package media2

import onvif "github.com/IOTechSystems/onvif/xsd/onvif"

// GetStreamUri requests the RTSP streaming URI for a media profile.
// Protocol should be "RTSP" for standard unicast streaming.
type GetStreamUri struct {
	XMLName      string               `xml:"tr2:GetStreamUri"`
	Protocol     string               `xml:"tr2:Protocol,omitempty"`
	ProfileToken onvif.ReferenceToken `xml:"tr2:ProfileToken"`
}

// GetStreamUriResponse holds the RTSP URI returned by the device.
// Credentials are not included by the device; the caller injects them separately.
type GetStreamUriResponse struct {
	Uri string `xml:"Uri"`
}

// GetStreamUriFunction implements the Function interface for Media2FunctionMap.
type GetStreamUriFunction struct{}

func (_ *GetStreamUriFunction) Request() interface{}  { return &GetStreamUri{} }
func (_ *GetStreamUriFunction) Response() interface{} { return &GetStreamUriResponse{} }
