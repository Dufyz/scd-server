package interfaces

import (
	"github.com/Dufyz/scd-server/internal/domain/entities"
	"github.com/Dufyz/scd-server/internal/shared/dtos"
)

type MessageRepositoryInterface interface {
	FindById(id int64) (*entities.Message, error)
	ListByChatId(chatId int64) ([]entities.Message, error)
	Create(body dtos.CreateMessage) (entities.Message, error)
	Update(id int64, body dtos.UpdateMessage) (entities.Message, error)
	UpdateLanguage(id int64, language string) (entities.Message, error)
	Delete(id int64) error
}
