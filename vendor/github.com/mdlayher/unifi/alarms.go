package unifi

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// Alarms returns all of the Alarms for a specified site name.
func (c *Client) Alarms(siteName string) ([]*Alarm, error) {
	var v struct {
		Alarms []*Alarm `json:"data"`
	}

	req, err := c.newRequest(
		"GET",
		fmt.Sprintf("/api/s/%s/list/alarm", siteName),
		nil,
	)
	if err != nil {
		return nil, err
	}

	_, err = c.do(req, &v)
	return v.Alarms, err
}

// An Alarm is an alert which is triggered when a Device becomes
// unavailable.
type Alarm struct {
	ID        string
	APMAC     net.HardwareAddr
	APName    string
	Archived  bool
	DateTime  time.Time
	Key       string
	Message   string
	SiteID    string
	Subsystem string
}

// UnmarshalJSON unmarshals the raw JSON representation of an Alarm.
func (a *Alarm) UnmarshalJSON(b []byte) error {
	var al alarm
	if err := json.Unmarshal(b, &al); err != nil {
		return err
	}

	mac, err := net.ParseMAC(al.AP)
	if err != nil {
		return err
	}

	t, err := time.Parse(time.RFC3339, al.DateTime)
	if err != nil {
		return err
	}

	*a = Alarm{
		ID:        al.ID,
		APMAC:     mac,
		APName:    al.APName,
		Archived:  al.Archived,
		DateTime:  t,
		Key:       al.Key,
		Message:   al.Msg,
		SiteID:    al.SiteID,
		Subsystem: al.Subsystem,
	}

	return nil
}

// An alarm is the raw structure of an Alarm returned from the UniFi Controller
// API.
type alarm struct {
	ID        string `json:"_id"`
	AP        string `json:"ap"`
	APName    string `json:"ap_name"`
	Archived  bool   `json:"archived"`
	DateTime  string `json:"datetime"`
	Key       string `json:"key"`
	Msg       string `json:"msg"`
	SiteID    string `json:"site_id"`
	Subsystem string `json:"subsystem"`
	// A UNIX timestamp field "time" exists here, but seems
	// redundant with DateTime
}
