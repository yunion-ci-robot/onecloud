package core

import (
	"yunion.io/x/jsonutils"

	schedapi "yunion.io/x/onecloud/pkg/apis/scheduler"
	"yunion.io/x/onecloud/pkg/compute/baremetal"
	computemodels "yunion.io/x/onecloud/pkg/compute/models"
	"yunion.io/x/onecloud/pkg/scheduler/api"
	"yunion.io/x/onecloud/pkg/scheduler/core/score"
)

const (
	PriorityStep int = 1
)

type FailedCandidate struct {
	Stage     string
	Candidate Candidater
	Reasons   []PredicateFailureReason
}

type FailedCandidates struct {
	Candidates []FailedCandidate
}

type SelectPlugin interface {
	Name() string
	OnPriorityEnd(*Unit, Candidater)
	OnSelectEnd(u *Unit, c Candidater, count int64)
}

type Kind int

const (
	KindFree Kind = iota
	KindRaw
	KindReserved
)

type CandidatePropertyGetter interface {
	Id() string
	Name() string
	Zone() *computemodels.SZone
	Host() *computemodels.SHost
	Cloudprovider() *computemodels.SCloudprovider
	Region() *computemodels.SCloudregion
	HostType() string
	HostSchedtags() []computemodels.SSchedtag
	Storages() []*api.CandidateStorage
	Networks() []computemodels.SNetwork
	Status() string
	HostStatus() string
	Enabled() bool
	IsEmpty() bool
	ResourceType() string
	NetInterfaces() map[string][]computemodels.SNetInterface
	ProjectGuests() map[string]int64
	CreatingGuestCount() int

	RunningCPUCount() int64
	TotalCPUCount(useRsvd bool) int64
	FreeCPUCount(useRsvd bool) int64

	RunningMemorySize() int64
	TotalMemorySize(useRsvd bool) int64
	FreeMemorySize(useRsvd bool) int64

	StorageInfo() []*baremetal.BaremetalStorage
	GetFreeStorageSizeOfType(storageType string, useRsvd bool) int64

	GetFreePort(netId string) int
}

// Candidater replace host Candidate resource info
type Candidater interface {
	Getter() CandidatePropertyGetter
	// IndexKey return candidate cache item's ident, usually host ID
	IndexKey() string
	Type() int

	GetSchedDesc() *jsonutils.JSONDict
	GetGuestCount() int64
	GetResourceType() string
}

// HostPriority represents the priority of scheduling to particular host, higher priority is better.
type HostPriority struct {
	// Name of the host
	Host string
	// Score associated with the host
	Score Score
	// Resource wraps Candidate host info
	Candidate Candidater
}

type HostPriorityList []HostPriority

func (h HostPriorityList) Len() int {
	return len(h)
}

func (h HostPriorityList) Less(i, j int) bool {
	s1 := h[i].Score.ScoreBucket
	s2 := h[j].Score.ScoreBucket
	preferLess := score.PreferLess(s1, s2)
	avoidLess := score.AvoidLess(s1, s2)
	normalLess := score.NormalLess(s1, s2)

	if !(preferLess || avoidLess || normalLess) {
		return h[i].Host < h[j].Host
	}

	if preferLess {
		return true
	}
	if avoidLess {
		return false
	}
	return normalLess
}

func (h HostPriorityList) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

type FitPredicate interface {
	// Get filter's name
	Name() string
	Clone() FitPredicate
	PreExecute(*Unit, []Candidater) (bool, error)
	Execute(*Unit, Candidater) (bool, []PredicateFailureReason, error)
}

type PredicateFailureReason interface {
	GetReason() string
}

type PriorityPreFunction func(*Unit, []Candidater) (bool, []PredicateFailureReason, error)

// PriorityMapFunction is a function that computes per-resource results for a given resource.
type PriorityMapFunction func(*Unit, Candidater) (HostPriority, error)

// PriorityReduceFunction is a function that aggregated per-resource results and computes
// final scores for all hosts.
type PriorityReduceFunction func(*Unit, []Candidater, HostPriorityList) error

type PriorityConfig struct {
	Name   string
	Pre    PriorityPreFunction
	Map    PriorityMapFunction
	Reduce PriorityReduceFunction
	Weight int
}

type Priority interface {
	Name() string
	Clone() Priority
	Map(*Unit, Candidater) (HostPriority, error)
	Reduce(*Unit, []Candidater, HostPriorityList) error
	PreExecute(*Unit, []Candidater) (bool, []PredicateFailureReason, error)

	// Score intervals
	ScoreIntervals() score.Intervals
}

type AllocatedResource struct {
	Disks []*schedapi.CandidateDisk `json:"disks"`
}

func NewAllocatedResource() *AllocatedResource {
	return &AllocatedResource{
		Disks: make([]*schedapi.CandidateDisk, 0),
	}
}
