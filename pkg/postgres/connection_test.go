package postgres

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestDBValues(t *testing.T) {
	testCases := []struct {
		name   string
		config Config
		want   map[string]string
	}{
		{
			name: "empty configs",
			want: make(map[string]string),
		},
		{
			name: "some config",
			config: Config{
				Name:              "myDatabase",
				User:              "superuser",
				Password:          "notAG00DP@ssword",
				Port:              "1234",
				ConnectionTimeout: 5,
				PoolHealthCheck:   5 * time.Minute,
			},
			want: map[string]string{
				"dbname":                   "myDatabase",
				"password":                 "notAG00DP@ssword",
				"port":                     "1234",
				"user":                     "superuser",
				"connect_timeout":          "5",
				"pool_health_check_period": "5m0s",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := dbValues(&tc.config)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}
