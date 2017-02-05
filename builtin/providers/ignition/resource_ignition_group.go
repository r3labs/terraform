package ignition

import (
	"github.com/coreos/ignition/config/types"
	"github.com/r3labs/terraform/helper/schema"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Delete: resourceGroupDelete,
		Exists: resourceGroupExists,
		Read:   resourceGroupRead,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"gid": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"password_hash": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	id, err := buildGroup(d, meta.(*cache))
	if err != nil {
		return err
	}

	d.SetId(id)
	return nil
}

func resourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

func resourceGroupExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id, err := buildGroup(d, meta.(*cache))
	if err != nil {
		return false, err
	}

	return id == d.Id(), nil
}

func resourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func buildGroup(d *schema.ResourceData, c *cache) (string, error) {
	return c.addGroup(&types.Group{
		Name:         d.Get("name").(string),
		PasswordHash: d.Get("password_hash").(string),
		Gid:          getUInt(d, "gid"),
	}), nil
}
