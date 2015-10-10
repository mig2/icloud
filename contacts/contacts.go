package contacts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"

	"github.com/mig2/icloud/engine"
)

type ICloudContactsResponse struct {
	HeaderPositions map[string]int            `json:"headerPositions"`
	SyncToken       string                    `json:"syncToken"`
	contactsOrder   []string                  `json:"contactsOrder"`
	MeCardId        string                    `json:"meCardId"`
	Collections     []ICloudContactCollection `json:"collections"`
	PrefToken       string                    `json:"prefToken"`
	Groups          []ICloudContactGroup      `json:"groups"`
	Contacts        []ICloudContact           `json:"contacts"`
}

type ICloudContact struct {
	FirstName  string               `json:"firstName"`
	MiddleName string               `json:"middleName"`
	LastName   string               `json:"lastName"`
	Suffix     string               `json:"suffix"`
	ContactId  string               `json:"contactId"`
	Prefix     string               `json:"prefix"`
	Phones     []ICloudContactPhone `json:"phones"`
	Etag       string               `json:"etag"`
	IsCompany  bool                 `json:"isCompany"`
}

type ICloudContactPhone struct {
	Field string `json:"field"`
	Label string `json:"label"`
}

type ICloudContactCollection struct {
	GroupsOrder  []string `json:"groupsOrder"`
	Etag         string   `json:"etag"`
	CollectionId string   `json:"collectionId"`
}

type ICloudContactGroup struct {
}

const (
	ContactsUrl string = "%v/co/startup"
)

func parseLocale(locale string) string {
	if len(locale) < 5 {
		return "en_US"
	}

	return locale[0:2] + "_" + locale[len(locale)-2:]
}

func Get(cloud *engine.ICloudEngine) (*ICloudContactsResponse, error) {
	v := url.Values{}
	v.Add("clientBuildNumber", cloud.ReportedVersion.BuildNumber)
	v.Add("clientId", cloud.ClientID)
	v.Add("clientVersion", "2.1")
	v.Add("dsid", cloud.User.Dsid)
	v.Add("locale", parseLocale(cloud.User.Locale))
	v.Add("order", "last,first")

	var req *http.Request
	var e error
	host, _, _ := net.SplitHostPort(cloud.Webservices["contacts"].Url)
	curl := fmt.Sprintf("%v?%v", host+"/co/startup", v.Encode())
	if req, e = http.NewRequest("GET", curl, nil); e != nil {
		return nil, e
	}

	// headers and stuff
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Origin", "https://www.icloud.com")
	req.Header.Set("Referer", "https://www.icloud.com/")
	req.Header.Set("User-Agent", "Opera/9.52 (X11; Linux i686; U; en)")
	req.Header.Set("Host", host[8:])

	var resp *http.Response
	if resp, e = cloud.Client.Do(req); e != nil {
		return nil, e
	}
	defer resp.Body.Close()

	var body []byte
	if body, e = ioutil.ReadAll(resp.Body); e != nil {
		return nil, e
	}

	iresp := new(ICloudContactsResponse)
	json.Unmarshal(body, iresp)
	return iresp, nil
}
