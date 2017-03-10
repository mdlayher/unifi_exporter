package unifi

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// Stations returns all of the Stations for a specified site name.
func (c *Client) Stations(siteName string) ([]*Station, error) {
	var v struct {
		Stations []*Station `json:"data"`
	}

	req, err := c.newRequest(
		"GET",
		fmt.Sprintf("/api/s/%s/stat/sta", siteName),
		nil,
	)
	if err != nil {
		return nil, err
	}

	_, err = c.do(req, &v)
	return v.Stations, err
}

// A Station is a client connected to a UniFi access point.
type Station struct {
	ID              string
	APMAC           net.HardwareAddr
	AssociationTime time.Time
	Channel         int
	FirstSeen       time.Time
	Hostname        string // Device-provided name
	IdleTime        time.Duration
	IP              net.IP
	IsWired         bool
	LastSeen        time.Time
	MAC             net.HardwareAddr
	RoamCount       int
	Name            string // Unifi-set name
	Noise           int
	RSSI            int
	SiteID          string
	Stats           *StationStats
	Uptime          time.Duration
	UserID          string
}

// StationStats contains station network activity statistics.
type StationStats struct {
	ReceiveBytes    int64
	ReceivePackets  int64
	ReceiveRate     int
	TransmitBytes   int64
	TransmitPackets int64
	TransmitPower   int
	TransmitRate    int
}

// UnmarshalJSON unmarshals the raw JSON representation of a Station.
func (s *Station) UnmarshalJSON(b []byte) error {
	var sta station
	if err := json.Unmarshal(b, &sta); err != nil {
		return err
	}

	apMAC, err := net.ParseMAC(sta.ApMac)
	if !sta.IsWired && err != nil {
		return err
	}

	mac, err := net.ParseMAC(sta.Mac)
	if err != nil {
		return err
	}

	*s = Station{
		ID:              sta.ID,
		APMAC:           apMAC,
		AssociationTime: time.Unix(int64(sta.AssocTime), 0),
		Channel:         sta.Channel,
		FirstSeen:       time.Unix(int64(sta.FirstSeen), 0),
		Hostname:        sta.Hostname,
		IdleTime:        time.Duration(time.Duration(sta.Idletime) * time.Second),
		IP:              net.ParseIP(sta.IP),
		IsWired:         sta.IsWired,
		LastSeen:        time.Unix(int64(sta.LastSeen), 0),
		MAC:             mac,
		Name:            sta.Name,
		Noise:           sta.Noise,
		RSSI:            sta.RSSI,
		RoamCount:       sta.RoamCount,
		SiteID:          sta.SiteID,
		Stats: &StationStats{
			ReceiveBytes:    sta.RxBytes,
			ReceivePackets:  sta.RxPackets,
			ReceiveRate:     sta.RxRate,
			TransmitBytes:   sta.TxBytes,
			TransmitPackets: sta.TxPackets,
			TransmitPower:   sta.TxPower,
			TransmitRate:    sta.TxRate,
		},
		Uptime: time.Duration(time.Duration(sta.Uptime) * time.Second),
		UserID: sta.UserID,
	}

	return nil
}

// A station is the raw structure of a Station returned from the UniFi Controller
// API.
type station struct {
	// TODO(mdlayher): give all fields appropriate names and data types.
	ID               string `json:"_id"`
	IsGuestByUap     bool   `json:"_is_guest_by_uap"`
	LastSeenByUap    int    `json:"_last_seen_by_uap"`
	UptimeByUap      int    `json:"_uptime_by_uap"`
	ApMac            string `json:"ap_mac"`
	AssocTime        int    `json:"assoc_time"`
	Authorized       bool   `json:"authorized"`
	Bssid            string `json:"bssid"`
	BytesR           int64  `json:"bytes-r"`
	Ccq              int    `json:"ccq"`
	Channel          int    `json:"channel"`
	Essid            string `json:"essid"`
	FirstSeen        int    `json:"first_seen"`
	Hostname         string `json:"hostname"`
	Idletime         int    `json:"idletime"`
	IP               string `json:"ip"`
	IsGuest          bool   `json:"is_guest"`
	IsWired          bool   `json:"is_wired"`
	LastSeen         int    `json:"last_seen"`
	Mac              string `json:"mac"`
	Name             string `json:"name"`
	Noise            int    `json:"noise"`
	Oui              string `json:"oui"`
	PowersaveEnabled bool   `json:"powersave_enabled"`
	QosPolicyApplied bool   `json:"qos_policy_applied"`
	Radio            string `json:"radio"`
	RadioProto       string `json:"radio_proto"`
	RoamCount        int    `json:"roam_count"`
	RSSI             int    `json:"rssi"`
	RxBytes          int64  `json:"rx_bytes"`
	RxBytesR         int64  `json:"rx_bytes-r"`
	RxPackets        int64  `json:"rx_packets"`
	RxRate           int    `json:"rx_rate"`
	Signal           int    `json:"signal"`
	SiteID           string `json:"site_id"`
	TxBytes          int64  `json:"tx_bytes"`
	TxBytesR         int64  `json:"tx_bytes-r"`
	TxPackets        int64  `json:"tx_packets"`
	TxPower          int    `json:"tx_power"`
	TxRate           int    `json:"tx_rate"`
	Uptime           int    `json:"uptime"`
	UserID           string `json:"user_id"`
}
