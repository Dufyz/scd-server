package repositories

import (
	"database/sql"

	db "github.com/Dufyz/scd-server/infra/database"
	"github.com/Dufyz/scd-server/internal/domain/entities"
	"github.com/Dufyz/scd-server/internal/shared/dtos"
	"github.com/Dufyz/scd-server/internal/shared/interfaces"
	"go.uber.org/zap"
)

var _ interfaces.MessageRepositoryInterface = &MessageRepository{}

type MessageRepository struct {
	connection *db.ReplicatedDB
}

func NewMessageRepository(connection *db.ReplicatedDB) interfaces.MessageRepositoryInterface {
	return &MessageRepository{
		connection: connection,
	}
}

func (r *MessageRepository) FindById(id int64) (*entities.Message, error) {
	var message entities.Message
	err := r.connection.QueryRow(`
		SELECT id, chat_id, message, user_name, created_at, updated_at, language
		FROM "messages"
		WHERE id = $1
	`, id).Scan(
		&message.ID,
		&message.ChatID,
		&message.Message,
		&message.UserName,
		&message.CreatedAt,
		&message.UpdatedAt,
		&message.Language,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		zap.L().Error("Error on scan row Message/Repository/FindById", zap.Error(err))
		return nil, err
	}

	return &message, nil
}

func (r *MessageRepository) ListByChatId(chatId int64) ([]entities.Message, error) {
	rows, err := r.connection.Query(`
		SELECT id, chat_id, message, user_name, created_at, updated_at, language
		FROM "messages"
		WHERE chat_id = $1
		ORDER BY created_at ASC
	`, chatId)
	if err != nil {
		zap.L().Error("Error on query rows Message/Repository/ListByAgentId", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var messages []entities.Message
	for rows.Next() {
		var message entities.Message
		err := rows.Scan(
			&message.ID,
			&message.ChatID,
			&message.Message,
			&message.UserName,
			&message.CreatedAt,
			&message.UpdatedAt,
			&message.Language,
		)
		if err != nil {
			zap.L().Error("Error on scan row Message/Repository/ListByAgentId", zap.Error(err))
			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (r *MessageRepository) Create(body dtos.CreateMessage) (entities.Message, error) {
	tx, err := r.connection.Begin()
	if err != nil {
		zap.L().Error("Error starting transaction Message/Repository/Create", zap.Error(err))
		return entities.Message{}, err
	}
	defer tx.Rollback()

	var message entities.Message
	err = tx.QueryRow(`
		INSERT INTO "messages" (chat_id, message, user_name)
		VALUES ($1, $2, $3)
		RETURNING id, chat_id, message, user_name, created_at, updated_at, language
	`, body.ChatID, body.Message, body.UserName).Scan(
		&message.ID,
		&message.ChatID,
		&message.Message,
		&message.UserName,
		&message.CreatedAt,
		&message.UpdatedAt,
		&message.Language,
	)

	if err != nil {
		zap.L().Error("Error on INSERT INTO message Message/Repository/Create", zap.Error(err))
		return entities.Message{}, err
	}

	err = tx.Commit()
	if err != nil {
		zap.L().Error("Error committing transaction Message/Repository/Create", zap.Error(err))
		return entities.Message{}, err
	}

	return message, nil
}

func (r *MessageRepository) Update(id int64, body dtos.UpdateMessage) (entities.Message, error) {
	tx, err := r.connection.Begin()
	if err != nil {
		zap.L().Error("Error starting transaction Message/Repository/UpdateBasicInfo", zap.Error(err))
		return entities.Message{}, err
	}
	defer tx.Rollback()

	var message entities.Message
	err = tx.QueryRow(`
		UPDATE "messages"
		SET message = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, chat_id, message, user_name, created_at, updated_at, language
	`, body.Message, id).Scan(
		&message.ID,
		&message.ChatID,
		&message.Message,
		&message.UserName,
		&message.CreatedAt,
		&message.UpdatedAt,
		&message.Language,
	)

	if err != nil {
		zap.L().Error("Error on UPDATE message Message/Repository/UpdateBasicInfo", zap.Error(err))
		return entities.Message{}, err
	}

	err = tx.Commit()
	if err != nil {
		zap.L().Error("Error committing transaction Message/Repository/UpdateBasicInfo", zap.Error(err))
		return entities.Message{}, err
	}

	return message, nil
}

func (r *MessageRepository) UpdateLanguage(id int64, language string) (entities.Message, error) {
	tx, err := r.connection.Begin()
	if err != nil {
		zap.L().Error("Error starting transaction Message/Repository/UpdateLanguage", zap.Error(err))
		return entities.Message{}, err
	}
	defer tx.Rollback()

	var message entities.Message
	err = tx.QueryRow(`
		UPDATE "messages"
		SET language = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, chat_id, message, user_name, created_at, updated_at, language
	`, language, id).Scan(
		&message.ID,
		&message.ChatID,
		&message.Message,
		&message.UserName,
		&message.CreatedAt,
		&message.UpdatedAt,
		&message.Language,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("No message updated with language", zap.Int64("id", id), zap.String("language", language))
			return entities.Message{}, err
		}
		zap.L().Error("Error on UPDATE language Message/Repository/UpdateLanguage", zap.Error(err))
		return entities.Message{}, err
	}

	if err = tx.Commit(); err != nil {
		zap.L().Error("Error committing transaction Message/Repository/UpdateLanguage", zap.Error(err))
		return entities.Message{}, err
	}

	return message, nil
}

func (r *MessageRepository) Delete(id int64) error {
	result, err := r.connection.Exec(`DELETE FROM "messages" WHERE id = $1`, id)
	if err != nil {
		zap.L().Error("Error on DELETE message Message/Repository/Delete", zap.Error(err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.L().Error("Error getting rows affected Message/Repository/Delete", zap.Error(err))
		return err
	}

	if rowsAffected == 0 {
		zap.L().Warn("No message deleted", zap.Int64("id", id))
	}

	return nil
}
