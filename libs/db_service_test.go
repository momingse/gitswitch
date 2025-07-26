package libs_test

import (
	"errors"
	"gs/libs"
	mocks "gs/mocks/libs"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestDBService_Add(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		value     string
		setupMock func(*mocks.MockDB, *mocks.MockTx, *mocks.MockBucket)
		wantErr   bool
	}{
		{
			name:  "successful add",
			key:   "test-key",
			value: "test-value",
			setupMock: func(mockDB *mocks.MockDB, mockTx *mocks.MockTx, mockBucket *mocks.MockBucket) {
				// Mock the Update call - this is where the magic happens
				mockDB.EXPECT().Update(gomock.Any()).DoAndReturn(func(fn func(libs.Tx) error) error {
					// Create a mock transaction and simulate the function call
					// The fn parameter is the anonymous function passed to Update()
					// We call it with our mock transaction to simulate real behavior
					return fn(mockTx)
				})

				// Set expectations for what happens inside the transaction function
				mockTx.EXPECT().Bucket([]byte("test-bucket")).Return(mockBucket)
				mockBucket.EXPECT().Put([]byte("test-key"), []byte("test-value")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "database update error",
			key:   "test-key",
			value: "test-value",
			setupMock: func(mockDB *mocks.MockDB, mockTx *mocks.MockTx, mockBucket *mocks.MockBucket) {
				mockDB.EXPECT().Update(gomock.Any()).Return(errors.New("database error"))
				// No need to set up tx/bucket expectations since Update fails immediately
			},
			wantErr: true,
		},
		{
			name:  "bucket not found",
			key:   "test-key",
			value: "test-value",
			setupMock: func(mockDB *mocks.MockDB, mockTx *mocks.MockTx, mockBucket *mocks.MockBucket) {
				mockDB.EXPECT().Update(gomock.Any()).DoAndReturn(func(fn func(libs.Tx) error) error {
					return fn(mockTx)
				})
				mockTx.EXPECT().Bucket([]byte("test-bucket")).Return(nil) // Bucket not found
			},
			wantErr: true,
		},
		{
			name:  "bucket put error",
			key:   "test-key",
			value: "test-value",
			setupMock: func(mockDB *mocks.MockDB, mockTx *mocks.MockTx, mockBucket *mocks.MockBucket) {
				mockDB.EXPECT().Update(gomock.Any()).DoAndReturn(func(fn func(libs.Tx) error) error {
					return fn(mockTx)
				})
				mockTx.EXPECT().Bucket([]byte("test-bucket")).Return(mockBucket)
				mockBucket.EXPECT().Put([]byte("test-key"), []byte("test-value")).Return(errors.New("put error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mocks.NewMockDB(ctrl)
			mockTx := mocks.NewMockTx(ctrl)
			mockBucket := mocks.NewMockBucket(ctrl)

			tt.setupMock(mockDB, mockTx, mockBucket)

			service := libs.NewDBService(mockDB, "test-bucket")
			err := service.Add(tt.key, tt.value)

			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDBService_Get(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		setupMock func(*mocks.MockDB, *mocks.MockTx, *mocks.MockBucket)
		wantValue string
		wantErr   bool
	}{
		{
			name: "successful get",
			key:  "test-key",
			setupMock: func(mockDB *mocks.MockDB, mockTx *mocks.MockTx, mockBucket *mocks.MockBucket) {
				mockDB.EXPECT().View(gomock.Any()).DoAndReturn(func(fn func(libs.Tx) error) error {
					// Simulate calling the view function with our mock transaction
					return fn(mockTx)
				})
				mockTx.EXPECT().Bucket([]byte("test-bucket")).Return(mockBucket)
				mockBucket.EXPECT().Get([]byte("test-key")).Return([]byte("test-value"))
			},
			wantValue: "test-value",
			wantErr:   false,
		},
		{
			name: "key not found",
			key:  "nonexistent-key",
			setupMock: func(mockDB *mocks.MockDB, mockTx *mocks.MockTx, mockBucket *mocks.MockBucket) {
				mockDB.EXPECT().View(gomock.Any()).DoAndReturn(func(fn func(libs.Tx) error) error {
					return fn(mockTx)
				})
				mockTx.EXPECT().Bucket([]byte("test-bucket")).Return(mockBucket)
				mockBucket.EXPECT().Get([]byte("nonexistent-key")).Return(nil)
			},
			wantValue: "",
			wantErr:   false,
		},
		{
			name: "bucket not found",
			key:  "test-key",
			setupMock: func(mockDB *mocks.MockDB, mockTx *mocks.MockTx, mockBucket *mocks.MockBucket) {
				mockDB.EXPECT().View(gomock.Any()).DoAndReturn(func(fn func(libs.Tx) error) error {
					return fn(mockTx)
				})
				mockTx.EXPECT().Bucket([]byte("test-bucket")).Return(nil)
			},
			wantValue: "",
			wantErr:   true,
		},
		{
			name: "database view error",
			key:  "test-key",
			setupMock: func(mockDB *mocks.MockDB, mockTx *mocks.MockTx, mockBucket *mocks.MockBucket) {
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
			mockTx := mocks.NewMockTx(ctrl)
			mockBucket := mocks.NewMockBucket(ctrl)

			tt.setupMock(mockDB, mockTx, mockBucket)

			service := libs.NewDBService(mockDB, "test-bucket")
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
