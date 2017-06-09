package metrics

type Manager interface {
	AddRelay(relay *Relay)
	RemoveRelay(relayName string)
	GetRelay(relayName string) *Relay
}
