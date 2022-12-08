package dialogs

import (
	"encoding/json"
	"errors"
)

const (
	SimpleUtterance = "SimpleUtterance"
)

type Request struct {
	Meta struct {
		Locale     string `json:"locale"`
		Timezone   string `json:"timezone"`
		ClientID   string `json:"client_id"`
		Interfaces struct {
			AccountLinking *struct{} `json:"account_linking"`
			Screen         *struct{} `json:"screen"`
		} `json:"interfaces"`
	} `json:"meta"`

	LinkingComplete *struct{} `json:"account_linking_complete_event,omitempty"`

	Request struct {
		Command           string `json:"command"`
		OriginalUtterance string `json:"original_utterance"`
		Type              string `json:"type"`
		Markup            struct {
			DangerousContext *bool `json:"dangerous_context,omitempty"`
		} `json:"markup,omitempty"`
		Payload interface{} `json:"payload,omitempty"`
		NLU     struct {
			Tokens   []string `json:"tokens"`
			Entities []Entity `json:"entities,omitempty"`
		} `json:"nlu"`
	} `json:"request"`

	Session struct {
		New       bool   `json:"new"`
		MessageID int    `json:"message_id"`
		SessionID string `json:"session_id"`
		SkillID   string `json:"skill_id"`
		UserID    string `json:"user_id"`
	} `json:"session"`

	State struct {
		Session interface{} `json:"session,omitempty"`
	} `json:"state,omitempty"`

	Version string `json:"version"`
	Bearer  string
}

// Type тип запроса (реплика или нажатие кнопки).
func (req *Request) Type() string {
	return req.Request.Type
}

// Text реплика пользователя без изменений.
func (req *Request) Text() string {
	return req.Request.OriginalUtterance
}

// IsNewSession отправлена реплика в рамках нового разговора или уже начатого.
func (req *Request) IsNewSession() bool {
	return req.Session.New
}

// ClientID идентификатор клиентского устройства или приложения. Не рекомендуется использовать.
func (req *Request) ClientID() string {
	return req.Meta.ClientID
}

// Command реплика пользователя, преобразованная Алисой. В частности, текст очищается от знаков препинания, а числительные преобразуются в числа.
func (req *Request) Command() string {
	return req.Request.Command
}

// DangerousContext флаг опасной реплики.
func (req *Request) DangerousContext() bool {
	if req.Request.Markup.DangerousContext != nil {
		return *req.Request.Markup.DangerousContext
	}
	return false
}

// Payload возвращает map[string]interface{} с данными, которые были переданы в Payload кнопки. Подходит для Payload, оформленного в виде json-объекта. Если Payload простая строка следует использовать метод PayloadString(). Если в запросе нет Payload возвращается ошибка.
func (req *Request) Payload() (map[string]interface{}, error) {
	if req.Request.Payload != nil {
		return req.Request.Payload.(map[string]interface{}), nil
	}
	return nil, errors.New("Payload is nil")
}

// SessionID идентификатор сессии.
func (req *Request) SessionID() string {
	return req.Session.SessionID
}

// MessageID счетчик сообщений в рамках сессии.
func (req *Request) MessageID() int {
	return req.Session.MessageID
}

// UserID идентификатор пользователя.
func (req *Request) UserID() string {
	return req.Session.UserID
}

// StateSession Состояние сессии.
func (req *Request) StateSession(key string) interface{} {
	if req.State.Session == nil {
		return nil
	}
	session := req.State.Session.(map[string]interface{})

	return session[key]
}

// State.Session Состояние сессии json строкой
func (req *Request) StateSessionAsJson() (string, error) {
	data, err := json.Marshal(req.State.Session)

	return string(data), err
}

func (req *Request) clean() *Request {
	req.Meta.Interfaces = struct {
		AccountLinking *struct{} `json:"account_linking"`
		Screen         *struct{} `json:"screen"`
	}{
		nil,
		nil,
	}
	req.LinkingComplete = nil
	req.Request.Command = ""
	req.Request.OriginalUtterance = ""
	req.Request.Payload = nil
	req.Request.Markup = struct {
		DangerousContext *bool `json:"dangerous_context,omitempty"`
	}{
		nil,
	}
	req.Request.NLU = struct {
		Tokens   []string `json:"tokens"`
		Entities []Entity `json:"entities,omitempty"`
	}{
		[]string{},
		[]Entity{},
	}
	req.Bearer = ""
	return req
}
