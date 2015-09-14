package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	gip, err := GetGlobalIP()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	fmt.Println(gip)
}

func GetGlobalIP() (string, error) {
	response, err := http.Get("http://whatismyip.akamai.com")
	if err != nil {
		return "", err
	}

	if response.StatusCode != 200 {
		err := fmt.Errorf("could not get global IP: %s", response.Status)
		return "", err
	}

	globalIP, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(globalIP), nil
}
