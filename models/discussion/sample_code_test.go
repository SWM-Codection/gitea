package discussion_test

import (
	"code.gitea.io/gitea/models/discussion"
	"code.gitea.io/gitea/models/unittest"
	"context"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateDiscussionAiSampleCode(t *testing.T) {

	tests := []struct {
		name    string
		opts    *discussion.CreateDiscussionAiCommentOpt
		wantErr bool
	}{
		{
			name: "Valid creation",
			opts: &discussion.CreateDiscussionAiCommentOpt{
				TargetCommentId: 1,
				GenearaterId:    1,
				Type: 			 "pull",
				Content:         stringPtr("Test content"),
			},
			wantErr: false,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			assert.NoError(t, unittest.PrepareTestDatabase())
			got, err := discussion.CreateAiSampleCode(context.Background(), tt.opts)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.opts.TargetCommentId, got.TargetCommentId)
				assert.Equal(t, tt.opts.GenearaterId, got.GenearaterId)
				assert.Equal(t, *tt.opts.Content, got.Content)
			}
		})
	}
}

func TestDeleteDiscussionAiSampleCodeByID(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	tests := []struct {
		name    string
		id      int64
		Type string
		wantErr bool
	}{
		{
			name:    "Valid deletion",
			id:      1,
			Type: "pull",
			wantErr: false,
		},
		{
			name:    "Invalid deletion - non-existent ID",
			id:      99999,
			Type: "pull",
			wantErr: true,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NoError(t, unittest.PrepareTestDatabase())
			err := discussion.DeleteAiSampleCodeByID(context.Background(), tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Verify deletion
				_, err := discussion.GetAiSampleCodeByCommentID(context.Background(), tt.id, tt.Type)
				assert.Error(t, err) // Should not find the deleted item
			}
		})
	}
}

func TestGetAiSampleCodeByCommentID(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	tests := []struct {
		name    string
		id      int64
		want    int
		Type	string
		wantErr bool
	}{
		{
			name:    "Valid retrieval",
			id:      1,
			want:    2, // Assuming there are 2 sample codes for comment ID 1
			wantErr: false,
		},
		{
			name:    "No sample codes",
			id:      99999,
			want:    0,
			wantErr: false,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NoError(t, unittest.PrepareTestDatabase())
			got, err := discussion.GetAiSampleCodeByCommentID(context.Background(), tt.id, tt.Type)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, got, tt.want)
				for _, sampleCode := range got {
					assert.Equal(t, tt.id, sampleCode.TargetCommentId)
				}
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
