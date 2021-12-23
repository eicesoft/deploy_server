package gen

import "testing"

func TestCapitalize(t *testing.T) {
	s := capitalize("Role")
	t.Logf("capitalize: %s", s)
}
