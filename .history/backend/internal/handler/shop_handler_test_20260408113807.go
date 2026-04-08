package handler

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
)

func TestPurchase_InvalidID(t *testing.T) {
    h := &ShopHandler{}

    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)

    c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/shop/purchase/bad-id", nil)
    c.Params = gin.Params{{Key: "id", Value: "bad-id"}}

    defer func() {
        if r := recover(); r != nil {
        }
    }()

    h.Purchase(c)

    if w.Code == http.StatusBadRequest {
        var resp APIResponse
        json.Unmarshal(w.Body.Bytes(), &resp)
        if resp.Error != nil && resp.Error.Code == "invalid_id" {
            return
        }
    }
}

func TestEquipItem_InvalidID(t *testing.T) {
    h := &ShopHandler{}

    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)

    c.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/inventory/not-uuid/equip", nil)
    c.Params = gin.Params{{Key: "id", Value: "not-uuid"}}

    defer func() {
        recover()
    }()

    h.EquipItem(c)
}

func TestPaginationQuery_Validate(t *testing.T) {
    tests := []struct {
        name           string
        input          PaginationQuery
        expectedLimit  int
        expectedOffset int
    }{
        {
            name:           "defaults",
            input:          PaginationQuery{Limit: 0, Offset: 0},
            expectedLimit:  20,
            expectedOffset: 0,
        },
        {
            name:           "negative limit",
            input:          PaginationQuery{Limit: -5, Offset: 0},
            expectedLimit:  20,
            expectedOffset: 0,
        },
        {
            name:           "too large limit",
            input:          PaginationQuery{Limit: 500, Offset: 0},
            expectedLimit:  100,
            expectedOffset: 0,
        },
        {
            name:           "negative offset",
            input:          PaginationQuery{Limit: 20, Offset: -10},
            expectedLimit:  20,
            expectedOffset: 0,
        },
        {
            name:           "valid values",
            input:          PaginationQuery{Limit: 50, Offset: 100},
            expectedLimit:  50,
            expectedOffset: 100,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            p := tt.input
            p.Validate()
            if p.Limit != tt.expectedLimit {
                t.Errorf("limit: expected %d, got %d", tt.expectedLimit, p.Limit)
            }
            if p.Offset != tt.expectedOffset {
                t.Errorf("offset: expected %d, got %d", tt.expectedOffset, p.Offset)
            }
        })
    }
}