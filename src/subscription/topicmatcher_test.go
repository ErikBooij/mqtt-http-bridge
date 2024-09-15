package subscription

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTopicMatcher(t *testing.T) {
	tm := newTopicMatcher()

	tt := []struct {
		topic         string
		subscription  string
		expectedMatch bool
	}{
		{"zigbee2mqtt/shortcut-button-001", "zigbee2mqtt/shortcut-button-001", true},
		{"zigbee2mqtt/shortcut-button-001", "zigbee2mqtt/shortcut-button-002", false},
		{"zigbee2mqtt/shortcut-button-001", "zigbee2mqtt/shortcut-button-00+", true},
		{"zigbee2mqtt/shortcut-button-001", "zigbee2mqtt/shortcut-button-00#", true},
		{"zigbee2mqtt/shortcut-button-001", "#", true},
		{"zigbee2mqtt/shortcut-button-001", "+", false},
		{"zigbee2mqtt/shortcut-button-001", "zigbee2mqtt/+", true},
		{"zigbee2mqtt/shortcut-button-001/test", "zigbee2mqtt/+/test", true},
		{"zigbee2mqtt/shortcut-button-001/test", "zigbee2mqtt/#", true},
	}

	for n, tc := range tt {
		t.Run(fmt.Sprintf("Topic Matcher Test Case #%d", n+1), func(t *testing.T) {
			assert.Equal(t, tm.match(tc.topic, tc.subscription), tc.expectedMatch, "Expected %s matching against %s to be %t", tc.topic, tc.subscription, tc.expectedMatch)
		})
	}
}
