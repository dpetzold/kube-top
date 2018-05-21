package main

import (
	"fmt"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/api/resource"
)

func cmp_fields(f1, f2 interface{}, reverse bool, field string) bool {

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

	if c1, ok := f1.(*CpuResource); ok {
		c2 := f2.(*CpuResource)
		v := c2.Quantity
		if reverse {
			return c1.Cmp(*v) < 0
		}
		return c1.Cmp(*v) > 0
	}

	if m1, ok := f1.(*MemoryResource); ok {
		m2 := f2.(*MemoryResource)
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

	panic(fmt.Sprintf("Unknown type: cmp_fields %s", field))
}

func cmp_struct(t interface{}, field string, i, j int, reverse bool) bool {

	if ra, ok := t.([]*ContainerResources); ok {
		return cmp_fields(getField(ra[i], field), getField(ra[j], field), reverse, field)
	}

	if ci, ok := t.([]*ContainerInfo); ok {
		return cmp_fields(getField(ci[i], field), getField(ci[j], field), reverse, field)
	}

	if nr, ok := t.([]*NodeResources); ok {
		return cmp_fields(getField(nr[i], field), getField(nr[j], field), reverse, field)
	}

	panic("Unknown type: cmp_struct")
}
