package gosoap

import "testing"

func TestAddStringHeaderContentRejectsRootlessFragments(t *testing.T) {
	for _, frag := range []string{"", "   ", "<!-- c -->", "<?pi x?>"} {
		msg := NewEmptySOAP()
		if err := msg.AddStringHeaderContent(frag); err == nil {
			t.Fatalf("fragment %q: want error, got nil", frag)
		}
	}
	msg := NewEmptySOAP()
	if err := msg.AddStringHeaderContent(`<a:Id xmlns:a="urn:x">7</a:Id>`); err != nil {
		t.Fatalf("well-formed element: %v", err)
	}
}
