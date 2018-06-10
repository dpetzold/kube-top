package kube

import (
	"fmt"

	humanize "github.com/dustin/go-humanize"
	"k8s.io/apimachinery/pkg/api/resource"
)

func calcPercentage(dividend, divisor int64) int {
	return int(float64(dividend) / float64(divisor) * 100)
}

func NewMemoryResource(value int64) *MemoryResource {
	return &MemoryResource{resource.NewQuantity(value, resource.BinarySI)}
}

func (r *MemoryResource) calcPercentage(divisor *resource.Quantity) int {
	return calcPercentage(r.Value(), divisor.Value())
}

func (r *MemoryResource) String() string {
	// XXX: Support more units
	// return fmt.Sprintf("%vMi", r.Value()/(1024*1024))
	return humanize.Bytes(uint64(r.Value()))
}

func NewCpuResource(value int64) *CpuResource {
	r := resource.NewMilliQuantity(value, resource.DecimalSI)
	return &CpuResource{r}
}

func (r *CpuResource) String() string {
	// XXX: Support more units
	return fmt.Sprintf("%vm", r.MilliValue())
}

func (r *CpuResource) calcPercentage(divisor *resource.Quantity) int {
	return calcPercentage(r.MilliValue(), divisor.MilliValue())
}
