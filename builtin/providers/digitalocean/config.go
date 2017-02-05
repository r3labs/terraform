package digitalocean

import (
	"log"
	"time"

	"github.com/digitalocean/godo"
	"github.com/r3labs/terraform/helper/resource"
	"golang.org/x/oauth2"
)

type Config struct {
	Token string
}

// Client() returns a new client for accessing digital ocean.
func (c *Config) Client() (*godo.Client, error) {
	tokenSrc := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: c.Token,
	})

	client := godo.NewClient(oauth2.NewClient(oauth2.NoContext, tokenSrc))

	log.Printf("[INFO] DigitalOcean Client configured for URL: %s", client.BaseURL.String())

	return client, nil
}

// waitForAction waits for the action to finish using the resource.StateChangeConf.
func waitForAction(client *godo.Client, action *godo.Action) error {
	var (
		pending   = "in-progress"
		target    = "completed"
		refreshfn = func() (result interface{}, state string, err error) {
			a, _, err := client.Actions.Get(action.ID)
			if err != nil {
				return nil, "", err
			}
			if a.Status == "errored" {
				return a, "errored", nil
			}
			if a.CompletedAt != nil {
				return a, target, nil
			}
			return a, pending, nil
		}
	)
	_, err := (&resource.StateChangeConf{
		Pending: []string{pending},
		Refresh: refreshfn,
		Target:  []string{target},

		Delay:      10 * time.Second,
		Timeout:    60 * time.Minute,
		MinTimeout: 3 * time.Second,

		// This is a hack around DO API strangeness.
		// https://github.com/r3labs/terraform/issues/481
		//
		NotFoundChecks: 60,
	}).WaitForState()
	return err
}
