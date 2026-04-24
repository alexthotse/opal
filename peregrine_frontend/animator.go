package main

import (
	"github.com/charmbracelet/harmonica"
)

type Animator struct {
	spring   harmonica.Spring
	position float64
	velocity float64
	target   float64
}

func NewAnimator() *Animator {
	return &Animator{
		spring:   harmonica.NewSpring(harmonica.FPS(60), 6.0, 0.5),
		position: 0,
		velocity: 0,
		target:   0,
	}
}

func (a *Animator) Update() {
	a.position, a.velocity = a.spring.Update(a.position, a.velocity, a.target)
}

func (a *Animator) SetTarget(t float64) {
	a.target = t
}
