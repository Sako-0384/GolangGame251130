package game

// Scene はゲームの各シーン（タイトル、メインゲームなど）が満たすべきインターフェース
type Scene interface {
	Update(dt float32)
	Draw()
	OnEnter()
}

// SceneManager は現在のシーンを管理し、遷移を制御する
type SceneManager struct {
	currentScene Scene
	enteredScene bool
}

// NewSceneManager は新しいシーンマネージャーを作成する
func NewSceneManager() *SceneManager {
	return &SceneManager{}
}

// ChangeScene は現在のシーンを変更する
func (sm *SceneManager) ChangeScene(scene Scene) {
	sm.currentScene = scene
	sm.enteredScene = true
}

// Update は現在のシーンのUpdateを呼び出す
func (sm *SceneManager) Update(dt float32) {
	if sm.enteredScene {
		sm.enteredScene = false
		sm.currentScene.OnEnter()
	}

	if sm.currentScene != nil {
		sm.currentScene.Update(dt)
	}
}

// Draw は現在のシーンのDrawを呼び出す
func (sm *SceneManager) Draw() {
	if sm.currentScene != nil {
		sm.currentScene.Draw()
	}
}
