package livechat

import "strconv"

const (
	HeaderVersion    = "X-LIVEAGENT-API-VERSION"
	HeaderAffinity   = "X-LIVEAGENT-AFFINITY"
	HeaderSessionKey = "X-LIVEAGENT-SESSION-KEY"
	HeaderSequence   = "X-LIVEAGENT-SEQUENCE"
)

type Header struct {
	Version    string
	Affinity   string
	SessionKey string
	Sequence   int
}

func (h *Header) Values() map[string]string {
	m := make(map[string]string, 4)
	m[HeaderVersion] = h.Version

	if h.Affinity == "" {
		h.Affinity = "null"
	}
	m[HeaderAffinity] = h.Affinity

	if h.SessionKey != "" {
		m[HeaderSessionKey] = h.SessionKey
	}

	if h.Sequence != 0 {
		m[HeaderSequence] = strconv.Itoa(h.Sequence)
	}

	return m
}
