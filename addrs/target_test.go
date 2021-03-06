package addrs

import (
	"fmt"
	"testing"
)

func TestTargetContains(t *testing.T) {
	for _, test := range []struct {
		addr, other Targetable
		expect      bool
	}{
		{
			mustParseTarget("module.foo"),
			mustParseTarget("module.bar"),
			false,
		},
		{
			mustParseTarget("module.foo"),
			mustParseTarget("module.foo"),
			true,
		},
		{
			// module.foo is an unkeyed module instance here, so it cannot
			// contain another instance
			mustParseTarget("module.foo"),
			mustParseTarget("module.foo[0]"),
			false,
		},
		{
			RootModuleInstance,
			mustParseTarget("module.foo"),
			true,
		},
		{
			mustParseTarget("module.foo"),
			RootModuleInstance,
			false,
		},
		{
			mustParseTarget("module.foo"),
			mustParseTarget("module.foo.module.bar[0]"),
			true,
		},
		{
			mustParseTarget("module.foo"),
			mustParseTarget("module.foo.module.bar[0]"),
			true,
		},
		{
			mustParseTarget("module.foo[2]"),
			mustParseTarget("module.foo[2].module.bar[0]"),
			true,
		},
		{
			mustParseTarget("module.foo"),
			mustParseTarget("module.foo.test_resource.bar"),
			true,
		},
		{
			mustParseTarget("module.foo"),
			mustParseTarget("module.foo.test_resource.bar[0]"),
			true,
		},

		// Resources
		{
			mustParseTarget("test_resource.foo"),
			mustParseTarget("test_resource.foo[\"bar\"]"),
			true,
		},
		{
			mustParseTarget(`test_resource.foo["bar"]`),
			mustParseTarget(`test_resource.foo["bar"]`),
			true,
		},
		{
			mustParseTarget("test_resource.foo"),
			mustParseTarget("test_resource.foo[2]"),
			true,
		},
		{
			mustParseTarget("test_resource.foo"),
			mustParseTarget("module.bar.test_resource.foo[2]"),
			false,
		},
		{
			mustParseTarget("module.bar.test_resource.foo"),
			mustParseTarget("module.bar.test_resource.foo[2]"),
			true,
		},
		{
			mustParseTarget("module.bar.test_resource.foo"),
			mustParseTarget("module.bar[0].test_resource.foo[2]"),
			false,
		},

		// Config paths, while never returned from parsing a target, must still be targetable
		{
			ConfigResource{
				Module: []string{"bar"},
				Resource: Resource{
					Mode: ManagedResourceMode,
					Type: "test_resource",
					Name: "foo",
				},
			},
			mustParseTarget("module.bar.test_resource.foo[2]"),
			true,
		},
		{
			ConfigResource{
				Resource: Resource{
					Mode: ManagedResourceMode,
					Type: "test_resource",
					Name: "foo",
				},
			},
			mustParseTarget("module.bar.test_resource.foo[2]"),
			false,
		},
		{
			ConfigResource{
				Module: []string{"bar"},
				Resource: Resource{
					Mode: ManagedResourceMode,
					Type: "test_resource",
					Name: "foo",
				},
			},
			mustParseTarget("module.bar[0].test_resource.foo"),
			true,
		},
	} {
		t.Run(fmt.Sprintf("%s-in-%s", test.other, test.addr), func(t *testing.T) {
			got := test.addr.TargetContains(test.other)
			if got != test.expect {
				t.Fatalf("expected %q.TargetContains(%q) == %t", test.addr, test.other, test.expect)
			}
		})
	}
}

func TestResourceContains(t *testing.T) {
	for _, test := range []struct {
		in, other Targetable
		expect    bool
	}{} {
		t.Run(fmt.Sprintf("%s-in-%s", test.other, test.in), func(t *testing.T) {
			got := test.in.TargetContains(test.other)
			if got != test.expect {
				t.Fatalf("expected %q.TargetContains(%q) == %t", test.in, test.other, test.expect)
			}
		})
	}
}

func mustParseTarget(str string) Targetable {
	t, diags := ParseTargetStr(str)
	if diags != nil {
		panic(fmt.Sprintf("%s: %s", str, diags.ErrWithWarnings()))
	}
	return t.Subject
}
