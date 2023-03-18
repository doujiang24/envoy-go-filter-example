package main

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"net/http"

	"github.com/envoyproxy/envoy/contrib/golang/filters/http/source/go/pkg/api"
)

type filter struct {
	callbacks api.FilterCallbackHandler
	config    *config
}

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

// checkRemote checks the username and password against a remote server, through http request.
func checkRemote(config *config, username, password string) bool {
	body := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password)
	remoteAddr := "http://" + config.host + ":" + strconv.Itoa(int(config.port)) + "/check"
	resp, err := http.Post(remoteAddr, "application/json", strings.NewReader(body))
	if err != nil {
		fmt.Printf("check error: %v\n", err)
		return false
	}
	if resp.StatusCode != 200 {
		return false
	}
	return true
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
	fmt.Printf("got username: %v, password: %v\n", username, password)

	if ok := checkRemote(f.config, username, password); !ok {
		return false, "invalid username or password"
	}
	return true, ""
}

func (f *filter) DecodeHeaders(header api.RequestHeaderMap, endStream bool) api.StatusType {
	go func() {
		if ok, msg := f.verify(header); !ok {
			// TODO: set the WWW-Authenticate response header
			f.callbacks.SendLocalReply(401, msg, map[string]string{}, 0, "bad-request")
			return
		}
		f.callbacks.Continue(api.Continue)
	}()
	return api.Running
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
