package ignition

import (
	"github.com/coreos/ignition/config/types"
	"github.com/r3labs/terraform/helper/schema"
)

func resourceRaid() *schema.Resource {
	return &schema.Resource{
		Create: resourceRaidCreate,
		Delete: resourceRaidDelete,
		Exists: resourceRaidExists,
		Read:   resourceRaidRead,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"level": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"devices": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"spares": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceRaidCreate(d *schema.ResourceData, meta interface{}) error {
	id, err := buildRaid(d, meta.(*cache))
	if err != nil {
		return err
	}

	d.SetId(id)
	return nil
}

func resourceRaidDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

func resourceRaidExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id, err := buildRaid(d, meta.(*cache))
	if err != nil {
		return false, err
	}

	return id == d.Id(), nil
}

func resourceRaidRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func buildRaid(d *schema.ResourceData, c *cache) (string, error) {
	var devices []types.Path
	for _, value := range d.Get("devices").([]interface{}) {
		devices = append(devices, types.Path(value.(string)))
	}

	return c.addRaid(&types.Raid{
		Name:    d.Get("name").(string),
		Level:   d.Get("level").(string),
		Devices: devices,
		Spares:  d.Get("spares").(int),
	}), nil
}
