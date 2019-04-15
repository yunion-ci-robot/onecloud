package scheduler

import (
	"yunion.io/x/jsonutils"

	"yunion.io/x/onecloud/pkg/apis"
	"yunion.io/x/onecloud/pkg/apis/compute"
)

type ScheduleBaseConfig struct {
	BestEffort      bool     `json:"best_effort"`
	SuggestionLimit int64    `json:"suggestion_limit"`
	SuggestionAll   bool     `json:"suggestion_all"`
	IgnoreFilters   []string `json:"ignore_filters"`
	SessionId       string   `json:"session_id"`

	// usedby test api
	RecordLog bool `json:"record_to_history"`
	Details   bool `json:"details"`
}

type ForGuest struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GroupRelation struct {
	GroupId  string `json:"group_id"`
	Strategy string `json:"strategy"`
	Scope    string `json:"scope"`
}

type ServerConfig struct {
	*compute.ServerConfigs

	Memory      int    `json:"vmem_size"`
	Ncpu        int    `json:"vcpu_count"`
	Name        string `json:"name"`
	GuestStatus string `json:"guest_status"`

	// HostId used by migrate
	HostId string `json:"host_id"`

	// DEPRECATED
	Metadata       map[string]string `json:"__meta__"`
	ForGuests      []*ForGuest       `json:"for_guests"`
	GroupRelations []*GroupRelation  `json:"group_releations"`
	Groups         interface{}       `json:"groups"`
	Id             string            `json:"id"`
}

// ScheduleInput used by scheduler sync-schedule/test/forecast api
type ScheduleInput struct {
	apis.Meta

	ScheduleBaseConfig

	ServerConfig
}

func (input ScheduleInput) ToConditionInput() *jsonutils.JSONDict {
	ret := input.JSON(input)
	// old condition compatible
	ret.Add(jsonutils.NewString(input.Project), "owner_tenant_id")
	return ret
}

type CandidateDisk struct {
	Index     int    `json:"index"`
	StorageId string `json:"storage_id"`
}

type CandidateResource struct {
	HostId string           `json:"host_id"`
	Name   string           `json:"name"`
	Disks  []*CandidateDisk `json:"disks"`

	// used by backup schedule
	BackupCandidate *CandidateResource `json:"backup_candidate"`

	// Error means no candidate found, include reasons
	Error string `json:"error"`
}

type ScheduleOutput struct {
	apis.Meta

	Candidates []*CandidateResource `json:"candidates"`
}
