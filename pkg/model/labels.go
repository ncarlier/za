package model

import (
	"fmt"
	"sort"
	"strings"
)

// Labels used for Loki Labels
type Labels map[string]string

func (l Labels) String() string {
	labels := make([]string, 0, len(l))
	for k, v := range l {
		labels = append(labels, fmt.Sprintf("%s=%q", k, v))
	}
	sort.Strings(labels)
	return fmt.Sprintf("{%s}", strings.Join(labels, ", "))
}
