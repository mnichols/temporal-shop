package session

import (
	"context"
	"github.com/temporalio/temporal-shop/api/temporal_shop/values/v1"
	"github.com/temporalio/temporal-shop/services/go/internal/encrypt"
	"go.temporal.io/sdk/client"
	"google.golang.org/protobuf/proto"
)

type Sessioner interface {
	SignalWorkflow(ctx context.Context, workflowID string, runID string, signalName string, arg interface{}) error
	ExecuteWorkflow(ctx context.Context, options client.StartWorkflowOptions, workflow interface{}, args ...interface{}) (client.WorkflowRun, error)
}

type ID struct {
	inner          *values.SessionID
	sessionBytes   []byte
	key            []byte
	associatedData []byte
	encrypted      string
}

func NewID(key, associatedData []byte, sessionID *values.SessionID) (*ID, error) {
	sessBytes, err := proto.Marshal(sessionID)
	if err != nil {
		return nil, err
	}

	id := &ID{
		inner:          sessionID,
		sessionBytes:   sessBytes,
		key:            key,
		associatedData: associatedData,
	}
	sid, err := encrypt.EncryptDeterministically(key, sessBytes, associatedData)
	if err != nil {
		return nil, err
	}
	id.encrypted = sid
	return id, nil
}
func NewFromID(key, value, associatedData []byte) (*ID, error) {
	sessBytes, err := encrypt.DecryptDeterministically(key, []byte(value), associatedData)
	if err != nil {
		return nil, err
	}
	sessionID := &values.SessionID{}
	err = proto.Unmarshal(sessBytes, sessionID)
	if err != nil {
		return nil, err
	}
	return &ID{
		inner:          sessionID,
		sessionBytes:   sessBytes,
		key:            key,
		associatedData: associatedData,
		encrypted:      string(value),
	}, nil
}
func (i ID) SessionID() *values.SessionID {
	return i.inner
}
func (i ID) String() string {
	return i.encrypted
}
