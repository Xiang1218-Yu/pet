package appearance

import (
	"image/color"
	"math"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"

	"pet/internal/pet"
)

// SizeForStage 根据成长阶段返回宠物尺寸（越大的阶段越大，模拟成长）
func SizeForStage(s pet.GrowthStage) (w, h float32) {
	switch s {
	case pet.StageBaby:
		return 110, 110
	case pet.StageChild:
		return 140, 140
	case pet.StageTeen:
		return 170, 170
	case pet.StageAdult:
		return 200, 200
	default:
		return 140, 140
	}
}

type palette struct {
	Body  color.Color
	Inner color.Color
	Cheek color.Color
}

func paletteFor(stage pet.GrowthStage) palette {
	switch stage {
	case pet.StageBaby:
		return palette{rgba(255, 213, 181, 255), rgba(255, 180, 140, 255), rgba(255, 179, 193, 255)}
	case pet.StageChild:
		return palette{rgba(255, 184, 138, 255), rgba(255, 149, 99, 255), rgba(255, 159, 177, 255)}
	case pet.StageTeen:
		return palette{rgba(245, 158, 92, 255), rgba(232, 131, 63, 255), rgba(255, 143, 166, 255)}
	case pet.StageAdult:
		return palette{rgba(229, 131, 66, 255), rgba(200, 101, 40, 255), rgba(255, 127, 156, 255)}
	default:
		return palette{rgba(255, 184, 138, 255), rgba(255, 149, 99, 255), rgba(255, 159, 177, 255)}
	}
}

func rgba(r, g, b, a uint8) color.Color { return color.RGBA{r, g, b, a} }

// BuildPet 返回一个可放入任意容器的宠物形象对象。
// 优先尝试加载用户自定义PNG（抠好的透明图）；找不到就用Fyne原生Canvas矢量绘制。
func BuildPet(s *pet.PetState) fyne.CanvasObject {
	if png := tryLoadCustomPNG(s); png != nil {
		return png
	}
	return buildDrawnPet(s)
}

// buildDrawnPet 用Fyne的基础几何对象组装一只治愈系橘猫
func buildDrawnPet(s *pet.PetState) fyne.CanvasObject {
	w, _ := SizeForStage(s.Stage)
	pal := paletteFor(s.Stage)
	k := float32(w) / 200.0

	bodyRX, bodyRY := stageBodySize(s)
	bodyX, bodyY := 100*k, 130*k
	body := &canvas.Ellipse{FillColor: pal.Body}
	body.Resize(fyne.NewSize(bodyRX*k*2, bodyRY*k*2))
	body.Move(fyne.NewPos(bodyX-bodyRX*k, bodyY-bodyRY*k))

	headR := stageHeadR(s)
	head := &canvas.Circle{FillColor: pal.Body}
	head.Resize(fyne.NewSize(headR*2*k, headR*2*k))
	head.Move(fyne.NewPos(100*k-headR*k, 85*k-headR*k))

	earScale := headR / 42.0
	earW, earH := 28*earScale*k, 40*earScale*k
	earInW, earInH := 18*earScale*k, 26*earScale*k
	le, re := &canvas.Ellipse{FillColor: pal.Body}, &canvas.Ellipse{FillColor: pal.Body}
	le.Resize(fyne.NewSize(earW, earH))
	re.Resize(fyne.NewSize(earW, earH))
	le.Move(fyne.NewPos(100*k-40*earScale*k-earW/2, 85*k-36*earScale*k-earH/2))
	re.Move(fyne.NewPos(100*k+40*earScale*k-earW/2, 85*k-36*earScale*k-earH/2))
	li, ri := &canvas.Ellipse{FillColor: pal.Inner}, &canvas.Ellipse{FillColor: pal.Inner}
	li.Resize(fyne.NewSize(earInW, earInH))
	ri.Resize(fyne.NewSize(earInW, earInH))
	li.Move(fyne.NewPos(100*k-38*earScale*k-earInW/2, 85*k-28*earScale*k-earInH/2))
	ri.Move(fyne.NewPos(100*k+38*earScale*k-earInW/2, 85*k-28*earScale*k-earInH/2))

	tailRX := float32(16+int(s.Stage)*2) * k
	tX := 100*k + bodyRX*k*0.95
	tY := 150 * k
	tail := &canvas.Ellipse{FillColor: pal.Body}
	tail.Resize(fyne.NewSize(tailRX*2, 32*k))
	tail.Move(fyne.NewPos(tX-tailRX, tY-16*k))
	tailIn := &canvas.Ellipse{FillColor: pal.Inner}
	tailIn.Resize(fyne.NewSize(16*k, 20*k))
	tailIn.Move(fyne.NewPos(tX+tailRX-8*k, tY-10*k))

	cl := &canvas.Ellipse{FillColor: rgba(255, 179, 193, 170)}
	cr := &canvas.Ellipse{FillColor: rgba(255, 179, 193, 170)}
	cl.Resize(fyne.NewSize(18*k, 10*k))
	cr.Resize(fyne.NewSize(18*k, 10*k))
	cl.Move(fyne.NewPos(70*k-9*k, 95*k-5*k))
	cr.Move(fyne.NewPos(130*k-9*k, 95*k-5*k))

	eyes, mouth := drawFace(s.Mood(), k)

	fl := &canvas.Ellipse{FillColor: pal.Inner}
	fr := &canvas.Ellipse{FillColor: pal.Inner}
	fl.Resize(fyne.NewSize(26*k, 14*k))
	fr.Resize(fyne.NewSize(26*k, 14*k))
	fl.Move(fyne.NewPos(75*k-13*k, 180*k-7*k))
	fr.Move(fyne.NewPos(125*k-13*k, 180*k-7*k))

	root := container.NewWithoutLayout(
		tail, tailIn, body, fl, fr, le, re, li, ri, head, cl, cr,
	)
	for _, o := range eyes {
		root.Add(o)
	}
	for _, o := range mouth {
		root.Add(o)
	}
	root.Resize(fyne.NewSize(w, w))
	return container.New(layout.NewCenterLayout(), root)
}

func stageBodySize(s *pet.PetState) (float32, float32) {
	factor := 1.0 + float64(s.Stage)*0.18
	return float32(40.0 * factor), float32(32.0 * factor)
}
func stageHeadR(s *pet.PetState) float32 {
	switch s.Stage {
	case pet.StageBaby:
		return 38
	case pet.StageChild:
		return 42
	case pet.StageTeen:
		return 45
	case pet.StageAdult:
		return 46
	default:
		return 42
	}
}

// drawFace 根据心情绘制眼睛+嘴
func drawFace(m pet.MoodState, k float32) (eyes, mouth []fyne.CanvasObject) {
	brown := rgba(45, 27, 14, 255)
	white := rgba(255, 255, 255, 255)
	sickRed := rgba(138, 59, 59, 255)
	tongue := rgba(255, 142, 161, 255)
	sickGray := rgba(107, 114, 128, 255)
	switch m {
	case pet.MoodHappy:
		eyes = []fyne.CanvasObject{
			makeArc(70*k, 77*k, 10*k, -170, -10, 4, brown),
			makeArc(110*k, 77*k, 10*k, -170, -10, 4, brown),
		}
		mouth = []fyne.CanvasObject{
			makeArc(100*k, 100*k, 10*k, -160, -20, 3, brown),
			makeTongue(100*k, 106*k, 6*k, tongue),
		}
	case pet.MoodNormal:
		eyes = []fyne.CanvasObject{
			makeCircle(80*k, 80*k, 6*k, brown),
			makeCircle(82*k, 78*k, 2*k, white),
			makeCircle(120*k, 80*k, 6*k, brown),
			makeCircle(122*k, 78*k, 2*k, white),
		}
		mouth = []fyne.CanvasObject{
			makeArc(100*k, 102*k, 6*k, -165, -15, 3, brown),
		}
	case pet.MoodSad:
		eyes = []fyne.CanvasObject{
			makeArc(80*k, 82*k, 9*k, 15, 155, 4, brown),
			makeArc(120*k, 82*k, 9*k, 15, 155, 4, brown),
			makeCircle(80*k, 88*k, 3*k, brown),
			makeCircle(120*k, 88*k, 3*k, brown),
		}
		mouth = []fyne.CanvasObject{
			makeArc(100*k, 110*k, 10*k, 10, 170, 3, brown),
		}
	case pet.MoodSick:
		fallthrough
	default:
		eyes = []fyne.CanvasObject{
			makeLine(74*k, 74*k, 86*k, 86*k, 3, sickRed),
			makeLine(86*k, 74*k, 74*k, 86*k, 3, sickRed),
			makeLine(114*k, 74*k, 126*k, 86*k, 3, sickRed),
			makeLine(126*k, 74*k, 114*k, 86*k, 3, sickRed),
		}
		mouth = []fyne.CanvasObject{
			func() fyne.CanvasObject {
				e := &canvas.Ellipse{FillColor: sickGray}
				e.Resize(fyne.NewSize(16*k, 10*k))
				e.Move(fyne.NewPos(100*k-8*k, 105*k-5*k))
				return e
			}(),
		}
	}
	return eyes, mouth
}

// ===== 绘图辅助 =====

func makeCircle(x, y, r float32, fill color.Color) *canvas.Circle {
	c := &canvas.Circle{FillColor: fill}
	c.Resize(fyne.NewSize(2*r, 2*r))
	c.Move(fyne.NewPos(x-r, y-r))
	return c
}

func makeLine(x1, y1, x2, y2, width float32, stroke color.Color) *canvas.Line {
	return &canvas.Line{
		StrokeColor: stroke,
		StrokeWidth: width,
		Position1:   fyne.NewPos(x1, y1),
		Position2:   fyne.NewPos(x2, y2),
	}
}

// makeArc 用多段Line近似一段圆弧。角度: 0度=x轴正方向，逆时针为正
// 画点：P(θ) = (cx + r·cosθ, cy - r·sinθ) （屏幕y轴向下）
func makeArc(cx, cy, r float32, startDeg, endDeg int, width float32, stroke color.Color) fyne.CanvasObject {
	const steps = 24
	segments := make([]fyne.CanvasObject, 0, steps)
	delta := endDeg - startDeg
	for i := 0; i < steps; i++ {
		a1 := degToRad(float64(startDeg + (delta*i)/steps))
		a2 := degToRad(float64(startDeg + (delta*(i+1))/steps))
		x1 := cx + r*float32(math.Cos(a1))
		y1 := cy - r*float32(math.Sin(a1))
		x2 := cx + r*float32(math.Cos(a2))
		y2 := cy - r*float32(math.Sin(a2))
		segments = append(segments, makeLine(x1, y1, x2, y2, width, stroke))
	}
	return container.NewWithoutLayout(segments...)
}

func degToRad(d float64) float64 { return d * math.Pi / 180 }

func makeTongue(cx, cy, r float32, fill color.Color) *canvas.Ellipse {
	e := &canvas.Ellipse{FillColor: fill}
	e.Resize(fyne.NewSize(2*r, r*1.2))
	e.Move(fyne.NewPos(cx-r, cy-r*0.3))
	return e
}

// ======== 用户自定义PNG加载（抠图透明PNG，优先级最高） ========
// 文件名约定： assets/pets/{baby|child|teen|adult}_{happy|normal|sad|sick}.png
// 搜索路径顺序：  ~/.healing-pet/assets/pets/  →  {exeDir}/assets/pets/  →  {cwd}/assets/pets/
func tryLoadCustomPNG(s *pet.PetState) *canvas.Image {
	stageName := []string{"baby", "child", "teen", "adult"}[s.Stage]
	moodName := []string{"happy", "normal", "sad", "sick"}[s.Mood()]
	fileName := stageName + "_" + moodName + ".png"

	candidates := make([]string, 0, 3)
	if home, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates, filepath.Join(home, ".healing-pet", "assets", "pets", fileName))
	}
	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), "assets", "pets", fileName))
	}
	if cwd, err := os.Getwd(); err == nil {
		candidates = append(candidates, filepath.Join(cwd, "assets", "pets", fileName))
	}
	for _, c := range candidates {
		if st, err := os.Stat(c); err == nil && !st.IsDir() {
			u := storage.NewFileURI(c)
			img := canvas.NewImageFromURI(u)
			img.FillMode = canvas.ImageFillContain
			img.ScaleMode = canvas.ImageScaleSmooth
			return img
		}
	}
	return nil
}
