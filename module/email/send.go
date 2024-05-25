// Package email provides a Starlark module that sends email using Resend API.
package email

import (
	"fmt"

	"github.com/1set/gut/ystring"
	"github.com/1set/starlet"
	"github.com/1set/starlet/dataconv"
	"github.com/1set/starlet/dataconv/types"
	"github.com/PureMature/starcli/util"
	"github.com/resend/resend-go/v2"
	"github.com/samber/lo"
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

// NewModule creates a new Module with the given config loaders.
func NewModule(resendAPIKey, senderDomain ConfigGetter) starlet.ModuleLoader {
	m := &Module{
		resendAPIKey: resendAPIKey,
		senderDomain: senderDomain,
	}
	return m.LoadModule()
}

// LoadModule returns the Starlark module with the given config loaders.
func (m *Module) LoadModule() starlet.ModuleLoader {
	sd := starlark.StringDict{
		"send": m.genSendFunc(),
	}
	return dataconv.WrapModuleData(ModuleName, sd)
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

		// do the actual work here
		fmt.Println("API Key:", resendAPIKey)
		fmt.Println("Sender Domain:", senderDomain)

		// parse args
		newOneOrListStr := func() *util.OneOrMany[starlark.String] { return util.NewOneOrManyNoDefault[starlark.String]() }
		var (
			subject            types.StringOrBytes // must be set
			bodyHTML           types.StringOrBytes // one of the three must be set
			bodyText           types.StringOrBytes
			bodyMarkdown       types.StringOrBytes
			toAddresses        = newOneOrListStr() // one of the three must be set
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
			"to?", toAddresses, "cc?", ccAddresses, "bcc?", bccAddresses,
			"from?", &fromAddress, "from_id?", &fromNameID,
			"reply_to?", &replyAddress, "reply_id?", &replyNameID,
			"attachment_files?", attachmentFiles, "attachment?", attachmentContents); err != nil {
			return starlark.None, err
		}

		// validate args
		if body := []string{bodyHTML.GoString(), bodyText.GoString(), bodyMarkdown.GoString()}; lo.EveryBy(body, ystring.IsBlank) {
			return starlark.None, fmt.Errorf("one of body_html, body_text, or body_markdown must be non-blank")
		}
		if recv := []int{toAddresses.Len(), ccAddresses.Len(), bccAddresses.Len()}; lo.Sum(recv) == 0 {
			return starlark.None, fmt.Errorf("one of to, cc, or bcc must be set")
		}
		if from := []string{fromAddress.GoString(), fromNameID.GoString()}; lo.EveryBy(from, ystring.IsBlank) {
			return starlark.None, fmt.Errorf("one of from or from_id must be non-blank")
		}

		// prepare request
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
		// TODO: reply to, attachments --- https://resend.com/docs/api-reference/emails/send-email

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
			Html:    bodyHTML.GoString(),
			Text:    bodyText.GoString(), // TODO: markdown
			// TODO: reply to, attachments
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
