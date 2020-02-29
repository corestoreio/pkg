package dmlgen

import (
	"bytes"
	"crypto/md5"
	"fmt"
)

type tables []*Table

// hasFeature returns false when any of the tables does not have the feature/s.
func (ts tables) hasFeature(g *Generator, feature FeatureToggle) bool {
	for _, tbl := range ts {
		if g.hasFeature(tbl.featuresInclude, tbl.featuresExclude, feature) {
			return true
		}
	}
	return false
}

func (ts tables) names() []string {
	names := make([]string, len(ts))
	for i, tbl := range ts {
		names[i] = tbl.Table.Name
	}
	return names
}

// nameID returns a consistent md5 hash of the table names.
func (ts tables) nameID() string {
	var buf bytes.Buffer
	for _, tbl := range ts {
		buf.WriteString(tbl.Table.Name)
	}
	return fmt.Sprintf("%x", md5.Sum(buf.Bytes()))
}
