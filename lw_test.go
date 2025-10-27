package lw_api

import (
	"reflect"
	"testing"
)

func TestNewLinnworksAPIBuilder(t *testing.T) {
	builder := NewLinnworksAPIBuilder()
	if builder == nil {
		t.Error("Expected builder to not be nil")
	}
}

func TestLinnworksAPIBuilder_Build(t *testing.T) {
	t.Run("should return an error if baseURL is not provided", func(t *testing.T) {
		builder := NewLinnworksAPIBuilder().
			Token("test-token").
			AppID("test-app-id").
			AppSecret("test-app-secret")

		_, err := builder.Build()
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("should return an error if token is not provided", func(t *testing.T) {
		builder := NewLinnworksAPIBuilder().
			BaseURL("https://example.com").
			AppID("test-app-id").
			AppSecret("test-app-secret")

		_, err := builder.Build()
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("should return an error if appID is not provided", func(t *testing.T) {
		builder := NewLinnworksAPIBuilder().
			BaseURL("https://example.com").
			Token("test-token").
			AppSecret("test-app-secret")

		_, err := builder.Build()
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("should return an error if appSecret is not provided", func(t *testing.T) {
		builder := NewLinnworksAPIBuilder().
			BaseURL("https://example.com").
			Token("test-token").
			AppID("test-app-id")

		_, err := builder.Build()
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("should return a LinnworksAPI instance on success", func(t *testing.T) {
		builder := NewLinnworksAPIBuilder().
			BaseURL("https://eu-ext.linnworks.net").
			Token("test-token").
			AppID("test-app-id").
			AppSecret("test-app-secret")

		api, err := builder.Build()
		if err != nil {
			t.Errorf("Expected nil error, got %v", err)
		}
		if api == nil {
			t.Error("Expected api to not be nil")
		}
	})

	t.Run("all api branches have to be set", func(t *testing.T) {
		builder := NewLinnworksAPIBuilder().
			BaseURL("https://eu-ext.linnworks.net").
			Token("test-token").
			AppID("test-app-id").
			AppSecret("test-app-secret")

		api, err := builder.Build()
		if err != nil {
			t.Fatalf("Expected nil error, got %v", err)
		}

		val := reflect.ValueOf(api).Elem()
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			if field.IsNil() {
				t.Errorf("Expected %s to not be nil", val.Type().Field(i).Name)
			}
		}
	})
}
