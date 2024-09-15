package subscription

import (
	"regexp"
	"strings"
	"sync"
)

func newTopicMatcher() *topicMatcher {
	return &topicMatcher{
		topicRegexes: make(map[string]*regexp.Regexp),
	}
}

type topicMatcher struct {
	topicRegexes   map[string]*regexp.Regexp
	topicRegexesMu sync.RWMutex
}

func (tm *topicMatcher) match(topic, subscription string) bool {
	if topic == subscription {
		return true
	}

	if !hasWildcard(subscription) {
		return false
	}

	tm.topicRegexesMu.RLock()
	regex, ok := tm.topicRegexes[subscription]
	tm.topicRegexesMu.RUnlock()

	if !ok {
		cleaned := regexp.QuoteMeta(subscription)
		cleaned = strings.ReplaceAll(cleaned, "\\+", "[^/]+") // The plus is escaped in the quote meta, so we need to unescape it
		cleaned = strings.ReplaceAll(cleaned, "#", ".+")

		regex = regexp.MustCompile("^" + cleaned + "$")

		tm.topicRegexesMu.Lock()
		tm.topicRegexes[subscription] = regex
		tm.topicRegexesMu.Unlock()
	}

	if regex == nil {
		return false
	}

	return regex.MatchString(topic)
}

func (tm *topicMatcher) reset() {
	tm.topicRegexesMu.Lock()
	defer tm.topicRegexesMu.Unlock()

	tm.topicRegexes = make(map[string]*regexp.Regexp)
}

func hasWildcard(subscription string) bool {
	if strings.Contains(subscription, "+") {
		return true
	}

	if strings.Contains(subscription, "#") {
		return true
	}

	return false
}
