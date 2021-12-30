package client

type postResponseType string

const (
	Success = "success"
	Danger  = "danger"
	Error   = "error"
)

type postResponse struct {
	Type    postResponseType `json:"type"`
	Message string           `json:"msg"`
}

type AliasResponse struct {
	ID      int64  `json:"id"`
	Domain  string `json:"domain"`
	GoTo    string `json:"goto"`
	Address string `json:"address"`
	Active  int64  `json:"active"`
}

type DomainResponse struct {
	Name                    string `json:"domain_name"`
	Description             string `json:"description"`
	Active                  int64  `json:"active"`
	QuotaBytes              int64  `json:"max_quota_for_domain"`
	Mailboxes               int64  `json:"max_num_mboxes_for_domain"`
	MailboxDefaultSizeBytes int64  `json:"def_new_mailbox_quota"`
	MailboxMaxSizeBytes     int64  `json:"max_quota_for_mbox"`
	Aliases                 int64  `json:"max_num_aliases_for_domain"`
}

type MailboxResponse struct {
	Username string `json:"local_part"`
	Domain   string `json:"domain"`
	Email    string `json:"username"`
	Active   int    `json:"active"`
	Name     string `json:"name"`
	Quota    int64  `json:"quota"`
}
