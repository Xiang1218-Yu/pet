package pet

import "time"

// 成长阶段枚举
type GrowthStage int

const (
	StageBaby GrowthStage = iota // 0: 幼崽期
	StageChild                   // 1: 成长期
	StageTeen                    // 2: 少年期
	StageAdult                   // 3: 成年期
)

func (s GrowthStage) String() string {
	switch s {
	case StageBaby:
		return "幼崽"
	case StageChild:
		return "幼年"
	case StageTeen:
		return "少年"
	case StageAdult:
		return "成年"
	default:
		return "未知"
	}
}

// 心情状态枚举（由综合指标推导，用于选择外观表情）
type MoodState int

const (
	MoodHappy MoodState = iota // 开心
	MoodNormal                 // 普通
	MoodSad                    // 不开心
	MoodSick                   // 生病/状态很差
)

func (m MoodState) String() string {
	switch m {
	case MoodHappy:
		return "开心"
	case MoodNormal:
		return "普通"
	case MoodSad:
		return "不开心"
	case MoodSick:
		return "难受"
	default:
		return "未知"
	}
}

// 食物类型
type FoodType struct {
	Name     string  // 名称，如"小鱼干"
	Hunger   float64 // 饱腹度增加量
	Happiness float64 // 好感度增加量（好吃的食物额外加好感）
	Energy   float64 // 精力增加量（如营养食物加精力，零食可能减）
}

var (
	FoodFish    = FoodType{Name: "小鱼干", Hunger: 25, Happiness: 5, Energy: 8}
	FoodMilk    = FoodType{Name: "牛奶", Hunger: 15, Happiness: 3, Energy: 15}
	FoodSnack   = FoodType{Name: "小零食", Hunger: 10, Happiness: 10, Energy: -3}
	FoodHealthy = FoodType{Name: "营养餐", Hunger: 30, Happiness: 2, Energy: 20}
)

// 互动历史记录（用于成长评估与可解释性）
type InteractionRecord struct {
	Type      string    // "feed"/"pet"/"clean"
	SubType   string    // 食物名或空
	Quality   string    // "good"/"normal"/"bad"
	Timestamp time.Time `json:"t"`
}

// PetState 宠物全部持久化状态
type PetState struct {
	Name string `json:"name"`

	// 四项核心状态 (0-100)
	Hunger    float64 `json:"hunger"`    // 饱腹度，100=饱，0=饿
	Happiness float64 `json:"happiness"` // 好感度，100=非常喜欢主人
	Hygiene   float64 `json:"hygiene"`   // 卫生值，100=干净（方法名Clean()保留做清洁动作，避免字段/方法同名）
	Energy    float64 `json:"energy"`    // 精力值，100=精力充沛

	// 成长系统
	Stage           GrowthStage `json:"stage"`
	CareScore       float64     `json:"care_score"`        // 照顾质量累计分
	CareStreak      int         `json:"care_streak"`       // 连续照顾天数
	DailyCareCount  int         `json:"daily_care_count"`  // 今日互动次数
	LastCareDay     string      `json:"last_care_day"`     // 上次照顾日期(YYYY-MM-DD)
	BirthTime       time.Time   `json:"birth_time"`        // 出生时间
	AgeHours        float64     `json:"age_hours"`         // 累计年龄小时

	// 时间追踪（用于离线补偿）
	LastUpdate time.Time `json:"last_update"`

	// 互动历史
	History []InteractionRecord `json:"history"`
}
