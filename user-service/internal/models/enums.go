package models

import (
	"database/sql/driver"
	"fmt"
)

type UserStatus int

const (
	UserStatusActive UserStatus = iota
	UserStatusInactive
	UserStatusSuspended
	UserStatusPendingVerification
	UserStatusBanned
)

func (s UserStatus) String() string {
	switch s {
	case UserStatusActive:
		return "active"
	case UserStatusInactive:
		return "inactive"
	case UserStatusSuspended:
		return "suspended"
	case UserStatusPendingVerification:
		return "pending_verification"
	case UserStatusBanned:
		return "banned"
	default:
		return "unknown"
	}
}

func (s *UserStatus) FromString(str string) error {
	switch str {
	case "active":
		*s = UserStatusActive
	case "inactive":
		*s = UserStatusInactive
	case "suspended":
		*s = UserStatusSuspended
	case "pending_verification":
		*s = UserStatusPendingVerification
	case "banned":
		*s = UserStatusBanned
	default:
		return fmt.Errorf("invalid user status: %s", str)
	}
	return nil
}

func (s UserStatus) Value() (driver.Value, error) {
	return s.String(), nil
}

func (s *UserStatus) Scan(value interface{}) error {
	if value == nil {
		*s = UserStatusActive
		return nil
	}

	switch v := value.(type) {
	case string:
		return s.FromString(v)
	case []byte:
		return s.FromString(string(v))
	default:
		return fmt.Errorf("cannot scan %T into UserStatus", value)
	}
}

type Gender int

const (
	GenderNotSpecified Gender = iota
	GenderMale
	GenderFemale
	GenderOther
)

func (g Gender) String() string {
	switch g {
	case GenderMale:
		return "male"
	case GenderFemale:
		return "female"
	case GenderOther:
		return "other"
	default:
		return "not_specified"
	}
}

func (g *Gender) FromString(str string) error {
	switch str {
	case "male":
		*g = GenderMale
	case "female":
		*g = GenderFemale
	case "other":
		*g = GenderOther
	case "not_specified", "":
		*g = GenderNotSpecified
	default:
		return fmt.Errorf("invalid gender: %s", str)
	}
	return nil
}

func (g Gender) Value() (driver.Value, error) {
	return g.String(), nil
}

func (g *Gender) Scan(value interface{}) error {
	if value == nil {
		*g = GenderNotSpecified
		return nil
	}

	switch v := value.(type) {
	case string:
		return g.FromString(v)
	case []byte:
		return g.FromString(string(v))
	default:
		return fmt.Errorf("cannot scan %T into Gender", value)
	}
}

type Theme int

const (
	ThemeLight Theme = iota
	ThemeDark
	ThemeAuto
)

func (t Theme) String() string {
	switch t {
	case ThemeLight:
		return "light"
	case ThemeDark:
		return "dark"
	case ThemeAuto:
		return "auto"
	default:
		return "light"
	}
}

func (t *Theme) FromString(str string) error {
	switch str {
	case "light":
		*t = ThemeLight
	case "dark":
		*t = ThemeDark
	case "auto":
		*t = ThemeAuto
	default:
		return fmt.Errorf("invalid theme: %s", str)
	}
	return nil
}

func (t Theme) Value() (driver.Value, error) {
	return t.String(), nil
}

func (t *Theme) Scan(value interface{}) error {
	if value == nil {
		*t = ThemeLight
		return nil
	}

	switch v := value.(type) {
	case string:
		return t.FromString(v)
	case []byte:
		return t.FromString(string(v))
	default:
		return fmt.Errorf("cannot scan %T into Theme", value)
	}
}
