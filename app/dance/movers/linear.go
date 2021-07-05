package movers

import (
	"github.com/wieku/danser-go/app/beatmap/objects"
	"github.com/wieku/danser-go/app/bmath"
	"github.com/wieku/danser-go/app/settings"
	"github.com/wieku/danser-go/framework/math/animation/easing"
	"github.com/wieku/danser-go/framework/math/curves"
	"github.com/wieku/danser-go/framework/math/vector"
	"math"
)

type LinearMover struct {
	*basicMover

	line    curves.Linear
	startTime float64
	simple  bool
}

func NewLinearMover() MultiPointMover {
	return &LinearMover{basicMover: &basicMover{}}
}

func NewLinearMoverSimple() MultiPointMover {
	return &LinearMover{
		basicMover: &basicMover{},
		simple:     true,
	}
}

func (mover *LinearMover) SetObjects(objs []objects.IHitObject) int {
	start, end := objs[0], objs[1]
	startPos := start.GetStackedEndPositionMod(mover.diff.Mods)
	startTime := start.GetEndTime()
	endPos := end.GetStackedStartPositionMod(mover.diff.Mods)
	endTime := end.GetStartTime()

	mover.line = curves.NewLinear(startPos, endPos)

	mover.startTime = startTime
	mover.endTime = endTime

	if mover.simple {
		mover.startTime = math.Max(startTime, end.GetStartTime()-(mover.diff.Preempt-100*mover.diff.Speed))
	} else {
		config := settings.CursorDance.MoverSettings.Linear[mover.id%len(settings.CursorDance.MoverSettings.Linear)]

		if config.WaitForPreempt {
			mover.startTime = math.Max(startTime, end.GetStartTime()-(mover.diff.Preempt-config.ReactionTime*mover.diff.Speed))
		}
	}

	return 2
}

func (mover *LinearMover) Update(time float64) vector.Vector2f {
	t := bmath.ClampF64((time-mover.startTime)/(mover.endTime-mover.startTime), 0, 1)
	return mover.line.PointAt(float32(easing.OutQuad(t)))
}

func (mover *LinearMover) GetObjectsPosition(time float64, object objects.IHitObject) vector.Vector2f {
	config := settings.CursorDance.MoverSettings.Linear[mover.id%len(settings.CursorDance.MoverSettings.Linear)]

	if !config.ChoppyLongObjects || mover.simple || object.GetType() == objects.CIRCLE {
		return mover.basicMover.GetObjectsPosition(time, object)
	}

	const sixtyTime = 1000.0 / 60

	timeDiff := math.Mod(time-object.GetStartTime(), sixtyTime)

	time1 := time - timeDiff
	time2 := time1 + sixtyTime

	pos1 := object.GetStackedPositionAtMod(time1, mover.diff.Mods)
	pos2 := object.GetStackedPositionAtMod(time2, mover.diff.Mods)

	return pos1.Lerp(pos2, float32((time-time1)/sixtyTime))
}
