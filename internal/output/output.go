package output

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/tob/mailcowctl/internal/client"
	"gopkg.in/yaml.v3"
)

func Render(w io.Writer, format string, value any) error {
	switch format {
	case "json":
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(value)
	case "yaml":
		return yaml.NewEncoder(w).Encode(value)
	case "table":
		return renderTable(w, value)
	default:
		return fmt.Errorf("unsupported output format %q", format)
	}
}

func renderTable(w io.Writer, value any) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	switch v := value.(type) {
	case []client.Domain:
		fmt.Fprintln(tw, "DOMAIN\tACTIVE\tMAILBOXES\tQUOTA")
		for _, d := range v {
			fmt.Fprintf(tw, "%s\t%v\t%d/%d\t%s/%s\n", d.Name(), d.Active, d.MailboxesInDomain, d.MaxNumMailboxesForDomain, d.QuotaUsedInDomain, d.MaxQuotaForDomain)
		}
	case client.Domain:
		return renderTable(w, []client.Domain{v})
	case []client.Mailbox:
		fmt.Fprintln(tw, "ADDRESS\tNAME\tACTIVE\tQUOTA")
		for _, m := range v {
			fmt.Fprintf(tw, "%s\t%s\t%v\t%v/%s\n", m.Username, m.Name, m.Active, m.QuotaUsed, m.Quota)
		}
	case client.Mailbox:
		return renderTable(w, []client.Mailbox{v})
	case []client.Alias:
		fmt.Fprintln(tw, "ID\tADDRESS\tGOTO\tACTIVE")
		for _, a := range v {
			fmt.Fprintf(tw, "%d\t%s\t%s\t%v\n", a.ID, a.Address, a.Goto, a.Active)
		}
	case client.Alias:
		return renderTable(w, []client.Alias{v})
	case []client.Dkim:
		fmt.Fprintln(tw, "DOMAIN\tSELECTOR\tLENGTH\tDNS-TXT")
		for _, d := range v {
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", d.Domain, d.Selector, d.Length, truncate(d.DkimTxt, 50))
		}
	case client.Dkim:
		return renderTable(w, []client.Dkim{v})
	case []client.SyncJob:
		fmt.Fprintln(tw, "ID\tMAILBOX\tHOST\tUSER1\tINTERVAL\tACTIVE\tLAST-RUN")
		for _, j := range v {
			fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%v\t%v\t%s\n", j.ID, j.Mailbox, j.Host1, j.User1, j.MinsInterval, j.Active, j.LastRun)
		}
	case client.SyncJob:
		return renderTable(w, []client.SyncJob{v})
	case []client.AppPassword:
		fmt.Fprintln(tw, "ID\tMAILBOX\tNAME\tPROTOCOLS\tACTIVE")
		for _, a := range v {
			fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%v\n", a.ID, a.Mailbox, a.Name, a.Protocols(), a.Active)
		}
	case client.AppPassword:
		return renderTable(w, []client.AppPassword{v})
	case []client.Filter:
		fmt.Fprintln(tw, "ID\tMAILBOX\tDESC\tTYPE\tACTIVE")
		for _, f := range v {
			fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%v\n", f.ID, f.Username, f.ScriptDesc, f.FilterType, f.Active)
		}
	case client.Filter:
		return renderTable(w, []client.Filter{v})
	case []client.PolicyItem:
		fmt.Fprintln(tw, "PREFID\tKIND\tDOMAIN\tVALUE")
		for _, p := range v {
			fmt.Fprintf(tw, "%d\t%s\t%s\t%s\n", p.PrefID, p.Kind, p.Object, p.Value)
		}
	case client.PolicyItem:
		return renderTable(w, []client.PolicyItem{v})
	case []client.Transport:
		fmt.Fprintln(tw, "ID\tDESTINATION\tNEXTHOP\tUSERNAME\tPASSWORD\tACTIVE")
		for _, t := range v {
			fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\t%v\n", t.ID, t.Destination, t.Nexthop, t.Username, maskSecret(t.PasswordShort, t.Password), t.Active)
		}
	case client.Transport:
		return renderTable(w, []client.Transport{v})
	case []client.Bcc:
		fmt.Fprintln(tw, "ID\tTYPE\tLOCAL-DEST\tBCC-DEST\tACTIVE")
		for _, b := range v {
			fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%v\n", b.ID, b.Type, b.LocalDest, b.BccDest, b.Active)
		}
	case client.Bcc:
		return renderTable(w, []client.Bcc{v})
	default:
		return json.NewEncoder(w).Encode(value)
	}
	return tw.Flush()
}

// truncate slices by bytes and is intended for ASCII content (e.g. DKIM TXT records).
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

// maskSecret renders a password column for table output without leaking the
// full secret: it prefers the API-provided short form, falls back to a fixed
// mask when a password is set, and is empty otherwise.
func maskSecret(short, full string) string {
	if short != "" {
		return short
	}
	if full != "" {
		return "***"
	}
	return ""
}
