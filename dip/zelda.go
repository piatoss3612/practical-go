package dip

type Hero interface {
	SetLeftHand(Weapon)
	SetRightHand(Weapon)
	Attack() int
}

type Link struct {
	LeftHand  Weapon
	RightHand Weapon
}

func (l *Link) SetLeftHand(w Weapon) {
	l.LeftHand = w
}

func (l *Link) SetRightHand(w Weapon) {
	l.RightHand = w
}

func (l *Link) Attack() int {
	return l.LeftHand.Attack() + l.RightHand.Attack()
}

var _ Hero = (*Link)(nil)

type Weapon interface {
	Attack() int
}

type MasterSword struct {
	AttackPower int
}

func (m *MasterSword) Attack() int {
	return m.AttackPower
}

type Bow struct {
	AttackPower int
}

func (b *Bow) Attack() int {
	return b.AttackPower
}

type Arrow struct {
	AttackPower int
}

func (a *Arrow) Attack() int {
	return a.AttackPower
}

var _ Weapon = (*MasterSword)(nil)
var _ Weapon = (*Bow)(nil)
var _ Weapon = (*Arrow)(nil)
