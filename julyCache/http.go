package julyCache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_julycache/"

// HttpPool implements PeerPicker for a pool of HTTP peers.
type HttpPool struct {
	self     string
	basePath string
}

// NewHttpPool initializes an HTTP pool of peers.
func NewHttpPool(self string) *HttpPool {
	return &HttpPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// Log info with server name
func (p *HttpPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// ServeHTTP handle all http requests
func (p *HttpPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HttpPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)
	// /<basepath>/<groupname>/<key> required
	params := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(params) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupName := params[0]
	key := params[1]
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}
	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-type", "application/octet-stream")
	_, _ = w.Write(view.ByteSlice())
}
