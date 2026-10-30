package user

import (
	"auth/internal/model"
	desc "auth/pkg/user_v1"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func ToUserCreateFromProto(data *desc.UserCreate) model.UserCreate {

	return model.UserCreate{
		Name:     data.Info.Name,
		Password: data.Password,
		Email:    data.Info.Email,
		Roles:    data.Info.Roles.String(),
	}

}

func ToUserPatchFromProto(req *desc.UpdateRequest) model.UserPatch {
	return model.UserPatch{
		ID:    req.Id,
		Name:  optionalString(req.Name),
		Email: optionalString(req.Email),
		Roles: optionsRole(req.Roles),
	}
}

func ToProtoUserFromService(data model.User) *desc.User {

	return &desc.User{
		Id: data.ID,
		Info: &desc.UserData{
			Name:  data.Name,
			Email: data.Email,
			Roles: ToProtoRoles(data.Roles),
		},
		CreatedAt: toProtoTimestamp(data.CreatedAt),
		UpdatedAt: toProtoTimestamp(data.UpdatedAt),
	}

}

// model → proto
func ToProtoRoles(role string) desc.Roles {
	switch role {
	case "ADMIN":
		return desc.Roles_ADMIN
	default:
		return desc.Roles_USER
	}
}

// proto → model
func ToModelRoles(role desc.Roles) string {
	switch role {
	case desc.Roles_ADMIN:
		return "ADMIN"
	default:
		return "USER"
	}
}

func optionalString(v *wrapperspb.StringValue) *string {
	if v == nil {
		return nil
	}

	s := v.GetValue()
	return &s
}

func optionsRole(r *desc.Roles) *string {
	if r == nil {
		return nil
	}

	s := ToModelRoles(*r)
	return &s
}

func toProtoTimestamp(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}
