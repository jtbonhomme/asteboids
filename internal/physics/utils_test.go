package physics_test

import (
	"testing"

	"github.com/jtbonhomme/asteboids/internal/physics"
	"github.com/jtbonhomme/asteboids/internal/vector"
)

type testBody struct {
	physics.Body
}

func newTestBody(position, velocity vector.Vector2D) *physics.Body {
	b := &physics.Body{}
	b.Init(velocity)
	b.Move(position)
	return b
}

func TestFuturePosition(t *testing.T) {
	type TestCase struct {
		name string
		t    float64
		body *physics.Body
		want vector.Vector2D
	}

	tests := []TestCase{
		{
			name: "test1",
			t:    1,
			body: newTestBody(
				vector.Vector2D{X: 0, Y: 0},
				vector.Vector2D{X: 1, Y: 1},
			),
			want: vector.Vector2D{X: 1, Y: 1},
		},
		{
			name: "test2",
			t:    1,
			body: newTestBody(
				vector.Vector2D{X: 1, Y: 1},
				vector.Vector2D{X: 1, Y: 1},
			),
			want: vector.Vector2D{X: 2, Y: 2},
		},
		{
			name: "test3",
			t:    2,
			body: newTestBody(
				vector.Vector2D{X: 1, Y: 1},
				vector.Vector2D{X: 1, Y: 1},
			),
			want: vector.Vector2D{X: 3, Y: 3},
		},
	}
	for _, tt := range tests {
		tt := tt // NOTE: https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // marks each test case as capable of running in parallel with each other
			t.Log(tt.name)
			got := tt.body.FuturePosition(tt.t)
			if !got.IsEqual(tt.want) {
				t.Errorf("test %s expected %+v got %+v", tt.name, tt.want, got)
			}
		})
	}

}
