package storage

import (
	"testing"
)

func TestCASPathTransform(t *testing.T) {
	key := "test"
	pathKey := CASPathTransform(key)
	expectedPath := "a94a8/fe5cc/b19ba/61c4c/0873d/391e9/87982/fbbd3"
	expectedFileName := "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3"
	if pathKey.PathName != expectedPath {
		t.Errorf("Expected %s, got %s", expectedPath, pathKey.PathName)
	}
	if pathKey.FileName != expectedFileName {
		t.Errorf("Expected %s, got %s", expectedFileName, pathKey.FileName)
	}
}
