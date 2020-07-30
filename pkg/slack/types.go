package slack

//TODO: this pattern in go?

// ChannelID - strict type for slack channel ID
type ChannelID struct {
	value string
}

// Get returns string value for the channelID
func (cID *ChannelID) Get() string {
	return cID.value
}

// NewChannelID creates a strict type of channelID
func NewChannelID(id string) *ChannelID {
	return &ChannelID{
		value: id,
	}
}
