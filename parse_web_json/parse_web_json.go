package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	// http://docs.aws.amazon.com/ja_jp/general/latest/gr/aws-ip-ranges.html
	AWS_IP_RANGE_URL = "https://ip-ranges.amazonaws.com/ip-ranges.json"
)

type AWSIPRangeResp struct {
	SyncToken  string       `json:"syncToken"`
	CreateDate string       `json:"createDate"`
	Prefixes   []AWSIPRange `json:"prefixes"`
}

type AWSIPRange struct {
	IPPrefix string `json:"ip_prefix"`
	Region   string `json:"region"`
	Service  string `json:"service"`
}

func main() {
	awsIPs, err := requestAWSGlobalIP()
	if err != nil {
		fmt.Printf("fail parse response json: %v\n", err)
		os.Exit(1)
	}

	if awsIPs != nil {
		fmt.Printf("sync token: %s, created at %s\n", awsIPs.SyncToken, awsIPs.CreateDate)
		for _, ip := range awsIPs.Prefixes {
			fmt.Printf("%s (%s) created at %s\n", ip.IPPrefix, ip.Region, ip.Service)
		}
	}
}

func requestAWSGlobalIP() (*AWSIPRangeResp, error) {
	resp, err := http.Get(AWS_IP_RANGE_URL)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(bufio.NewReader(resp.Body))

	var awsIPs AWSIPRangeResp
	if err := dec.Decode(&awsIPs); err == io.EOF {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &awsIPs, nil
}
