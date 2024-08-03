package sender

import "net/url"

type RCManagerServer struct {
	Url string
}

// Create a new RCManagerServer and verify the url
func (n *RCManagerServer) SetConn(u, token string) error {

	if _, err := url.ParseRequestURI(u); err != nil {
		return err
	}
	n.Url = u

	return nil
}
