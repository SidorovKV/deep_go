package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		for i, char := range unsafe.Slice(unsafe.StringData(name), len(name)) {
			person.name[i] = char
		}
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = int32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.mana = int16(mana)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.healthPart = byte(health)
		person.houseGunFamTypeHealthPart |= byte(health >> 8 << 6)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.respectStr |= byte(respect)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.respectStr |= byte(strength) << 4
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.explvl |= byte(experience)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.explvl |= byte(level) << 4
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.houseGunFamTypeHealthPart |= 1 << 5
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.houseGunFamTypeHealthPart |= 1 << 4
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.houseGunFamTypeHealthPart |= 1 << 3
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.houseGunFamTypeHealthPart |= 1 << personType
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	x    int32
	y    int32
	z    int32
	gold int32
	// still 6 unused bits in mana
	mana int16
	name [42]byte
	// 1-4 bits - respect, 5-8 bits - strength
	respectStr byte
	// 1-4 bits - experience, 5-8 bits - level
	explvl byte
	// 1-3 bits - type, 4 bit - has house, 5 bit - has gun, 6 bit - has family
	// 7-8 bits - health part
	houseGunFamTypeHealthPart byte
	healthPart                byte
}

func NewGamePerson(options ...Option) GamePerson {
	p := &GamePerson{}

	for _, option := range options {
		option(p)
	}

	return *p
}

func (p *GamePerson) Name() string {
	chars := p.name[:]

	for i := 0; i < 42; i++ {
		if p.name[i] == 0 {
			chars = p.name[:i]

			break
		}
	}

	return unsafe.String(unsafe.SliceData(chars), len(chars))
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	return int(p.mana)
}

func (p *GamePerson) Health() int {
	health := int(p.healthPart)
	health += int(p.houseGunFamTypeHealthPart>>6) << 8

	return health
}

func (p *GamePerson) Respect() int {
	return int(p.respectStr << 4 >> 4)
}

func (p *GamePerson) Strength() int {
	return int(p.respectStr >> 4)
}

func (p *GamePerson) Experience() int {
	return int(p.explvl << 4 >> 4)
}

func (p *GamePerson) Level() int {
	return int(p.explvl >> 4)
}

func (p *GamePerson) HasHouse() bool {
	return p.houseGunFamTypeHealthPart>>5&1 == 1
}

func (p *GamePerson) HasGun() bool {
	return p.houseGunFamTypeHealthPart>>4&1 == 1
}

func (p *GamePerson) HasFamily() bool {
	return p.houseGunFamTypeHealthPart>>3&1 == 1
}

func (p *GamePerson) Type() int {
	val := p.houseGunFamTypeHealthPart
	for i := BuilderGamePersonType; i <= WarriorGamePersonType; i++ {
		if val&1 == 1 {
			return i
		}

		val >>= 1
	}

	return 0
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamily())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
