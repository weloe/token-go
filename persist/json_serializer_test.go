package persist

import (
	"testing"
)

func TestJsonSerializer_Serialize(t *testing.T) {
	serializer := NewJsonSerializer()

	data := struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Email string `json:"email"`
	}{
		Name:  "John Doe",
		Age:   30,
		Email: "john.doe@example.com",
	}

	serializedData, err := serializer.Serialize(data)
	if err != nil {
		t.Fatalf("Failed to serialize data: %v", err)
	}

	expectedData := `{"name":"John Doe","age":30,"email":"john.doe@example.com"}`
	if string(serializedData) != expectedData {
		t.Errorf("Serialized data does not match expected value. Got: %s, Want: %s", serializedData, expectedData)
	}
}

func TestJsonSerializer_UnSerialize(t *testing.T) {
	serializer := NewJsonSerializer()

	serializedData := []byte(`{"name":"John Doe","age":30,"email":"john.doe@example.com"}`)

	var result struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Email string `json:"email"`
	}

	err := serializer.UnSerialize(serializedData, &result)
	if err != nil {
		t.Fatalf("Failed to unserialize data: %v", err)
	}

	expectedResult := struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Email string `json:"email"`
	}{
		Name:  "John Doe",
		Age:   30,
		Email: "john.doe@example.com",
	}
	if result != expectedResult {
		t.Errorf("Unserialized data does not match expected value. Got: %+v, Want: %+v", result, expectedResult)
	}
}
