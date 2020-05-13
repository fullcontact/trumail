package api

import (
	"encoding/xml"
	"github.com/labstack/echo"
	"github.com/sdwolfe32/trumail/verifier"
	"net/http"
	"strings"
)

// Lookup contains all output data for an email verification Lookup
type Lookup struct {
	XMLName     xml.Name `json:"-" xml:"lookup"`
	Address     string   `json:"address" xml:"address"`
	Username    string   `json:"username" xml:"username"`
	Domain      string   `json:"domain" xml:"domain"`
	MD5Hash     string   `json:"md5Hash" xml:"md5Hash"`
	ValidFormat bool     `json:"validFormat" xml:"validFormat"`
	Deliverable string    `json:"deliverable" xml:"deliverable"`
	FullInbox   bool     `json:"fullInbox" xml:"fullInbox"`
	HostExists  bool     `json:"hostExists" xml:"hostExists"`
	CatchAll    bool     `json:"catchAll" xml:"catchAll"`
}

// LookupHandler performs a single email verification and returns
// a fully populated lookup or an error
func LookupHandler(v *verifier.Verifier) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Perform the unlimited verification
		lookup, err := v.Verify(c.Param("email"))
		deliverable := "true"
		if err != nil {
			if strings.HasPrefix(err.Error(), "Fatal") {
				return FormatEncoder(c, http.StatusInternalServerError, err)
			}
			deliverable = ""
		} else if !lookup.Deliverable {
			deliverable = "false"
		}
		return FormatEncoder(c, http.StatusOK, &Lookup{
			Address:     lookup.Address.Address,
			Username:    lookup.Username,
			Domain:      lookup.Domain,
			MD5Hash:     lookup.MD5Hash,
			ValidFormat: lookup.ValidFormat,
			Deliverable: deliverable,
			FullInbox:   lookup.FullInbox,
			HostExists:  lookup.HostExists,
			CatchAll:    lookup.CatchAll,
		})
	}
}
