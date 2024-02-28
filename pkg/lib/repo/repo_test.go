package repo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"oras.land/oras-go/v2/registry"
)

func TestParseReference(t *testing.T) {
	tests := []struct {
		input        string
		expectedRef  *registry.Reference
		expectedTags []string
		expectErr    bool
	}{
		{
			input:     "",
			expectErr: true,
		},
		{
			input:        "testregistry.io/test-organization/test-repository:test-tag",
			expectedRef:  reference("testregistry.io", "test-organization/test-repository", "test-tag"),
			expectedTags: []string{},
		},
		{
			input:        "testregistry.io/test-organization/test-repository:test-tag,extraTag1,extraTag2",
			expectedRef:  reference("testregistry.io", "test-organization/test-repository", "test-tag"),
			expectedTags: []string{"extraTag1", "extraTag2"},
		},
		{
			input:        "test-repository:test-tag,extraTag1,extraTag2",
			expectedRef:  reference("localhost", "test-repository", "test-tag"),
			expectedTags: []string{"extraTag1", "extraTag2"},
		},
		{
			input:        "localhost:5000/test-organization/test-repository:test-tag,extraTag1,extraTag2",
			expectedRef:  reference("localhost:5000", "test-organization/test-repository", "test-tag"),
			expectedTags: []string{"extraTag1", "extraTag2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			actualRef, actualTags, actualErr := ParseReference(tt.input)
			if tt.expectErr {
				assert.Error(t, actualErr)
				assert.Nil(t, actualRef)
				assert.Nil(t, actualTags)
			} else {
				if !assert.NoError(t, actualErr) {
					return
				}
				assert.Equal(t, tt.expectedRef, actualRef)
				assert.Equal(t, tt.expectedTags, actualTags)
			}
		})
	}
}

func reference(reg, repo, ref string) *registry.Reference {
	return &registry.Reference{
		Registry:   reg,
		Repository: repo,
		Reference:  ref,
	}
}
