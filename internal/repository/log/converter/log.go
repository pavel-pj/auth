package converter

import (
	modelService "auth/internal/model"
	modelRepo "auth/internal/repository/log/model"
)

func ToRepoLogFromService(data modelService.Log) modelRepo.Log {

	return modelRepo.Log{
		Level:   modelRepo.Level(data.Level),
		Message: data.Message,
	}

}
