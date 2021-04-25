package physics_test

import (
	"testing"

	"github.com/jtbonhomme/asteboids/internal/physics"
)

func TestIntersect(t *testing.T) {
	type TestCase struct {
		name      string
		body      *physics.Body
		candidate *physics.Body
		intersect bool
	}

	tests := []TestCase{
		{
			name:      "intersection",
			body:      physics.NewBody(100, 100, 50, 50),
			candidate: physics.NewBody(130, 130, 60, 60),
			intersect: true,
		},
		{
			name:      "no intersection",
			body:      physics.NewBody(100, 100, 50, 50),
			candidate: physics.NewBody(100, 150, 50, 40),
			intersect: false,
		},
		{
			name:      "no intersection",
			body:      physics.NewBody(100, 100, 50, 50),
			candidate: physics.NewBody(62.5, 87.5, 25, 25),
			intersect: false,
		},
		{
			name:      "no intersection",
			body:      physics.NewBody(369, 156, 50, 50),
			candidate: physics.NewBody(404, 56, 100, 100),
			intersect: false,
		},
		{
			name:      "no intersection",
			body:      physics.NewBody(416, 417, 50, 50),
			candidate: physics.NewBody(316, 334, 100, 100),
			intersect: false,
		},
	}
	for _, tt := range tests {
		tt := tt // NOTE: https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // marks each test case as capable of running in parallel with each other
			t.Log(tt.name)
			res := tt.body.Intersect(tt.candidate)
			if res != tt.intersect {
				t.Errorf("test %s expected %t got %t", tt.name, tt.intersect, res)
			}
		})
	}
}
