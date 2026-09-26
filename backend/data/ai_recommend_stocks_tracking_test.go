package data

import (
	"math"
	"testing"

	"go-stock/backend/models"
)

func TestParseRecommendationPrices(t *testing.T) {
	got := parseRecommendationPrices("¥ 10.50 - 12.80 元")
	if len(got) != 2 || got[0] != 10.5 || got[1] != 12.8 {
		t.Fatalf("unexpected parsed prices: %#v", got)
	}
}

func TestApplyRecommendationTracking(t *testing.T) {
	tests := []struct {
		name    string
		current string
		want    string
		label   string
	}{
		{name: "waiting entry", current: "90", want: recommendTrackingWatching, label: "等待入场"},
		{name: "entry reached", current: "102", want: recommendTrackingEntryReached, label: "到达买入区"},
		{name: "tracking above entry", current: "110", want: recommendTrackingTracking, label: "跟踪中"},
		{name: "take profit reached", current: "121", want: recommendTrackingTakeProfitReached, label: "已到止盈"},
		{name: "stop loss reached", current: "79", want: recommendTrackingStopLossReached, label: "已到止损"},
		{name: "missing price", current: "", want: recommendTrackingWatching, label: "观察中"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := &models.AiRecommendStocks{
				StockCurrentPrice:        tt.current,
				RecommendBuyPrice:        "100-105",
				RecommendStopProfitPrice: "120",
				RecommendStopLossPrice:   "80",
			}
			applyRecommendationTracking(item)
			if item.TrackingState != tt.want {
				t.Fatalf("state=%q, want %q", item.TrackingState, tt.want)
			}
			if item.TrackingLabel != tt.label {
				t.Fatalf("label=%q, want %q", item.TrackingLabel, tt.label)
			}
		})
	}
}

func TestApplyRecommendationTrackingDistances(t *testing.T) {
	item := &models.AiRecommendStocks{
		StockCurrentPrice:             "100",
		RecommendBuyPriceMin:         95,
		RecommendBuyPriceMax:         98,
		RecommendStopProfitPriceMin:  120,
		RecommendStopLossPrice:       "90",
	}
	applyRecommendationTracking(item)

	if !item.HasTakeProfitTarget || !item.HasStopLossTarget {
		t.Fatalf("expected both targets to be present: %#v", item)
	}
	if math.Abs(item.TakeProfitDistancePct-20) > 0.0001 {
		t.Fatalf("take-profit distance=%f, want 20", item.TakeProfitDistancePct)
	}
	if math.Abs(item.StopLossDistancePct-(-10)) > 0.0001 {
		t.Fatalf("stop-loss distance=%f, want -10", item.StopLossDistancePct)
	}
	if item.TrackingState != recommendTrackingTracking {
		t.Fatalf("state=%q, want %q", item.TrackingState, recommendTrackingTracking)
	}
}
