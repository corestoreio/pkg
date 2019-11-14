package dmlgen

import (
	"context"
	"testing"

	"github.com/corestoreio/pkg/util/assert"
)

func TestGenerator_isAllowedRelationship(t *testing.T) {
	ctx := context.TODO()
	t.Run("exclude wildcards01", func(t *testing.T) {
		g, err := NewGenerator("",
			WithForeignKeyRelationships(ctx, nil, ForeignKeyOptions{
				// IncludeRelationShips: []string{"what are the names?"},
				ExcludeRelationships: []string{
					"athlete.athlete_id", "athlete_team_member.athlete_id",
					"athlete_team.team_id", "athlete_team_member.team_id",

					"athlete_team_member.*", "*.*", // do not print relations for the relation table itself.
					"athlete_team_member2.column2", "*.*", // do not print relations for the relation table itself.
					"athlete_team_member3.column3", "athlete.*", // do not print relations for the relation table itself.
				},
			},
			),
		)
		assert.NoError(t, err)
		assert.False(t, g.isAllowedRelationship("athlete", "athlete_id", "athlete_team_member", "athlete_id"))
		assert.False(t, g.isAllowedRelationship("athlete_team_member", "column1", "athlete_team", "column1"))
		assert.False(t, g.isAllowedRelationship("athlete_team_member2", "column2", "athlete_team", "column1"))
		assert.False(t, g.isAllowedRelationship("athlete_team_member3", "column3", "athlete", "column3"))
		assert.True(t, g.isAllowedRelationship("athlete_team_member4", "column4", "athlete", "column3"))
	})

	t.Run("include wildcards01", func(t *testing.T) {
		g, err := NewGenerator("",
			WithForeignKeyRelationships(ctx, nil, ForeignKeyOptions{
				// IncludeRelationShips: []string{"what are the names?"},
				IncludeRelationShips: []string{
					"athlete.athlete_id", "athlete_team_member.athlete_id",
					"athlete_team.team_id", "athlete_team_member.team_id",

					"athlete_team_member.*", "*.*", // do not print relations for the relation table itself.
					"athlete_team_member2.column2", "*.*", // do not print relations for the relation table itself.
					"athlete_team_member3.column3", "athlete.*", // do not print relations for the relation table itself.

				},
			},
			),
		)
		assert.NoError(t, err)
		assert.True(t, g.isAllowedRelationship("athlete", "athlete_id", "athlete_team_member", "athlete_id"))
		assert.True(t, g.isAllowedRelationship("athlete_team_member", "column1", "athlete_team", "column1"))
		assert.True(t, g.isAllowedRelationship("athlete_team_member2", "column2", "athlete_team", "column1"))
		assert.True(t, g.isAllowedRelationship("athlete_team_member3", "column3", "athlete", "column3"))
		assert.True(t, g.isAllowedRelationship("athlete_team", "team_id", "athlete_team_member", "team_id"))
	})
}
