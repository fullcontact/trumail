package verifier

// Verifier contains all dependencies needed to perform educated email
// verification lookups
type Verifier struct {
	hostname, sourceAddr string
	disposabler          *Disposabler
}

// Lookup contains all output data for an email verification Lookup
type Lookup struct {
	Address
	ValidFormat, Deliverable, FullInbox, HostExists, CatchAll, Disposable bool
	ErrorDetails                                                     string
}

// NewVerifier generates a new Verifier using the passed hostname and
// source email address
func NewVerifier(hostname, sourceAddr string) *Verifier {
	return &Verifier{hostname, sourceAddr, NewDisposabler()}
}

// Verify performs an email verification on the passed email address
func (v *Verifier) Verify(email string) (*Lookup, error) {
	// Allocate memory for the Lookup
	var l Lookup
	l.Address.Address = email

	// First parse the email address passed
	address, err := ParseAddress(email)
	if err != nil {
		l.ValidFormat = false
		return &l, nil
	}
	l.ValidFormat = true
	l.Address = *address

	// Attempt to form an SMTP Connection
	del, err := NewDeliverabler(address.Domain, v.hostname, v.sourceAddr)
	if err != nil {
		l.ErrorDetails = err.Error()
		if le := ParseSMTPError(err); le != nil && le.Fatal {
			return nil, le
		} else {
			return &l, nil
		}
	}
	defer del.Close() // Defer close the SMTP connection

	// Host exists if we've successfully formed a connection
	l.HostExists = true

	if v.disposabler.IsDisposable(address.Domain) {
		l.Disposable = true
		return &l, nil
	}

	// Retrieve the catchall status and check deliverability
	if del.HasCatchAll(3) {
		l.CatchAll = true
		l.Deliverable = true
	} else {
		if err := del.IsDeliverable(address.Address, 3); err != nil {
			l.ErrorDetails = err.Error()
			if le := ParseSMTPError(err); le != nil {
				if le.Message == ErrFullInbox {
					l.FullInbox = true // set FullInbox and return no error
					return &l, nil
				}
				return &l, le // Return if there's a true error
			}
		} else {
			l.Deliverable = true
		}
	}
	return &l, nil
}
