package scheduler

import (
	"context"

	"github.com/levifitzpatrick1/guild_tracker/internal/battlenet/updates"
)

func (s *Scheduler) dailyGuildSync() {
	s.Logger.Info("Starting daily sync...")
	ctx := context.Background()

	if err := s.Battlenet.RefreshToken(); err != nil {
		s.Logger.Error("Failed to refresh bnet token on daily sync: %v", err)
		return
	}

	if err := updates.SyncGuild(ctx, s.GuildName, s.Server, s.Battlenet, s.Queries, s.Logger); err != nil {
		s.Logger.Error("Daily sync failed: %v", err)
		return
	}

	s.Logger.Info("Daily sync complete!")
}
