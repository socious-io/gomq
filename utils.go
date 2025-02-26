package gomq

import "fmt"

func categorizeChannel(channelDir, channel string) string {
	return fmt.Sprintf("%s/%s", channelDir, channel)
}
