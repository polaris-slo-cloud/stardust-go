package types

type LinkNodeProtocol interface {
	Mount(s Node)
	ConnectLink(link Link) error
	DisconnectLink(link Link) error
	UpdateLinks() ([]Link, error)
	Established() []Link
	Links() []Link
}
