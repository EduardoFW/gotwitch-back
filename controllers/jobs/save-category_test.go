package jobs

import (
	"testing"

	"api.gotwitch.tk/models"
	"api.gotwitch.tk/settings"
)

// Test that runs before all tests
func TestMain(m *testing.M) {
	settings.Setup("../../.env")
	m.Run()
}

func TestAppendUniqueCategories(t *testing.T) {
	categoriesA := []models.Category{
		{
			Id:   "1",
			Name: "test",
		},
		{
			Id:   "2",
			Name: "test2",
		},
	}

	categoriesB := []models.Category{
		{
			Id:   "2",
			Name: "test2",
		},
		{
			Id:   "3",
			Name: "test3",
		},
		{
			Id:   "4",
			Name: "test4",
		},
	}

	categories := appendUniqueCategories(categoriesA, categoriesB)

	if len(categories) != 4 {
		t.Errorf("Expected 4 categories, got %d", len(categories))
	}

	if categories[0].Id != "1" {
		t.Errorf("Expected first category to be 1, got %s", categories[0].Id)
	}

	if categories[3].Id != "4" {
		t.Errorf("Expected last category to be 4, got %s", categories[3].Id)
	}
}

func TestContainsCategory(t *testing.T) {
	categories := []models.Category{
		{
			Id:   "1",
			Name: "test",
		},
		{
			Id:   "2",
			Name: "test2",
		},
	}

	if !containsCategory(categories, categories[0]) {
		t.Errorf("Expected to find category %s in categories", categories[0].Id)
	}

	if containsCategory(categories, models.Category{
		Id:   "3",
		Name: "test3",
	}) {
		t.Errorf("Expected not to find category 3 in categories")
	}
}

func TestGenerateAlphabet(t *testing.T) {
	alphabet := generateAlphabet()
	if len(alphabet) != 26 {
		t.Errorf("Expected alphabet to be 26, got %d", len(alphabet))
	}

	if alphabet[0] != "a" {
		t.Errorf("Expected alphabet to start with a, got %s", alphabet[0])
	}

	if alphabet[25] != "z" {
		t.Errorf("Expected alphabet to end with z, got %s", alphabet[25])
	}
}
