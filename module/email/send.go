// Package email provides a Starlark module that sends email using Resend API.
package email

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/1set/gut/ystring"
	"github.com/1set/starlet"
	"github.com/1set/starlet/dataconv"
	"github.com/1set/starlet/dataconv/types"
	"github.com/PureMature/starcli/util"
	"github.com/resend/resend-go/v2"
	"github.com/samber/lo"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	renderer "github.com/yuin/goldmark/renderer/html"
	"go.starlark.net/starlark"
)

// ModuleName defines the expected name for this Module when used in starlark's load() function, eg: load('email', 'send')
const ModuleName = "email"

// ConfigGetter is a function type that returns a string.
type ConfigGetter func() string

// Module wraps the Starlark module with the given config loaders.
type Module struct {
	resendAPIKey ConfigGetter
	senderDomain ConfigGetter
}

// NewModule creates a new bare Module.
func NewModule() starlet.ModuleLoader {
	m := &Module{}
	return m.LoadModule()
}

// NewModuleWithConfig creates a new Module with the given config.
func NewModuleWithConfig(resendAPIKey, senderDomain string) starlet.ModuleLoader {
	m := &Module{
		resendAPIKey: func() string { return resendAPIKey },
		senderDomain: func() string { return senderDomain },
	}
	return m.LoadModule()
}

// NewModuleWithGetter creates a new Module with the given config loaders.
func NewModuleWithGetter(resendAPIKey, senderDomain ConfigGetter) starlet.ModuleLoader {
	m := &Module{
		resendAPIKey: resendAPIKey,
		senderDomain: senderDomain,
	}
	return m.LoadModule()
}

// LoadModule returns the Starlark module with the given config loaders.
func (m *Module) LoadModule() starlet.ModuleLoader {
	sd := starlark.StringDict{
		"set_resend_api_key": m.genSetConfig("resend_api_key"),
		"set_sender_domain":  m.genSetConfig("sender_domain"),
		"send":               m.genSendFunc(),
	}
	return dataconv.WrapModuleData(ModuleName, sd)
}

func (m *Module) genSetConfig(name string) starlark.Callable {
	return starlark.NewBuiltin(ModuleName+".set_"+name, func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var v starlark.String
		if err := starlark.UnpackArgs(b.Name(), args, kwargs, name, &v); err != nil {
			return starlark.None, err
		}
		switch name {
		case "resend_api_key":
			m.resendAPIKey = func() string { return v.GoString() }
		case "sender_domain":
			m.senderDomain = func() string { return v.GoString() }
		}
		return starlark.None, nil
	})
}

func (m *Module) genSendFunc() starlark.Callable {
	return starlark.NewBuiltin(ModuleName+".send", func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		// load config: resend_api_key is required, sender_domain is optional
		var (
			resendAPIKey string
			senderDomain string
		)
		if m.resendAPIKey != nil {
			resendAPIKey = m.resendAPIKey()
		}
		if ystring.IsBlank(resendAPIKey) {
			return starlark.None, fmt.Errorf("resend_api_key is not set")
		}
		if m.senderDomain != nil {
			senderDomain = m.senderDomain()
		}

		// parse args
		newOneOrListStr := func() *util.OneOrMany[starlark.String] { return util.NewOneOrManyNoDefault[starlark.String]() }
		var (
			subject            types.StringOrBytes         // must be set
			bodyHTML           types.NullableStringOrBytes // one of the three must be set
			bodyText           types.NullableStringOrBytes
			bodyMarkdown       types.NullableStringOrBytes
			toAddresses        = newOneOrListStr() // must be set
			ccAddresses        = newOneOrListStr()
			bccAddresses       = newOneOrListStr()
			fromAddress        types.StringOrBytes // one of the two must be set
			fromNameID         types.StringOrBytes
			replyAddress       types.StringOrBytes // two of them are optional
			replyNameID        types.StringOrBytes
			attachmentFiles    = newOneOrListStr()
			attachmentContents = util.NewOneOrManyNoDefault[*starlark.Dict]()
		)
		if err := starlark.UnpackArgs(b.Name(), args, kwargs,
			"subject", &subject,
			"body_html?", &bodyHTML, "body_text?", &bodyText, "body_markdown?", &bodyMarkdown,
			"to", toAddresses, "cc?", ccAddresses, "bcc?", bccAddresses,
			"from?", &fromAddress, "from_id?", &fromNameID,
			"reply_to?", &replyAddress, "reply_id?", &replyNameID,
			"attachment_files?", attachmentFiles, "attachment?", attachmentContents); err != nil {
			return starlark.None, err
		}

		// validate args
		if body := []string{bodyHTML.GoString(), bodyText.GoString(), bodyMarkdown.GoString()}; lo.EveryBy(body, ystring.IsBlank) {
			return starlark.None, fmt.Errorf("one of body_html, body_text, or body_markdown must be non-blank")
		}
		if toAddresses.Len() == 0 {
			return starlark.None, fmt.Errorf("to must be set and non-empty")
		}
		if from := []string{fromAddress.GoString(), fromNameID.GoString()}; lo.EveryBy(from, ystring.IsBlank) {
			return starlark.None, fmt.Errorf("one of from or from_id must be non-blank")
		}

		// convert from to send address
		var sendAddr string
		if fromAddr := fromAddress.GoString(); ystring.IsNotBlank(fromAddr) {
			sendAddr = fromAddr
		} else if fromID := fromNameID.GoString(); ystring.IsNotBlank(fromID) {
			if ystring.IsNotBlank(senderDomain) {
				sendAddr = fromID + "@" + senderDomain
			} else {
				return starlark.None, fmt.Errorf("sender_domain should be set when from_id is used")
			}
		} else {
			return starlark.None, fmt.Errorf("no valid from or from_id found")
		}

		// prepare request
		convGoString := func(v []starlark.String) []string {
			l := make([]string, len(v))
			for i, vv := range v {
				l[i] = dataconv.StarString(vv)
			}
			return l
		}
		req := &resend.SendEmailRequest{
			From:    sendAddr,
			To:      convGoString(toAddresses.Slice()),
			Cc:      convGoString(ccAddresses.Slice()),
			Bcc:     convGoString(bccAddresses.Slice()),
			Subject: subject.GoString(),
		}

		// for body content
		if !bodyHTML.IsNullOrEmpty() {
			// directly use HTML content
			req.Html = bodyHTML.GoString()
		} else if !bodyText.IsNullOrEmpty() {
			// directly use text content
			req.Text = bodyText.GoString()
		} else if !bodyMarkdown.IsNullOrEmpty() {
			// convert markdown to HTML
			markdown := goldmark.New(
				goldmark.WithRendererOptions(
					renderer.WithUnsafe(),
				),
				goldmark.WithExtensions(
					extension.Strikethrough,
					extension.Table,
					extension.Linkify,
				),
			)
			html := bytes.NewBufferString("")
			_ = markdown.Convert([]byte(bodyMarkdown.GoString()), html)
			req.Html = html.String()
		}

		// for attachments
		if fps := attachmentFiles.Slice(); len(fps) > 0 {
			// load file content and attach
			for _, r := range fps {
				fp := r.GoString()
				c, err := ioutil.ReadFile(fp)
				if err != nil {
					return starlark.None, err
				}
				n := filepath.Base(fp)
				req.Attachments = append(req.Attachments, &resend.Attachment{
					Filename: n,
					Content:  c,
				})
			}
		}
		if dcts := attachmentContents.Slice(); len(dcts) > 0 {
			// convert dict to attachment and attach
			for _, r := range dcts {
				fn, ok, err := r.Get(starlark.String("name"))
				if !ok || err != nil {
					return starlark.None, fmt.Errorf("attachment must have a name")
				}
				ct, ok, err := r.Get(starlark.String("content"))
				if !ok || err != nil {
					return starlark.None, fmt.Errorf("attachment must have content")
				}
				req.Attachments = append(req.Attachments, &resend.Attachment{
					Filename: dataconv.StarString(fn),
					Content:  []byte(dataconv.StarString(ct)),
				})
			}
		}

		// send it
		client := resend.NewClient(resendAPIKey)
		sent, err := client.Emails.Send(req)
		if err != nil {
			return starlark.None, err
		}
		return starlark.String(sent.Id), nil
	})
}
