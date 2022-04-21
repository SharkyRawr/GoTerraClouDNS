package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/SharkyRawr/GoTerraClouDNS/cloudns"
)

var authID = flag.Int("sub-auth-id", 0, "Sub Auth ID")
var authPassword = flag.String("auth-password", "", "Auth Password")
var zoneFlag = flag.String("zone", "", "Zone name")
var outPath = flag.String("out", "import.tf", "Output file path")

func main() {
	flag.Parse()
	if *authID <= 0 || *authPassword == "" || *zoneFlag == "" || *outPath == "" {
		flag.Usage()
		return
	}

	fd, err := os.OpenFile(*outPath, os.O_WRONLY|os.O_CREATE, 0644)
	defer fd.Close()
	writer := bufio.NewWriter(fd)
	defer writer.Flush()

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

	tpl, err := template.New("tf_dns_record").Parse(`resource "cloudns_dns_record" "{{.SafeID}}_{{.ID}}" {
  # {{.Host}}.{{.Zone}} {{.TTL}} in {{.Type}} {{.Record}}
  name  = "{{.Host}}"
  zone  = "{{.Zone}}"
  type  = "{{.Type}}"
  value = "{{.SafeRecord}}"
  ttl   = "{{.TTL}}"
}`)
	if err != nil {
		panic(err)
	}

	records, err := api.GetRecordsForZone(zone.Name)
	for _, record := range records {
		safeID := strings.Replace(zone.Name, ".", "_", -1)
		vals := struct {
			SafeID     string
			ID         string
			Zone       string
			Host       string
			Type       string
			Record     string
			SafeRecord string
			TTL        string
		}{
			SafeID:     safeID,
			ID:         record.ID,
			Zone:       zone.Name,
			Host:       record.Host,
			Type:       record.Type,
			Record:     record.Record,
			SafeRecord: strings.Replace(record.Record, "\"", "\\\"", -1),
			TTL:        record.TTL,
		}
		err := tpl.Execute(writer, vals)
		writer.WriteString("\n\n")
		fmt.Printf("terraform import cloudns_dns_record.%s_%s cloudns_dns_record/%s_%s\n", safeID, record.ID, safeID, record.ID)
		if err != nil {
			panic(err)
		}

	}
}
