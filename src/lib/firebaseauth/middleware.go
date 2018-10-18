package firebaseauth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aikizoku/beego/src/lib/log"
	"github.com/unrolled/render"
	"google.golang.org/appengine"
)

// Middleware ... JSONRPC2に準拠したミドルウェア
type Middleware struct {
	Svc Service
}

// Handle ... Firebase認証をする
func (m *Middleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)

		userID, claims, err := m.Svc.Authentication(ctx, r)
		if err != nil {
			m.renderError(ctx, w, http.StatusForbidden, "authentication: "+err.Error())
			return
		}

		ctx = setUserID(ctx, userID)
		log.Debugf(ctx, "UserID: %s", userID)

		ctx = setClaims(ctx, claims)
		log.Debugf(ctx, "Claims: %v", claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) renderError(ctx context.Context, w http.ResponseWriter, status int, msg string) {
	log.Warningf(ctx, msg)
	render.New().Text(w, status, fmt.Sprintf("%d %s", status, msg))
}

// NewMiddleware ... Middlewareを作成する
func NewMiddleware(svc Service) *Middleware {
	return &Middleware{
		Svc: svc,
	}
}
