package main

import (
	"flag"
	"fmt"
	"git.catgirl.biz/sophie/goterracloudns/v2/cloudns"
	"os"
	"strconv"
	"strings"
	"text/template"
)

var authID = flag.Int("sub-auth-id", 0, "Sub Auth ID")
var authPassword = flag.String("auth-password", "", "Auth Password")
var zoneFlag = flag.String("zone", "", "Zone name")

func main() {
	flag.Parse()
	if *authID <= 0 || *authPassword == "" || *zoneFlag == "" {
		flag.Usage()
		return
	}

	api := cloudns.New(strconv.Itoa(*authID), strings.Clone(*authPassword))
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
	fmt.Printf("Found %d zones...\n", len(zones))

	var zone cloudns.ZoneListResponse
	for _, z := range zones {
		if z.Name == *zoneFlag {
			zone = z
			break
		}
	}
	fmt.Printf("Reading zone records for %s and converting to Terraform...\n", zone.Name)

	tpl, err := template.New("tf_dns_record").Parse(`resource "cloudns_dns_record" "terraform_managed_record_{{.ID}}" {
  # {{.Host}}.{{.Zone}} {{.TTL}} in {{.Type}} {{.Record}}
  name  = "{{.Host}}"
  zone  = "{{.Zone}}"
  type  = "{{.Type}}"
  value = "{{.Record}}"
  ttl   = "{{.TTL}}"
}`)
	if err != nil {
		panic(err)
	}

	records, err := api.GetRecordsForZone(zone.Name)
	for _, record := range records {
		vals := struct {
			ID     string
			Zone   string
			Host   string
			Type   string
			Record string
			TTL    string
		}{
			ID:     record.ID,
			Zone:   zone.Name,
			Host:   record.Host,
			Type:   record.Type,
			Record: record.Record,
			TTL:    record.TTL,
		}
		err := tpl.Execute(os.Stdout, vals)
		os.Stdout.WriteString("\n\n")
		if err != nil {
			panic(err)
		}

	}
}
