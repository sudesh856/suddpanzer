package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (

	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)

	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type WorkSpec struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ScenarioYaml  string                 `protobuf:"bytes,1,opt,name=scenario_yaml,json=scenarioYaml,proto3" json:"scenario_yaml,omitempty"`
	AgentId       string                 `protobuf:"bytes,2,opt,name=agent_id,json=agentId,proto3" json:"agent_id,omitempty"`
	Region        string                 `protobuf:"bytes,3,opt,name=region,proto3" json:"region,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *WorkSpec) Reset() {
	*x = WorkSpec{}
	mi := &file_proto_suddpanzer_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *WorkSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorkSpec) ProtoMessage() {}

func (x *WorkSpec) ProtoReflect() protoreflect.Message {
	mi := &file_proto_suddpanzer_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*WorkSpec) Descriptor() ([]byte, []int) {
	return file_proto_suddpanzer_proto_rawDescGZIP(), []int{0}
}

func (x *WorkSpec) GetScenarioYaml() string {
	if x != nil {
		return x.ScenarioYaml
	}
	return ""
}

func (x *WorkSpec) GetAgentId() string {
	if x != nil {
		return x.AgentId
	}
	return ""
}

func (x *WorkSpec) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

type MetricsSnapshot struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	AgentId       string                 `protobuf:"bytes,1,opt,name=agent_id,json=agentId,proto3" json:"agent_id,omitempty"`
	Region        string                 `protobuf:"bytes,2,opt,name=region,proto3" json:"region,omitempty"`
	Timestamp     int64                  `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Rps           float64                `protobuf:"fixed64,4,opt,name=rps,proto3" json:"rps,omitempty"`
	P50Ms         int64                  `protobuf:"varint,5,opt,name=p50_ms,json=p50Ms,proto3" json:"p50_ms,omitempty"`
	P95Ms         int64                  `protobuf:"varint,6,opt,name=p95_ms,json=p95Ms,proto3" json:"p95_ms,omitempty"`
	P99Ms         int64                  `protobuf:"varint,7,opt,name=p99_ms,json=p99Ms,proto3" json:"p99_ms,omitempty"`
	ErrorRatePct  float64                `protobuf:"fixed64,8,opt,name=error_rate_pct,json=errorRatePct,proto3" json:"error_rate_pct,omitempty"`
	TotalRequests int64                  `protobuf:"varint,9,opt,name=total_requests,json=totalRequests,proto3" json:"total_requests,omitempty"`
	ErrorCount    int64                  `protobuf:"varint,10,opt,name=error_count,json=errorCount,proto3" json:"error_count,omitempty"`
	VusActive     int64                  `protobuf:"varint,11,opt,name=vus_active,json=vusActive,proto3" json:"vus_active,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MetricsSnapshot) Reset() {
	*x = MetricsSnapshot{}
	mi := &file_proto_suddpanzer_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MetricsSnapshot) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetricsSnapshot) ProtoMessage() {}

func (x *MetricsSnapshot) ProtoReflect() protoreflect.Message {
	mi := &file_proto_suddpanzer_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*MetricsSnapshot) Descriptor() ([]byte, []int) {
	return file_proto_suddpanzer_proto_rawDescGZIP(), []int{1}
}

func (x *MetricsSnapshot) GetAgentId() string {
	if x != nil {
		return x.AgentId
	}
	return ""
}

func (x *MetricsSnapshot) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

func (x *MetricsSnapshot) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *MetricsSnapshot) GetRps() float64 {
	if x != nil {
		return x.Rps
	}
	return 0
}

func (x *MetricsSnapshot) GetP50Ms() int64 {
	if x != nil {
		return x.P50Ms
	}
	return 0
}

func (x *MetricsSnapshot) GetP95Ms() int64 {
	if x != nil {
		return x.P95Ms
	}
	return 0
}

func (x *MetricsSnapshot) GetP99Ms() int64 {
	if x != nil {
		return x.P99Ms
	}
	return 0
}

func (x *MetricsSnapshot) GetErrorRatePct() float64 {
	if x != nil {
		return x.ErrorRatePct
	}
	return 0
}

func (x *MetricsSnapshot) GetTotalRequests() int64 {
	if x != nil {
		return x.TotalRequests
	}
	return 0
}

func (x *MetricsSnapshot) GetErrorCount() int64 {
	if x != nil {
		return x.ErrorCount
	}
	return 0
}

func (x *MetricsSnapshot) GetVusActive() int64 {
	if x != nil {
		return x.VusActive
	}
	return 0
}

type AgentStatus struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	AgentId       string                 `protobuf:"bytes,1,opt,name=agent_id,json=agentId,proto3" json:"agent_id,omitempty"`
	Status        string                 `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	Message       string                 `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AgentStatus) Reset() {
	*x = AgentStatus{}
	mi := &file_proto_suddpanzer_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AgentStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AgentStatus) ProtoMessage() {}

func (x *AgentStatus) ProtoReflect() protoreflect.Message {
	mi := &file_proto_suddpanzer_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}


func (*AgentStatus) Descriptor() ([]byte, []int) {
	return file_proto_suddpanzer_proto_rawDescGZIP(), []int{2}
}

func (x *AgentStatus) GetAgentId() string {
	if x != nil {
		return x.AgentId
	}
	return ""
}

func (x *AgentStatus) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *AgentStatus) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_proto_suddpanzer_proto protoreflect.FileDescriptor

const file_proto_suddpanzer_proto_rawDesc = "" +
	"\n" +
	"\x16proto/suddpanzer.proto\x12\n" +
	"suddpanzer\"b\n" +
	"\bWorkSpec\x12#\n" +
	"\rscenario_yaml\x18\x01 \x01(\tR\fscenarioYaml\x12\x19\n" +
	"\bagent_id\x18\x02 \x01(\tR\aagentId\x12\x16\n" +
	"\x06region\x18\x03 \x01(\tR\x06region\"\xc6\x02\n" +
	"\x0fMetricsSnapshot\x12\x19\n" +
	"\bagent_id\x18\x01 \x01(\tR\aagentId\x12\x16\n" +
	"\x06region\x18\x02 \x01(\tR\x06region\x12\x1c\n" +
	"\ttimestamp\x18\x03 \x01(\x03R\ttimestamp\x12\x10\n" +
	"\x03rps\x18\x04 \x01(\x01R\x03rps\x12\x15\n" +
	"\x06p50_ms\x18\x05 \x01(\x03R\x05p50Ms\x12\x15\n" +
	"\x06p95_ms\x18\x06 \x01(\x03R\x05p95Ms\x12\x15\n" +
	"\x06p99_ms\x18\a \x01(\x03R\x05p99Ms\x12$\n" +
	"\x0eerror_rate_pct\x18\b \x01(\x01R\ferrorRatePct\x12%\n" +
	"\x0etotal_requests\x18\t \x01(\x03R\rtotalRequests\x12\x1f\n" +
	"\verror_count\x18\n" +
	" \x01(\x03R\n" +
	"errorCount\x12\x1d\n" +
	"\n" +
	"vus_active\x18\v \x01(\x03R\tvusActive\"Z\n" +
	"\vAgentStatus\x12\x19\n" +
	"\bagent_id\x18\x01 \x01(\tR\aagentId\x12\x16\n" +
	"\x06status\x18\x02 \x01(\tR\x06status\x12\x18\n" +
	"\amessage\x18\x03 \x01(\tR\amessage2U\n" +
	"\x0fSuddpanzerAgent\x12B\n" +
	"\vRunScenario\x12\x14.suddpanzer.WorkSpec\x1a\x1b.suddpanzer.MetricsSnapshot0\x01B'Z%github.com/sudesh856/suddpanzer/protob\x06proto3"

var (
	file_proto_suddpanzer_proto_rawDescOnce sync.Once
	file_proto_suddpanzer_proto_rawDescData []byte
)

func file_proto_suddpanzer_proto_rawDescGZIP() []byte {
	file_proto_suddpanzer_proto_rawDescOnce.Do(func() {
		file_proto_suddpanzer_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_suddpanzer_proto_rawDesc), len(file_proto_suddpanzer_proto_rawDesc)))
	})
	return file_proto_suddpanzer_proto_rawDescData
}

var file_proto_suddpanzer_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proto_suddpanzer_proto_goTypes = []any{
	(*WorkSpec)(nil),       
	(*MetricsSnapshot)(nil), 
	(*AgentStatus)(nil),     
}
var file_proto_suddpanzer_proto_depIdxs = []int32{
	0, 
	1, 
	1, 
	0, 
	0, 
	0, 
	0, 
}

func init() { file_proto_suddpanzer_proto_init() }
func file_proto_suddpanzer_proto_init() {
	if File_proto_suddpanzer_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_suddpanzer_proto_rawDesc), len(file_proto_suddpanzer_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_suddpanzer_proto_goTypes,
		DependencyIndexes: file_proto_suddpanzer_proto_depIdxs,
		MessageInfos:      file_proto_suddpanzer_proto_msgTypes,
	}.Build()
	File_proto_suddpanzer_proto = out.File
	file_proto_suddpanzer_proto_goTypes = nil
	file_proto_suddpanzer_proto_depIdxs = nil
}
