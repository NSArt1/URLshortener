package handler

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/you/url-shortener/internal/model"
    "github.com/you/url-shortener/internal/storage"
    "github.com/you/url-shortener/internal/util"
    "github.com/prometheus/client_golang/prometheus"
)

var (
    reqCounter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "url_shortener_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"path", "method", "status"},
    )
)

func init() {
    prometheus.MustRegister(reqCounter)
}

type Handler struct {
    store *storage.Store
}

func New(db *storage.Store) *Handler {
    return &Handler{store: db}
}

type shortenReq struct {
    URL string `json:"url" binding:"required,url"`
}

func (h *Handler) Shorten(c *gin.Context) {
    var req shortenReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    code := util.GenerateCode(8)
    u := &model.URL{Code: code, OriginalURL: req.URL}
    if err := h.store.Create(u); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot save URL"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"code": code, "short_url": c.Request.Host + "/" + code})
    reqCounter.WithLabelValues("/shorten", "POST", "200").Inc()
}

func (h *Handler) Redirect(c *gin.Context) {
    code := c.Param("code")
    u, err := h.store.FindByCode(code)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        reqCounter.WithLabelValues("/"+code, "GET", "404").Inc()
        return
    }
    // инкремент счётчика
    _ = h.store.IncrementClicks(code)
    reqCounter.WithLabelValues("/"+code, "GET", "302").Inc()
    c.Redirect(http.StatusFound, u.OriginalURL)
}

func (h *Handler) Stats(c *gin.Context) {
    code := c.Param("code")
    u, err := h.store.FindByCode(code)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "code":         u.Code,
        "original_url": u.OriginalURL,
        "clicks":       u.Clicks,
        "created_at":   u.CreatedAt.Format(time.RFC3339),
    })
    reqCounter.WithLabelValues("/stats/"+code, "GET", "200").Inc()
}