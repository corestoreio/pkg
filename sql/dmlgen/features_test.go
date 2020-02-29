package dmlgen

import "testing"

func Test_hasFeature(t *testing.T) {
	type args struct {
		tableInclude FeatureToggle
		tableExclude FeatureToggle
		featureReq   FeatureToggle
		mode         rune // a=and or o=or default is or
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "OR include and exclude zero", args: args{
			featureReq: FeatureEntityStruct,
		}, want: 1},
		{name: "AND includes only, one feature ok", args: args{
			tableInclude: FeatureDBUpdate | FeatureDBInsert | FeatureEntityStruct,
			tableExclude: 0,
			featureReq:   FeatureEntityStruct,
			mode:         'a',
		}, want: 2},
		{name: "AND includes only, one feature nok", args: args{
			tableInclude: FeatureDBUpdate | FeatureDBInsert | FeatureEntityStruct,
			tableExclude: 0,
			featureReq:   FeatureCollectionEach,
			mode:         'a',
		}, want: -2},
		{name: "AND includes only, two feature ok", args: args{
			tableInclude: FeatureDBUpdate | FeatureDBInsert | FeatureEntityStruct,
			tableExclude: 0,
			featureReq:   FeatureEntityStruct | FeatureDBUpdate,
			mode:         'a',
		}, want: 2},
		{name: "OR includes only, two feature ok", args: args{
			tableInclude: FeatureDBUpdate | FeatureDBInsert | FeatureEntityStruct,
			tableExclude: 0,
			featureReq:   FeatureEntityStruct | FeatureDBUpdate,
			mode:         'o',
		}, want: 4},
		{name: "OR includes only, one feature ok", args: args{
			tableInclude: FeatureDBUpdate | FeatureDBInsert | FeatureEntityStruct,
			tableExclude: 0,
			featureReq:   FeatureEntityStruct | FeatureCollectionEach,
			mode:         'o',
		}, want: 4},
		{name: "OR includes only, zero feature ok", args: args{
			tableInclude: FeatureDBUpdate | FeatureDBInsert | FeatureEntityStruct,
			tableExclude: 0,
			featureReq:   FeatureCollectionStruct | FeatureCollectionEach,
			mode:         'o',
		}, want: -2},
		{name: "AND includes only, two feature nok", args: args{
			tableInclude: FeatureDBUpdate | FeatureDBInsert | FeatureEntityStruct,
			tableExclude: 0,
			featureReq:   FeatureEntityStruct | FeatureDBUpsert,
			mode:         'a',
		}, want: -2},
		{name: "AND excludes only, two feature nok", args: args{
			tableInclude: 0,
			tableExclude: FeatureEntityStruct | FeatureDBUpdate | FeatureDBInsert,
			featureReq:   FeatureEntityStruct | FeatureDBUpsert,
			mode:         'a',
		}, want: -1},
		{name: "AND excludes only, two feature ok", args: args{
			tableInclude: 0,
			tableExclude: FeatureEntityStruct | FeatureDBUpdate | FeatureDBInsert,
			featureReq:   FeatureDBSelect | FeatureDBUpsert,
			mode:         'a',
		}, want: 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasFeature(tt.args.tableInclude, tt.args.tableExclude, tt.args.featureReq, tt.args.mode); got != tt.want {
				t.Errorf("hasFeature() = %v, want %v", got, tt.want)
			}
		})
	}
}
