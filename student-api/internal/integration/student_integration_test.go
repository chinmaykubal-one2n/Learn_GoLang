package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	handler "student-api/internal/handler"
	logging "student-api/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

func setupLoggerForIntegrationTests() {
	zapLogger := zap.NewExample()
	logging.Logger = otelzap.New(zapLogger)
}

func TestHealthCheckIntegration(t *testing.T) {
	setupLoggerForIntegrationTests()

	tests := []struct {
		name           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "health check returns 200 and ok",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"ok","timestamp":"`,
		},
	}

	gin.SetMode(gin.TestMode)

	h := handler.NewHandler(nil)

	router := gin.Default()
	router.GET("/healthz", h.HealthCheck)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tc.expectedStatus, resp.Code)

			assert.Contains(t, resp.Body.String(), tc.expectedBody)
		})
	}
}
