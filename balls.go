package main

type Ball struct {
	r Vector // position
	v Vector // velocity

	radius float64
	mass float64
	id int
}

func (b *Ball) addGravity(level float64) {
	b.v.y += ag*level
}

func (b *Ball) updatePosition() {
	b.r = b.r.add(b.v)
}

func (b *Ball) addVelocity(v Vector) {
	b.v = b.v.add(v)
}

func (b *Ball) isHit(b2 *Ball) bool {
	dist := b.r.distance(b2.r)
	return dist <= (b.radius + b2.radius)
}


func (b *Ball) ballHitVelocity(b2 *Ball, Cr float64) Vector {
	mRatio := (Cr+1)*b2.mass/(b.mass+b2.mass)
	vDiff := b.v.subtract(b2.v)
	rDiff := b.r.subtract(b2.r)
	proj := vDiff.projection(rDiff)
	return b.v.subtract(proj.multiply(mRatio))
}

func (b *Ball) handleWallCollision(padding float64, Cr float64) bool{
	isHit := false
	if b.r.x-b.radius <= padding || b.r.x+b.radius >= screenWidth-padding {
		b.v.x = -Cr*b.v.x
		isHit = true
	}
	if b.r.y-b.radius <= padding || b.r.y+b.radius >= screenHeight-padding {
		b.v.y = -Cr*b.v.y
		isHit = true
	}
	//Move out of wall
	if b.r.x-b.radius < padding {
		b.r.x = b.radius + padding + 1
	}
	if b.r.x+b.radius > screenWidth-padding {
		b.r.x = screenWidth - b.radius - padding -1
	}
	if b.r.y-b.radius < padding {
		b.r.y = b.radius + padding + 1
	}
	if b.r.y+b.radius > screenHeight-padding {
		b.r.y = screenHeight - b.radius - padding -1
	}
	return isHit
}













































