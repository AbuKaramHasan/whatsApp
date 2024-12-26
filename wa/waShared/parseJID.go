package waShared

import (
	"app/definitions"
	"strings"

	"go.mau.fi/whatsmeow/types"
)

func ParseJID(arg string) (types.JID, bool) {
	if arg == "" {
		definitions.Log.Errorf("ParseJID called with an empty argument")
		return types.JID{}, false
	}

	// Remove any leading '+' (common in phone numbers)
	if arg[0] == '+' {
		arg = arg[1:]
	}

	// If no '@' symbol, assume it's a user JID
	if !strings.ContainsRune(arg, '@') {
		jid := types.NewJID(arg, types.DefaultUserServer)
		if jid.User == "" {
			definitions.Log.Errorf("Invalid JID: empty user for input %s", arg)
			return types.JID{}, false
		}
		return jid, true
	}

	// Handle complete JID with '@'
	recipient, err := types.ParseJID(arg)
	if err != nil {
		definitions.Log.Errorf("Invalid JID %s: %v", arg, err)
		return types.JID{}, false
	} else if recipient.User == "" {
		definitions.Log.Errorf("Invalid JID %s: no user specified", arg)
		return types.JID{}, false
	}

	return recipient, true
}
