package imported

import (
	"github.com/sidkurella/gunion/internal/loader/testdata/imported/inner"
)

type myUnion struct {
	a int
	b inner.MyValue
}
