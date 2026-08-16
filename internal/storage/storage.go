package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"pet/internal/pet"
)

// AppDirName 应用在用户目录下的文件夹名
const AppDirName = ".healing-pet"

// SaveFileName 存档文件名
const SaveFileName = "pet_save.json"

// candidateDirs 返回候选存储目录，按优先级从高到低：
// 1. ~/.healing-pet  (真实用户环境的标准位置)
// 2. ${EXECUTABLE_DIR}/.healing-pet  (绿色便携)
// 3. ${CWD}/.healing-pet-data  (任何环境都能写)
func candidateDirs() []string {
	var dirs []string
	if home, err := os.UserHomeDir(); err == nil {
		dirs = append(dirs, filepath.Join(home, AppDirName))
	}
	if exe, err := os.Executable(); err == nil {
		dirs = append(dirs, filepath.Join(filepath.Dir(exe), AppDirName))
	}
	if cwd, err := os.Getwd(); err == nil {
		dirs = append(dirs, filepath.Join(cwd, AppDirName+"-data"))
	}
	return dirs
}

// EnsureAppDir 依次尝试候选目录，返回第一个可成功创建/已存在的目录
func EnsureAppDir() (string, error) {
	var lastErr error
	for _, d := range candidateDirs() {
		if err := os.MkdirAll(d, 0o755); err == nil {
			// 再做一次可写测试
			testFile := filepath.Join(d, ".wtest")
			if err := os.WriteFile(testFile, []byte("ok"), 0o644); err == nil {
				_ = os.Remove(testFile)
				return d, nil
			} else {
				lastErr = err
			}
		} else {
			lastErr = err
		}
	}
	if lastErr == nil {
		lastErr = errors.New("no writable directory found")
	}
	return "", lastErr
}

// SavePath 返回完整的存档文件路径
func SavePath() (string, error) {
	dir, err := EnsureAppDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, SaveFileName), nil
}

// Save 将宠物状态以 JSON 持久化到磁盘
// 写入采用"临时文件 + rename"模式，防止中途崩溃导致存档损坏
func Save(state *pet.PetState) error {
	if state == nil {
		return errors.New("nil pet state")
	}
	path, err := SavePath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	// Windows 下直接覆盖可能失败，先尝试；若失败则删除目标再改名
	if err := os.Rename(tmp, path); err != nil {
		if runtime.GOOS == "windows" {
			_ = os.Remove(path)
			if err2 := os.Rename(tmp, path); err2 != nil {
				return err2
			}
		} else {
			return err
		}
	}
	return nil
}

// Load 从磁盘加载宠物状态；如果文件不存在或损坏返回 nil, nil
// 尝试所有候选目录，返回找到的第一个存在的存档
func Load() (*pet.PetState, error) {
	dirs := candidateDirs()
	var lastErr error
	for _, d := range dirs {
		path := filepath.Join(d, SaveFileName)
		st, statErr := os.Stat(path)
		if statErr != nil {
			if os.IsNotExist(statErr) {
				continue
			}
			lastErr = statErr
			continue
		}
		if st.IsDir() {
			continue
		}
		data, err := os.ReadFile(path)
		if err != nil {
			lastErr = err
			continue
		}
		var s pet.PetState
		if err := json.Unmarshal(data, &s); err != nil {
			// 存档损坏：备份损坏文件后返回 nil，让应用重置
			_ = os.Rename(path, path+".broken")
			return nil, nil
		}
		return &s, nil
	}
	// 没找到存档是正常的
	if lastErr != nil && !os.IsNotExist(lastErr) {
		// 如果只有真正的IO错误才返回错误（否则"找不到"算正常nil返回）
	}
	// 所有候选都没有找到 -> nil, nil
	return nil, nil
}
