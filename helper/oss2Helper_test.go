package helper

import (
	"testing"
	"time"
)

func TestGetOss2SavedKey(t *testing.T) {
	date := time.Now().Format("20060102")
	uuid := "foo"
	actual := GetOss2SavedKey(uuid)
	expected := "image/" + date + "/" + uuid
	if actual != expected {
		t.Errorf("GetOss2SavedKey(\"%s\") expected %s, actual %s", uuid, expected, actual)
	}
}
