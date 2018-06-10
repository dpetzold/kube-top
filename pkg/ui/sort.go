package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/dpetzold/kube-top/pkg/kube"
	"k8s.io/apimachinery/pkg/api/resource"
)

func compareFields(f1, f2 interface{}, reverse bool, field string) bool {

	if t1, ok := f1.(time.Time); ok {
		t2 := f2.(time.Time)
		if reverse {
			return t1.Before(t2)
		}
		return t1.After(t2)
	}

	if q1, ok := f1.(*resource.Quantity); ok {
		q2 := f2.(*resource.Quantity)
		if reverse {
			return q1.Cmp(*q2) < 0
		}
		return q1.Cmp(*q2) > 0
	}

	if c1, ok := f1.(*kube.CpuResource); ok {
		c2 := f2.(*kube.CpuResource)
		v := c2.Quantity
		if reverse {
			return c1.Cmp(*v) < 0
		}
		return c1.Cmp(*v) > 0
	}

	if m1, ok := f1.(*kube.MemoryResource); ok {
		m2 := f2.(*kube.MemoryResource)
		v := m2.Quantity
		if reverse {
			return m1.Cmp(*v) < 0
		}
		return m1.Cmp(*v) > 0
	}

	if v1, ok := f1.(int64); ok {
		v2 := f2.(int64)
		if reverse {
			return v1 < v2
		}
		return v1 > v2

	}

	if s1, ok := f1.(string); ok {
		s2 := f2.(string)
		if reverse {
			return strings.Compare(s1, s2) > 0
		}
		return strings.Compare(s1, s2) < 0
	}

	panic(fmt.Sprintf("Unknown type to compareFields: %s", field))
}

func CompareStruct(t interface{}, field string, i, j int, reverse bool) bool {

	if ra, ok := t.([]*kube.ContainerResources); ok {
		return compareFields(getField(ra[i], field), getField(ra[j], field), reverse, field)
	}

	if ci, ok := t.([]*kube.ContainerInfo); ok {
		return compareFields(getField(ci[i], field), getField(ci[j], field), reverse, field)
	}

	if nr, ok := t.([]*kube.NodeResources); ok {
		return compareFields(getField(nr[i], field), getField(nr[j], field), reverse, field)
	}

	panic("Unknown type to CompareStruct")
}
