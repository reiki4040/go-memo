package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"log"
)

func main() {
	region := "ap-northeast-1"

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	svc := route53.New(sess)
	resp, err := svc.ListHostedZones(&route53.ListHostedZonesInput{})
	if err != nil {
		log.Fatalf("failed list hosted zones: %v", err)
	}
	for _, z := range resp.HostedZones {
		log.Printf("hosted zone name: %s (%s)", aws.StringValue(z.Name), aws.StringValue(z.Config.Comment))
		in := &route53.ListResourceRecordSetsInput{
			HostedZoneId: z.Id,
		}
		r, err := svc.ListResourceRecordSets(in)
		if err != nil {
			log.Printf("failed list hosted zones: %v", err)
			continue
		}
		for _, rs := range r.ResourceRecordSets {
			log.Printf("record set %s %s TTL %d", aws.StringValue(rs.Type), aws.StringValue(rs.Name), aws.Int64Value(rs.TTL))
		}
	}
}
