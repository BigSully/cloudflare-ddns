package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"flag"

	"strings"

	"github.com/joho/godotenv"

	"cloudflare-ddns/util"

	"github.com/cloudflare/cloudflare-go"
)

type Cloudflare struct {
	Token    string
	ZoneID   string
	ZoneName string
}

func (cf Cloudflare) listAllRecords(name string) {
	api, _ := cloudflare.NewWithAPIToken(cf.Token)

	// Fetch all records for a zone
	recs, err := api.DNSRecords(context.Background(), cf.ZoneID, cloudflare.DNSRecord{})
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range recs {
		if name != "" {
			if strings.Contains(r.Name, name) {
				fmt.Printf("%s: %s, %s\n", r.Name, r.Type, r.ID)
			}
		} else {
			fmt.Printf("%s: %s, %s\n", r.Name, r.Type, r.ID)
		}
	}
}

func (cf Cloudflare) findByRecord(query cloudflare.DNSRecord) (records []cloudflare.DNSRecord, err error) {
	api, _ := cloudflare.NewWithAPIToken(cf.Token)

	// Fetch all records for a zone
	records, err = api.DNSRecords(context.Background(), cf.ZoneID, query)
	if err != nil {
		log.Fatal(err)
	}

	return
}

func (cf Cloudflare) findOne(query cloudflare.DNSRecord) (record cloudflare.DNSRecord, err error) {
	recs, err := cf.findByRecord(query)
	if err != nil || recs == nil || len(recs) == 0 {
		return
	}

	return recs[0], nil
}

func (cf Cloudflare) bindInterface() (err error) {
	api, _ := cloudflare.NewWithAPIToken(cf.Token)
	ctx := context.Background()

	ifName := util.Getenv("IF_NAME", "eth0")
	ipv4, err := util.GetInterfaceIpv4Addr(ifName)
	if err != nil {
		return
	}

	query := cloudflare.DNSRecord{
		Name: fmt.Sprintf("%s.%s", os.Getenv("IF_BIND"), cf.ZoneName),
		Type: util.Getenv("IF_BIND_TYPE", "A"),
	}
	record, err := cf.findOne(query)
	if err != nil {
		return
	}

	// create new
	if record.ID == "" {
		newRecord := cloudflare.DNSRecord{
			Name:    query.Name,
			Type:    query.Type,
			Content: ipv4,
		}
		if _, err = api.CreateDNSRecord(ctx, cf.ZoneID, newRecord); err != nil {
			log.Fatal(err)
		}
		return
	}

	// update
	record.Content = ipv4

	if err := api.UpdateDNSRecord(ctx, cf.ZoneID, record.ID, record); err != nil {
		log.Fatal(err)
	}
	return
}

func (cf Cloudflare) bindPublicIP() (err error) {
	api, _ := cloudflare.NewWithAPIToken(cf.Token)
	ctx := context.Background()

	ipv4, err := util.GetPublicIP()
	if err != nil {
		return
	}

	query := cloudflare.DNSRecord{
		Name: fmt.Sprintf("%s.%s", os.Getenv("PUB_BIND"), cf.ZoneName),
		Type: util.Getenv("PUB_BIND_TYPE", "A"),
	}
	record, err := cf.findOne(query)
	if err != nil {
		return
	}

	// create new
	if record.ID == "" {
		newRecord := cloudflare.DNSRecord{
			Name:    query.Name,
			Type:    query.Type,
			Content: ipv4,
		}
		if _, err = api.CreateDNSRecord(ctx, cf.ZoneID, newRecord); err != nil {
			log.Fatal(err)
		}
		return
	}

	// update
	record.Content = ipv4

	if err = api.UpdateDNSRecord(ctx, cf.ZoneID, record.ID, record); err != nil {
		log.Fatal(err)
	}
	return
}

func NewCloudflare(token, zoneName string) (cf Cloudflare, err error) {
	api, _ := cloudflare.NewWithAPIToken(token)
	zoneID, err := api.ZoneIDByName(zoneName)
	if err != nil {
		return
	}
	cf = Cloudflare{token, zoneID, zoneName}
	return
}

func main() {

	// Load the .env file in the current directory
	godotenv.Load()
	token := os.Getenv("CF_TOKEN")
	if token == "" {
		fmt.Println("Please set token in .env file")
		return
	}
	domain := os.Getenv("CF_DOMAIN")
	if domain == "" {
		fmt.Println("Please set domain in .env file")
		return
	}

	cf, err := NewCloudflare(token, domain)
	if err != nil {
		log.Fatal(err)
		return
	}

	flag.Parse()
	args := flag.Args()
	subcmd := ""
	if len(args) > 0 {
		subcmd = args[0]
		args = args[1:]
	}

	switch subcmd {
	case "list":
		name := ""
		if len(args) > 0 {
			name = args[0]
		}
		cf.listAllRecords(name)
		os.Exit(0)
	default:
	}

	if ifBind := util.Getenv("IF_BIND", ""); ifBind != "" {
		if err := cf.bindInterface(); err != nil {
			log.Fatal(err)
		}
	}

	if pubBind := util.Getenv("PUB_BIND", ""); pubBind != "" {
		if err := cf.bindPublicIP(); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("%s - complete.\n", time.Now().Format("2006-01-02 15:04:05"))

}
