package main

import (
	"net/http"
	"strings"

	"github.com/laoggu/cabbage-vehicle-backend/internal/pkg/jwt"
	"github.com/laoggu/cabbage-vehicle-backend/internal/pkg/tenant"
)

func withIdempotency(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. 幂等 key 校验
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			if r.Header.Get("X-Idempotency-Key") == "" {
				http.Error(w, "missing X-Idempotency-Key", http.StatusBadRequest)
				return
			}
		}

		// 2. 从 JWT 提取 openID 作为 tenant_id 注入 ctx
		auth := r.Header.Get("Authorization")
		if strings.HasPrefix(auth, "Bearer ") {
			token := strings.TrimPrefix(auth, "Bearer ")
			if cl, _ := jwt.Parse(token); cl != nil {
				r = r.WithContext(tenant.WithID(r.Context(), cl.Sub))
			}
		}

		next.ServeHTTP(w, r)
	})
}
