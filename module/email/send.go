// Package email provides a Starlark module that sends email using Resend API.
package email

import (
	"fmt"

	"github.com/1set/gut/ystring"
	"github.com/1set/starlet"
	"github.com/1set/starlet/dataconv"
	"github.com/1set/starlet/dataconv/types"
	"github.com/PureMature/starcli/util"
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
		newOneOrList := func() starlark.Unpacker { return util.NewOneOrManyNoDefault[starlark.String]() }
		var (
			subject            types.StringOrBytes // must be set
			bodyHTML           types.StringOrBytes // one of the three must be set
			bodyText           types.StringOrBytes
			bodyMarkdown       types.StringOrBytes
			toAddresses        = newOneOrList() // one of the three must be set
			ccAddresses        = newOneOrList()
			bccAddresses       = newOneOrList()
			fromAddress        types.StringOrBytes // one of the two must be set
			fromName           types.StringOrBytes
			replyAddress       types.StringOrBytes
			replyName          types.StringOrBytes
			attachmentFiles    = newOneOrList()
			attachmentContents = util.NewOneOrManyNoDefault[*starlark.Dict]()
		)
		if err := starlark.UnpackArgs(b.Name(), args, kwargs,
			"subject", &subject,
			"body_html?", &bodyHTML, "body_text?", &bodyText, "body_markdown?", &bodyMarkdown,
			"to?", &toAddresses, "cc?", &ccAddresses, "bcc?", &bccAddresses,
			"from?", &fromAddress, "from_name?", &fromName,
			"reply_to?", &replyAddress, "reply_name?", &replyName,
			"attachment_files?", &attachmentFiles, "attachment?", &attachmentContents); err != nil {
			return starlark.None, err
		}

		return starlark.None, nil
	})
}
