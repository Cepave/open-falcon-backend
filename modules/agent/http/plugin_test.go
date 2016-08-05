package http

import "testing"

func TestDeleteAndCloneRepo(t *testing.T) {
	out := deleteAndCloneRepo("./plugin", "https://github.com/humorless/openfalcon-plugin.git")
	t.Log("test deleteAndCloneRepo: ", out)
	out = deleteAndCloneRepo("/", "https://github.com/humorless/openfalcon-plugin.git")
	t.Log("test deleteAndCloneRepo: ", out)
}
