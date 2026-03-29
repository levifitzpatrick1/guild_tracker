package updates

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/levifitzpatrick1/guild_tracker/generated/dbstore"
	"github.com/levifitzpatrick1/guild_tracker/internal/battlenet"
	"github.com/levifitzpatrick1/guild_tracker/internal/battlenet/responseStructs"
)

func UpdateCharactersInGuild(ctx context.Context, name, server string, b *battlenet.Battlenet, q *dbstore.Queries) ([]responseStructs.Member, error) {
	roster, err := b.GetGuildRoster(name, server)
	if err != nil {
		return nil, fmt.Errorf("failed to get guild info: %w", err)
	}

	for _, member := range roster.Members {

		params := dbstore.UpsertCharacterParams{
			ID:     int64(member.Character.ID),
			Name:   member.Character.Name,
			Server: member.Character.Realm.Name,
			Level:  int64(member.Character.Level),
			Guild: sql.NullString{
				String: name,
				Valid:  true,
			},
			Score: sql.NullFloat64{
				Valid: false,
			},
		}

		_, err := q.UpsertCharacter(ctx, params)
		if err != nil {
			b.Logger.Error("failed to upsert character %s: %v", member.Character.Name, err)
			continue
		}
		b.Logger.Info("upserted character: %s", member.Character.Name)
	}

	return roster.Members, nil

}
