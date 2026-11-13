// Copyright IBM Corp. 2024, 2026
// SPDX-License-Identifier: BUSL-1.1

package agent

import (
	"bytes"
	"testing"

	"github.com/hashicorp/consul/sdk/testutil"
	"github.com/stretchr/testify/require"
)

func TestAgentRetryNewDiscover(t *testing.T) {
	d, err := newDiscover()
	require.NoError(t, err)
	// This build does not contain any go-discover provider.
	require.Empty(t, d.Names())
}

func TestAgentRetryJoinAddrs(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	d, err := newDiscover()
	require.NoError(t, err)

	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{"handles nil", nil, []string{}},
		{"handles empty input", []string{}, []string{}},
		{"handles one element",
			[]string{"192.168.0.12"},
			[]string{"192.168.0.12"},
		},
		{"handles two elements",
			[]string{"192.168.0.12", "192.168.0.13"},
			[]string{"192.168.0.12", "192.168.0.13"},
		},
		{"rejects a go-discover expression, because this build does not contain go-discover",
			[]string{"192.168.0.12", "provider=aws region=eu-west-1 tag_key=consul tag_value=tag access_key_id=a secret_access_key=a"},
			[]string{"192.168.0.12"},
		},
	}
	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := testutil.LoggerWithOutput(t, &buf)

			output := retryJoinAddrs(d, retryJoinSerfVariant, "LAN", test.input, logger)
			bufout := buf.String()
			require.Equal(t, test.expected, output, bufout)
			if i == 4 {
				require.Contains(t, bufout, "Cannot discover address")
			}
		})
	}
	t.Run("handles nil discover", func(t *testing.T) {
		require.Equal(t, []string{}, retryJoinAddrs(nil, retryJoinSerfVariant, "LAN", []string{"a"}, nil))
	})
}
