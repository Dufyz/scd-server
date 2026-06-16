package repositories

import (
	"database/sql"
	"fmt"

	db "github.com/Dufyz/scd-server/infra/database"
	"github.com/Dufyz/scd-server/internal/domain/entities"
	"github.com/Dufyz/scd-server/internal/shared/dtos"
	"github.com/Dufyz/scd-server/internal/shared/interfaces"
	"go.uber.org/zap"
)

var _ interfaces.ChatRepositoryInterface = &ChatRepository{}

type ChatRepository struct {
	connection *db.ReplicatedDB
}

func NewChatRepository(connection *db.ReplicatedDB) interfaces.ChatRepositoryInterface {
	return &ChatRepository{
		connection: connection,
	}
}

func (r *ChatRepository) FindById(id int64) (*entities.Chat, error) {
	var chat entities.Chat
	err := r.connection.QueryRow(`
		SELECT id, name, category, created_at, updated_at
		FROM "chats"
		WHERE id = $1
	`, id).Scan(
		&chat.ID,
		&chat.Name,
		&chat.Category,
		&chat.CreatedAt,
		&chat.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		zap.L().Error("Error on scan row Chat/Repository/FindById", zap.Error(err))
		return nil, err
	}

	return &chat, nil
}

func (r *ChatRepository) List(filters dtos.ChatFilters) ([]entities.Chat, error) {
	query := `
		SELECT id, name, category, created_at, updated_at
		FROM chats
		WHERE 1 = 1
	`

	args := []interface{}{}
	argPos := 1

	if filters.Name != nil && *filters.Name != "" {
		query += fmt.Sprintf(" AND unaccent(LOWER(name)) ILIKE '%%' || unaccent(LOWER($%d)) || '%%'", argPos)
		args = append(args, *filters.Name)
		argPos++
	}

	if filters.Category != nil && *filters.Category != "" {
		query += fmt.Sprintf(" AND unaccent(LOWER(category)) ILIKE '%%' || unaccent(LOWER($%d)) || '%%'", argPos)
		args = append(args, *filters.Category)
		argPos++
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.connection.Query(query, args...)
	if err != nil {
		zap.L().Error("Error on SELECT Chat/Repository/List", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var chats []entities.Chat

	for rows.Next() {
		var chat entities.Chat

		err := rows.Scan(
			&chat.ID,
			&chat.Name,
			&chat.Category,
			&chat.CreatedAt,
			&chat.UpdatedAt,
		)
		if err != nil {
			zap.L().Error("Error on scan row Chat/Repository/List", zap.Error(err))
			return nil, err
		}

		chats = append(chats, chat)
	}

	if err := rows.Err(); err != nil {
		zap.L().Error("Error iterating rows Chat/Repository/List", zap.Error(err))
		return nil, err
	}

	return chats, nil
}

func (r *ChatRepository) Create(body dtos.CreateChat) (entities.Chat, error) {
	tx, err := r.connection.Begin()
	if err != nil {
		zap.L().Error("Error starting transaction Chat/Repository/Create", zap.Error(err))
		return entities.Chat{}, err
	}
	defer tx.Rollback()

	var chat entities.Chat
	err = tx.QueryRow(`
		INSERT INTO "chats" (name, category)
		VALUES ($1, $2)
		RETURNING id, name, category, created_at, updated_at
	`, body.Name, body.Category).Scan(
		&chat.ID,
		&chat.Name,
		&chat.Category,
		&chat.CreatedAt,
		&chat.UpdatedAt,
	)

	if err != nil {
		zap.L().Error("Error on INSERT INTO chat Chat/Repository/Create", zap.Error(err))
		return entities.Chat{}, err
	}

	err = tx.Commit()
	if err != nil {
		zap.L().Error("Error committing transaction Chat/Repository/Create", zap.Error(err))
		return entities.Chat{}, err
	}

	return chat, nil
}

func (r *ChatRepository) Update(id int64, body dtos.UpdateChat) (entities.Chat, error) {
	tx, err := r.connection.Begin()
	if err != nil {
		zap.L().Error("Error starting transaction Chat/Repository/UpdateBasicInfo", zap.Error(err))
		return entities.Chat{}, err
	}
	defer tx.Rollback()

	var chat entities.Chat
	err = tx.QueryRow(`
		UPDATE "chats"
		SET name = $1, category = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, name, category, created_at, updated_at
	`, body.Name, body.Category, id).Scan(
		&chat.ID,
		&chat.Name,
		&chat.Category,
		&chat.CreatedAt,
		&chat.UpdatedAt,
	)

	if err != nil {
		zap.L().Error("Error on UPDATE chat Chat/Repository/UpdateBasicInfo", zap.Error(err))
		return entities.Chat{}, err
	}

	err = tx.Commit()
	if err != nil {
		zap.L().Error("Error committing transaction Chat/Repository/UpdateBasicInfo", zap.Error(err))
		return entities.Chat{}, err
	}

	return chat, nil
}

func (r *ChatRepository) Delete(id int64) error {
	result, err := r.connection.Exec(`DELETE FROM "chats" WHERE id = $1`, id)
	if err != nil {
		zap.L().Error("Error on DELETE chats Chat/Repository/Delete", zap.Error(err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.L().Error("Error getting rows affected Chat/Repository/Delete", zap.Error(err))
		return err
	}

	if rowsAffected == 0 {
		zap.L().Warn("No chat deleted", zap.Int64("id", id))
	}

	return nil
}
