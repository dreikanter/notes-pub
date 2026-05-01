package deploy

import (
	"os"
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
	repo, err := RepoDir("https://github.com/user/site.git")
	require.NoError(t, err)
	assert.Equal(t, filepath.Dir(build), filepath.Dir(repo), "build and repo share a parent")
	assert.Equal(t, "build", filepath.Base(build))
	assert.Equal(t, "repo", filepath.Base(repo))
	assert.Equal(t, "site", filepath.Base(filepath.Dir(build)))
}

func TestSyncMirrorsBuildIntoRepoPreservingGit(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "build")
	dst := filepath.Join(root, "repo")
	require.NoError(t, os.MkdirAll(filepath.Join(src, "sub"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(src, "index.html"), []byte("new"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(src, "sub", "page.html"), []byte("p"), 0o644))

	require.NoError(t, os.MkdirAll(filepath.Join(dst, ".git", "objects"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dst, ".git", "config"), []byte("git"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dst, "stale.html"), []byte("old"), 0o644))
	require.NoError(t, os.MkdirAll(filepath.Join(dst, "stale-dir"), 0o755))

	require.NoError(t, Sync(src, dst))

	// .git survives.
	_, err := os.Stat(filepath.Join(dst, ".git", "config"))
	assert.NoError(t, err, ".git/config should be preserved")

	// New files copied.
	got, err := os.ReadFile(filepath.Join(dst, "index.html"))
	require.NoError(t, err)
	assert.Equal(t, "new", string(got))
	_, err = os.Stat(filepath.Join(dst, "sub", "page.html"))
	assert.NoError(t, err)

	// Stale entries gone.
	_, err = os.Stat(filepath.Join(dst, "stale.html"))
	assert.True(t, os.IsNotExist(err), "stale file should be removed")
	_, err = os.Stat(filepath.Join(dst, "stale-dir"))
	assert.True(t, os.IsNotExist(err), "stale dir should be removed")
}

func TestSyncErrorsWhenSourceMissing(t *testing.T) {
	root := t.TempDir()
	dst := filepath.Join(root, "repo")
	require.NoError(t, os.MkdirAll(dst, 0o755))
	err := Sync(filepath.Join(root, "missing"), dst)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "run `npub build` first")
}
