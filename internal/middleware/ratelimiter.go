package middleware

import (
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"github.com/kisalto/Feel-The-Night-Swagger/internal/dto"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type IPRateLimiter struct {
	mu      sync.Mutex
	clients map[string]*client
	r       rate.Limit
	burst   int
}

func NewIpRateLimiter(r rate.Limit, burst int) *IPRateLimiter {
	l := &IPRateLimiter{
		clients: make(map[string]*client),
		r:       r,
		burst:   burst,
	}
	go l.cleanup(3 * time.Minute)
	return l
}

func (l *IPRateLimiter) get(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	c, ok := l.clients[ip]
	if !ok {
		c = &client{limiter: rate.NewLimiter(l.r, l.burst)}
		l.clients[ip] = c
	}
	c.lastSeen = time.Now()
	return c.limiter
}

func (l *IPRateLimiter) cleanup(maxIdle time.Duration) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		l.mu.Lock()
		for ip, c := range l.clients {
			if time.Since(c.lastSeen) > maxIdle {
				delete(l.clients, ip)
			}
		}
		l.mu.Unlock()
	}
}

func (l *IPRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		lim := l.get(c.ClientIP())

		res := lim.Reserve()
		if delay := res.Delay(); delay > 0 {
			res.Cancel()
			c.Header("Retry-After", strconv.Itoa(int(math.Ceil(delay.Seconds()))))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, dto.ErrorResponse{
				Error: "Muitas tentativas, tente novamente mais tarde",
			})
			return
		}
		c.Next()
	}
}
