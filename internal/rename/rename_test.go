package rename

import (
	"testing"

	"img-meta/internal/meta"
)

func TestPlan_GeneratesNames(t *testing.T) {
	metas := []meta.Meta{
		{Path: "a.png", Format: "png", Width: 1920, Height: 1080},
		{Path: "b.png", Format: "png", Width: 600, Height: 800},
	}
	plans := Plan(metas, "{aspect}-{w}x{h}-{i}")
	if len(plans) != 2 {
		t.Fatalf("expected 2 plans, got %d", len(plans))
	}
	if plans[0].To != "landscape-1920x1080-0.png" {
		t.Fatalf("unexpected name: %q", plans[0].To)
	}
	if plans[1].To != "portrait-600x800-1.png" {
		t.Fatalf("unexpected name: %q", plans[1].To)
	}
}

func TestPlan_Empty(t *testing.T) {
	plans := Plan(nil, "")
	if len(plans) != 0 {
		t.Fatalf("expected empty, got %d", len(plans))
	}
}

func TestPlan_Collision(t *testing.T) {
	metas := []meta.Meta{
		{Path: "a.png", Format: "png", Width: 100, Height: 100},
		{Path: "b.png", Format: "png", Width: 100, Height: 100},
	}
	plans := Plan(metas, "same")
	if plans[0].To != "same.png" {
		t.Fatalf("first should be unchanged: %q", plans[0].To)
	}
	if plans[1].To != "same_1.png" {
		t.Fatalf("collision should be disambiguated: %q", plans[1].To)
	}
}

func TestPlan_ReuseDoesNotAlias(t *testing.T) {
	first := Plan([]meta.Meta{
		{Path: "a.png", Format: "png", Width: 1920, Height: 1080},
	}, "{aspect}-{w}x{h}-{i}")
	if len(first) != 1 || first[0].To != "landscape-1920x1080-0.png" {
		t.Fatalf("first plan: %+v", first)
	}
	second := Plan([]meta.Meta{
		{Path: "b.png", Format: "png", Width: 600, Height: 800},
	}, "{aspect}-{w}x{h}-{i}")
	if len(second) != 1 || second[0].To != "portrait-600x800-0.png" {
		t.Fatalf("second plan: %+v", second)
	}
	if first[0].To != "landscape-1920x1080-0.png" {
		t.Fatalf("first plan changed after second Plan: %q", first[0].To)
	}
}
