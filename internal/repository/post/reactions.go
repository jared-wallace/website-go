package post

import (
	"context"
	"fmt"
)

const insertReaction = `
	INSERT INTO reactions (post_id, ip_hash)
	VALUES ($1, $2)
	ON CONFLICT (post_id, ip_hash) DO NOTHING`

const countReactions = `SELECT COUNT(*) FROM reactions WHERE post_id = $1`

func (r *postgresRepository) AddReaction(ctx context.Context, postID int64, ipHash string) (bool, error) {
	tag, err := r.pool.Exec(ctx, insertReaction, postID, ipHash)
	if err != nil {
		return false, fmt.Errorf("AddReaction: %w", err)
	}
	return tag.RowsAffected() == 0, nil
}

func (r *postgresRepository) CountReactions(ctx context.Context, postID int64) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, countReactions, postID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("CountReactions: %w", err)
	}
	return count, nil
}
