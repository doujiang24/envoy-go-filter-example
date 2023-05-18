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
	if f.config.username == username && f.config.password == password {
		return true, ""
	}
	return false, "invalid username or password"
}

func (f *filter) DecodeHeaders(header api.RequestHeaderMap, endStream bool) api.StatusType {
	if ok, msg := f.verify(header); !ok {
		// TODO: set the WWW-Authenticate response header
		f.callbacks.SendLocalReply(401, msg, map[string]string{}, 0, "bad-request")
		return api.LocalReply
	}
	return api.Continue
}

func (f *filter) DecodeData(buffer api.BufferInstance, endStream bool) api.StatusType {
	return api.Continue
}

func (f *filter) DecodeTrailers(trailers api.RequestTrailerMap) api.StatusType {
	return api.Continue
}

func (f *filter) EncodeHeaders(header api.ResponseHeaderMap, endStream bool) api.StatusType {
	return api.Continue
}

func (f *filter) EncodeData(buffer api.BufferInstance, endStream bool) api.StatusType {
	return api.Continue
}

func (f *filter) EncodeTrailers(trailers api.ResponseTrailerMap) api.StatusType {
	return api.Continue
}

func (f *filter) OnDestroy(reason api.DestroyReason) {
}

func main() {
}
