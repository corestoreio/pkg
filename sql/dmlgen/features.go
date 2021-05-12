package dmlgen

import "strings"

// FeatureToggle allows certain generated code blocks to be switched off or on.
type FeatureToggle uint64

// List of available features
const (
	FeatureCollectionAppend FeatureToggle = 1 << iota
	FeatureCollectionClear
	FeatureCollectionCut
	FeatureCollectionDelete
	FeatureCollectionEach
	FeatureCollectionFilter
	FeatureCollectionInsert
	FeatureCollectionStruct // creates the struct type
	FeatureCollectionSwap
	FeatureCollectionUniqueGetters
	FeatureCollectionUniquifiedGetters
	FeatureCollectionValidate
	FeatureDB
	FeatureDBAssignLastInsertID
	FeatureDBDelete
	FeatureDBInsert
	FeatureDBMapColumns
	FeatureDBSelect
	FeatureDBTracing // opentelemetry tracing
	FeatureDBUpdate
	FeatureDBUpsert
	FeatureDBTableColumnNames
	FeatureEntityCopy
	FeatureEntityEmpty
	FeatureEntityGetSetPrivateFields
	FeatureEntityIsSet
	FeatureEntityRelationships
	FeatureEntityStruct // creates the struct type
	FeatureEntityValidate
	FeatureEntityWriteTo
	featureMax
)

var featureNames = map[FeatureToggle]string{
	FeatureCollectionAppend:            "FeatureCollectionAppend",
	FeatureCollectionCut:               "FeatureCollectionCut",
	FeatureCollectionDelete:            "FeatureCollectionDelete",
	FeatureCollectionEach:              "FeatureCollectionEach",
	FeatureCollectionFilter:            "FeatureCollectionFilter",
	FeatureCollectionInsert:            "FeatureCollectionInsert",
	FeatureCollectionStruct:            "FeatureCollectionStruct",
	FeatureCollectionSwap:              "FeatureCollectionSwap",
	FeatureCollectionUniqueGetters:     "FeatureCollectionUniqueGetters",
	FeatureCollectionUniquifiedGetters: "FeatureCollectionUniquifiedGetters",
	FeatureCollectionValidate:          "FeatureCollectionValidate",
	FeatureDB:                          "FeatureDB",
	FeatureDBAssignLastInsertID:        "FeatureDBAssignLastInsertID",
	FeatureDBDelete:                    "FeatureDBDelete",
	FeatureDBInsert:                    "FeatureDBInsert",
	FeatureDBMapColumns:                "FeatureDBMapColumns",
	FeatureDBSelect:                    "FeatureDBSelect",
	FeatureDBTracing:                   "FeatureDBTracing",
	FeatureDBUpdate:                    "FeatureDBUpdate",
	FeatureDBUpsert:                    "FeatureDBUpsert",
	FeatureEntityCopy:                  "FeatureEntityCopy",
	FeatureEntityEmpty:                 "FeatureEntityEmpty",
	FeatureEntityGetSetPrivateFields:   "FeatureEntityGetSetPrivateFields",
	FeatureEntityIsSet:                 "FeatureEntityIsSet",
	FeatureEntityRelationships:         "FeatureEntityRelationships",
	FeatureEntityStruct:                "FeatureEntityStruct",
	FeatureEntityValidate:              "FeatureEntityValidate",
	FeatureEntityWriteTo:               "FeatureEntityWriteTo",
}

func (f FeatureToggle) String() string {
	var buf strings.Builder
	j := 0
	for i, k := 0, 0; i <= int(featureMax); i, k = 1<<k, k+1 {
		if int(f)&i > 0 {
			if j > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(featureNames[FeatureToggle(i)])
			j++
		}
	}
	return buf.String()
}

func hasFeature(include, exclude, features FeatureToggle, mode ...rune) int {
	featureFlagCount := 0
	tableIncludeFlagCount := 0
	for i, j := FeatureToggle(0), FeatureToggle(0); j <= featureMax; i++ {
		j = 1 << i
		if j&features != 0 {
			featureFlagCount++
		}
		if j&include != 0 && j&features != 0 {
			tableIncludeFlagCount++
		}
	}
	mode2 := 'o' // or is the default mode
	if len(mode) == 1 && mode[0] > 0 {
		mode2 = mode[0]
	}

	ret := -2
	switch {
	case mode2 != 'a' && mode2 != 'o':
		panic("mode can only be character a (=and) or o (=or)")
	case include == 0 && exclude == 0:
		// everything allowed
		ret = 1
	case mode2 == 'a' && include > 0 && featureFlagCount == tableIncludeFlagCount && exclude == 0:
		// all features must be set in include
		ret = 2
	case mode2 == 'o' && include > 0 && (include&features) != 0 && exclude == 0:
		// at least one feature must be set in include
		ret = 4
	case include == 0 && exclude > 0 && (exclude&features) != 0:
		// exclude all features
		ret = -1
	case include == 0 && exclude > 0 && (exclude&features) == 0:
		// features in exclude flags not found
		ret = 3
	}
	return ret
}
