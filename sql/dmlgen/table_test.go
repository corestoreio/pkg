package dmlgen

import "testing"

func TestFeatureToggle_String(t *testing.T) {
	tests := []struct {
		name string
		f    FeatureToggle
		want string
	}{
		{"one", FeatureEntityDBUpsert, "FeatureEntityDBUpsert"},
		{"two", FeatureEntityDBUpsert | FeatureEntityDBSelect, "FeatureEntityDBSelect,FeatureEntityDBUpsert"},
		{"none", 0, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
