//go:build !linux


package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"pet/internal/appearance"
	"pet/internal/pet"
	"pet/internal/storage"
)

// ======== 全局应用状态 ========

type appState struct {
	state *pet.PetState
	// binding 字段，UI自动刷新
	name      binding.String
	stage     binding.String
	mood      binding.String
	hunger    binding.Float
	happiness binding.Float
	hygiene   binding.Float
	energy    binding.Float
	summary   binding.String

	petImage    fyne.CanvasObject // 可能是矢量绘制Container，也可能是PNG Image；统一抽象
	petRoot     *fyne.Container   // 桌宠窗口的根容器
	petStage    fyne.Size         // 内部活动区大小
	petOffsetX  float32   // 宠物在活动区的x偏移
	petOffsetY  float32   // 宠物在活动区的y偏移
	walkVx      float32
	walkVy      float32
	walkTick    *time.Ticker

	petMoveOverlay *petMoveLayer // 内部可拖动的透明覆盖层
}

func newAppState(s *pet.PetState) *appState {
	as := &appState{state: s}
	as.name = binding.NewString()
	as.stage = binding.NewString()
	as.mood = binding.NewString()
	as.hunger = binding.NewFloat()
	as.happiness = binding.NewFloat()
	as.hygiene = binding.NewFloat()
	as.energy = binding.NewFloat()
	as.summary = binding.NewString()
	as.refreshBindings()
	return as
}

func (as *appState) refreshBindings() {
	s := as.state
	_ = as.name.Set(s.Name)
	_ = as.stage.Set(fmt.Sprintf("%s期", s.Stage.String()))
	_ = as.mood.Set(s.Mood().String())
	_ = as.hunger.Set(s.Hunger)
	_ = as.happiness.Set(s.Happiness)
	_ = as.hygiene.Set(s.Hygiene)
	_ = as.energy.Set(s.Energy)
	_ = as.summary.Set(s.StatusSummary())
}

var lastStageMoodKey = ""

// rebuildPetImageIfNeeded 重建/刷新宠物形象；容器内对象整体替换
func (as *appState) rebuildPetImageIfNeeded() {
	cur := fmt.Sprintf("%d_%d", as.state.Stage, as.state.Mood())
	if cur == lastStageMoodKey && as.petImage != nil {
		return
	}
	lastStageMoodKey = cur
	newObj := appearance.BuildPet(as.state)
	w, h := appearance.SizeForStage(as.state.Stage)
	newObj.Resize(fyne.NewSize(w, h))

	// 先移除旧的再添加新的（如果在同一个petRoot里）
	if as.petImage != nil {
		// 找到父容器并移除
		if as.petRoot != nil {
			as.petRoot.Remove(as.petImage)
		}
	}
	as.petImage = newObj
	if as.petRoot != nil {
		// 插在overlay之前（保证点击层在最上面）
		objects := as.petRoot.Objects
		insertIdx := len(objects)
		// 如果最后一个是overlay(>moveOverlay count)，就插在它前面
		if insertIdx > 0 {
			insertIdx = len(objects) - 1
		}
		// 简单做法：直接append，然后再把moveOverlay加到最后
		as.petRoot.Objects = append(objects, newObj)
		if as.petMoveOverlay != nil {
			// 确保overlay在最顶层
			as.petRoot.Remove(as.petMoveOverlay)
			as.petRoot.Add(as.petMoveOverlay)
		}
	}
	positionPet(as)
	as.petImage.Refresh()
}

func positionPet(as *appState) {
	if as.petImage == nil {
		return
	}
	sz := as.petImage.Size()
	stageW := as.petStage.Width
	stageH := as.petStage.Height
	if stageW <= 0 {
		stageW = 260
	}
	if stageH <= 0 {
		stageH = 260
	}
	if as.petOffsetX == 0 && as.petOffsetY == 0 {
		as.petOffsetX = (stageW - sz.Width) / 2
		as.petOffsetY = stageH - sz.Height - 10
	}
	if as.petOffsetX < 0 {
		as.petOffsetX = 0
	}
	if as.petOffsetX+sz.Width > stageW {
		as.petOffsetX = stageW - sz.Width
	}
	if as.petOffsetY < 0 {
		as.petOffsetY = 0
	}
	if as.petOffsetY+sz.Height > stageH {
		as.petOffsetY = stageH - sz.Height
	}
	as.petImage.Move(fyne.NewPos(as.petOffsetX, as.petOffsetY))
}

// ======== 桌宠窗口（带可爱毛玻璃背景的小透明窗口，内部宠物自由活动） ========

const (
	petWinW = 300
	petWinH = 300
)

func buildPetWindow(a fyne.App, as *appState) fyne.Window {
	w := a.NewWindow(as.state.Name)
	// 先尝试桌面端能力：置顶
	if dw, ok := w.(desktop.Window); ok {
		dw.RequestAlwaysOnTop()
	}
	w.Resize(fyne.NewSize(petWinW, petWinH))
	w.SetPadded(false)

	// 毛玻璃背景
	bg := canvas.NewRectangle(color.RGBA{255, 248, 240, 210})
	bg.CornerRadius = 24

	// 阴影描边
	stroke := canvas.NewRectangle(color.RGBA{255, 220, 180, 255})
	stroke.CornerRadius = 24
	stroke.StrokeColor = color.RGBA{255, 200, 150, 255}
	stroke.StrokeWidth = 1.5

	as.petStage = fyne.NewSize(petWinW, petWinH)
	root := container.NewWithoutLayout(stroke, bg)

	// 把宠物加到内部
	as.rebuildPetImageIfNeeded()
	root.Add(as.petImage)

	// 覆盖整个内部活动区的拖动+点击透明层
	as.petMoveOverlay = &petMoveLayer{
		as:       as,
		onLeft:   func() { doPet(as) },
		onRight:  func(pos fyne.Position) { showPetMenu(as, w, pos) },
	}
	as.petMoveOverlay.ExtendBaseWidget(as.petMoveOverlay)
	root.Add(as.petMoveOverlay)

	// 每次窗口大小改变，重排背景+overlay+更新舞台大小
	w.SetContent(container.NewStack(root))
	// Fyne没有SetOnResized，但Stack会自适应；我们监听内容布局变化
	// 这里用一个简单的做法：每当心跳时同步舞台尺寸。
	updateStageSize := func() {
		csize := w.Canvas().Size()
		as.petStage = csize
		stroke.Resize(csize)
		bg.Resize(csize)
		stroke.Move(fyne.NewPos(0, 0))
		bg.Move(fyne.NewPos(0, 0))
		as.petMoveOverlay.Resize(csize)
		as.petMoveOverlay.Move(fyne.NewPos(0, 0))
		positionPet(as)
	}
	w.SetOnClosed(func() {})
	_ = updateStageSize
	// 放到heartbeat里调用
	as.petStage = fyne.NewSize(petWinW, petWinH)
	// 初始调用一次updateStageSize需要canvas尺寸，Show之后才可靠，
	// 所以在startHeartbeat里每10秒更新一次

	// 请求初始位置：右下角
	// 如果拿不到屏幕尺寸就不动了，让WM默认居中
	if screenSize, ok := tryGetPrimaryScreenSize(); ok {
		x := int(float32(screenSize.Width) - petWinW - 80)
		y := int(float32(screenSize.Height) - petWinH - 150)
		if dw, ok := w.(desktop.Window); ok {
			dw.RequestPosition(x, y)
		}
	}

	// 把 updateStageSize 导出到外部调用：通过hook
	stageSyncHook = updateStageSize

	return w
}

var stageSyncHook func()

// tryGetPrimaryScreenSize 获取主屏幕尺寸；Fyne v2.8 未公开此API；返回常见笔记本分辨率兜底
func tryGetPrimaryScreenSize() (fyne.Size, bool) {
	// 兜底：常见笔记本分辨率
	return fyne.NewSize(1440, 900), false
}

// ======== 拖动+点击覆盖层（Fyne自定义控件） ========

type petMoveLayer struct {
	widget.BaseWidget
	as      *appState
	onLeft  func()
	onRight func(pos fyne.Position)

	pressed    bool
	pressPos   fyne.Position
	startOffX  float32
	startOffY  float32
	moved      bool
}

func (t *petMoveLayer) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(canvas.NewRectangle(color.Transparent))
}

var _ desktop.Mouseable = (*petMoveLayer)(nil)

func (t *petMoveLayer) Tapped(ev *fyne.PointEvent) {
	if !t.moved {
		if t.onLeft != nil {
			t.onLeft()
		}
	}
}

func (t *petMoveLayer) TappedSecondary(ev *fyne.PointEvent) {
	if t.onRight != nil {
		t.onRight(ev.Position)
	}
}

func (t *petMoveLayer) MouseDown(ev *desktop.MouseEvent) {
	t.pressed = true
	t.moved = false
	t.pressPos = ev.Position
	t.startOffX = t.as.petOffsetX
	t.startOffY = t.as.petOffsetY
}

func (t *petMoveLayer) MouseUp(ev *desktop.MouseEvent) {
	t.pressed = false
}

func (t *petMoveLayer) MouseMoved(ev *desktop.MouseEvent) {
	if !t.pressed {
		return
	}
	dx := ev.Position.X - t.pressPos.X
	dy := ev.Position.Y - t.pressPos.Y
	if !t.moved {
		if dx*dx+dy*dy <= 9.0 {
			return
		}
		t.moved = true
	}
	t.as.petOffsetX = t.startOffX + dx
	t.as.petOffsetY = t.startOffY + dy
	positionPet(t.as)
	t.as.petImage.Refresh()
}

// ======== 操作菜单（用PopUpMenu，不依赖desktop.ShowMenuAt） ========

func showPetMenu(as *appState, w fyne.Window, pos fyne.Position) {
	items := []*fyne.MenuItem{
		{Label: "🐟 喂小鱼干", Action: func() { doFeed(as, pet.FoodFish) }},
		{Label: "🥛 喂牛奶", Action: func() { doFeed(as, pet.FoodMilk) }},
		{Label: "🍬 喂小零食", Action: func() { doFeed(as, pet.FoodSnack) }},
		{Label: "🍱 喂营养餐", Action: func() { doFeed(as, pet.FoodHealthy) }},
		{IsSeparator: true, Label: ""},
		{Label: "🫶 抚摸一下", Action: func() { doPet(as) }},
		{Label: "🛁 洗澡清洁", Action: func() { doClean(as) }},
		{IsSeparator: true, Label: ""},
		{Label: "📊 打开控制面板", Action: func() { globalControlWin.Show() }},
		{Label: "💾 立即保存", Action: func() { doSave(as) }},
		{Label: "🚪 退出应用", Action: func() {
			doSave(as)
			fyne.CurrentApp().Quit()
		}},
	}
	menu := fyne.NewMenu("", items...)
	pop := widget.NewPopUpMenu(menu, w.Canvas())
	// PopUpMenu 默认在菜单下；pos 是活动区相对坐标，转换到 Canvas 坐标
	// 在 Stack 里 overlay 的位置就是 0,0，所以可以直接用 pos
	canvasPos := pos
	// 但菜单需要显示在鼠标点上
	pop.ShowAtPosition(canvasPos)
}

// ======== 控制面板窗口 ========

var globalControlWin fyne.Window

func buildControlWindow(a fyne.App, as *appState) fyne.Window {
	w := a.NewWindow("治愈系养宠 · 控制面板")
	globalControlWin = w
	w.Resize(fyne.NewSize(440, 600))

	// 头部
	titleText := canvas.NewText("", color.RGBA{150, 80, 30, 255})
	titleText.TextSize = 22
	titleText.TextStyle = fyne.TextStyle{Bold: true}
	stageText := canvas.NewText("", color.RGBA{80, 80, 80, 255})
	stageText.TextSize = 13
	moodText := canvas.NewText("", color.RGBA{80, 80, 80, 255})
	moodText.TextSize = 13

	updateHeader := func() {
		n, _ := as.name.Get()
		titleText.Text = n
		titleText.Refresh()
		st, _ := as.stage.Get()
		stageText.Text = "成长：" + st
		stageText.Refresh()
		md, _ := as.mood.Get()
		moodText.Text = "心情：" + md
		moodText.Refresh()
	}
	updateHeader()
	as.name.AddListener(binding.NewDataListener(updateHeader))
	as.stage.AddListener(binding.NewDataListener(updateHeader))
	as.mood.AddListener(binding.NewDataListener(updateHeader))

	header := container.NewVBox(
		titleText,
		container.NewHBox(stageText, layout.NewSpacer(), moodText),
		widget.NewSeparator(),
	)

	// 状态条
	barHunger := newBar(as.hunger, "🍗 饱腹度", color.RGBA{243, 156, 18, 255})
	barHappiness := newBar(as.happiness, "💖 好感度", color.RGBA{230, 76, 120, 255})
	barClean := newBar(as.hygiene, "✨ 卫生值", color.RGBA{52, 152, 219, 255})
	barEnergy := newBar(as.energy, "⚡ 精力值", color.RGBA{46, 204, 113, 255})
	barsBox := container.NewVBox(
		barHunger,
		barHappiness,
		barClean,
		barEnergy,
		widget.NewSeparator(),
	)

	// 状态摘要
	summaryLabel := widget.NewLabelWithData(as.summary)
	summaryLabel.Wrapping = fyne.TextWrapWord
	summaryCard := container.NewVBox(
		canvas.NewText("当前状态：", color.RGBA{100, 100, 100, 255}),
		summaryLabel,
		widget.NewSeparator(),
	)

	// 操作按钮
	btnFeed1 := widget.NewButton("🐟 小鱼干", func() { doFeed(as, pet.FoodFish) })
	btnFeed2 := widget.NewButton("🥛 牛奶", func() { doFeed(as, pet.FoodMilk) })
	btnFeed3 := widget.NewButton("🍬 小零食", func() { doFeed(as, pet.FoodSnack) })
	btnFeed4 := widget.NewButton("🍱 营养餐", func() { doFeed(as, pet.FoodHealthy) })
	feedRow := container.NewGridWithColumns(2, btnFeed1, btnFeed2, btnFeed3, btnFeed4)

	btnPet := widget.NewButton("🫶 抚摸 (加好感)", func() { doPet(as) })
	btnPet.Importance = widget.HighImportance
	btnClean := widget.NewButton("🛁 洗澡 (恢复卫生)", func() { doClean(as) })
	btnClean.Importance = widget.WarningImportance
	actRow := container.NewGridWithColumns(2, btnPet, btnClean)

	opBox := container.NewVBox(
		canvas.NewText("喂食", color.RGBA{100, 60, 20, 255}),
		feedRow,
		canvas.NewText("互动", color.RGBA{100, 60, 20, 255}),
		actRow,
		widget.NewSeparator(),
	)

	// 成长档案
	growthLabel := widget.NewLabel("")
	updateGrowth := func() {
		s := as.state
		growthLabel.SetText(fmt.Sprintf(
			"🎂 出生：%s\n⏳ 累计：%.1f 小时\n⭐ 照顾得分：%.0f\n🔥 连续照顾：%d 天\n📝 今日互动：%d 次\n📚 累计互动：%d 次",
			s.BirthTime.Format("2006-01-02 15:04"),
			s.AgeHours,
			s.CareScore,
			s.CareStreak,
			s.DailyCareCount,
			len(s.History),
		))
	}
	updateGrowth()
	as.mood.AddListener(binding.NewDataListener(func() { updateGrowth() }))

	growthBox := container.NewVBox(
		canvas.NewText("成长档案", color.RGBA{100, 60, 20, 255}),
		growthLabel,
	)

	// 底部
	saveBtn := widget.NewButton("💾 立即保存", func() { doSave(as) })
	resetBtn := widget.NewButton("🔄 重开新宠物", func() {
		dialog.ShowConfirm("确认重开", "当前宠物将被替换为新生幼崽，存档覆盖，确定吗？", func(ok bool) {
			if !ok {
				return
			}
			as.state = pet.NewBorn("小治愈")
			lastStageMoodKey = ""
			as.refreshBindings()
			as.rebuildPetImageIfNeeded()
			updateGrowth()
			doSave(as)
		}, w)
	})
	resetBtn.Importance = widget.DangerImportance
	bottomRow := container.NewGridWithColumns(2, saveBtn, resetBtn)

	content := container.NewVScroll(container.NewVBox(
		header,
		barsBox,
		summaryCard,
		opBox,
		growthBox,
		layout.NewSpacer(),
		bottomRow,
	))
	w.SetContent(content)
	return w
}

// 彩色状态条：标题 + 数字 + 颜色条（低进度变红）
func newBar(data binding.Float, title string, c color.Color) fyne.CanvasObject {
	bg := canvas.NewRectangle(color.RGBA{225, 225, 225, 255})
	bg.CornerRadius = 6
	fg := canvas.NewRectangle(c)
	fg.CornerRadius = 6
	bars := container.NewWithoutLayout(bg, fg)
	bars.MinSize()
	bg.Resize(fyne.NewSize(380, 14))
	fg.Resize(fyne.NewSize(100, 14))

	titleTxt := canvas.NewText(title, color.RGBA{80, 60, 20, 255})
	titleTxt.TextSize = 13
	numTxt := canvas.NewText("0 / 100", color.RGBA{80, 80, 80, 255})
	numTxt.TextSize = 12
	numTxt.Alignment = fyne.TextAlignTrailing

	// 用一行HBox放标题和数字，然后bars单独一行
	top := container.NewHBox(titleTxt, layout.NewSpacer(), numTxt)
	box := container.NewVBox(top, bars)
	refresh := func() {
		v, _ := data.Get()
		numTxt.Text = fmt.Sprintf("%.0f / 100", v)
		numTxt.Refresh()
		pct := float32(v) / 100
		if pct < 0 {
			pct = 0
		}
		if pct > 1 {
			pct = 1
		}
		fg.Resize(fyne.NewSize(380*pct, 14))
		if pct < 0.25 {
			fg.FillColor = color.RGBA{220, 60, 60, 255}
		} else {
			fg.FillColor = c
		}
		fg.Refresh()
	}
	data.AddListener(binding.NewDataListener(refresh))
	refresh()
	return box
}

// ======== 动作实现 ========

func doFeed(as *appState, f pet.FoodType) {
	q := as.state.Feed(f, time.Now())
	showToastMsg(fmt.Sprintf("喂了 %s（%s）", f.Name, zhQuality(q)))
	afterAction(as)
}
func doPet(as *appState) {
	q := as.state.Pet(time.Now())
	showToastMsg(fmt.Sprintf("摸了摸头（%s）", zhQuality(q)))
	afterAction(as)
}
func doClean(as *appState) {
	q := as.state.Clean(time.Now())
	showToastMsg(fmt.Sprintf("洗了个澡（%s）", zhQuality(q)))
	afterAction(as)
}

func zhQuality(q string) string {
	switch q {
	case "good":
		return "很棒！"
	case "normal":
		return "还行~"
	case "bad":
		return "不太好…"
	default:
		return q
	}
}

func afterAction(as *appState) {
	as.refreshBindings()
	as.rebuildPetImageIfNeeded()
	doSave(as)
}

func doSave(as *appState) {
	if err := storage.Save(as.state); err != nil {
		log.Printf("保存失败: %v", err)
	}
}

var toastHook func(string)

func showToastMsg(msg string) {
	if toastHook != nil {
		toastHook(msg)
	}
}

// ======== 内部随机游走（宠物在stage容器内走动） ========

func startWalkAnimation(as *appState) {
	as.walkVx = 0.4
	as.walkVy = 0
	as.walkTick = time.NewTicker(60 * time.Millisecond)
	go func() {
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		turnCounter := 0
		for range as.walkTick.C {
			// 精力低或被用户拖着时不自动走
			as.state.Tick(time.Now())
			if as.state.Energy < 15 || as.petMoveOverlay.pressed {
				continue
			}
			if as.petImage == nil {
				continue
			}
			sz := as.petImage.Size()
			stageW := as.petStage.Width
			stageH := as.petStage.Height
			if stageW <= 10 || stageH <= 10 {
				continue
			}
			newX := as.petOffsetX + as.walkVx
			newY := as.petOffsetY + as.walkVy
			moved := false
			if newX < 0 {
				newX = 0
				as.walkVx = -as.walkVx
				moved = true
			}
			if newX+sz.Width > stageW {
				newX = stageW - sz.Width
				as.walkVx = -as.walkVx
				moved = true
			}
			if newY < 0 {
				newY = 0
				as.walkVy = -as.walkVy
				moved = true
			}
			if newY+sz.Height > stageH {
				newY = stageH - sz.Height
				as.walkVy = -as.walkVy
				moved = true
			}
			_ = moved
			as.petOffsetX = newX
			as.petOffsetY = newY
			// UI操作必须切回主线程
			fyne.Do(func() {
				if as.petImage != nil {
					as.petImage.Move(fyne.NewPos(newX, newY))
					as.petImage.Refresh()
				}
			})

			turnCounter++
			if turnCounter > 100+rng.Intn(180) {
				turnCounter = 0
				speed := float32(0.3 + rng.Float64()*0.7)
				switch rng.Intn(6) {
				case 0:
					as.walkVx = speed
					as.walkVy = 0
				case 1:
					as.walkVx = -speed
					as.walkVy = 0
				case 2:
					as.walkVx = 0
					as.walkVy = 0
				case 3:
					as.walkVx = speed
					as.walkVy = -speed * 0.5
				case 4:
					as.walkVx = -speed
					as.walkVy = speed * 0.5
				case 5:
					as.walkVx = 0
					as.walkVy = speed
				}
			}
		}
	}()
}

// 每秒心跳：衰减状态、重绘、定期同步舞台尺寸、自动保存
func startHeartbeat(as *appState) {
	tick := time.NewTicker(1 * time.Second)
	saveCounter := 0
	sizeCounter := 0
	go func() {
		for range tick.C {
			as.state.Tick(time.Now())
			// 所有UI相关操作切主线程
			fyne.Do(func() {
				as.refreshBindings()
				as.rebuildPetImageIfNeeded()
			})
			saveCounter++
			sizeCounter++
			if saveCounter >= 30 {
				saveCounter = 0
				doSave(as)
			}
			if sizeCounter >= 10 {
				sizeCounter = 0
				if stageSyncHook != nil {
					fyne.Do(func() {
						if stageSyncHook != nil {
							stageSyncHook()
						}
					})
				}
			}
		}
	}()
}

// ======== 自定义主题（更可爱的治愈系配色） ========

type cuteTheme struct{}

var _ fyne.Theme = (*cuteTheme)(nil)

func (c *cuteTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	defaults := theme.DefaultTheme()
	return defaults.Color(n, v)
}
func (c *cuteTheme) Font(s fyne.TextStyle) fyne.Resource { return theme.DefaultTheme().Font(s) }
func (c *cuteTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}
func (c *cuteTheme) Size(n fyne.ThemeSizeName) float32 {
	s := theme.DefaultTheme().Size(n)
	if n == theme.SizeNamePadding {
		return s * 1.15
	}
	return s
}

// ======== main ========

func main() {
	// 1) 加载或创建宠物
	s, err := storage.Load()
	if err != nil {
		log.Fatalf("读取存档失败: %v", err)
	}
	if s == nil {
		log.Println("未发现存档，创建新宠物…")
		s = pet.NewBorn("小治愈")
		if err := storage.Save(s); err != nil {
			log.Printf("初始存档失败: %v", err)
		}
	} else {
		s.Tick(time.Now()) // 离线补偿
	}

	a := app.NewWithID("com.healing.pet.desktop")
	a.Settings().SetTheme(&cuteTheme{})

	as := newAppState(s)

	controlWin := buildControlWindow(a, as)
	petWin := buildPetWindow(a, as)

	// 气泡提示：把消息临时替换summary 2秒
	toastHook = func(msg string) {
		go func() {
			old, _ := as.summary.Get()
			_ = as.summary.Set("💬 " + msg)
			time.Sleep(2 * time.Second)
			_ = as.summary.Set(old)
		}()
	}

	startHeartbeat(as)
	startWalkAnimation(as)

	petWin.Show()
	controlWin.Show()

	a.Lifecycle().SetOnStopped(func() {
		doSave(as)
		if as.walkTick != nil {
			as.walkTick.Stop()
		}
	})

	a.Run()
}
