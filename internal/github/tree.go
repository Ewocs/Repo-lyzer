package github

type TreeEntry struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	Type string `json:"type"`
	Size int    `json:"size"`
	Sha  string `json:"sha"`
}

type TreeResponse struct {
	Sha       string      `json:"sha"`
	Url       string      `json:"url"`
	Tree      []TreeEntry `json:"tree"`
	Truncated bool        `json:"truncated"`
}

func (c *Client) GetFileTree(owner, repo, branch string) ([]TreeEntry, error) {
	var t TreeResponse
	// recursive=1 to get full tree
	err := c.get("https://api.github.com/repos/"+owner+"/"+repo+"/git/trees/"+branch+"?recursive=1", &t)
	return t.Tree, err
}
