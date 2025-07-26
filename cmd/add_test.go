package cmd_test

import (
	"gs/cmd"
	mocks "gs/mocks/cmd"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type MockCall[T any] struct {
	Times    int
	Response T
	Error    error
}

func TestAddCmd(t *testing.T) {
	currentPath := "currentPath"
	parentFolderName := "parentFolderName"
	tests := []struct {
		name                    string
		args                    []string
		expectedError           string
		mockGetCurrentPath      MockCall[string]
		mockGetParentFolderName MockCall[string]
		mockCheckIfPathExists   MockCall[bool]
		mockAdd                 MockCall[error]
	}{
		{
			name: "successful no args",
			args: []string{},
			mockGetCurrentPath: MockCall[string]{
				Times:    1,
				Response: currentPath,
			},
			mockGetParentFolderName: MockCall[string]{
				Times:    1,
				Response: parentFolderName,
			},
			mockAdd: MockCall[error]{
				Times: 1,
			},
		},
		{
			name: "failed no args due to fail to get current directory",
			args: []string{},
			mockGetCurrentPath: MockCall[string]{
				Times:    1,
				Response: "",
				Error:    assert.AnError,
			},
			expectedError: "failed to get current directory",
		},
		{
			name: "failed no args due to fail to add to database",
			args: []string{},
			mockGetCurrentPath: MockCall[string]{
				Times:    1,
				Response: currentPath,
			},
			mockGetParentFolderName: MockCall[string]{
				Times:    1,
				Response: parentFolderName,
			},
			mockAdd: MockCall[error]{
				Times: 1,
				Error: assert.AnError,
			},
			expectedError: "failed to add to database",
		},
		{
			name: "successful with alias arg",
			args: []string{"alias"},
			mockGetCurrentPath: MockCall[string]{
				Times:    1,
				Response: currentPath,
			},

			mockAdd: MockCall[error]{
				Times: 1,
			},
		},
		{
			name: "failed with alias arg due to fail to get current directory",
			args: []string{"alias"},
			mockGetCurrentPath: MockCall[string]{
				Times: 1,
				Error: assert.AnError,
			},
			expectedError: "failed to get current directory",
		},
		{
			name: "failed with alias arg due to fail to add to database",
			args: []string{"alias"},
			mockGetCurrentPath: MockCall[string]{
				Times:    1,
				Response: currentPath,
			},

			mockAdd: MockCall[error]{
				Times: 1,
				Error: assert.AnError,
			},
			expectedError: "failed to add to database",
		},
		{
			name: "successful with alias and path arg",
			args: []string{"alias", "path"},
			mockCheckIfPathExists: MockCall[bool]{
				Times:    1,
				Response: true,
			},
			mockAdd: MockCall[error]{
				Times: 1,
			},
		},
		{
			name: "failed with alias and path arg due to fail to check if path exists",
			args: []string{"alias", "path"},
			mockCheckIfPathExists: MockCall[bool]{
				Times:    1,
				Response: false,
			},
			expectedError: "path does not exist",
		},
		{
			name: "failed with alias and path arg due to fail to add to database",
			args: []string{"alias", "path"},
			mockCheckIfPathExists: MockCall[bool]{
				Times:    1,
				Response: true,
			},
			mockAdd: MockCall[error]{
				Times: 1,
				Error: assert.AnError,
			},
			expectedError: "failed to add to database",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDBService := mocks.NewMockDBService(ctrl)
			mockFileService := mocks.NewMockFileService(ctrl)
			cmd := cmd.NewAddCmd(mockDBService, mockFileService)

			var alias string
			if len(tt.args) > 0 {
				alias = tt.args[0]
			} else {
				alias = tt.mockGetParentFolderName.Response
			}

			var path string
			if len(tt.args) > 1 {
				path = tt.args[1]
			} else {
				path = tt.mockGetCurrentPath.Response
			}

			mockFileService.EXPECT().GetCurrentPath().Return(tt.mockGetCurrentPath.Response, tt.mockGetCurrentPath.Error).Times(tt.mockGetCurrentPath.Times)
			mockFileService.EXPECT().GetParentFolderName(gomock.Any()).Return(tt.mockGetParentFolderName.Response).Times(tt.mockGetParentFolderName.Times)
			mockFileService.EXPECT().CheckIfPathExists(path).Return(tt.mockCheckIfPathExists.Response).Times(tt.mockCheckIfPathExists.Times)
			mockDBService.EXPECT().Add(alias, path).Return(tt.mockAdd.Error).Times(tt.mockAdd.Times)

			cmd.SetArgs(tt.args)
			err := cmd.Execute()

			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError)
			}
		})
	}
}
