package v2

import "testing"

func Test_GetFIdentifierDescription(t *testing.T) {
	description := GetFIdentifierDescription("TIT2")

	if description != "Title/songname/content description" {
		t.Errorf("GetFIdentifierDescription failed: expected description for TIT2 to be %s, got %s",
			"Title/songname/content description", description)
	}

	description = GetFIdentifierDescription("TBP")
	if description != "BPM (Beats Per Minute)" {
		t.Errorf("GetFIdentifierDescription failed: expected description for TBP to be %s, got %s",
			"BPM (Beats Per Minute)", description)
	}

	description = GetFIdentifierDescription("SomeInvalidFrameIDName")
	if description != "" {
		t.Errorf("GetFIdentifierDescription failed: expected description for SomeInvalidFrameIDName to be \"\", got %s",
			description)
	}
}
