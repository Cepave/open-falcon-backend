package plugins

import (
	"testing"
)

func TestGitLsRemote(t *testing.T) {
	hash, _ := GitLsRemote("https://github.com/humorless/lonelyone.git", "refs/heads/master")
	if hash != "e0ab4048b87337a6b6d5b6ab051af2a00569a024" {
		t.Error("commit hash dismatched")
	}
}
