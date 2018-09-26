package gitlab

import (
	"fmt"
	"net/http"
	"os"

	gitlab "github.com/xanzy/go-gitlab"
)

// Client for Gitlab REST API
type Client struct {
	*gitlab.Client
}

// NewClientWithGitlabPrivateToken returns a new Client with a Gitlab's private token
func NewClientWithGitlabPrivateToken(client *http.Client, gitlabDomain string, privateToken string) *Client {
	gl := gitlab.NewClient(client, privateToken)
	gl.SetBaseURL(fmt.Sprintf("https://%s/api/v4", gitlabDomain))
	return &Client{gl}
}

func NewClientFromEnv(client *http.Client) *Client {
	return NewClientWithGitlabPrivateToken(client, os.Getenv("GITLAB_DOMAIN"), os.Getenv("GITLAB_PRIVATE_TOKEN"))
}
