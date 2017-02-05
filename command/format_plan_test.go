package command

import (
	"strings"
	"testing"

	"github.com/r3labs/terraform/terraform"
	"github.com/mitchellh/colorstring"
)

// Test that a root level data source gets a special plan output on create
func TestFormatPlan_destroyDeposed(t *testing.T) {
	plan := &terraform.Plan{
		Diff: &terraform.Diff{
			Modules: []*terraform.ModuleDiff{
				&terraform.ModuleDiff{
					Path: []string{"root"},
					Resources: map[string]*terraform.InstanceDiff{
						"aws_instance.foo": &terraform.InstanceDiff{
							DestroyDeposed: true,
						},
					},
				},
			},
		},
	}
	opts := &FormatPlanOpts{
		Plan: plan,
		Color: &colorstring.Colorize{
			Colors:  colorstring.DefaultColors,
			Disable: true,
		},
		ModuleDepth: 1,
	}

	actual := FormatPlan(opts)

	expected := strings.TrimSpace(`
- aws_instance.foo (deposed)
	`)
	if actual != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, actual)
	}
}

// Test that computed fields with an interpolation string get displayed
func TestFormatPlan_displayInterpolations(t *testing.T) {
	plan := &terraform.Plan{
		Diff: &terraform.Diff{
			Modules: []*terraform.ModuleDiff{
				&terraform.ModuleDiff{
					Path: []string{"root"},
					Resources: map[string]*terraform.InstanceDiff{
						"aws_instance.foo": &terraform.InstanceDiff{
							Attributes: map[string]*terraform.ResourceAttrDiff{
								"computed_field": &terraform.ResourceAttrDiff{
									New:         "${aws_instance.other.id}",
									NewComputed: true,
								},
							},
						},
					},
				},
			},
		},
	}
	opts := &FormatPlanOpts{
		Plan: plan,
		Color: &colorstring.Colorize{
			Colors:  colorstring.DefaultColors,
			Disable: true,
		},
		ModuleDepth: 1,
	}

	out := FormatPlan(opts)
	lines := strings.Split(out, "\n")
	if len(lines) != 2 {
		t.Fatal("expected 2 lines of output, got:\n", out)
	}

	actual := strings.TrimSpace(lines[1])
	expected := `computed_field: "" => "${aws_instance.other.id}"`

	if actual != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, actual)
	}
}

// Test that a root level data source gets a special plan output on create
func TestFormatPlan_rootDataSource(t *testing.T) {
	plan := &terraform.Plan{
		Diff: &terraform.Diff{
			Modules: []*terraform.ModuleDiff{
				&terraform.ModuleDiff{
					Path: []string{"root"},
					Resources: map[string]*terraform.InstanceDiff{
						"data.type.name": &terraform.InstanceDiff{
							Attributes: map[string]*terraform.ResourceAttrDiff{
								"A": &terraform.ResourceAttrDiff{
									New:         "B",
									RequiresNew: true,
								},
							},
						},
					},
				},
			},
		},
	}
	opts := &FormatPlanOpts{
		Plan: plan,
		Color: &colorstring.Colorize{
			Colors:  colorstring.DefaultColors,
			Disable: true,
		},
		ModuleDepth: 1,
	}

	actual := FormatPlan(opts)

	expected := strings.TrimSpace(`
 <= data.type.name
    A: "B"
	`)
	if actual != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, actual)
	}
}

// Test that data sources nested in modules get the same plan output
func TestFormatPlan_nestedDataSource(t *testing.T) {
	plan := &terraform.Plan{
		Diff: &terraform.Diff{
			Modules: []*terraform.ModuleDiff{
				&terraform.ModuleDiff{
					Path: []string{"root", "nested"},
					Resources: map[string]*terraform.InstanceDiff{
						"data.type.name": &terraform.InstanceDiff{
							Attributes: map[string]*terraform.ResourceAttrDiff{
								"A": &terraform.ResourceAttrDiff{
									New:         "B",
									RequiresNew: true,
								},
							},
						},
					},
				},
			},
		},
	}
	opts := &FormatPlanOpts{
		Plan: plan,
		Color: &colorstring.Colorize{
			Colors:  colorstring.DefaultColors,
			Disable: true,
		},
		ModuleDepth: 2,
	}

	actual := FormatPlan(opts)

	expected := strings.TrimSpace(`
 <= module.nested.data.type.name
    A: "B"
	`)
	if actual != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, actual)
	}
}
