package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"goMcp/internal/errors"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRequest 测试请求结构体
type TestRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Age   int    `json:"age" validate:"min=1,max=120"`
}

func setupBaseController() *BaseController {
	gin.SetMode(gin.TestMode)
	return NewBaseController()
}

func TestBaseController_NewBaseController(t *testing.T) {
	controller := NewBaseController()
	assert.NotNil(t, controller)
}

func TestBaseController_BindAndValidate(t *testing.T) {
	controller := setupBaseController()

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedError  bool
		setupContext   func(*gin.Context)
		validateResult func(*testing.T, *TestRequest)
	}{
		{
			name: "Success_ValidJSON",
			requestBody: TestRequest{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   30,
			},
			expectedError: false,
			setupContext:  func(c *gin.Context) {},
			validateResult: func(t *testing.T, req *TestRequest) {
				assert.Equal(t, "John Doe", req.Name)
				assert.Equal(t, "john@example.com", req.Email)
				assert.Equal(t, 30, req.Age)
			},
		},
		{
			name:          "Error_InvalidJSON",
			requestBody:   `{"invalid": json}`,
			expectedError: true,
			setupContext:  func(c *gin.Context) {},
			validateResult: func(t *testing.T, req *TestRequest) {
				// 不验证结果，因为应该有错误
			},
		},
		{
			name: "Success_WithValidatedData",
			requestBody: TestRequest{
				Name:  "Jane Doe",
				Email: "jane@example.com",
				Age:   25,
			},
			expectedError: false,
			setupContext: func(c *gin.Context) {
				validatedData := TestRequest{
					Name:  "Jane Doe",
					Email: "jane@example.com",
					Age:   25,
				}
				c.Set("validated_data", validatedData)
			},
			validateResult: func(t *testing.T, req *TestRequest) {
				assert.Equal(t, "Jane Doe", req.Name)
				assert.Equal(t, "jane@example.com", req.Email)
				assert.Equal(t, 25, req.Age)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// 设置请求体
			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				var err error
				body, err = json.Marshal(tt.requestBody)
				require.NoError(t, err)
			}

			c.Request = httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// 设置上下文
			tt.setupContext(c)

			var req TestRequest
			err := controller.BindAndValidate(c, &req)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				tt.validateResult(t, &req)
			}
		})
	}
}

func TestBaseController_HandleValidationError(t *testing.T) {
	controller := setupBaseController()

	tests := []struct {
		name               string
		inputError         error
		expectedStatusCode int
	}{
		{
			name:               "OtherError",
			inputError:         fmt.Errorf("some other error"),
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			controller.HandleValidationError(c, tt.inputError)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
		})
	}
}

func TestBaseController_CopyValidatedData(t *testing.T) {
	controller := setupBaseController()

	tests := []struct {
		name          string
		source        interface{}
		target        interface{}
		expectedError bool
		validateResult func(*testing.T, interface{})
	}{
		{
			name:   "Success_SameType",
			source: TestRequest{Name: "John", Email: "john@example.com", Age: 30},
			target: &TestRequest{},
			expectedError: false,
			validateResult: func(t *testing.T, target interface{}) {
				req := target.(*TestRequest)
				assert.Equal(t, "John", req.Name)
				assert.Equal(t, "john@example.com", req.Email)
				assert.Equal(t, 30, req.Age)
			},
		},
		{
			name:   "Success_PointerSource",
			source: &TestRequest{Name: "Jane", Email: "jane@example.com", Age: 25},
			target: &TestRequest{},
			expectedError: false,
			validateResult: func(t *testing.T, target interface{}) {
				req := target.(*TestRequest)
				assert.Equal(t, "Jane", req.Name)
				assert.Equal(t, "jane@example.com", req.Email)
				assert.Equal(t, 25, req.Age)
			},
		},
		{
			name:          "Error_TargetNotPointer",
			source:        TestRequest{Name: "John", Email: "john@example.com", Age: 30},
			target:        TestRequest{},
			expectedError: true,
			validateResult: func(t *testing.T, target interface{}) {},
		},
		{
			name:          "Error_TypeMismatch",
			source:        "string value",
			target:        &TestRequest{},
			expectedError: true,
			validateResult: func(t *testing.T, target interface{}) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := controller.CopyValidatedData(tt.source, tt.target)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				tt.validateResult(t, tt.target)
			}
		})
	}
}

func TestBaseController_ParseIDParam(t *testing.T) {
	controller := setupBaseController()

	tests := []struct {
		name        string
		paramValue  string
		paramName   string
		expectedID  int64
		expectedErr bool
	}{
		{
			name:        "Success_ValidID",
			paramValue:  "123",
			paramName:   "id",
			expectedID:  123,
			expectedErr: false,
		},
		{
			name:        "Success_ZeroID",
			paramValue:  "0",
			paramName:   "id",
			expectedID:  0,
			expectedErr: false,
		},
		{
			name:        "Error_InvalidID",
			paramValue:  "abc",
			paramName:   "id",
			expectedID:  0,
			expectedErr: true,
		},
		{
			name:        "Error_EmptyID",
			paramValue:  "",
			paramName:   "id",
			expectedID:  0,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{
				{Key: tt.paramName, Value: tt.paramValue},
			}

			id, err := controller.ParseIDParam(c, tt.paramName)

			if tt.expectedErr {
				assert.Error(t, err)
				var appErr *errors.AppError
				assert.ErrorAs(t, err, &appErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, id)
			}
		})
	}
}

func TestBaseController_ParsePaginationParams(t *testing.T) {
	controller := setupBaseController()

	tests := []struct {
		name           string
		queryParams    map[string]string
		expectedPage   int64
		expectedLimit  int64
		expectedOffset int64
		expectedErr    bool
	}{
		{
			name:           "Success_DefaultValues",
			queryParams:    map[string]string{},
			expectedPage:   1,
			expectedLimit:  10,
			expectedOffset: 0,
			expectedErr:    false,
		},
		{
			name: "Success_CustomValues",
			queryParams: map[string]string{
				"page":  "2",
				"limit": "20",
			},
			expectedPage:   2,
			expectedLimit:  20,
			expectedOffset: 20,
			expectedErr:    false,
		},
		{
			name: "Error_InvalidPage",
			queryParams: map[string]string{
				"page": "0",
			},
			expectedErr: true,
		},
		{
			name: "Error_InvalidLimit",
			queryParams: map[string]string{
				"limit": "0",
			},
			expectedErr: true,
		},
		{
			name: "Error_LimitTooLarge",
			queryParams: map[string]string{
				"limit": "101",
			},
			expectedErr: true,
		},
		{
			name: "Error_NonNumericPage",
			queryParams: map[string]string{
				"page": "abc",
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// 构建查询字符串
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			q := req.URL.Query()
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()
			c.Request = req

			page, limit, offset, err := controller.ParsePaginationParams(c)

			if tt.expectedErr {
				assert.Error(t, err)
				var appErr *errors.AppError
				assert.ErrorAs(t, err, &appErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPage, page)
				assert.Equal(t, tt.expectedLimit, limit)
				assert.Equal(t, tt.expectedOffset, offset)
			}
		})
	}
}

func TestBaseController_HandleError(t *testing.T) {
	controller := setupBaseController()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testErr := fmt.Errorf("test error")
	controller.HandleError(c, testErr)

	assert.Len(t, c.Errors, 1)
	assert.Equal(t, testErr, c.Errors[0].Err)
}