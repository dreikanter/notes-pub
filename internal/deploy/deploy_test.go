package deploy

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepoSlug(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"https with .git", "https://github.com/user/repo.git", "repo"},
		{"https without .git", "https://github.com/user/repo", "repo"},
		{"ssh with .git", "git@github.com:user/repo.git", "repo"},
		{"trailing slash", "https://github.com/user/repo/", "repo"},
		{"only basename", "repo.git", "repo"},
		{"empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, RepoSlug(tt.in))
		})
	}
}

func TestCacheLayout(t *testing.T) {
	build, err := BuildDir("https://github.com/user/site.git")
	require.NoError(t, err)
	gitDir, err := GitDir("https://github.com/user/site.git")
	require.NoError(t, err)
	assert.Equal(t, filepath.Dir(build), filepath.Dir(gitDir), "build and git share a parent")
	assert.Equal(t, "build", filepath.Base(build))
	assert.Equal(t, "git", filepath.Base(gitDir))
	assert.Equal(t, "site", filepath.Base(filepath.Dir(build)))
}
