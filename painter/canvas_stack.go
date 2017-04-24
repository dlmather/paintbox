package painter

type CanvasStack struct {
	internal []*Canvas
}

func NewCanvasStack() *CanvasStack {
	return &CanvasStack{
		internal: make([]*Canvas, 0),
	}
}

func (cs *CanvasStack) Push(can *Canvas) {
	cs.internal = append(cs.internal, can)
}

func (cs *CanvasStack) Pop() *Canvas {
	if len(cs.internal) == 0 {
		return nil
	}
	last := cs.internal[len(cs.internal)-1]
	if len(cs.internal) == 1 {
		return last
	}
	cs.internal = cs.internal[0 : len(cs.internal)-1]
	return last
}
