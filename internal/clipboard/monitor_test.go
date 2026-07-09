package clipboard

import "testing"

func TestValidateUTF8RemovesInvalidSequences(t *testing.T) {
	t.Parallel()

	got := validateUTF8("ok\xff中文")
	if got != "ok中文" {
		t.Fatalf("expected invalid UTF-8 bytes to be removed, got %q", got)
	}
}
