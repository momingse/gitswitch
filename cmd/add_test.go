package cmd_test

import (
	"fmt"
	"gs/cmd"
	mocks "gs/mocks/cmd"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type MockCall[T any] struct {
	args     []string
	Times    int
	Response T
	Error    error
}

func TestAddCmd(t *testing.T) {
	currentPath := "currentPath"
	parentFolderName := "parentFolderName"
	pathValue := "pathValue"
	aliasValue := "aliasValue"

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
				args:     []string{currentPath},
				Times:    1,
				Response: parentFolderName,
			},
			mockAdd: MockCall[error]{
				args:  []string{parentFolderName, currentPath},
				Times: 1,
			},
		},
		{
			name: "failed no args due to fail to get current path",
			args: []string{},
			mockGetCurrentPath: MockCall[string]{
				Times: 1,
				Error: assert.AnError,
			},
			expectedError: "failed to get current path",
		},
		{
			name: "failed no args due to fail to add to database",
			args: []string{},
			mockGetCurrentPath: MockCall[string]{
				Times:    1,
				Response: currentPath,
			},
			mockGetParentFolderName: MockCall[string]{
				args:     []string{currentPath},
				Times:    1,
				Response: parentFolderName,
			},
			mockAdd: MockCall[error]{
				args:  []string{parentFolderName, currentPath},
				Times: 1,
				Error: assert.AnError,
			},
			expectedError: fmt.Sprintf("failed to add %s with alias %s", currentPath, parentFolderName),
		},
		{
			name: "successful with alias arg",
			args: []string{aliasValue},
			mockGetCurrentPath: MockCall[string]{
				Times:    1,
				Response: currentPath,
			},
			mockAdd: MockCall[error]{
				args:  []string{aliasValue, currentPath},
				Times: 1,
			},
		},
		{
			name: "failed with alias arg due to fail to get current path",
			args: []string{aliasValue},
			mockGetCurrentPath: MockCall[string]{
				Times: 1,
				Error: assert.AnError,
			},
			expectedError: "failed to get current path",
		},
		{
			name: "failed with alias arg due to fail to add to database",
			args: []string{aliasValue},
			mockGetCurrentPath: MockCall[string]{
				Times:    1,
				Response: currentPath,
			},
			mockAdd: MockCall[error]{
				args:  []string{aliasValue, currentPath},
				Times: 1,
				Error: assert.AnError,
			},
			expectedError: fmt.Sprintf("failed to add %s with alias %s", currentPath, aliasValue),
		},
		{
			name: "successful with alias and path arg",
			args: []string{aliasValue, pathValue},
			mockCheckIfPathExists: MockCall[bool]{
				args:     []string{pathValue},
				Times:    1,
				Response: true,
			},
			mockAdd: MockCall[error]{
				args:  []string{aliasValue, pathValue},
				Times: 1,
			},
		},
		{
			name: "failed with alias and path arg due to fail to check if path exists",
			args: []string{aliasValue, pathValue},
			mockCheckIfPathExists: MockCall[bool]{
				args:     []string{pathValue},
				Times:    1,
				Response: false,
			},
			expectedError: "path does not exist",
		},
		{
			name: "failed with alias and path arg due to fail to add to database",
			args: []string{aliasValue, pathValue},
			mockCheckIfPathExists: MockCall[bool]{
				args:     []string{pathValue},
				Times:    1,
				Response: true,
			},
			mockAdd: MockCall[error]{
				args:  []string{aliasValue, pathValue},
				Times: 1,
				Error: assert.AnError,
			},
			expectedError: fmt.Sprintf("failed to add %s with alias %s", pathValue, aliasValue),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDBService := mocks.NewMockDBService(ctrl)
			mockFileService := mocks.NewMockFileService(ctrl)
			cmd := cmd.NewAddCmd(mockDBService, mockFileService)

			// Set up mock expectations only when they should be called
			if tt.mockGetCurrentPath.Times > 0 {
				mockFileService.EXPECT().GetCurrentPath().Return(tt.mockGetCurrentPath.Response, tt.mockGetCurrentPath.Error).Times(tt.mockGetCurrentPath.Times)
			}

			if tt.mockGetParentFolderName.Times > 0 && len(tt.mockGetParentFolderName.args) > 0 {
				mockFileService.EXPECT().GetParentFolderName(tt.mockGetParentFolderName.args[0]).Return(tt.mockGetParentFolderName.Response).Times(tt.mockGetParentFolderName.Times)
			}

			if tt.mockCheckIfPathExists.Times > 0 && len(tt.mockCheckIfPathExists.args) > 0 {
				mockFileService.EXPECT().CheckIfPathExists(tt.mockCheckIfPathExists.args[0]).Return(tt.mockCheckIfPathExists.Response, tt.mockCheckIfPathExists.Error).Times(tt.mockCheckIfPathExists.Times)
			}

			if tt.mockAdd.Times > 0 && len(tt.mockAdd.args) >= 2 {
				mockDBService.EXPECT().Add(tt.mockAdd.args[0], tt.mockAdd.args[1]).Return(tt.mockAdd.Error).Times(tt.mockAdd.Times)
			}

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

