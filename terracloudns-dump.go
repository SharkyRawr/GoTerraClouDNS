package main

import (
	"fmt"
	"git.catgirl.biz/sophie/goterracloudns/v2/cloudns"
)

func main() {
	api := cloudns.New("16535", "Sculptor-Jumble8-Chastise")
	loginStatus, err := api.Login()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", loginStatus)

	fmt.Println("Listing zones...")

	zones, err := api.GetZones()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", zones)

	fmt.Println("Listing zone records for first zonename:")
	firstZone := zones[0]
	records, err := api.GetRecordsForZone(firstZone.Name)
	fmt.Printf("%+v\n", records)
}
