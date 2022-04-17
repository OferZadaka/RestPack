package RestPack

import (
	"fmt"
	"net/http"
	"strings"
)

//creates a request for a rest api
func NewRequest(endpoint string, token_key string, token_value string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://%s", endpoint), nil)
	req.Header.Add(token_key, token_value)
	req.Header.Add("accept", "application/json")
	if err != nil {
		return nil, err
	}
	return req, err
}

type Outlet struct {
	Ip           string
	Outlet       string
	Manufacturer string
	Model        string
}

//creates a new outlet struct
func newOutlet(outlet string, ip string, manufacturer string, model string) Outlet {
	o := Outlet{Outlet: outlet, Ip: ip, Manufacturer: manufacturer, Model: model}
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
func GetOutlet(response []byte) []Outlet {
	var outlet string
	var ip string
	var ip_leg string
	var outlet_leg string
	var outlet_s Outlet

	outlet_slice := make([]Outlet, 0)
	//connects to the netbox api and gets the PDU details
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
				outlet_s = newOutlet(outlet, ip, "", "")
				outlet_slice = append(outlet_slice, outlet_s)
			}
		}
	}
	return outlet_slice
}

func GetManufacturer(response []byte) []string {
	var manufacturer string
	var model string
	var area bool
	var slug_index int = 0
	manufacturer_slice := make([]string, 0)

	splitted := strings.Split(string(response), ",")
	for _, v := range splitted {
		if strings.Contains(v, "manufacturer") {
			area = true
			continue
		}
		if strings.Contains(v, "slug") && area && slug_index == 0 {
			manufacturer = strings.Split(v, ":")[1]
			slug_index = 1
			continue

		}
		if strings.Contains(v, "slug") && area && slug_index == 1 {
			model = strings.Split(v, ":")[1]
		}
		if model != "" && manufacturer != "" {
			manufacturer_slice = append(manufacturer_slice, manufacturer, model)
			return manufacturer_slice
		}
	}
	return []string{"", ""}
}

//////////////////////////////////////////////////////////////
// func getRequest(url string, name string) []byte {
// 	url = url + name
// 	token := "Token 7e5bce0cee12c654ea0c209c50defa49e5e22d4b"

// 	req, _ := NewRequest(url, "Authorization", token)
// 	res, _ := http.DefaultClient.Do(req)

// 	responseData, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return responseData
// }

// func appendManufacturerDetails(outlet []Outlet) []Outlet {
// 	for _, out := range outlet {
// 		out.Ip = strings.Replace(out.Ip, "\"", "", -1)
// 		manufacturerDetails := getRequest("netbox.habana-labs.com/api/dcim/devices/?name=", out.Ip)
// 		manufacturer_slice := GetManufacturer(manufacturerDetails)
// 		out.Manufacturer = manufacturer_slice[0]
// 		out.Model = manufacturer_slice[1]
// 	}
// 	return outlet
// }

// func main() {
// 	name := "ofer-test-hls2"
// 	pdu_details := getRequest("netbox.habana-labs.com/api/dcim/power-ports/?device=", name)
// 	pdu_slice := GetOutlet(pdu_details)
// 	pdu_slice = appendManufacturerDetails(pdu_slice)

// }
