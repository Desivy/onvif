package onvif

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// A SOAP call to an unresponsive device must abort as soon as the caller's
// context is cancelled — not block until the http.Client timeout fires.
func TestDevice_SendSoapWithContext_CancelAbortsInFlightRequest(t *testing.T) {
	for _, mode := range []string{DigestAuth, UsernameTokenAuth} {
		t.Run(mode, func(t *testing.T) {
			release := make(chan struct{})

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				<-release // block until test signals us to stop
			}))
			// Defers execute LIFO: release is closed first (unblocking the handler),
			// then srv.Close() can drain cleanly.
			defer srv.Close()
			defer close(release)

			dev := NewDeviceRaw(DeviceParams{
				Xaddr:    srv.Listener.Addr().String(),
				AuthMode: mode,
				Username: "user",
				Password: "pass",
			})

			ctx, cancel := context.WithCancel(context.Background())
			go func() {
				time.Sleep(100 * time.Millisecond)
				cancel()
			}()

			start := time.Now()
			resp, err := dev.SendSoapWithContext(ctx, srv.URL, "<test/>")
			elapsed := time.Since(start)
			if resp != nil {
				resp.Body.Close()
			}

			if err == nil {
				t.Fatal("SendSoapWithContext should fail when ctx is cancelled mid-request")
			}
			if !errors.Is(err, context.Canceled) {
				t.Fatalf("err = %v, want context.Canceled", err)
			}
			if elapsed > 2*time.Second {
				t.Fatalf("SendSoapWithContext took %v after cancel, want prompt return", elapsed)
			}
		})
	}
}
