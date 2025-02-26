package gomq

// The global config variable holds the configuration for the application.
var config = new(Config)

// Config represents the configuration structure for the Nats MQ.
type Config struct {
	Url        string // Url represents the Nats host.
	Token      string // Token represents the Nats token.
	ChannelDir string // ChannelDir is a token that put the stream on each platform separately like ChannelDir/[channel].
	Consumers  map[string]func(interface{})
}

// Setup initializes the MQ SDK with the provided configuration.
func Setup(cfg Config) error {

	// Set the global configuration to the provided config.
	config = &cfg
	return nil // Return nil to indicate successful setup.
}
