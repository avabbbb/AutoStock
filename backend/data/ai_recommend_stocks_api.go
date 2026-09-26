// Package data ai_recommend_stocks_api.go
package data

import (
	"go-stock/backend/db"
	"go-stock/backend/models"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/duke-git/lancet/v2/datetime"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/duke-git/lancet/v2/strutil"
)

type AiRecommendStocksService struct{}

func NewAiRecommendStocksService() *AiRecommendStocksService {
	return &AiRecommendStocksService{}
}

const (
	recommendTrackingWatching          = "watching"
	recommendTrackingEntryReached      = "entry_reached"
	recommendTrackingTracking          = "tracking"
	recommendTrackingTakeProfitReached = "take_profit_reached"
	recommendTrackingStopLossReached   = "stop_loss_reached"
)

var recommendationPriceNumberRE = regexp.MustCompile(`\d+(?:\.\d+)?`)

func parseRecommendationPrices(value string) []float64 {
	matches := recommendationPriceNumberRE.FindAllString(strings.ReplaceAll(value, ",", ""), -1)
	prices := make([]float64, 0, len(matches))
	for _, match := range matches {
		price, err := strconv.ParseFloat(match, 64)
		if err == nil && price > 0 {
			prices = append(prices, price)
		}
	}
	return prices
}

func recommendationEntryRange(item *models.AiRecommendStocks) (float64, float64, bool) {
	minPrice, maxPrice := item.RecommendBuyPriceMin, item.RecommendBuyPriceMax
	if minPrice > 0 && maxPrice > 0 {
		if minPrice > maxPrice {
			minPrice, maxPrice = maxPrice, minPrice
		}
		return minPrice, maxPrice, true
	}

	prices := parseRecommendationPrices(item.RecommendBuyPrice)
	if len(prices) == 0 {
		return 0, 0, false
	}
	if len(prices) == 1 {
		return prices[0], prices[0], true
	}
	minPrice, maxPrice = prices[0], prices[1]
	if minPrice > maxPrice {
		minPrice, maxPrice = maxPrice, minPrice
	}
	return minPrice, maxPrice, true
}

func recommendationTakeProfit(item *models.AiRecommendStocks) (float64, bool) {
	if item.RecommendStopProfitPriceMin > 0 {
		return item.RecommendStopProfitPriceMin, true
	}
	prices := parseRecommendationPrices(item.RecommendStopProfitPrice)
	if len(prices) == 0 {
		return 0, false
	}
	return prices[0], true
}

func recommendationStopLoss(item *models.AiRecommendStocks) (float64, bool) {
	prices := parseRecommendationPrices(item.RecommendStopLossPrice)
	if len(prices) == 0 {
		return 0, false
	}
	return prices[0], true
}

func applyRecommendationTracking(item *models.AiRecommendStocks) {
	if item == nil {
		return
	}

	current, err := strconv.ParseFloat(strings.TrimSpace(item.StockCurrentPrice), 64)
	if err != nil || current <= 0 {
		item.TrackingState = recommendTrackingWatching
		item.TrackingLabel = "观察中"
		return
	}

	entryMin, entryMax, hasEntry := recommendationEntryRange(item)
	takeProfit, hasTakeProfit := recommendationTakeProfit(item)
	stopLoss, hasStopLoss := recommendationStopLoss(item)

	item.HasTakeProfitTarget = hasTakeProfit
	item.HasStopLossTarget = hasStopLoss

	if hasTakeProfit {
		item.TakeProfitDistancePct = (takeProfit - current) / current * 100
	}
	if hasStopLoss {
		item.StopLossDistancePct = (stopLoss - current) / current * 100
	}

	switch {
	case hasStopLoss && current <= stopLoss:
		item.TrackingState = recommendTrackingStopLossReached
		item.TrackingLabel = "已到止损"
	case hasTakeProfit && current >= takeProfit:
		item.TrackingState = recommendTrackingTakeProfitReached
		item.TrackingLabel = "已到止盈"
	case hasEntry && current >= entryMin && current <= entryMax:
		item.TrackingState = recommendTrackingEntryReached
		item.TrackingLabel = "到达买入区"
	case hasEntry && current > entryMax:
		item.TrackingState = recommendTrackingTracking
		item.TrackingLabel = "跟踪中"
	default:
		item.TrackingState = recommendTrackingWatching
		item.TrackingLabel = "等待入场"
	}
}

// CreateAiRecommendStocks 创建AI推荐股票记录
func (s *AiRecommendStocksService) CreateAiRecommendStocks(recommend *models.AiRecommendStocks) error {
	result := db.Dao.Create(recommend)
	return result.Error
}

func (s *AiRecommendStocksService) BatchCreateAiRecommendStocks(recommends []*models.AiRecommendStocks) error {
	result := db.Dao.Create(recommends)
	return result.Error
}

// GetAiRecommendStocksList 分页查询AI推荐股票记录
func (s *AiRecommendStocksService) GetAiRecommendStocksList(query *models.AiRecommendStocksQuery) (*models.AiRecommendStocksPageData, error) {
	var list []models.AiRecommendStocks
	var total int64

	q := db.Dao.Model(&models.AiRecommendStocks{})

	// 构建关键词搜索条件（股票代码、股票名称、板块名称使用 OR 关系）
	keyword := query.StockCode
	if keyword == "" {
		keyword = query.StockName
	}
	if keyword == "" {
		keyword = query.BkName
	}
	if keyword == "" {
		keyword = query.ModelName
	}

	if keyword != "" {
		q = q.Where("(stock_code LIKE ? OR stock_name LIKE ? OR bk_name LIKE ? OR model_name LIKE ?)",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 日期范围查询
	if query.StartDate != "" && query.EndDate != "" {
		query.StartDate = strutil.ReplaceWithMap(query.StartDate, map[string]string{
			"T": " ",
			"Z": "",
		})
		query.EndDate = strutil.ReplaceWithMap(query.EndDate, map[string]string{
			"T": " ",
			"Z": "",
		})
		startDate, err := time.Parse("2006-01-02 15:04:05", query.StartDate)
		if err != nil {
			startDate, _ = time.Parse("2006-01-02", query.StartDate)
		}

		endDate, err := time.Parse("2006-01-02 15:04:05", query.EndDate)
		if err != nil {
			endDate, _ = time.Parse("2006-01-02", query.EndDate)
		}

		q = q.Where("data_time BETWEEN ? AND ?", datetime.BeginOfDay(startDate), datetime.EndOfDay(endDate))
	} else if query.StartDate == "" && query.EndDate == "" && keyword == "" {
		// 只有在没有关键词时才默认查询今天的数据
		q = q.Where("data_time BETWEEN ? AND ?", datetime.BeginOfDay(time.Now()), datetime.EndOfDay(time.Now()))
	} else if query.StartDate != "" && query.EndDate == "" {
		query.StartDate = strutil.ReplaceWithMap(query.StartDate, map[string]string{
			"T": " ",
			"Z": "",
		})
		startDate, _ := time.Parse("2006-01-02", query.StartDate)
		q = q.Where("data_time BETWEEN ? AND ?", datetime.BeginOfDay(startDate), datetime.EndOfDay(startDate))
	}

	// 预警状态筛选
	if query.EnableAlert != nil {
		q = q.Where("enable_alert = ?", *query.EnableAlert)
	}

	// 计算总数
	err := q.Count(&total).Error
	if err != nil {
		return nil, err
	}

	// 设置默认分页参数
	page := query.Page
	pageSize := query.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// 执行分页查询
	offset := (page - 1) * pageSize
	err = q.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&list).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	stockCodes := slice.Map(list, func(index int, item models.AiRecommendStocks) string {
		return ConvertTushareCodeToStockCode(item.StockCode)
	})
	stockData, _ := NewStockDataApi().GetStockCodeRealTimeData(stockCodes...)
	for _, info := range *stockData {
		for idx, item := range list {
			if ConvertTushareCodeToStockCode(item.StockCode) == ConvertTushareCodeToStockCode(info.Code) {
				list[idx].StockCurrentPrice = info.Price
				list[idx].StockPrePrice = info.PreClose
				list[idx].StockCurrentPriceTime = info.Date + " " + info.Time
			}
		}
	}
	for idx := range list {
		applyRecommendationTracking(&list[idx])
	}

	return &models.AiRecommendStocksPageData{
		List:       list,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetAiRecommendStocksTodayStats 统计截止指定日期（默认今天）近 days 天的推荐股池：按股票聚合推荐次数、评级、开仓价、目标价、止损价，并补充实时行情
// days <= 1 表示只统计该日期当天，排序按推荐次数倒序
func (s *AiRecommendStocksService) GetAiRecommendStocksTodayStats(date string, days int) (*models.AiRecommendStocksTodayStatsData, error) {
	var list []models.AiRecommendStocks

	day := time.Now()
	if d := strings.TrimSpace(date); d != "" {
		if parsed, err := time.Parse("2006-01-02", d); err == nil {
			day = parsed
		}
	}
	if days <= 0 {
		days = 1
	}
	startDay := day.AddDate(0, 0, -(days - 1))

	err := db.Dao.Model(&models.AiRecommendStocks{}).
		Where("data_time BETWEEN ? AND ?", datetime.BeginOfDay(startDay), datetime.EndOfDay(day)).
		Order("created_at ASC").
		Find(&list).Error
	if err != nil {
		return nil, err
	}

	result := &models.AiRecommendStocksTodayStatsData{
		Date:  day.Format("2006-01-02"),
		Items: make([]models.AiRecommendStocksTodayStat, 0),
	}
	if days > 1 {
		result.Date = startDay.Format("2006-01-02") + " ~ " + day.Format("2006-01-02")
	}
	// 多日区间下时间需带日期，否则无法区分是哪一天推荐的
	timeLayout := "15:04"
	if days > 1 {
		timeLayout = "01-02 15:04"
	}

	// 归一化股票代码（000001.SZ -> sz000001）后聚合
	type todayStatAgg struct {
		item     models.AiRecommendStocksTodayStat
		lastTime time.Time
	}
	aggs := make([]todayStatAgg, 0)
	index := make(map[string]int)
	modelSet := make(map[string]bool)
	for _, item := range list {
		code := ConvertTushareCodeToStockCode(item.StockCode)
		key := code
		if key == "" {
			key = item.StockName
		}
		idx, ok := index[key]
		if !ok {
			aggs = append(aggs, todayStatAgg{
				item: models.AiRecommendStocksTodayStat{
					StockCode:  code,
					StockName:  item.StockName,
					FirstTime:  item.CreatedAt.Format(timeLayout),
					ModelNames: []string{},
				},
			})
			idx = len(aggs) - 1
			index[key] = idx
		}
		stat := &aggs[idx].item
		stat.Count++
		// 按创建时间升序遍历，后写入的记录覆盖前值，最终保留区间内最后一次推荐的评级与价位
		stat.StockName = item.StockName
		stat.BkName = item.BkName
		stat.Rating = item.Rating
		stat.RecommendBuyPrice = item.RecommendBuyPrice
		stat.RecommendBuyPriceMin = item.RecommendBuyPriceMin
		stat.RecommendBuyPriceMax = item.RecommendBuyPriceMax
		stat.RecommendStopProfitPrice = item.RecommendStopProfitPrice
		stat.RecommendStopProfitPriceMin = item.RecommendStopProfitPriceMin
		stat.RecommendStopProfitPriceMax = item.RecommendStopProfitPriceMax
		stat.RecommendStopLossPrice = item.RecommendStopLossPrice
		stat.StockPrice = item.StockPrice
		stat.LastTime = item.CreatedAt.Format(timeLayout)
		aggs[idx].lastTime = item.CreatedAt
		if item.ModelName != "" {
			modelSet[item.ModelName] = true
			if !slice.Contain(stat.ModelNames, item.ModelName) {
				stat.ModelNames = append(stat.ModelNames, item.ModelName)
			}
		}
	}

	// 推荐次数多的优先，次数相同按最近推荐时间倒序
	sort.SliceStable(aggs, func(i, j int) bool {
		if aggs[i].item.Count != aggs[j].item.Count {
			return aggs[i].item.Count > aggs[j].item.Count
		}
		return aggs[i].lastTime.After(aggs[j].lastTime)
	})
	for _, agg := range aggs {
		result.Items = append(result.Items, agg.item)
	}

	// 补充实时行情
	if len(result.Items) > 0 {
		codes := slice.Map(result.Items, func(_ int, item models.AiRecommendStocksTodayStat) string {
			return item.StockCode
		})
		if stockData, err := NewStockDataApi().GetStockCodeRealTimeData(codes...); err == nil && stockData != nil {
			for _, info := range *stockData {
				infoCode := ConvertTushareCodeToStockCode(info.Code)
				for i := range result.Items {
					if result.Items[i].StockCode == infoCode {
						result.Items[i].StockCurrentPrice = info.Price
						result.Items[i].StockPrePrice = info.PreClose
						result.Items[i].StockCurrentPriceTime = info.Date + " " + info.Time
					}
				}
			}
		}
	}

	result.StockCount = len(result.Items)
	result.TotalCount = len(list)
	result.ModelCount = len(modelSet)

	return result, nil
}

// GetAiRecommendStocksByID 根据ID获取AI推荐股票记录
func (s *AiRecommendStocksService) GetAiRecommendStocksByID(id uint) (*models.AiRecommendStocks, error) {
	var recommend models.AiRecommendStocks
	err := db.Dao.First(&recommend, id).Error
	if err != nil {
		return nil, err
	}
	return &recommend, nil
}

// UpdateAiRecommendStocks 更新AI推荐股票记录
func (s *AiRecommendStocksService) UpdateAiRecommendStocks(id uint, recommend *models.AiRecommendStocks) error {
	result := db.Dao.Model(&models.AiRecommendStocks{}).Where("id = ?", id).Updates(recommend)
	return result.Error
}

// DeleteAiRecommendStocks 根据ID删除AI推荐股票记录
func (s *AiRecommendStocksService) DeleteAiRecommendStocks(id uint) error {
	// 使用软删除
	result := db.Dao.Where("id = ?", id).Delete(&models.AiRecommendStocks{})
	return result.Error
}

// UpdateAiRecommendStocksAlert 更新AI推荐股票的预警状态
func (s *AiRecommendStocksService) UpdateAiRecommendStocksAlert(id uint, enableAlert bool) error {
	result := db.Dao.Model(&models.AiRecommendStocks{}).Where("id = ?", id).Update("enable_alert", enableAlert)
	return result.Error
}

// BatchDeleteAiRecommendStocks 批量删除AI推荐股票记录
func (s *AiRecommendStocksService) BatchDeleteAiRecommendStocks(ids []uint) error {
	// 使用软删除
	result := db.Dao.Where("id IN ?", ids).Delete(&models.AiRecommendStocks{})
	return result.Error
}
