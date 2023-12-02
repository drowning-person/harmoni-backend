package conf

import (
	"time"
)

type Database struct {
	Driver          string        `protobuf:"bytes,1,opt,name=driver,proto3" json:"driver,omitempty"`
	Source          string        `protobuf:"bytes,2,opt,name=source,proto3" json:"source,omitempty"`
	ConnMaxLifeTime time.Duration `protobuf:"bytes,3,opt,name=conn_max_life_time,json=connMaxLifeTime,proto3" json:"conn_max_life_time,omitempty"`
	MaxOpenConn     int32         `protobuf:"varint,4,opt,name=max_open_conn,json=maxOpenConn,proto3" json:"max_open_conn,omitempty"`
	MaxIdleConn     int32         `protobuf:"varint,5,opt,name=max_idle_conn,json=maxIdleConn,proto3" json:"max_idle_conn,omitempty"`
}

type ServerCommon struct {
	Http *HTTPCommon `protobuf:"bytes,1,opt,name=http,proto3" json:"http,omitempty"`
	Grpc *GRPCCommon `protobuf:"bytes,2,opt,name=grpc,proto3" json:"grpc,omitempty"`
}

type HTTPCommon struct {
	Addr string `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
}

type GRPCCommon struct {
	Addr    string        `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
	Timeout time.Duration `protobuf:"bytes,3,opt,name=timeout,proto3" json:"timeout,omitempty"`
}
