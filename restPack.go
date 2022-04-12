package restPack

import (
	"fmt"
	"net/http"
	"strings"
)

//creates a request for a rest api
func newRequest(endpoint string, token_key string, token_value string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://%s", endpoint), nil)
	req.Header.Add(token_key, token_value)
	req.Header.Add("accept", "application/json")
	if err != nil {
		return nil, err
	}
	return req, err
}

type Outlet struct {
	Ip     string
	Outlet string
}

//creates a new outlet struct
func newOutlet(outlet string, ip string) Outlet {
	o := Outlet{Outlet: outlet, Ip: ip}
	return o
}

//checks if the outlet is already in the slice
func contains(o []Outlet, e string) bool {
	for _, a := range o {
		if a.Outlet == e {
			return true
		}
	}
	return false
}

//gets the outlet number and pdu ip from the rest api
func GetOutlet(response []byte, value string) []Outlet {
	var outlet string
	var ip string
	var ip_leg string
	var outlet_leg string
	var outlet_s Outlet
	outlet_slice := make([]Outlet, 0)

	splitted := strings.Split(string(response), ",")
	for _, v := range splitted {
		if strings.Contains(v, "PO") {
			outlet_leg = outlet
			outlet = strings.Split(v, ":")[1]
		}
		if strings.Contains(v, "192") {
			ip_leg = ip
			ip = strings.Split(v, ":")[1]
		}
		if outlet != "" && ip != "" && ip != ip_leg && outlet != outlet_leg {
			if !contains(outlet_slice, outlet) {
				outlet_s = newOutlet(outlet, ip)
				outlet_slice = append(outlet_slice, outlet_s)
			}
		}
	}
	return outlet_slice
}

//api_key := "7e5bce0cee12c654ea0c209c50defa49e5e22d4b"
//response, err := http.Get("https://netbox.habana-labs.com/api/dcim/power-ports/?device=ofer-test-hls2")
// url := "netbox.habana-labs.com/api/dcim/power-ports/?device=ofer-test-hls2"

// req, _ := newRequest(url, "Authorization", "Token 7e5bce0cee12c654ea0c209c50defa49e5e22d4b")
// log.Println(req)
// res, _ := http.DefaultClient.Do(req)
// responseData, err := ioutil.ReadAll(res.Body)
// if err != nil {
// 	log.Fatal(err)
// }
// GetOutlet(responseData, "PSU")
