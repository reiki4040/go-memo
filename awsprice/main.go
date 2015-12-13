package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/reiki4040/cstore"
)

type EC2Price struct {
	FormatVersion   string
	Disclaimer      string
	Offercode       string
	Version         string
	Publicationdate string
	Products        map[string]SKU
	Terms           Term
}

type SKU struct {
	Sku           string
	ProfuctFamiry string
	Attributes    Attribute
}

type Attribute struct {
	Servicecode           string
	Location              string
	LocationType          string
	InstanceType          string
	CurrentGeneration     string
	InstanceFamily        string
	Vcpu                  string
	PhysicalProcessor     string
	ClockSpeed            string
	Memory                string
	Storage               string
	NetworkPerformance    string
	ProcessorArchitecture string
	Tenancy               string
	OperatingSystem       string
	LicenseModel          string
	Usagetype             string
	Operation             string
	PreInstalledSw        string
}

type Term struct {
	Ondemand map[string]map[string]Offer
	Reserved map[string]map[string]Offer
}

type Offer struct {
	OfferTermCode   string
	SKU             string
	EffectiveDate   string
	PriceDimensions map[string]PriceDimensions
	TermAttributes  TermAttribute
}

type TermAttribute struct {
	LeaseContractLength string
	PurchaseOption      string
}

type PriceDimensions struct {
	RateCode     string
	Description  string
	BeginRange   string
	EndRange     string
	Unit         string
	PricePerUnit PricePerUnit
	appliesTo    []string
}

type PricePerUnit struct {
	USD string
}

func main() {
	cs, err := cstore.NewCStore("ec2", "ec2.json", cstore.JSON)
	if err != nil {
		log.Fatalf("%v", err)
	}
	ec2price := EC2Price{}
	err = cs.GetWithoutValidate(&ec2price)
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Printf("Products: %d, Ondemand: %d, Reserved: %d\n", len(ec2price.Products), len(ec2price.Terms.Ondemand), len(ec2price.Terms.Reserved))

	pCode := ""
	for _, sku := range ec2price.Products {
		a := sku.Attributes
		if a.Location != "Asia Pacific (Tokyo)" {
			continue
		}

		if a.InstanceType != "c3.large" {
			continue
		}

		if a.OperatingSystem != "Linux" {
			continue
		}

		if a.Tenancy != "Shared" {
			continue
		}

		pCode = sku.Sku
	}

	t := ec2price.Terms.Ondemand[pCode]
	for _, offer := range t {
		fmt.Printf("%s ", offer.OfferTermCode)
		fmt.Printf("%s ", offer.SKU)
		for _, p := range offer.PriceDimensions {
			fmt.Printf("%s ", p.Description)
			fmt.Printf("$%s\n", p.PricePerUnit.USD)

			usd, err := strconv.ParseFloat(p.PricePerUnit.USD, 32)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("hour: $%f, day: $%f, month(30days): $%f\n", usd, usd*24, usd*24*30)

			yen := 125.0
			fmt.Printf("[yen]hour: %.0fyen, day: %.0fyen, month(30days): %.0fyen\n", usd*yen, usd*24*yen, usd*24*30*yen)
		}

	}
}
