package types

// LinkNodeProtocol defines methods which link node protocols have to implement
type LinkNodeProtocol interface {

	// Mount associates this protocol instance with the given node
	Mount(node Node)

	// ConnectLink establishes the given link
	ConnectLink(link Link) error

	// DisconnectLink removes the given link from established links
	DisconnectLink(link Link) error

	// UpdateLinks calculates and returns the current list of established links
	UpdateLinks() ([]Link, error)

	// Established returns only established links
	Established() []Link

	// Links returns all links managed by this protocol instance (including unestablished links)
	Links() []Link
}
