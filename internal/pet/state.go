package pet

import (
	"math"
	"time"
)

// ========== 衰减速率配置（每分钟扣减量） ==========
// 所有值都是"每分钟扣减多少点"；上限100，下限0
const (
	DecayHungerPerMin   = 0.35 // 饱腹度：1.5小时左右从100降到约70
	DecayHappinessPerMin = 0.18 // 好感度衰减较慢
	DecayCleanPerMin    = 0.25 // 卫生值
	DecayEnergyPerMin   = 0.22 // 精力值

	// 离线补偿最大窗口（防止离线一年瞬间数值全0）
	MaxOfflineCompensateHours = 12

	// 单条记录的最大careScore衰减上限（保护上限值不被快速归零）
	CareScoreDecayPerDay = 15.0
)

// NewBorn 创建一只新生幼崽
func NewBorn(name string) *PetState {
	now := time.Now()
	return &PetState{
		Name:        name,
		Hunger:      80,
		Happiness:   70,
		Hygiene:     90,
		Energy:      90,
		Stage:       StageBaby,
		CareScore:   0,
		CareStreak:  0,
		BirthTime:   now,
		LastUpdate:  now,
		LastCareDay: now.Format("2006-01-02"),
		History:     []InteractionRecord{},
	}
}

// clamp 将数值约束在 [0,100]
func clamp(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 100 {
		return 100
	}
	return v
}

// Tick 推进时间：根据 lastUpdate 和 now 的差值对状态进行衰减，并检查成长阶段
// 注意：此函数会修改 LastUpdate 为 now
func (p *PetState) Tick(now time.Time) {
	if now.Before(p.LastUpdate) {
		// 时钟回拨，保护
		p.LastUpdate = now
		return
	}
	delta := now.Sub(p.LastUpdate)
	minutes := delta.Minutes()

	// 限制最大离线补偿窗口
	maxMinutes := float64(MaxOfflineCompensateHours) * 60.0
	if minutes > maxMinutes {
		minutes = maxMinutes
	}

	if minutes <= 0 {
		return
	}

	p.Hunger = clamp(p.Hunger - DecayHungerPerMin*minutes)
	p.Happiness = clamp(p.Happiness - DecayHappinessPerMin*minutes)
	p.Hygiene = clamp(p.Hygiene - DecayCleanPerMin*minutes)
	p.Energy = clamp(p.Energy - DecayEnergyPerMin*minutes)

	// 年龄累计
	p.AgeHours += delta.Hours()
	if p.AgeHours < 0 {
		p.AgeHours = 0
	}

	// 照顾分随日数衰减（按天算）
	days := delta.Hours() / 24.0
	p.CareScore = math.Max(0, p.CareScore-CareScoreDecayPerDay*days)

	p.checkGrowth()
	p.LastUpdate = now
}

// ========== 成长系统 ==========
// 阶段晋升门槛（年龄小时 + 照顾质量分 + 互动次数）
// 这样保证"必须真的好好照顾才能长大"，而不是纯挂机
var stageThresholds = []struct {
	MinAgeHours   float64
	MinCareScore  float64
	MinCareEvents int
}{
	// StageBaby -> StageChild
	{MinAgeHours: 6, MinCareScore: 30, MinCareEvents: 5},
	// StageChild -> StageTeen
	{MinAgeHours: 24, MinCareScore: 120, MinCareEvents: 20},
	// StageTeen -> StageAdult
	{MinAgeHours: 72, MinCareScore: 350, MinCareEvents: 60},
}

func (p *PetState) checkGrowth() {
	totalEvents := len(p.History)
	for i := len(stageThresholds) - 1; i >= 0; i-- {
		th := stageThresholds[i]
		targetStage := GrowthStage(i + 1)
		if p.Stage < targetStage &&
			p.AgeHours >= th.MinAgeHours &&
			p.CareScore >= th.MinCareScore &&
			totalEvents >= th.MinCareEvents {
			p.Stage = targetStage
			// 晋升奖励
			p.Happiness = clamp(p.Happiness + 15)
			break
		}
	}
}

// ========== 互动质量评估 ==========
// 用于给成长系统提供可解释的质量标签
func (p *PetState) qualityOfFeed(f FoodType) string {
	// 饱腹度高于90还喂就是"不好"（吃撑了），低于30喂且吃健康食物=很好
	if p.Hunger >= 90 {
		return "bad"
	}
	if f.Name == FoodHealthy.Name && p.Hunger < 40 {
		return "good"
	}
	if p.Hunger < 50 {
		return "good"
	}
	return "normal"
}

func (p *PetState) qualityOfClean() string {
	if p.Hygiene <= 25 {
		return "good" // 及时清洁
	}
	if p.Hygiene >= 85 {
		return "normal" // 本来就很干净
	}
	return "normal"
}

func (p *PetState) qualityOfPet() string {
	if p.Happiness <= 30 {
		return "good"
	}
	if p.Happiness >= 95 {
		return "normal"
	}
	return "normal"
}

// updateCareDaily 每天第一次互动重置计数，连续天数+1
func (p *PetState) updateCareDaily(now time.Time) {
	today := now.Format("2006-01-02")
	yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")

	if p.LastCareDay == today {
		p.DailyCareCount++
	} else if p.LastCareDay == yesterday {
		p.CareStreak++
		p.DailyCareCount = 1
		p.LastCareDay = today
	} else {
		// 断签
		p.CareStreak = 1
		p.DailyCareCount = 1
		p.LastCareDay = today
	}
}

func (p *PetState) recordInteraction(t, sub, quality string, now time.Time) {
	p.History = append(p.History, InteractionRecord{
		Type:      t,
		SubType:   sub,
		Quality:   quality,
		Timestamp: now,
	})
	// 只保留最近500条，避免文件无限大
	if len(p.History) > 500 {
		p.History = p.History[len(p.History)-500:]
	}
}

// ========== 对外互动 API ==========

// Feed 喂食，返回质量标签
func (p *PetState) Feed(f FoodType, now time.Time) string {
	p.Tick(now)
	quality := p.qualityOfFeed(f)

	p.Hunger = clamp(p.Hunger + f.Hunger)
	p.Happiness = clamp(p.Happiness + f.Happiness)
	p.Energy = clamp(p.Energy + f.Energy)

	// careScore 加成
	switch quality {
	case "good":
		p.CareScore += 6
	case "normal":
		p.CareScore += 3
	case "bad":
		p.CareScore += 0.5
		p.Happiness = clamp(p.Happiness - 3) // 吃撑了反而不开心
	}

	p.updateCareDaily(now)
	p.recordInteraction("feed", f.Name, quality, now)
	p.checkGrowth()
	p.LastUpdate = now
	return quality
}

// Pet 抚摸
func (p *PetState) Pet(now time.Time) string {
	p.Tick(now)
	quality := p.qualityOfPet()

	// 精力不足时抚摸可能被嫌弃
	if p.Energy < 20 {
		p.Happiness = clamp(p.Happiness + 1)
		quality = "normal"
	} else {
		p.Happiness = clamp(p.Happiness + 8)
		// 抚摸略微消耗精力
		p.Energy = clamp(p.Energy - 1)
	}

	switch quality {
	case "good":
		p.CareScore += 5
	case "normal":
		p.CareScore += 2
	}

	p.updateCareDaily(now)
	p.recordInteraction("pet", "", quality, now)
	p.checkGrowth()
	p.LastUpdate = now
	return quality
}

// Clean 清洁
func (p *PetState) Clean(now time.Time) string {
	p.Tick(now)
	quality := p.qualityOfClean()

	cleanGain := 50.0
	if p.Stage >= StageAdult {
		cleanGain = 45 // 成犬/成猫体量更大，恢复略少
	}
	p.Hygiene = clamp(p.Hygiene + cleanGain)
	// 清洁过程略微消耗精力
	p.Energy = clamp(p.Energy - 3)
	// 有些宠物不喜欢洗澡，好感度小幅变化
	if p.Hygiene < 20 {
		p.Happiness = clamp(p.Happiness + 4) // 本来很脏，洗完舒服
	} else {
		p.Happiness = clamp(p.Happiness - 2) // 不想洗硬被洗
	}

	switch quality {
	case "good":
		p.CareScore += 6
	case "normal":
		p.CareScore += 2.5
	}

	p.updateCareDaily(now)
	p.recordInteraction("clean", "", quality, now)
	p.checkGrowth()
	p.LastUpdate = now
	return quality
}

// ========== 派生属性：心情 ==========

// Mood 根据当前四项核心状态给出综合心情
func (p *PetState) Mood() MoodState {
	avg := (p.Hunger + p.Happiness + p.Hygiene + p.Energy) / 4.0
	switch {
	case avg >= 75:
		return MoodHappy
	case avg >= 50:
		return MoodNormal
	case avg >= 25:
		return MoodSad
	default:
		return MoodSick
	}
}

// StatusSummary 给UI显示用的简短描述
func (p *PetState) StatusSummary() string {
	mood := p.Mood()
	var parts []string
	if p.Hunger < 30 {
		parts = append(parts, "好饿…")
	} else if p.Hunger >= 90 {
		parts = append(parts, "吃撑了")
	}
	if p.Hygiene < 30 {
		parts = append(parts, "脏脏的")
	}
	if p.Energy < 25 {
		parts = append(parts, "困困")
	}
	if p.Happiness < 30 {
		parts = append(parts, "想你摸摸")
	}
	if len(parts) == 0 {
		switch mood {
		case MoodHappy:
			parts = append(parts, "超级开心！")
		case MoodNormal:
			parts = append(parts, "嗯，还行~")
		case MoodSad:
			parts = append(parts, "有点蔫儿…")
		case MoodSick:
			parts = append(parts, "难受…快照顾我！")
		}
	}
	result := ""
	for i, s := range parts {
		if i > 0 {
			result += "，"
		}
		result += s
	}
	return result
}
