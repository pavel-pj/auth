package converter

import (
	modelService "github.com/pavel-pj/auth/internal/model"
	modelRepo "github.com/pavel-pj/auth/internal/repository/log/model"
)

func ToRepoLogFromService(data modelService.Log) modelRepo.Log {

	return modelRepo.Log{
		Level:   modelRepo.Level(data.Level),
		Message: data.Message,
	}

}
