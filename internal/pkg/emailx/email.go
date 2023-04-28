package emailx

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/smtp"
	"time"

	"github.com/jordan-wright/email"
)

const (
	defaultCodeTemplate = `<p>您好，欢迎注册。</p>您的验证码为：`
	codeMap             = "0123456789"
	defaultCodeLength   = 6
	defaultPoolSize     = 4
	defaultSubject      = "验证码"
)

var (
	KeyPrefix     = "mail:code:"
	KeyTTL        = time.Minute
	defaultConfig = Config{
		KeyPrefix:    KeyPrefix,
		KeyTTL:       KeyTTL,
		CodeLength:   uint8(defaultCodeLength),
		CodeTemplate: defaultCodeTemplate,
		Subject:      defaultSubject,
		PoolSize:     defaultPoolSize,
	}
)

type (
	Config struct {
		KeyPrefix    string
		KeyTTL       time.Duration
		CodeTemplate string
		Subject      string
		UserEmail    string
		CodeLength   uint8
		PoolSize     uint8
	}

	Option func(c *Config)

	EmailManager struct {
		*email.Pool
	}
)

func WithSubject(subject string) Option {
	return func(c *Config) {
		c.Subject = subject
	}
}

func WithKeyPrefix(keyPrefix string) Option {
	return func(c *Config) {
		c.KeyPrefix = keyPrefix
	}
}

func WithKeyTTL(keyTTL time.Duration) Option {
	return func(c *Config) {
		c.KeyTTL = keyTTL
	}
}

func WithCodeTemplate(codeTemplate string) Option {
	return func(c *Config) {
		c.CodeTemplate = codeTemplate
	}
}

func WithCodeLength(codeLength uint8) Option {
	return func(c *Config) {
		c.CodeLength = codeLength
	}
}

func WithPoolSize(poolsize uint8) Option {
	return func(c *Config) {
		c.PoolSize = poolsize
	}
}

func NewEmailManager(host string, port int, username string, password string, opts ...Option) (*EmailManager, error) {
	for _, opt := range opts {
		opt(&defaultConfig)
	}

	defaultConfig.UserEmail = username
	KeyPrefix = defaultConfig.KeyPrefix
	KeyTTL = defaultConfig.KeyTTL

	address := fmt.Sprintf("%s:%d", host, port)
	pool, err := email.NewPool(address, int(defaultConfig.PoolSize),
		smtp.PlainAuth("", username, password, host))
	if err != nil {
		return nil, err
	}

	return &EmailManager{pool}, nil
}

func RandomCode() string {
	rand.Seed(time.Now().Unix())
	code := bytes.Buffer{}
	code.Grow(6)
	codeMapLen := len(codeMap)
	for i := 0; i < int(defaultConfig.CodeLength); i++ {
		code.WriteByte(codeMap[rand.Intn(codeMapLen)])
	}
	return code.String()
}

func (m *EmailManager) SendCodeEmail(code, em string, timeout time.Duration) error {
	e := email.NewEmail()
	e.From = defaultConfig.UserEmail
	e.To = []string{em}
	e.Subject = defaultConfig.Subject
	e.HTML = []byte(defaultConfig.CodeTemplate + code)
	return m.Send(e, timeout)
}

func (m *EmailManager) SendEmail(content []byte, to string, timeout time.Duration) error {
	e := email.NewEmail()
	e.From = defaultConfig.UserEmail
	e.To = []string{to}
	e.Subject = defaultConfig.Subject
	e.HTML = content
	return m.Send(e, timeout)
}
