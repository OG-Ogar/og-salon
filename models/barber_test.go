package models

import "testing"

func TestProviderNounLabelsStylistsAndBarbers(t *testing.T) {
	stylist := Barber{ProviderType: ProviderHairstylist}
	barber := Barber{ProviderType: ProviderBarber}
	legacy := Barber{}
	if got := stylist.ProviderNoun(); got != "stylist" {
		t.Fatalf("stylist noun = %q", got)
	}
	if got := barber.ProviderNoun(); got != "barber" {
		t.Fatalf("barber noun = %q", got)
	}
	if got := legacy.ProviderNoun(); got != "barber" {
		t.Fatalf("legacy noun = %q", got)
	}
}
