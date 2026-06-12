package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"chaseapp.tv/api/internal/config"
)

// MetaHandler handles Meta Marketing API endpoints.
type MetaHandler struct {
	cfg    config.MetaConfig
	client *http.Client
	logger *slog.Logger
}

// NewMetaHandler creates a new MetaHandler.
func NewMetaHandler(cfg config.MetaConfig, logger *slog.Logger) *MetaHandler {
	return &MetaHandler{
		cfg:    cfg,
		client: http.DefaultClient,
		logger: logger,
	}
}

type metaAdAccount struct {
	ID            string `json:"id"`
	AccountID     string `json:"account_id,omitempty"`
	Name          string `json:"name"`
	AccountStatus int    `json:"account_status,omitempty"`
	Currency      string `json:"currency,omitempty"`
	TimezoneName  string `json:"timezone_name,omitempty"`
}

type metaBusiness struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type metaAdAccountCollection struct {
	Data   []metaAdAccount `json:"data"`
	Paging struct {
		Next string `json:"next"`
	} `json:"paging"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    int    `json:"code"`
	} `json:"error"`
}

type metaBusinessCollection struct {
	Data   []metaBusiness `json:"data"`
	Paging struct {
		Next string `json:"next"`
	} `json:"paging"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    int    `json:"code"`
	} `json:"error"`
}

// ListAdAccounts returns agency-accessible Meta ad accounts.
// GET /api/v1/meta/ad-accounts
func (h *MetaHandler) ListAdAccounts(w http.ResponseWriter, r *http.Request) {
	if h.cfg.SystemUserAccessToken == "" {
		Error(w, http.StatusServiceUnavailable, "Meta agency connection is not configured")
		return
	}

	fields := "fields=id,account_id,name,account_status,currency,timezone_name&limit=100"
	accounts, err := h.fetchAdAccounts(r.Context(), "me/adaccounts?"+fields)
	if err != nil {
		h.logger.Error("meta ad account fetch failed", slog.Any("error", err))
		Error(w, http.StatusBadGateway, "Failed to load Meta ad accounts")
		return
	}

	if h.cfg.SystemUserID != "" {
		assigned, err := h.fetchAdAccounts(r.Context(), fmt.Sprintf("%s/assigned_ad_accounts?%s", h.cfg.SystemUserID, fields))
		if err != nil {
			h.logger.Warn("meta assigned ad account fetch failed", slog.String("system_user_id", h.cfg.SystemUserID), slog.Any("error", err))
		}
		accounts = append(accounts, assigned...)
	}

	businesses := []metaBusiness{}
	if h.cfg.AgencyBusinessID != "" {
		businesses = append(businesses, metaBusiness{ID: h.cfg.AgencyBusinessID})
	} else {
		businesses, err = h.fetchBusinesses(r.Context(), "me/businesses?fields=id,name&limit=100")
		if err != nil {
			h.logger.Error("meta business fetch failed", slog.Any("error", err))
			Error(w, http.StatusBadGateway, "Failed to load Meta businesses")
			return
		}
	}

	for _, business := range businesses {
		owned, err := h.fetchAdAccounts(r.Context(), fmt.Sprintf("%s/owned_ad_accounts?%s", business.ID, fields))
		if err != nil {
			h.logger.Warn("meta owned ad account fetch failed", slog.String("business_id", business.ID), slog.Any("error", err))
		}
		accounts = append(accounts, owned...)

		client, err := h.fetchAdAccounts(r.Context(), fmt.Sprintf("%s/client_ad_accounts?%s", business.ID, fields))
		if err != nil {
			h.logger.Warn("meta client ad account fetch failed", slog.String("business_id", business.ID), slog.Any("error", err))
		}
		accounts = append(accounts, client...)
	}

	JSON(w, http.StatusOK, map[string]any{
		"data": dedupeMetaAdAccounts(accounts),
	})
}

func (h *MetaHandler) fetchAdAccounts(ctx context.Context, path string) ([]metaAdAccount, error) {
	items := []metaAdAccount{}
	nextURL := h.metaURL(path)

	for nextURL != "" {
		var payload metaAdAccountCollection
		err := h.fetchMetaPage(ctx, nextURL, &payload)
		if err != nil {
			return nil, err
		}

		items = append(items, payload.Data...)
		nextURL = payload.Paging.Next
	}

	return items, nil
}

func (h *MetaHandler) fetchBusinesses(ctx context.Context, path string) ([]metaBusiness, error) {
	items := []metaBusiness{}
	nextURL := h.metaURL(path)

	for nextURL != "" {
		var payload metaBusinessCollection
		err := h.fetchMetaPage(ctx, nextURL, &payload)
		if err != nil {
			return nil, err
		}

		items = append(items, payload.Data...)
		nextURL = payload.Paging.Next
	}

	return items, nil
}

func (h *MetaHandler) fetchMetaPage(ctx context.Context, pageURL string, payload any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pageURL, nil)
	if err != nil {
		return err
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(payload); err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("meta api status: %d", resp.StatusCode)
	}

	return nil
}

func (h *MetaHandler) metaURL(path string) string {
	base := strings.TrimRight(h.cfg.GraphBaseURL, "/")
	version := strings.Trim(h.cfg.GraphVersion, "/")
	path = strings.TrimLeft(path, "/")
	separator := "?"
	if strings.Contains(path, "?") {
		separator = "&"
	}

	return fmt.Sprintf("%s/%s/%s%saccess_token=%s", base, version, path, separator, url.QueryEscape(h.cfg.SystemUserAccessToken))
}

func dedupeMetaAdAccounts(accounts []metaAdAccount) []metaAdAccount {
	seen := map[string]metaAdAccount{}
	for _, account := range accounts {
		if account.ID == "" {
			continue
		}
		seen[account.ID] = account
	}

	out := make([]metaAdAccount, 0, len(seen))
	for _, account := range seen {
		out = append(out, account)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})
	return out
}
