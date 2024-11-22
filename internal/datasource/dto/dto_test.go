package dto

import (
	"database/sql"
	"testing"

	"github.com/EFG/api"
)

func TestGetUsersArgs_FromAPI(t *testing.T) {
	type fields struct {
		Page            sql.NullInt32
		PageSize        sql.NullInt32
		FilterID        sql.NullString
		FilterFirstName sql.NullString
		FilterLastName  sql.NullString
		FilterNickname  sql.NullString
		FilterEmail     sql.NullString
		FilterCountry   sql.NullString
	}
	type args struct {
		req *api.GetUsersRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "empty request",
			fields: fields{
				Page:            sql.NullInt32{Valid: false},
				PageSize:        sql.NullInt32{Valid: false},
				FilterID:        sql.NullString{Valid: false},
				FilterFirstName: sql.NullString{Valid: false},
				FilterLastName:  sql.NullString{Valid: false},
				FilterNickname:  sql.NullString{Valid: false},
				FilterEmail:     sql.NullString{Valid: false},
				FilterCountry:   sql.NullString{Valid: false},
			},
			args: args{
				req: &api.GetUsersRequest{},
			},
		},
		{
			name: "non-empty request",
			fields: fields{
				Page:            sql.NullInt32{Valid: true, Int32: 1},
				PageSize:        sql.NullInt32{Valid: true, Int32: 10},
				FilterID:        sql.NullString{Valid: true, String: "1"},
				FilterFirstName: sql.NullString{Valid: true, String: "John"},
				FilterLastName:  sql.NullString{Valid: true, String: "Doe"},
				FilterNickname:  sql.NullString{Valid: true, String: "johndoe"},
				FilterEmail: sql.NullString{
					Valid: true, String: "john.doe@example.com",
				},
				FilterCountry: sql.NullString{Valid: true, String: "US"},
			},
			args: args{
				req: &api.GetUsersRequest{
					Page:            1,
					PageSize:        10,
					FilterId:        "1",
					FilterFirstName: "John",
					FilterLastName:  "Doe",
					FilterNickname:  "johndoe",
					FilterEmail:     "john.doe@example.com",
					FilterCountry:   "US",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GetUsersArgs{
				Page:            tt.fields.Page,
				PageSize:        tt.fields.PageSize,
				FilterID:        tt.fields.FilterID,
				FilterFirstName: tt.fields.FilterFirstName,
				FilterLastName:  tt.fields.FilterLastName,
				FilterNickname:  tt.fields.FilterNickname,
				FilterEmail:     tt.fields.FilterEmail,
				FilterCountry:   tt.fields.FilterCountry,
			}
			g.FromAPI(tt.args.req)
		})
	}
}
