package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("domain, hasMX, hasSPF, sprRecord, hasDMARC, dmarcRecord\n")

	fmt.Println("Enter the email Address below: ")
	for scanner.Scan() {
		checkDomain(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error: could not read from input: %v\n", err)
	}

}

func checkDomain(domain string) {
	var (
		hasMX, hasSPF, hasDMARC bool
		spfRecord, dmarcRecord  string
	)

	mxRecords, err := net.LookupMX(domain)

	if err != nil {
		log.Println("Error: %v\n", err)
	}

	if len(mxRecords) > 0 {
		hasMX = true
	}

	txtRecords, err := net.LookupTXT(domain)

	if err != nil {
		log.Println("Error: %v\n", err)
	}

	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}
	}

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = record
			break
		}
	}

	fmt.Print("\n")
	fmt.Printf("Domain: %v\n", domain)
	fmt.Printf("hasMX: %v\n", hasMX)
	fmt.Printf("hasSPF: %v\n", hasSPF)
	fmt.Printf("spfRecord: %v\n", spfRecord)
	fmt.Printf("hasDMARC: %v\n", hasDMARC)
	fmt.Printf("dmarcRecord: %v\n", dmarcRecord)
	fmt.Print("\n\n")

	return
}
