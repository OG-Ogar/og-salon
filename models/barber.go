package models

import "strings"

// Barber account statuses. A suspended barber cannot log in and is
// hidden from customers.
const (
	BarberActive    = "Active"
	BarberSuspended = "Suspended"
	BarberArchived  = "Archived"
)

// Provider types. Barbers and hair stylists share the same account and
// workspace, so the type is only used to label them correctly in public
// places such as the home feed.
const (
	ProviderBarber      = "barber"
	ProviderHairstylist = "hairstylist"
)

// NormalizeProviderType keeps stored and submitted values inside the two
// supported provider types. Anything unknown is treated as a barber so old
// accounts created before the distinction continue to work.
func NormalizeProviderType(value string) string {
	if strings.EqualFold(strings.TrimSpace(value), ProviderHairstylist) {
		return ProviderHairstylist
	}
	return ProviderBarber
}

// Barber is a shop owner who publishes hairstyles and takes bookings.
type Barber struct {
	ID int

	FullName    string
	PhoneNumber string
	Email       string

	// PasswordHash holds the encoded PBKDF2 digest produced by the
	// security package. The plain password is never stored.
	PasswordHash string

	ShopName    string
	ShopAddress string
	State       string
	LocalArea   string

	// The two photos are separate on purpose. ProfileImage is the shop
	// itself, shown wide at the top of the page; PersonalImage is the
	// barber, shown beside their name so a customer knows who will be
	// cutting their hair. Either may be empty.
	ProfileImage  string
	PersonalImage string

	// ProviderType distinguishes barbers from hair stylists. Both use the
	// same account and workspace; the type only changes the label shown to
	// clients.
	ProviderType string

	Status string
}

// IsActive reports whether the barber may log in and receive bookings.
func (b Barber) IsActive() bool {
	return b.Status == BarberActive
}

// ProviderLabel is the public role label for this account.
func (b Barber) ProviderLabel() string {
	if NormalizeProviderType(b.ProviderType) == ProviderHairstylist {
		return "Hair Stylist"
	}
	return "Barber"
}

// ProviderNoun is the short role word used in calls to action, such as
// "Message stylist" for a stylist and "Message barber" for a barber.
func (b Barber) ProviderNoun() string {
	if NormalizeProviderType(b.ProviderType) == ProviderHairstylist {
		return "stylist"
	}
	return "barber"
}
