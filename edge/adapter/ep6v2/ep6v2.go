package ep6v2

import (
	"context"

	Adapter "github.com/maritimusj/chuanyan/gate/adapter"
	AdapterContract "github.com/maritimusj/chuanyan/gate/adapter/contract"
	L "github.com/maritimusj/chuanyan/gate/lang"
)

type adapter struct{}

func (adapter *adapter) Open(ctx context.Context, option AdapterContract.Option) (AdapterContract.Client, error) {
	client := New()
	if err := client.Open(ctx, option); err != nil {
		return nil, err
	}

	return client, nil
}

func init() {
	Adapter.Register("ep6v2", L.Str(L.EP6v2Desc), &adapter{})
}
