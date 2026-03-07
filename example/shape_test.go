package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShapeUnion(t *testing.T) {
	t.Run("constructors create correct variants", func(t *testing.T) {
		circle := NewShapeUnion_circle(3.14)
		assert.True(t, circle.Is_circle())
		assert.False(t, circle.Is_rectangle())
		assert.False(t, circle.Is_triangle())
		assert.False(t, circle.Is_Invalid())

		rect := NewShapeUnion_rectangle([2]float64{3.0, 4.0})
		assert.True(t, rect.Is_rectangle())
		assert.False(t, rect.Is_circle())

		tri := NewShapeUnion_triangle([3]float64{3.0, 4.0, 5.0})
		assert.True(t, tri.Is_triangle())

		invalid := NewShapeUnion_Invalid()
		assert.True(t, invalid.Is_Invalid())
	})

	t.Run("unwrap returns value for correct variant", func(t *testing.T) {
		circle := NewShapeUnion_circle(3.14)
		assert.Equal(t, 3.14, circle.Unwrap_circle())

		rect := NewShapeUnion_rectangle([2]float64{3.0, 4.0})
		assert.Equal(t, [2]float64{3.0, 4.0}, rect.Unwrap_rectangle())
	})

	t.Run("unwrap panics for wrong variant", func(t *testing.T) {
		circle := NewShapeUnion_circle(3.14)
		assert.Panics(t, func() {
			circle.Unwrap_rectangle()
		})
	})

	t.Run("get returns value and true for correct variant", func(t *testing.T) {
		circle := NewShapeUnion_circle(3.14)
		val, ok := circle.Get_circle()
		assert.True(t, ok)
		assert.Equal(t, 3.14, val)
	})

	t.Run("get returns zero and false for wrong variant", func(t *testing.T) {
		circle := NewShapeUnion_circle(3.14)
		val, ok := circle.Get_rectangle()
		assert.False(t, ok)
		assert.Equal(t, [2]float64{}, val)
	})

	t.Run("match is exhaustive", func(t *testing.T) {
		circle := NewShapeUnion_circle(3.14)
		result := Match_ShapeUnion(
			&circle,
			func(radius float64) string { return "circle" },
			func(dims [2]float64) string { return "rectangle" },
			func(sides [3]float64) string { return "triangle" },
			func() string { return "invalid" },
		)
		assert.Equal(t, "circle", result)
	})

	t.Run("zero value is invalid variant", func(t *testing.T) {
		var s ShapeUnion
		assert.True(t, s.Is_Invalid())
		assert.False(t, s.Is_circle())
	})

	t.Run("stringer works", func(t *testing.T) {
		assert.Equal(t, "Invalid", _shapeVariant_Invalid.String())
		assert.Equal(t, "circle", _shapeVariant_circle.String())
		assert.Equal(t, "rectangle", _shapeVariant_rectangle.String())
		assert.Equal(t, "triangle", _shapeVariant_triangle.String())
	})
}
