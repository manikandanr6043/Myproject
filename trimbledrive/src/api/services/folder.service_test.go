package services

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"

	"trimble.com/common/constants"
	"trimble.com/common/repository"

	"trimble.com/tdrive/api/model/folder"
)

func TestConvertToResponse(t *testing.T) {
	var testString = uuid.NewString()
	var testUserId = uuid.NewString()
	var testId = uuid.NewString()
	var now = time.Now()

	latestEntity := &repository.Latest{
		Id:             &testId,
		MajorVersion:   1,
		MinorVersion:   1,
		Name:           testString,
		Type:           constants.Folder,
		SpaceId:        testId,
		ParentFolderId: &testId,
		CreatedBy:      testUserId,
		ModifiedBy:     testUserId,
		Created:        now,
		Modified:       now,
	}
	want := &folder.Response{
		Id:         testId,
		Version:    "1.1",
		Type:       constants.Folder,
		SpaceId:    testId,
		CreatedBy:  testUserId,
		ModifiedBy: testUserId,
		CreatedOn:  now,
		ModifiedOn: now,
	}

	if got := convertToResponse(latestEntity, nil, nil); !reflect.DeepEqual(got, want) {
		t.Errorf("convertToResponse() = %v, want %v", got, want)
	}
}

func TestConvertToResponseWithSelectedFields(t *testing.T) {
	var testString = uuid.NewString()
	var testUserId = uuid.NewString()
	var testId = uuid.NewString()
	var now = time.Now()

	latestEntity := &repository.Latest{
		Id:             &testId,
		MajorVersion:   1,
		MinorVersion:   1,
		Name:           testString,
		Type:           constants.Folder,
		SpaceId:        testId,
		ParentFolderId: &testId,
		CreatedBy:      testUserId,
		ModifiedBy:     testUserId,
		Created:        now,
		Modified:       now,
	}
	want := &folder.Response{
		Id:                      testId,
		Version:                 "1.1",
		Name:                    &testString,
		Type:                    constants.Folder,
		SpaceId:                 testId,
		ParentFolderId:          &testId,
		CreatedBy:               testUserId,
		ModifiedBy:              testUserId,
		CreatedOn:               now,
		ModifiedOn:              now,
		TopActiveParentFolderId: &testId,
	}
	fields := []string{"name", "parentFolderId", "topActiveParentFolderId"}

	if got := convertToResponse(latestEntity, &testId, &fields); !reflect.DeepEqual(got, want) {
		t.Errorf("convertToResponse() = %v, want %v", got, want)
	}
}

func TestConvertToListResponse(t *testing.T) {
	var testString = uuid.NewString()
	var testUserId = uuid.NewString()
	var testId = uuid.NewString()
	var now = time.Now()
	var testDbId = testId + "#FOLDER#" + testId
	var latestEntityList []*repository.Latest
	latestEntityList = append(latestEntityList,
		&repository.Latest{
			Id:             &testDbId,
			MajorVersion:   1,
			MinorVersion:   2,
			Name:           testString,
			Type:           constants.Folder,
			SpaceId:        testId,
			ParentFolderId: &testDbId,
			CreatedBy:      testUserId,
			ModifiedBy:     testUserId,
			Created:        now,
			Modified:       now,
		},
	)

	want := &folder.ListResponse{
		Items: []folder.Response{
			{
				Id:         testId,
				Version:    "1.2",
				Type:       constants.Folder,
				SpaceId:    testId,
				CreatedBy:  testUserId,
				ModifiedBy: testUserId,
				CreatedOn:  now,
				ModifiedOn: now,
			},
		},
	}

	if got := convertToListResponse(latestEntityList, &testId, nil); !reflect.DeepEqual(got, want) {
		t.Errorf("convertToListResponse() = %v, want %v", got, want)
	}
}
