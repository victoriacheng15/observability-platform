package main

import (
	"os"
	"testing"
)

func TestLoadData(t *testing.T) {
	// Create a temporary data file
	tmpContent := `
landing:
  page_title: "Test Landing"
  hero:
    title: "Test Hero"
    subtitle: "Test Subtitle"
    cta_text: "Click Me"
    cta_link: "/test.html"
    cta2_text: "Click Me 2"
    cta2_link: "/test2.html"
  features:
    - title: "Feat 1"
      description: "Desc 1"
      icon: "rocket"

evolution:
  page_title: "Test Evolution"
  intro_text: "Test Intro"
  events:
    - date: "2024-01-01"
      title: "Event 1"
      description: "Desc 1"

dashboards:
  page_title: "Test Dashboards"
  categories:
    - name: "Cat 1"
      description: "Desc 1"
      items:
        - caption: "Cap 1"
          image_path: "/img.png"
`
	tmpFile, err := os.CreateTemp("", "data-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(tmpContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	data, err := loadData(tmpFile.Name())
	if err != nil {
		t.Fatalf("loadData failed: %v", err)
	}

	if data.Landing.PageTitle != "Test Landing" {
		t.Errorf("Expected 'Test Landing', got '%s'", data.Landing.PageTitle)
	}
	if len(data.Landing.Features) != 1 {
		t.Errorf("Expected 1 feature, got %d", len(data.Landing.Features))
	}
	if data.Evolution.PageTitle != "Test Evolution" {
		t.Errorf("Expected 'Test Evolution', got '%s'", data.Evolution.PageTitle)
	}
	if len(data.Dashboards.Categories) != 1 {
		t.Errorf("Expected 1 category, got %d", len(data.Dashboards.Categories))
	}
}
