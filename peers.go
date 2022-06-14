package mycache

import "github.com/hfpublic/mycache/pd"

type PeerPicker interface {
	PickPeer(key string) (PeerGetter, bool)
}

type PeerGetter interface {
	Get(in *pd.Request, out *pd.Response) error
}
