// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package spc2000

// Interrupt messages
const (
	GetStatus = iota
	SetSkipAmount
	TriggerSleepCycle
	SetSkipUnit
)

// Status codes
const (
	Evacuate = iota
	NotInVacuum
	InsufficientFuel
	GravityImbalance
	AngularMomentum
	OpenDoors
	MechanicalError
	UnknownError = 0xffff
)

// Time units
const (
	Milliseconds = iota
	Minutes
	Days
	Years
)
