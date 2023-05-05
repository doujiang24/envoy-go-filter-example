package main

import (
	"encoding/base64"
	"strings"

	"github.com/envoyproxy/envoy/contrib/golang/filters/http/source/go/pkg/api"
)

type filter struct {
	callbacks api.FilterCallbackHandler
	config    *config
}

const secretKey = "secret"

// parseBasicAuth parses an HTTP Basic Authentication string.
// "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==" returns ("Aladdin", "open sesame", true).
func parseBasicAuth(auth string) (username, password string, ok bool) {
	const prefix = "Basic "
	// Case insensitive prefix match. See Issue 22736.
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return "", "", false
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return "", "", false
	}
	cs := string(c)
	username, password, ok = strings.Cut(cs, ":")
	if !ok {
		return "", "", false
	}
	return username, password, true
}

func (f *filter) verify(header api.RequestHeaderMap) (bool, string) {
	auth, ok := header.Get("authorization")
	if !ok {
		return false, "no Authorization"
	}
	username, password, ok := parseBasicAuth(auth)
	if !ok {
		return false, "invalid Authorization format"
	}
	// fmt.Printf("expected username: %v, password: %v; got username: %v, password: %v\n", f.config.username, f.config.password, username, password)
	if f.config.username == username && f.config.password == password {
		return true, ""
	}
	return false, "invalid username or password"
}

func (f *filter) DecodeHeaders(header api.RequestHeaderMap, endStream bool) api.StatusType {
	go func() {
		defer f.callbacks.RecoverPanic()
		f.callbacks.Continue(api.Continue)
	}()
	return api.Running
}

func (f *filter) DecodeData(buffer api.BufferInstance, endStream bool) api.StatusType {
	go func() {
		defer f.callbacks.RecoverPanic()
		f.callbacks.Continue(api.Continue)
	}()
	return api.Running
}

func (f *filter) DecodeTrailers(trailers api.RequestTrailerMap) api.StatusType {
	go func() {
		defer f.callbacks.RecoverPanic()
		f.callbacks.Continue(api.Continue)
	}()
	return api.Running
}

func (f *filter) EncodeHeaders(header api.ResponseHeaderMap, endStream bool) api.StatusType {
	go func() {
		defer f.callbacks.RecoverPanic()
		f.callbacks.Continue(api.Continue)
	}()
	return api.Running
}

func (f *filter) EncodeData(buffer api.BufferInstance, endStream bool) api.StatusType {
	go func() {
		defer f.callbacks.RecoverPanic()
		f.callbacks.Continue(api.Continue)
	}()
	return api.Running
}

func (f *filter) EncodeTrailers(trailers api.ResponseTrailerMap) api.StatusType {
	return api.Continue
}

func (f *filter) OnDestroy(reason api.DestroyReason) {
}

func main() {
}
