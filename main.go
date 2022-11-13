package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/bunnydnsupdater")
	viper.AddConfigPath("$XDG_CONFIG_HOME/bunnydnsupdater")
	viper.AddConfigPath("$HOME/.config/bunnydnsupdater")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("error while loading the configuration file: %w", err))
	}

	var C config

	err = viper.Unmarshal(&C)
	if err != nil {
		panic(fmt.Errorf("error while reading the configuration file: %v", err))
	}

	validate := validator.New()

	if err := validate.Struct(&C); err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	httpClient := &http.Client{}

	zonesRequest, err := http.NewRequest("GET", "https://api.bunny.net/dnszone", nil)
	zonesRequest.Header.Set("AccessKey", C.AccessKey)

	if err != nil {
		panic(fmt.Errorf("error while creating HTTP request for getting DNS zones: %v", err))
	}

	zonesResponse, err := httpClient.Do(zonesRequest)

	if err != nil {
		panic(fmt.Errorf("error while getting DNS zones: %v", err))
	}

	if zonesResponse.StatusCode == 401 {
		panic(fmt.Errorf("your access key did not work: %s", zonesResponse.Status))
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(fmt.Errorf("error: %v", err))
		}
	}(zonesResponse.Body)

	zonesBody, err := io.ReadAll(zonesResponse.Body)

	if err != nil {
		panic(fmt.Errorf("error while reading body of DNS zones request: %v", err))
	}

	bunnyDNSZones, err := UnmarshalBunnyDNSZones(zonesBody)

	if err != nil {
		panic(fmt.Errorf("error while parsing DNS zones: %v", err))
	}

	zone, ok := lo.Find[BunnyDNSZone](bunnyDNSZones.Items, func(z BunnyDNSZone) bool {
		return z.Domain == C.Zone
	})

	if ok == false {
		panic(fmt.Errorf("could not find zone named %s", C.Zone))
	}

	record, ok := lo.Find[Record](zone.Records, func(r Record) bool {
		return r.Name == C.Record
	})

	if ok == false {
		panic(fmt.Errorf("could not find record named %s in zone %s", C.Zone, C.Record))
	}

	publicIPAddressRequest, err := http.Get("https://api.ipify.org")

	if err != nil {
		panic(fmt.Errorf("error while getting current IP address: %v", err))
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(fmt.Errorf("error: %v", err))
		}
	}(publicIPAddressRequest.Body)

	b, err := io.ReadAll(publicIPAddressRequest.Body)

	if err != nil {
		panic(fmt.Errorf("error while getting current IP address: %v", err))
	}

	publicIPAddress := string(b)

	if record.Value == publicIPAddress {
		fmt.Printf("IP addresses are the same: public %s; current record %s", publicIPAddress, record.Value)
	} else {
		fmt.Printf("IP addresses are not the same: public %s; current record %s\nupdating...", publicIPAddress, record.Value)
		postBody := struct {
			ID    int64
			Value string
		}{
			ID:    record.ID,
			Value: publicIPAddress,
		}

		postJSON, err := json.Marshal(postBody)

		if err != nil {
			panic(fmt.Errorf("error while updating record: %v", err))
		}

		postJSONBuf := bytes.NewBuffer(postJSON)

		updateRequest, err := http.NewRequest("POST",
			fmt.Sprintf("https://api.bunny.net/dnszone/%v/records/%v", zone.ID, record.ID),
			postJSONBuf)

		if err != nil {
			panic(fmt.Errorf("error creating HTTP request for updating DNS record: %v", err))
		}
		updateRequest.Header.Set("AccessKey", C.AccessKey)

		updateResponse, err := httpClient.Do(zonesRequest)

		if err != nil {
			panic(fmt.Errorf("error while updating DNS record: %v", err))
		}

		fmt.Println(updateResponse.Status)

	}

}
