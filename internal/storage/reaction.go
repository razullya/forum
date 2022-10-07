package storage

import (
	"database/sql"
	"fmt"
	"forum/models"
)

type Reaction interface {
	//
	UpdateReaction(id int, action string, object string, username string) error
	CreateReaction(id int, action string, object string, username string) error
	GetReactionsById(id int, object string) ([]models.Reaction, error)
}
type ReactionStorage struct {
	db *sql.DB
}

func newReactionStorage(db *sql.DB) *ReactionStorage {
	return &ReactionStorage{
		db: db,
	}
}
func (r *ReactionStorage) UpdateReaction(id int, action string, object string, username string) error {
	query := `UPDATE reaction SET action=$1 WHERE creator = $2 AND id_object = $3 AND object=$4;`

	if _, err := r.db.Exec(query, action, username, id, object); err != nil {
		return fmt.Errorf("storage: update reaction: %w", err)
	}
	return nil
}
func (r *ReactionStorage) GetReactionsById(id int, object string) ([]models.Reaction, error) {

	query := `SELECT * FROM reaction WHERE id_object=$1 and object=$2;`

	row, err := r.db.Query(query, id, object)
	if err != nil {
		return nil, fmt.Errorf("storage: get reactions by id: %w", err)
	}

	reactions := []models.Reaction{}
	for row.Next() {
		var reaction models.Reaction
		if err := row.Scan(&reaction.Id, &reaction.PostId, &reaction.Reaction, &reaction.Username, &reaction.Object); err != nil {
			return nil, fmt.Errorf("storage: get reactions by id: %w", err)
		}
		reactions = append(reactions, reaction)
	}
	return reactions, nil
}

func (r *ReactionStorage) CreateReaction(id int, action string, object string, username string) error {
	query := `INSERT INTO reaction(id_object, action, creator, object) VALUES($1, $2, $3, $4);`
	if _, err := r.db.Exec(query, id, action, username, object); err != nil {
		return fmt.Errorf("storage: create reaction: %w", err)
	}
	return nil
}
