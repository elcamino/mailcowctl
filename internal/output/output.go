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
	default:
		return json.NewEncoder(w).Encode(value)
	}
	return tw.Flush()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
