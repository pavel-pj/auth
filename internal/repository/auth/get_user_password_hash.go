package auth

import (
	"auth/internal/model"
	converter "auth/internal/repository/auth/converter"
	repoModel "auth/internal/repository/auth/model"
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/pavel-pj/platform_lib/pkg/client/db"
)

func (r *AuthRepository) GetUserHash(ctx context.Context, login string) (model.UserLogin, error) {

	sqlStr, arg, err := psql.Select().
		From(tableName).
		Columns(columnName, columnPassword, columnRole).
		Where(sq.Eq{columnName: login}).
		ToSql()
	if err != nil {
		return model.UserLogin{}, err
	}

	q := db.Query{
		Name:     "auth_repository.get_user_password_hash",
		QueryRow: sqlStr,
	}

	var userLogin repoModel.UserLogin

	err = r.db.DB().ScanOneContext(ctx, &userLogin, q, arg...)

	if err != nil {
		return model.UserLogin{}, err
	}

	user := converter.UserLoginFromRepoToService(userLogin)

	return user, nil

}
