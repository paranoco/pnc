package main

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/urfave/cli/v2"
	"net/http"
)

func PublicIpCommand(c *cli.Context) error {
	client := retryablehttp.NewClient()
	client.RetryMax = 3

	req, err := retryablehttp.NewRequest(http.MethodGet, "https://api.ipify.org?format=json", nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	var response map[string]interface{}

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return err
	}

	fmt.Println(response["ip"])
	return nil
}