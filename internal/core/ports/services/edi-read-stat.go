package ports

import "backend/internal/core/models"

type EDIReadStatService interface {
	TrackReadService(req models.EDIReadStatReq) error
}
