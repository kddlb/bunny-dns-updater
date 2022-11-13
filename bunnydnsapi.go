package main

import "encoding/json"

func UnmarshalBunnyDNSZones(data []byte) (BunnyDNSZones, error) {
	var r BunnyDNSZones
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *BunnyDNSZones) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type BunnyDNSZones struct {
	Items        []BunnyDNSZone `json:"Items"`
	CurrentPage  int64          `json:"CurrentPage"`
	TotalItems   int64          `json:"TotalItems"`
	HasMoreItems bool           `json:"HasMoreItems"`
}

type BunnyDNSZone struct {
	ID                            int64    `json:"Id"`
	Domain                        string   `json:"Domain"`
	Records                       []Record `json:"Records"`
	DateModified                  string   `json:"DateModified"`
	DateCreated                   string   `json:"DateCreated"`
	NameserversDetected           bool     `json:"NameserversDetected"`
	CustomNameserversEnabled      bool     `json:"CustomNameserversEnabled"`
	Nameserver1                   string   `json:"Nameserver1"`
	Nameserver2                   string   `json:"Nameserver2"`
	SOAEmail                      string   `json:"SoaEmail"`
	NameserversNextCheck          string   `json:"NameserversNextCheck"`
	LoggingEnabled                bool     `json:"LoggingEnabled"`
	LoggingIPAnonymizationEnabled bool     `json:"LoggingIPAnonymizationEnabled"`
	LogAnonymizationType          int64    `json:"LogAnonymizationType"`
}

type Record struct {
	ID                    int64              `json:"Id"`
	Type                  int64              `json:"Type"`
	TTL                   int64              `json:"Ttl"`
	Value                 string             `json:"Value"`
	Name                  string             `json:"Name"`
	Weight                int64              `json:"Weight"`
	Priority              int64              `json:"Priority"`
	Port                  int64              `json:"Port"`
	Flags                 int64              `json:"Flags"`
	Tag                   *string            `json:"Tag"`
	Accelerated           bool               `json:"Accelerated"`
	AcceleratedPullZoneID int64              `json:"AcceleratedPullZoneId"`
	LinkName              string             `json:"LinkName"`
	IPGeoLocationInfo     *IPGeoLocationInfo `json:"IPGeoLocationInfo"`
	MonitorStatus         int64              `json:"MonitorStatus"`
	MonitorType           int64              `json:"MonitorType"`
	GeolocationLatitude   float32            `json:"GeolocationLatitude"`
	GeolocationLongitude  float32            `json:"GeolocationLongitude"`
	EnviromentalVariables []interface{}      `json:"EnviromentalVariables"`
	LatencyZone           string             `json:"LatencyZone"`
	SmartRoutingType      int64              `json:"SmartRoutingType"`
	Disabled              bool               `json:"Disabled"`
}

type IPGeoLocationInfo struct {
	CountryCode      string `json:"CountryCode"`
	Country          string `json:"Country"`
	Asn              int64  `json:"ASN"`
	OrganizationName string `json:"OrganizationName"`
	City             string `json:"City"`
}
