package repository

import (
	"testing"

	"github.com/google/uuid"
)

var testSpaceId = uuid.NewString()
var testVisibleId = uuid.NewString()
var testResourceType = "TestType"
var testDbId = testSpaceId + "#" + testResourceType + "#" + testVisibleId

func TestGetDbId(t *testing.T) {
	type args struct {
		spaceId      string
		resourceType string
		visibleId    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestGetDbId",
			args: args{
				spaceId:      testSpaceId,
				resourceType: testResourceType,
				visibleId:    testVisibleId,
			},
			want: testDbId,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetDbId(tt.args.spaceId, tt.args.resourceType, &tt.args.visibleId)
			if *got != tt.want {
				t.Errorf("GetDbId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetVisibleId(t *testing.T) {
	type args struct {
		dbId string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestGetVisibleId",
			args: args{dbId: testDbId},
			want: testVisibleId,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetVisibleId(&tt.args.dbId)
			if *got != tt.want {
				t.Errorf("GetVisibleId() = %v, want %v", got, tt.want)
			}
		})
	}
}
