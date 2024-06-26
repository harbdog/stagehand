package stagehand

import ebiten "github.com/hajimehoshi/ebiten/v2"

type SceneManager[T any] struct {
	current ProtoScene[T]
}

func NewSceneManager[T any](scene Scene[T], state T) *SceneManager[T] {
	s := &SceneManager[T]{current: scene}
	scene.Load(state, s)
	return s
}

// Scene Switching
func (s *SceneManager[T]) SwitchTo(scene Scene[T]) {
	if prevTransition, ok := s.current.(SceneTransition[T]); ok {
		// previous transition is still running, end it first
		prevTransition.End()
	}
	if c, ok := s.current.(Scene[T]); ok {
		scene.Load(c.Unload(), s)
		s.current = scene
	}
}

func (s *SceneManager[T]) SwitchWithTransition(scene Scene[T], transition SceneTransition[T]) {
	if prevTransition, ok := s.current.(SceneTransition[T]); ok {
		// previous transition is still running, end it to start the next transition
		prevTransition.End()
	}
	sc := s.current.(Scene[T])
	transition.Start(sc, scene, s)
	if c, ok := sc.(TransitionAwareScene[T]); ok {
		scene.Load(c.PreTransition(scene), s)
	} else {
		scene.Load(sc.Unload(), s)
	}
	s.current = transition
}

func (s *SceneManager[T]) ReturnFromTransition(scene, origin Scene[T]) {
	if c, ok := scene.(TransitionAwareScene[T]); ok {
		c.PostTransition(origin.Unload(), origin)
	} else {
		scene.Load(origin.Unload(), s)
	}
	s.current = scene
}

// Ebiten Interface
func (s *SceneManager[T]) Update() error {
	return s.current.Update()
}

func (s *SceneManager[T]) Draw(screen *ebiten.Image) {
	s.current.Draw(screen)
}

func (s *SceneManager[T]) Layout(w, h int) (int, int) {
	return s.current.Layout(w, h)
}
