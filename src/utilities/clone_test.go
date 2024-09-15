package utilities_test

import (
	"github.com/stretchr/testify/assert"
	"mqtt-http-bridge/src/utilities"
	"testing"
)

func TestDeepCopy(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		original := 42
		clone, err := utilities.DeepCopy(original)

		if err != nil {
			t.Fatalf("Failed to copy: %s", err)
		}

		if clone != original {
			t.Fatalf("Expected %d, got %d", original, clone)
		}
	})

	t.Run("string", func(t *testing.T) {
		original := "hello"
		clone, err := utilities.DeepCopy(original)

		if err != nil {
			t.Fatalf("Failed to copy: %s", err)
		}

		if clone != original {
			t.Fatalf("Expected %s, got %s", original, clone)
		}
	})

	t.Run("struct", func(t *testing.T) {
		type testStruct struct {
			Value int
		}

		original := testStruct{Value: 42}
		clone, err := utilities.DeepCopy(original)

		if err != nil {
			t.Fatalf("Failed to copy: %s", err)
		}

		if clone != original {
			t.Fatalf("Expected %v, got %v", original, clone)
		}
	})

	t.Run("slice", func(t *testing.T) {
		original := []int{1, 2, 3}
		clone, err := utilities.DeepCopy(original)

		if err != nil {
			t.Fatalf("Failed to copy: %s", err)
		}

		if len(clone) != len(original) {
			t.Fatalf("Expected %v, got %v", original, clone)
		}

		for i := range original {
			if clone[i] != original[i] {
				t.Fatalf("Expected %v, got %v", original, clone)
			}
		}
	})

	t.Run("map", func(t *testing.T) {
		original := map[string]int{
			"one":   1,
			"two":   2,
			"three": 3,
		}
		clone, err := utilities.DeepCopy(original)

		if err != nil {
			t.Fatalf("Failed to copy: %s", err)
		}

		if len(clone) != len(original) {
			t.Fatalf("Expected %v, got %v", original, clone)
		}

		for k, v := range original {
			if clone[k] != v {
				t.Fatalf("Expected %v, got %v", original, clone)
			}
		}
	})

	t.Run("nil", func(t *testing.T) {
		var original *int = nil
		clone, err := utilities.DeepCopy(original)

		if err != nil {
			t.Fatalf("Failed to copy: %s", err)
		}

		if clone != nil {
			t.Fatalf("Expected nil, got %v", clone)
		}
	})

	t.Run("struct with map property", func(t *testing.T) {
		type testStruct struct {
			Value map[string]int
		}

		original := testStruct{Value: map[string]int{
			"one": 1,
			"two": 2,
		}}
		clone, err := utilities.DeepCopy(original)

		if err != nil {
			t.Fatalf("Failed to copy: %s", err)
		}

		if len(clone.Value) != len(original.Value) {
			t.Fatalf("Expected %v, got %v", original, clone)
		}

		for k, v := range original.Value {
			if clone.Value[k] != v {
				t.Fatalf("Expected %v, got %v", original, clone)
			}
		}

		original.Value["two"] = 3 // Overwrite original

		assert.Equal(t, 2, clone.Value["two"]) // Clone should not be affected
	})
}
