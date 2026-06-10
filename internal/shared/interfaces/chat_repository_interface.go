package interfaces

import (
	"github.com/Dufyz/scd-server/internal/domain/entities"
	"github.com/Dufyz/scd-server/internal/shared/dtos"
)

type ChatRepositoryInterface interface {
	FindById(id int64) (*entities.Chat, error)
	List(filters dtos.ChatFilters) ([]entities.Chat, error)
	Create(body dtos.CreateChat) (entities.Chat, error)
	Update(id int64, body dtos.UpdateChat) (entities.Chat, error)
	Delete(id int64) error
}
