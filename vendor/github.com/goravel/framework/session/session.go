package session

import (
	"maps"
	"slices"

	"github.com/spf13/cast"

	"github.com/goravel/framework/contracts/foundation"
	sessioncontract "github.com/goravel/framework/contracts/session"
	"github.com/goravel/framework/support/color"
	supportmaps "github.com/goravel/framework/support/maps"
	"github.com/goravel/framework/support/str"
)

type Session struct {
	driver     sessioncontract.Driver
	json       foundation.Json
	attributes map[string]any
	id         string
	name       string
	started    bool
}

func NewSession(name string, driver sessioncontract.Driver, json foundation.Json, id ...string) *Session {
	store := &Session{
		name:       name,
		driver:     driver,
		started:    false,
		attributes: make(map[string]any),
		json:       json,
	}
	if len(id) > 0 {
		store.SetID(id[0])
	} else {
		store.SetID("")
	}

	return store
}

func (s *Session) All() map[string]any {
	return s.attributes
}

func (s *Session) Exists(key string) bool {
	return supportmaps.Exists(s.attributes, key)
}

func (s *Session) Flash(key string, value any) sessioncontract.Session {
	s.Put(key, value)

	old := s.Get("_flash.new", []any{}).([]any)
	s.Put("_flash.new", append(old, key))

	s.removeFromOldFlashData(key)

	return s
}

func (s *Session) Flush() sessioncontract.Session {
	s.attributes = make(map[string]any)
	return s
}

func (s *Session) Forget(keys ...string) sessioncontract.Session {
	supportmaps.Forget(s.attributes, keys...)

	return s
}

func (s *Session) Get(key string, defaultValue ...any) any {
	return supportmaps.Get(s.attributes, key, defaultValue...)
}

func (s *Session) GetID() string {
	return s.id
}

func (s *Session) GetName() string {
	return s.name
}

func (s *Session) Has(key string) bool {
	val, ok := s.attributes[key]
	if !ok {
		return false
	}

	return val != nil
}

func (s *Session) Invalidate() error {
	s.Flush()
	return s.migrate(true)
}

func (s *Session) Keep(keys ...string) sessioncontract.Session {
	s.mergeNewFlashes(keys...)
	s.removeFromOldFlashData(keys...)
	return s
}

func (s *Session) Missing(key string) bool {
	return !s.Exists(key)
}

func (s *Session) Now(key string, value any) sessioncontract.Session {
	s.Put(key, value)

	old := s.Get("_flash.old", []any{}).([]any)
	s.Put("_flash.old", append(old, key))

	return s
}

func (s *Session) Only(keys []string) map[string]any {
	return supportmaps.Only(s.attributes, keys...)
}

func (s *Session) Pull(key string, def ...any) any {
	return supportmaps.Pull(s.attributes, key, def...)
}

func (s *Session) Put(key string, value any) sessioncontract.Session {
	s.attributes[key] = value
	return s
}

func (s *Session) Reflash() sessioncontract.Session {
	old := toStringSlice(s.Get("_flash.old", []any{}).([]any))
	s.mergeNewFlashes(old...)
	s.Put("_flash.old", []any{})
	return s
}

func (s *Session) Regenerate(destroy ...bool) error {
	err := s.migrate(destroy...)
	if err != nil {
		return err
	}

	s.regenerateToken()
	return nil
}

func (s *Session) Remove(key string) any {
	return s.Pull(key)
}

func (s *Session) Save() error {
	s.ageFlashData()

	data, err := s.json.MarshalString(s.attributes)
	if err != nil {
		return err
	}

	if err = s.driver.Write(s.GetID(), data); err != nil {
		return err
	}

	s.started = false

	return nil
}

func (s *Session) SetDriver(driver sessioncontract.Driver) sessioncontract.Session {
	if driver == nil {
		return s
	}

	s.driver = driver
	return s
}

func (s *Session) SetID(id string) sessioncontract.Session {
	if s.isValidID(id) {
		s.id = id
	} else {
		s.id = s.generateSessionID()
	}

	return s
}

func (s *Session) SetName(name string) sessioncontract.Session {
	s.name = name

	return s
}

func (s *Session) Start() bool {
	s.loadSession()

	if !s.Has("_token") {
		s.regenerateToken()
	}

	s.started = true
	return s.started
}

func (s *Session) Token() string {
	return s.Get("_token").(string)
}

func (s *Session) generateSessionID() string {
	return str.Random(40)
}

func (s *Session) isValidID(id string) bool {
	return len(id) == 40
}

func (s *Session) loadSession() {
	data := s.readFromHandler()
	if data != nil {
		maps.Copy(s.attributes, data)
	}
}

func (s *Session) migrate(destroy ...bool) error {
	shouldDestroy := false
	if len(destroy) > 0 {
		shouldDestroy = destroy[0]
	}

	if shouldDestroy {
		if err := s.driver.Destroy(s.GetID()); err != nil {
			return err
		}
	}

	s.SetID(s.generateSessionID())

	return nil
}

func (s *Session) readFromHandler() map[string]any {
	value, err := s.driver.Read(s.GetID())
	if err != nil {
		color.Errorln(err)
		return nil
	}
	var data map[string]any
	if value != "" {
		if err := s.json.UnmarshalString(value, &data); err != nil {
			color.Errorln(err)
			return nil
		}
	}

	return data
}

func (s *Session) ageFlashData() {
	old := toStringSlice(s.Get("_flash.old", []any{}).([]any))
	s.Forget(old...)
	s.Put("_flash.old", s.Get("_flash.new", []any{}))
	s.Put("_flash.new", []any{})
}

func (s *Session) mergeNewFlashes(keys ...string) {
	values := s.Get("_flash.new", []any{}).([]any)
	for _, key := range keys {
		if !slices.Contains(values, any(key)) {
			values = append(values, key)
		}
	}

	s.Put("_flash.new", values)
}

func (s *Session) regenerateToken() sessioncontract.Session {
	return s.Put("_token", str.Random(40))
}

func (s *Session) removeFromOldFlashData(keys ...string) {
	old := s.Get("_flash.old", []any{}).([]any)
	for _, key := range keys {
		old = slices.DeleteFunc(old, func(i any) bool {
			return cast.ToString(i) == key
		})
	}
	s.Put("_flash.old", old)
}

// toStringSlice converts an interface slice to a string slice.
func toStringSlice(anySlice []any) []string {
	strSlice := make([]string, len(anySlice))
	for i, v := range anySlice {
		strSlice[i] = cast.ToString(v)
	}
	return strSlice
}
