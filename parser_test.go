package abc

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tune, err := Parse("T:My song\nC2 _EG c2 B^A^^GFED C4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tune.Title != "My song" {
		t.Errorf("wrong title %q", tune.Title)
	}

	exp := []Note{
		{Oct: 3, Name: 'C', Length: 16},
		{Oct: 3, Name: 'E', Length: 8, Acc: -1},
		{Oct: 3, Name: 'G', Length: 8},
		{Oct: 4, Name: 'C', Length: 16},
		{Oct: 3, Name: 'B', Length: 8},
		{Oct: 3, Name: 'A', Length: 8, Acc: 1},
		{Oct: 3, Name: 'G', Length: 8, Acc: 2},
		{Oct: 3, Name: 'F', Length: 8},
		{Oct: 3, Name: 'E', Length: 8},
		{Oct: 3, Name: 'D', Length: 8},
		{Oct: 3, Name: 'C', Length: 32},
	}

	if len(exp) != len(tune.Notes) {
		t.Errorf("expected %v, got %#v", exp, tune.Notes)
	}

	for i, e := range exp {
		if !reflect.DeepEqual(e, tune.Notes[i]) {
			t.Errorf("expected %v, got %v", e, tune.Notes[i])
		}
	}
}
