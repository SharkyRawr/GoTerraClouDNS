package cloudns

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

type RequestMethod int

const (
	METHOD_GET RequestMethod = iota
	METHOD_POST
)

const ApiHostname = "https://api.cloudns.net/"

const (
	EPLogin = "dns/login.json"
)

type StatusResponse struct {
	Status            string `json:"status"`
	StatusDescription string `json:"statusDescription"`
}

type ZoneListResponse struct {
	StatusResponse
	Name      string `json:"name"`
	Type      string `json:"type"`
	HasBulk   bool   `json:"hasBulk"`
	Zone      string `json:"zone"`
	Status    string `json:"status"`
	Serial    string `json:"serial"`
	IsUpdated int64  `json:"isUpdated"`
}

type ZoneRecord struct {
	StatusResponse
	ID               string  `json:"id"`
	Type             string  `json:"type"`
	Host             string  `json:"host"`
	Record           string  `json:"record"`
	Failover         string  `json:"failover"`
	TTL              string  `json:"ttl"`
	Status           int64   `json:"status"`
	DynamicurlStatus *int64  `json:"dynamicurl_status,omitempty"`
	Priority         *string `json:"priority,omitempty"`
	Weight           *string `json:"weight,omitempty"`
	Port             *string `json:"port,omitempty"`
}

type ClouDNSAPI struct {
	subAuthID    string
	authPassword string
}

func New(SubAuthID string, AuthPassword string) ClouDNSAPI {
	return ClouDNSAPI{
		subAuthID:    SubAuthID,
		authPassword: AuthPassword,
	}
}

func (c *ClouDNSAPI) doRequest(endpoint string, values url.Values, target interface{}) error {
	values.Set("sub-auth-id", c.subAuthID)
	values.Set("auth-password", c.authPassword)

	resp, err := http.PostForm(ApiHostname+endpoint, values)
	if err != nil {
		return err
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func (c *ClouDNSAPI) Login() (StatusResponse, error) {
	values := url.Values{}
	status := StatusResponse{}
	err := c.doRequest(EPLogin, values, &status)
	if err != nil {
		return status, err
	}
	if status.Status != "Success" {
		return status, errors.New("Login failed: " + status.StatusDescription)
	}
	return status, err
}

func (c *ClouDNSAPI) GetZones() ([]ZoneListResponse, error) {
	values := url.Values{}
	values.Set("page", "1")
	values.Set("rows-per-page", "100") // @todo needs to work with pages
	zoneList := []ZoneListResponse{}
	err := c.doRequest("dns/list-zones.json", values, &zoneList)
	return zoneList, err
}

func (c *ClouDNSAPI) GetRecordsForZone(zone string) ([]ZoneRecord, error) {
	values := url.Values{}
	values.Set("domain-name", zone)
	// @todo needs to work with pages, hostnames + fuzzy search, types
	var recordMap map[string]ZoneRecord
	err := c.doRequest("dns/records.json", values, &recordMap)

	zoneRecords := []ZoneRecord{}

	for _, element := range recordMap {
		zoneRecords = append(zoneRecords, element)
	}

	return zoneRecords, err
}
