package handler

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
)

func TestGetBenefit_InvalidID(t *testing.T) {
    h := &BenefitHandler{}

    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)

    c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/benefits/not-a-uuid", nil)
    c.Params = gin.Params{{Key: "id", Value: "not-a-uuid"}}

    h.GetBenefit(c)

    if w.Code != http.StatusBadRequest {
        t.Errorf("expected 400, got %d", w.Code)
    }

    var resp APIResponse
    json.Unmarshal(w.Body.Bytes(), &resp)

    if resp.Error == nil || resp.Error.Code != "invalid_id" {
        t.Errorf("expected error code 'invalid_id', got %+v", resp.Error)
    }
}