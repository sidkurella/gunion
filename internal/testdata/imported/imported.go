package imported

import (
	"github.com/sidkurella/gunion/internal/testdata/imported/inner"
)

type myUnion struct {
	a int
	b inner.MyValue
}
