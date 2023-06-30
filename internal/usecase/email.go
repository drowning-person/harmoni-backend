package usecase

import (
	"bytes"
	"context"
	"fmt"
	"harmoni/internal/conf"
	emailentity "harmoni/internal/entity/email"
	"harmoni/internal/pkg/common"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/reason"
	"html/template"
	"math/rand"
	"mime"
	"net/smtp"
	"runtime"
	"time"

	"github.com/jordan-wright/email"
	"go.uber.org/zap"
)

const (
	registerTemplate = `
<!DOCTYPE html>
	<html>
	  <head>
		<meta charset="UTF-8">
		<title>验证码</title>
		<style>
		  /* CSS 样式表 */
		  body {
			font-family: Arial, sans-serif;
			background-color: #f7f7f7;
			margin: 0;
			padding: 0;
		  }
		  
		  .container {
			width: 80%;
			margin: 0 auto;
			background-color: #fff;
			padding: 20px;
			box-shadow: 0px 2px 4px rgba(0, 0, 0, 0.1);
		  }
		  
		  h1 {
			font-size: 24px;
			margin-bottom: 20px;
		  }
		  
		  p {
			font-size: 16px;
			line-height: 1.5;
			margin-bottom: 20px;
		  }
		  
		  .button {
			display: inline-block;
			padding: 10px 20px;
			background-color: #2f80ed;
			color: #fff;
			font-size: 16px;
			text-decoration: none;
			border-radius: 5px;
		  }
		  
		  .button:hover {
			background-color: #1c5de7;
		  }
		</style>
	  </head>
	  <body>
		<div class="container">
		  <p>尊敬的用户，感谢您注册我们的服务。</p>
		  <p>以下是您的验证码：</p>
		  <h2 style="font-size: 36px; margin-top: 0;">{{ .Code }}</h2>
		  <p>请在注册页面输入上述验证码以完成注册。</p>
		</div>
	  </body>
	</html>`
	registerSubject = "注册验证码"

	resetPasswordTemplate = `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>重置密码</title>
	</head>
	<body>
		<div>
			<h1>重置密码</h1>
			<p>我们收到了一份针对您账户的密码重置请求。如果您没有请求此操作，请忽略此邮件。</p>
			<p>如果您确实请求了密码重置，请在页面输入以下代码：</p>
			<p>{{ .Code }}</p>
			<p>此代码将在 {{ .TTL }} 分钟后失效。</p>
		</div>
	</body>
	</html>
	`
	resetPasswordSubject = "重置密码"

	changeTemplate = `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>修改账户信息</title>
	</head>
	<body>
		<div>
			<h1>修改账户信息</h1>
			<p>我们收到了您修改{{ .ChangeName }}的请求。如果您没有进行此操作，请忽略此邮件。</p>
			<p>如果您确实请求了修改账户信息，请在页面输入以下代码：</p>
			<p>{{ .Code }}</p>
			<p>此代码将在 {{ .TTL }} 分钟后失效。</p>
			<p>如果您没有进行上述操作或有任何疑问，请及时联系我们。</p>
		</div>
	</body>
	</html>
	`
	changeSubject = "修改账户信息"
)

const (
	codeMap = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	mailcodePrefix = "mail.code:"
)

func emailcodeKey(email string, emailtype emailentity.EmailType) string {
	return fmt.Sprintf("%s%s:%d", mailcodePrefix, email, emailtype)
}

type EmailUsecase struct {
	emailRepo emailentity.EmailRepo
	logger    *zap.SugaredLogger
	conf      *conf.Email
	emailPool *email.Pool
}

func NewEmailUsecase(conf *conf.Email, emailRepo emailentity.EmailRepo, logger *zap.SugaredLogger) (*EmailUsecase, error) {
	pool, err := email.NewPool(fmt.Sprintf("%s:%s", conf.Host, conf.Port), runtime.NumCPU()*2,
		smtp.PlainAuth("", conf.UserName, conf.Password, conf.Host))
	if err != nil {
		return nil, err
	}

	return &EmailUsecase{
		conf:      conf,
		emailRepo: emailRepo,
		logger:    logger,
		emailPool: pool,
	}, nil
}

func (u *EmailUsecase) CheckBeforeSendCode(ctx context.Context, email string, emailType emailentity.EmailType) error {
	key := emailcodeKey(email, emailType)
	content, err := u.emailRepo.GetCode(ctx, key)
	if err != nil {
		return err
	}

	if content != "" {
		data := &emailentity.EmailCodeContent{}
		err = data.FromJSONString(content)
		u.logger.Debugf("email code content: %#v", data)
		if err != nil {
			u.logger.Error(err)
			return errorx.BadRequest(reason.EmailShouldRequestLater)
		} else if time.Now().Unix()-data.LastReqTime < int64(time.Minute/time.Second) {
			return errorx.BadRequest(reason.EmailShouldRequestLater)
		} else {
			err = u.emailRepo.DelCode(ctx, key)
			if err != nil {
				u.logger.Error(err)
				return errorx.BadRequest(reason.EmailShouldRequestLater)
			}
		}
	}

	return nil
}

// SendAndSaveCode send email and save code
func (u *EmailUsecase) SendAndSaveCode(ctx context.Context, toEmailAddr, subject, body, codeContent string, emailType emailentity.EmailType) error {
	key := emailcodeKey(toEmailAddr, emailType)
	_, err := u.emailRepo.SetCode(ctx, key, codeContent, u.conf.CodeTTL)
	if err != nil {
		u.logger.Error(err)
	}

	go u.Send(ctx, toEmailAddr, subject, body)
	return nil
}

// SendAndSaveCodeWithTime send email and save code
func (u *EmailUsecase) SendAndSaveCodeWithTime(
	ctx context.Context, toEmailAddr, subject, body, codeContent string, emailType emailentity.EmailType, duration time.Duration) error {
	key := emailcodeKey(toEmailAddr, emailType)
	_, err := u.emailRepo.SetCode(ctx, key, codeContent, duration)
	if err != nil {
		u.logger.Error(err)
	}

	go u.Send(ctx, toEmailAddr, subject, body)
	return nil
}

// Send email send
func (u *EmailUsecase) Send(ctx context.Context, toEmailAddr, subject, body string) {
	u.logger.Infof("try to send email to %s", toEmailAddr)

	m := email.NewEmail()
	fromName := mime.QEncoding.Encode("utf-8", u.conf.FromName)
	m.From = fmt.Sprintf("%s <%s>", fromName, u.conf.UserName)
	m.To = []string{toEmailAddr}
	m.Subject = subject
	m.HTML = common.StringToBytes(body)

	if err := u.emailPool.Send(m, -1); err != nil {
		u.logger.Errorf("send email to %s failed: %s", toEmailAddr, err)
	} else {
		u.logger.Infof("send email to %s success", toEmailAddr)
	}

}

func (u *EmailUsecase) VerifyCode(ctx context.Context, email, code string, emailType emailentity.EmailType) error {
	key := emailcodeKey(email, emailType)
	content, err := u.emailRepo.GetCode(ctx, key)
	if err != nil {
		return err
	}
	if content == "" {
		u.logger.Warn("email code content is null")
		return errorx.BadRequest(reason.EmailCodeExpired)
	}

	data := &emailentity.EmailCodeContent{}
	err = data.FromJSONString(content)
	if err != nil {
		u.logger.Errorf("unmarshal redis data to email code content failed: %s", err)
		return errorx.BadRequest(reason.EmailCodeExpired)
	}
	if data.Code != code {
		return errorx.BadRequest(reason.EmailCodeIncorrect)
	}

	return nil
}

func (u *EmailUsecase) GenCode(ctx context.Context) string {
	rand.Seed(time.Now().Unix())
	code := bytes.Buffer{}
	code.Grow(6)
	codeMapLen := len(codeMap)
	for i := 0; i < 6; i++ {
		code.WriteByte(codeMap[rand.Intn(codeMapLen)])
	}
	return code.String()
}

func (u *EmailUsecase) RegisterTemplate(ctx context.Context, code string) (string, string, error) {
	// 解析 HTML 模板
	t := template.Must(template.New("register").Parse(registerTemplate))

	// 构造 HTML 邮件正文
	data := struct {
		Code string
	}{
		Code: code,
	}
	var body bytes.Buffer
	err := t.Execute(&body, data)
	if err != nil {
		return "", "", errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return registerSubject, body.String(), nil
}

func (u *EmailUsecase) ChangeTemplate(ctx context.Context, code string, changeName string) (string, string, error) {
	// 解析 HTML 模板
	t := template.Must(template.New("change").Parse(changeTemplate))

	// 构造 HTML 邮件正文
	data := struct {
		Code       string
		ChangeName string
		TTL        time.Duration
	}{
		Code:       code,
		ChangeName: changeName,
		TTL:        u.conf.CodeTTL,
	}
	var body bytes.Buffer
	err := t.Execute(&body, data)
	if err != nil {
		return "", "", errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return changeSubject, body.String(), nil
}

func (u *EmailUsecase) ResetPasswordTemplate(ctx context.Context, code string) (string, string, error) {
	// 解析 HTML 模板
	t := template.Must(template.New("resetpassword").Parse(resetPasswordTemplate))

	// 构造 HTML 邮件正文
	data := struct {
		Code string
		TTL  time.Duration
	}{
		Code: code,
		TTL:  u.conf.CodeTTL,
	}
	var body bytes.Buffer
	err := t.Execute(&body, data)
	if err != nil {
		return "", "", errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return resetPasswordSubject, body.String(), nil
}
