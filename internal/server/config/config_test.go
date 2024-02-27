package config

import (
	"reflect"
	"testing"
)

func TestConfig_LogDebug(t *testing.T) {
	type fields struct {
		Server        Server
		Log           Log
		BackupStorage BackupStorage
		Database      Database
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Server:        tt.fields.Server,
				Log:           tt.fields.Log,
				BackupStorage: tt.fields.BackupStorage,
				Database:      tt.fields.Database,
			}
			cfg.LogDebug()
		})
	}
}

func TestDatabase_Use(t *testing.T) {
	type fields struct {
		ConnStr        string
		RetryIntervals []time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Database{
				ConnStr:        tt.fields.ConnStr,
				RetryIntervals: tt.fields.RetryIntervals,
			}
			if got := d.Use(); got != tt.want {
				t.Errorf("Use() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
