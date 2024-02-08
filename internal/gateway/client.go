package gateway

import (
	"github.com/eldersonjls/ms-wallet-core/internal/entity"
)

type ClientGateway interface {
	FindByID(id string) (*entity.Client, error)
	Save(client *entity.Client) error
}
