package gomq

import "fmt"

func categorizeChannel(channelDir, channel string) string {
	if channelDir == "" {
		return channel
	}
	return fmt.Sprintf("%s/%s", channelDir, channel)
}
