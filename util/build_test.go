package util

import "testing"

func TestGenerateUUID(t *testing.T) {
	uuid, err := GenerateUUID()
	if err != nil {
		t.Errorf("Error generating UUID: %v", err)
	}
	if len(uuid) != 36 {
		t.Errorf("Incorrect UUID length: expected 36, got %d", len(uuid))
	}
	t.Logf("GenerateUUID() = %s", uuid)
}

func TestGenerateSimpleUUID(t *testing.T) {
	uuid, err := GenerateSimpleUUID()
	if err != nil {
		t.Errorf("Error generating simple UUID: %v", err)
	}
	if len(uuid) != 32 {
		t.Errorf("Incorrect simple UUID length: expected 32, got %d", len(uuid))
	}
}

func TestGenerateRandomString32(t *testing.T) {
	randomString, err := GenerateRandomString32()
	if err != nil {
		t.Errorf("Error generating random string: %v", err)
	}
	if len(randomString) != 32 {
		t.Errorf("Incorrect random string length: expected 32, got %d", len(randomString))
	}
	t.Logf("GenerateRandomString32() = %s", randomString)

}

func TestGenerateRandomString64(t *testing.T) {
	randomString, err := GenerateRandomString64()
	if err != nil {
		t.Errorf("Error generating random string: %v", err)
	}
	if len(randomString) != 64 {
		t.Errorf("Incorrect random string length: expected 64, got %d", len(randomString))
	}
	t.Logf("GenerateRandomString64() = %s", randomString)

}

func TestGenerateRandomString128(t *testing.T) {
	randomString, err := GenerateRandomString128()
	if err != nil {
		t.Errorf("Error generating random string: %v", err)
	}
	if len(randomString) != 128 {
		t.Errorf("Incorrect random string length: expected 128, got %d", len(randomString))
	}
	t.Logf("GenerateRandomString128() = %s", randomString)
}
