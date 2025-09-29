package types

type LinkNodeProtocol interface {

	// Mount associates this protocol instance with the given node
	Mount(s Node)

	// ConnectLink adds a new link to the protocol's management
	ConnectLink(link Link) error

	// DisconnectLink removes the given link from the protocol's management
	DisconnectLink(link Link) error

	// UpdateLinks calculates and returns the current list of established links
	UpdateLinks() ([]Link, error)

	// Established returns only established links
	Established() []Link

	// Links returns all links managed by this protocol instance (including unestablished links)
	Links() []Link
}
