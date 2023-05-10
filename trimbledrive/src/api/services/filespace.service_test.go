package services

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"

	"trimble.com/common/repository"

	"trimble.com/tdrive/api/model/filespace"
)

var testString = uuid.NewString()
var testUserId = uuid.NewString()
var testId = uuid.NewString()
var changeToken = "0"

func getTestString() *string {
	return &testString
}

func getTestUserId() *string {
	return &testUserId
}

func getTestId() *string {
	return &testId
}

func Test_convertToFilespaceEntity(t *testing.T) {
	type args struct {
		createRequestBody filespace.CreateRequest
		userId            string
	}
	tests := []struct {
		name string
		args args
		want *repository.Filespace
	}{
		{
			name: "CreateRequestWithoutId",
			args: args{
				createRequestBody: filespace.CreateRequest{},
				userId:            testUserId,
			},
			want: &repository.Filespace{
				CreatedBy:  getTestString(),
				ModifiedBy: getTestString(),
			},
		},
		{
			name: "CreateRequestWithId",
			args: args{
				createRequestBody: filespace.CreateRequest{Id: getTestId()},
				userId:            testUserId,
			},
			want: &repository.Filespace{
				ID:         getTestId(),
				CreatedBy:  getTestString(),
				ModifiedBy: getTestString(),
			},
		},
		{
			name: "CreateRequestWithIdAndDescription",
			args: args{
				createRequestBody: filespace.CreateRequest{Id: getTestId(), Description: getTestString()},
				userId:            testUserId,
			},
			want: &repository.Filespace{
				ID:          getTestId(),
				Description: getTestString(),
				CreatedBy:   getTestUserId(),
				ModifiedBy:  getTestUserId(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertToFilespaceEntity(tt.args.createRequestBody, tt.args.userId)
			request := tt.args.createRequestBody
			// Validate Id
			if (request.Id == nil && got.ID == nil) || (request.Id != nil && got.ID != getTestId()) {
				t.Errorf("ID validation failed. convertToFilespaceEntity() = %v, want %v", got, tt.want)
			}
			// Validate UserId
			if *got.CreatedBy != testUserId || *got.ModifiedBy != testUserId {
				t.Errorf("User ID validation failed. convertToFilespaceEntity() = %v, want %v", got, tt.want)
			}
			// Validate Description
			if request.Description != nil && got.Description != getTestString() {
				t.Errorf("Description validation failed. convertToFilespaceEntity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertToFilespaceResponse(t *testing.T) {
	rootId := uuid.NewString()
	createdTime := time.Now().UTC()
	spaceEntity := repository.NewFilespace(getTestString(), createdTime, testUserId)
	spaceEntity.ID = getTestId()
	spaceEntity.RootId = &rootId

	want := &filespace.Response{
		Id:          getTestId(),
		Description: getTestString(),
		ChangeToken: &changeToken,
		RootId:      &rootId,
		Deleted:     nil,
		CreatedOn:   &createdTime,
		CreatedBy:   getTestUserId(),
		ModifiedOn:  &createdTime,
		ModifiedBy:  getTestUserId(),
	}

	if got := convertToFilespaceResponse(spaceEntity); !reflect.DeepEqual(got, want) {
		t.Errorf("convertToFilespaceResponse() = %v, want %v", got, want)
	}
}
