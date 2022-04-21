package main

import (
	"testing"

	"github.com/SharkyRawr/GoTerraClouDNS/cloudns"
)

/*
	Probably should move it to the correct directory.
*/

func setup() cloudns.ClouDNSAPI {
	api := cloudns.New("16535", "Sculptor-Jumble8-Chastise")
	return api
}

func TestLogin(t *testing.T) {
	api := setup()
	_, err := api.Login()
	if err != nil {
		t.Error(err)
	}
}

func TestGetZones(t *testing.T) {
	api := setup()
	_, err := api.GetZones()
	if err != nil {
		t.Error(err)
	}
}

func TestGetRecords(t *testing.T) {
	api := setup()
	zones, err := api.GetZones()
	if err != nil {
		t.Error(err)
	}
	firstZone := zones[0]
	_, err = api.GetRecordsForZone(firstZone.Name)
	if err != nil {
		t.Error(err)
	}
}
