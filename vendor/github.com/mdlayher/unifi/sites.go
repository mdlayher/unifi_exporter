package unifi

// A Site is a physical location with UniFi devices managed by a UniFi
// Controller.
type Site struct {
	ID          string `json:"_id"`
	Description string `json:"desc"`
	Name        string `json:"name"`
	NumAPs      int    `json:"num_ap"`
	NumStations int    `json:"num_sta"`
	Role        string `json:"role"`
}

// Sites returns all of the Sites managed by a UniFi Controller.
func (c *Client) Sites() ([]*Site, error) {
	var v struct {
		Sites []*Site `json:"data"`
	}

	req, err := c.newRequest(
		"GET",
		"/api/self/sites",
		nil,
	)
	if err != nil {
		return nil, err
	}

	_, err = c.do(req, &v)
	return v.Sites, err
}
