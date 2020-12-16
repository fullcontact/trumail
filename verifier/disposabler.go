package verifier

import (
	"s32x.com/httpclient"
	"strings"
	"sync"
)

// updateInterval is how often we should reach out to update
// the disposable address map
const disposableDomainsUrl = "https://raw.githubusercontent.com/fullcontact/trumail/v1.1.5/disposable_emails.txt"

// Disposabler contains the map of known disposable email domains
type Disposabler struct {
	client  *httpclient.Client
	dispMap *sync.Map
}

// NewDisposabler creates a new Disposabler and starts a domain farmer
// that retrieves all known disposable domains periodically
func NewDisposabler() *Disposabler {
	d := &Disposabler{httpclient.New(), &sync.Map{}}
	go d.farmDomains()
	return d
}

// IsDisposable tests whether a string is among the known set of disposable
// mailbox domains. Returns true if the address is disposable
func (d *Disposabler) IsDisposable(domain string) bool {
	_, ok := d.dispMap.Load(domain)
	return ok
}

// farmDomains retrieves new disposable domains every set interval
func (d *Disposabler) farmDomains() error {
	for {
		// Perform the request for the domain list
		body, err := d.client.Get(disposableDomainsUrl).String()
		if err != nil {
			continue
		}

		// Split
		for _, domain := range strings.Split(body, "\n") {
			d.dispMap.Store(strings.TrimSpace(domain), true)
		}
	}
}
