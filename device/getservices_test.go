package device

import (
	"encoding/xml"
	"strings"
	"testing"
)

func TestGetServices_MarshalIncludesElements(t *testing.T) {
	b, err := xml.Marshal(GetServices{IncludeCapability: true})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got := string(b)
	if !strings.Contains(got, "GetServices") {
		t.Errorf("want GetServices element, got: %s", got)
	}
	if !strings.Contains(got, "IncludeCapability") {
		t.Errorf("want IncludeCapability element, got: %s", got)
	}
	if !strings.Contains(got, "true") {
		t.Errorf("want 'true' value, got: %s", got)
	}
}

func TestGetServicesResponse_UnmarshalServicesSlice(t *testing.T) {
	input := `<GetServicesResponse>
        <Service>
            <Namespace>http://www.onvif.org/ver20/media/wsdl</Namespace>
            <XAddr>http://192.0.2.1/onvif/media2</XAddr>
            <Capabilities><Detail/></Capabilities>
        </Service>
        <Service>
            <Namespace>http://www.onvif.org/ver10/events/wsdl</Namespace>
            <XAddr>http://192.0.2.1/onvif/events</XAddr>
        </Service>
    </GetServicesResponse>`

	var resp GetServicesResponse
	if err := xml.Unmarshal([]byte(input), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp.Service) != 2 {
		t.Fatalf("len(Service) = %d, want 2", len(resp.Service))
	}
	if resp.Service[0].Namespace != "http://www.onvif.org/ver20/media/wsdl" {
		t.Errorf("Service[0].Namespace = %q", resp.Service[0].Namespace)
	}
	if resp.Service[0].XAddr != "http://192.0.2.1/onvif/media2" {
		t.Errorf("Service[0].XAddr = %q", resp.Service[0].XAddr)
	}
	if len(resp.Service[0].Capabilities.Inner) == 0 {
		t.Error("Service[0].Capabilities.Inner is empty, want non-empty")
	}
	if len(resp.Service[1].Capabilities.Inner) != 0 {
		t.Error("Service[1] has no Capabilities element, Inner should be empty")
	}
}
