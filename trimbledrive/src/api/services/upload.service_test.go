package services

import (
	"reflect"
	"testing"
	"time"

	"trimble.com/common/api_error"
	"trimble.com/common/constants"
	"trimble.com/common/repository"
	"trimble.com/tdrive/api/model/upload"
)

var testUpdatedOn = time.Now().UTC()
var testFormat = "TestFormat"
var testFileId = "0d168f42-6774-4bdf-a1b6-32811ca1eadf"
var testSpaceId = "cb978029-26fa-45ac-aa48-a178d7be9a61"
var testUploadId = "bb119fc5-1bc6-4bbb-87eb-696e1d8ac589"
var testFormatTrb = "TRB"
var testFormatPdf = "PDF"
var testFormatThumb = constants.FormatThumbnail
var testContentsWithNilFormat = []upload.FileUploadContent{{Format: nil}}
var testMaxContentsExceeded = []upload.FileUploadContent{{Format: &testFormat}, {Format: &testFileId}, {Format: &testSpaceId}, {Format: &testUploadId}, {Format: nil}, {Format: &testFormatTrb}, {Format: &testFormatPdf}}
var testContentsWithValidFormat = []upload.FileUploadContent{{Format: &testFormatTrb}}
var testContentsWithDuplicateFormat = []upload.FileUploadContent{{Format: nil}, {Format: nil}}
var testExistingFile = repository.Latest{MajorVersion: 2, MinorVersion: 1}
var testIfNoneMatchInvalid = "1"
var testIfNoneMatchValid = "*"
var minorVersion int64 = 1
var testIfMatchValid = repository.Version{MajorVersion: 2, MinorVersion: &minorVersion}
var testIfMatchOld = repository.Version{MajorVersion: 1, MinorVersion: &minorVersion}

func Test_createUploadContent(t *testing.T) {
	type args struct {
		format    *string
		updatedOn time.Time
	}
	tests := []struct {
		name string
		args args
		want repository.UploadContent
	}{
		{
			name: "TestCreateUploadContentWithNilFormat",
			args: args{
				format:    nil,
				updatedOn: testUpdatedOn,
			},
			want: repository.UploadContent{
				Format:     nil,
				Status:     constants.ContentUploadStatusUploadable,
				UploadMode: constants.ContentUploadModeSinglepart,
				UpdatedOn:  testUpdatedOn,
			},
		},
		{
			name: "TestCreateUploadContentWithFormat",
			args: args{
				format:    &testFormat,
				updatedOn: testUpdatedOn,
			},
			want: repository.UploadContent{
				Format:     &testFormat,
				Status:     constants.ContentUploadStatusUploadable,
				UploadMode: constants.ContentUploadModeSinglepart,
				UpdatedOn:  testUpdatedOn,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createUploadContent(tt.args.format, tt.args.updatedOn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createUploadContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateStoragePath(t *testing.T) {
	type args struct {
		spaceId  string
		fileId   *string
		uploadId string
		format   *string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestGeneratePathWithoutFormat",
			args: args{
				spaceId:  testSpaceId,
				fileId:   &testFileId,
				uploadId: testUploadId,
				format:   nil,
			},
			want: "orig/cb978029-26fa-45ac-aa48-a178d7be9a61/0d168f42-6774-4bdf-a1b6-32811ca1eadf/bb119fc5-1bc6-4bbb-87eb-696e1d8ac589/src",
		},
		{
			name: "TestGeneratePathWithFormat",
			args: args{
				spaceId:  testSpaceId,
				fileId:   &testFileId,
				uploadId: testUploadId,
				format:   &testFormat,
			},
			want: "orig/cb978029-26fa-45ac-aa48-a178d7be9a61/0d168f42-6774-4bdf-a1b6-32811ca1eadf/bb119fc5-1bc6-4bbb-87eb-696e1d8ac589/testformat",
		},
		{
			name: "TestGeneratePathWithFormatThumbnail",
			args: args{
				spaceId:  testSpaceId,
				fileId:   &testFileId,
				uploadId: testUploadId,
				format:   &testFormatThumb,
			},
			want: "thumb/cb978029-26fa-45ac-aa48-a178d7be9a61/0d168f42-6774-4bdf-a1b6-32811ca1eadf/bb119fc5-1bc6-4bbb-87eb-696e1d8ac589/thumbnail",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateStoragePath(tt.args.spaceId, *tt.args.fileId, tt.args.uploadId, tt.args.format); got != tt.want {
				t.Errorf("generateStoragePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateUploadContents(t *testing.T) {
	type args struct {
		contents *[]upload.FileUploadContent
	}
	tests := []struct {
		name string
		args args
		want *api_error.ApiError
	}{
		{name: "TestUploadContentNil", args: args{contents: nil}, want: nil},
		{name: "TestUploadContentFormatNil", args: args{contents: &testContentsWithNilFormat}, want: nil},
		{name: "TestUploadContentsMaxLimitExceeded", args: args{contents: &testMaxContentsExceeded}, want: api_error.FileContentsLimitExceeded},
		{name: "TestUploadContentFormatValid", args: args{contents: &testContentsWithValidFormat}, want: nil},
		{name: "TestUploadContentFormatDuplicate", args: args{contents: &testContentsWithDuplicateFormat}, want: api_error.DuplicateFormat},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateUploadContents(tt.args.contents); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validateUploadContents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateIfNoneMatchAndIfMatch(t *testing.T) {
	type args struct {
		existingFile *repository.Latest
		ifNoneMatch  *string
		ifMatch      *repository.Version
	}
	tests := []struct {
		name string
		args args
		want *api_error.ApiError
	}{
		{name: "TestExistingFileNil",
			args: args{
				existingFile: nil,
				ifNoneMatch:  nil,
				ifMatch:      nil,
			},
			want: nil,
		},
		{
			name: "TestIfNoneMatchAndIfMatchNil",
			args: args{
				existingFile: &testExistingFile,
				ifNoneMatch:  nil,
				ifMatch:      nil,
			},
			want: nil,
		},
		{
			name: "TestIfNoneMatchInvalidAndIfMatchNil",
			args: args{
				existingFile: &testExistingFile,
				ifNoneMatch:  &testIfNoneMatchInvalid,
				ifMatch:      nil,
			},
			want: nil,
		},
		{
			name: "TestIfNoneMatchInvalidAndIfMatchValid",
			args: args{
				existingFile: &testExistingFile,
				ifNoneMatch:  &testIfNoneMatchInvalid,
				ifMatch:      &testIfMatchValid,
			},
			want: nil,
		},
		{
			name: "TestIfNoneMatchInvalidAndIfMatchOld",
			args: args{
				existingFile: &testExistingFile,
				ifNoneMatch:  &testIfNoneMatchInvalid,
				ifMatch:      &testIfMatchOld,
			},
			want: api_error.InvalidVersion,
		},
		{
			name: "TestIfNoneMatchValidAndIfMatchNil",
			args: args{
				existingFile: &testExistingFile,
				ifNoneMatch:  &testIfNoneMatchValid,
				ifMatch:      nil,
			},
			want: api_error.FileWithSameNameExists,
		},
		{
			name: "TestIfNoneMatchValidAndIfMatchValid",
			args: args{
				existingFile: &testExistingFile,
				ifNoneMatch:  &testIfNoneMatchValid,
				ifMatch:      &testIfMatchValid,
			},
			want: api_error.FileWithSameNameExists,
		},
		{
			name: "TestIfNoneMatchValidAndIfMatchOld",
			args: args{
				existingFile: &testExistingFile,
				ifNoneMatch:  &testIfNoneMatchValid,
				ifMatch:      &testIfMatchOld,
			},
			want: api_error.FileWithSameNameExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := preCheckIfNoneMatchAndIfMatch(tt.args.existingFile, tt.args.ifNoneMatch, tt.args.ifMatch); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validateIfNoneMatchAndIfMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
