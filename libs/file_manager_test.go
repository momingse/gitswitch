package libs_test

import (
	"errors"
	"fs/libs"
	"fs/mocks"
	"testing"

	"go.etcd.io/bbolt"
	"go.uber.org/mock/gomock"
)

func TestService_Add(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		value     string
		setupMock func(*mocks.MockDB)
		wantErr   bool
	}{
		{
			name:  "successful add",
			key:   "test-key",
			value: "test-value",
			setupMock: func(mockDB *mocks.MockDB) {
				mockDB.EXPECT().Update(gomock.Any()).DoAndReturn(func(fn func(*bbolt.Tx) error) error {
					// Create a mock transaction and bucket for testing
					// Since we can't easily mock bbolt.Tx, we'll just return nil to simulate success
					return nil
				})
			},
			wantErr: false,
		},
		{
			name:  "database update error",
			key:   "test-key",
			value: "test-value",
			setupMock: func(mockDB *mocks.MockDB) {
				mockDB.EXPECT().Update(gomock.Any()).Return(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mocks.NewMockDB(ctrl)
			tt.setupMock(mockDB)

			service := libs.NewService(mockDB)
			err := service.Add(tt.key, tt.value)

			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_Get(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		setupMock func(*mocks.MockDB)
		wantValue string
		wantErr   bool
	}{
		{
			name: "successful get",
			key:  "test-key",
			setupMock: func(mockDB *mocks.MockDB) {
				mockDB.EXPECT().View(gomock.Any()).DoAndReturn(func(fn func(*bbolt.Tx) error) error {
					// Simulate successful view operation
					return nil
				})
			},
			wantValue: "",
			wantErr:   false,
		},
		{
			name: "database view error",
			key:  "test-key",
			setupMock: func(mockDB *mocks.MockDB) {
				mockDB.EXPECT().View(gomock.Any()).Return(errors.New("database error"))
			},
			wantValue: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mocks.NewMockDB(ctrl)
			tt.setupMock(mockDB)

			service := libs.NewService(mockDB)
			value, err := service.Get(tt.key)

			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if value != tt.wantValue {
				t.Errorf("Service.Get() value = %v, want %v", value, tt.wantValue)
			}
		})
	}
}
