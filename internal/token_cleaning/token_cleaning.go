package token_cleaning

import (
	"context"
	"documents/internal/store"
	"log"
	"time"
)

func StartTokenCleaning(ctx context.Context, tokens *store.TokenStore, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		log.Printf("[token_cleaning] token cleanup started, interval=%s", interval)

		for {
			select {
			case <-ctx.Done():
				log.Println("[token_cleaning] token cleaning stopped")

				return
			case <-ticker.C:
				deletedTokens, err := tokens.DeleteExpiredTokens(ctx)
				if err != nil {
					log.Printf("[token_cleaning] failed to delete tokens: %v", err)
					continue
				}

				if deletedTokens > 0 {
					log.Printf("[token_cleaning] deleted %d tokens", deletedTokens)
				}
			}
		}
	}()
}
