package cbor_refmt

// From: https://github.com/ipfs/go-ipld-cbor/blob/821d2db12599a4c79963e2c7988f2d77c8e19c7e/refmt.go

import (
	"math/big"
	"sync"

	"github.com/ipfs/go-cid"
	cbornode "github.com/ipfs/go-ipld-cbor"

	"github.com/ipfs/go-ipld-cbor/encoding"

	"github.com/polydawn/refmt/obj/atlas"
)

// This atlas describes the CBOR Tag (42) for IPLD links, such that refmt can marshal and unmarshal them
var cidAtlasEntry = atlas.BuildEntry(cid.Cid{}).
	UseTag(cbornode.CBORTagLink).
	Transform().
	TransformMarshal(atlas.MakeMarshalTransformFunc(
		castCidToBytes,
	)).
	TransformUnmarshal(atlas.MakeUnmarshalTransformFunc(
		castBytesToCid,
	)).
	Complete()

// BigIntAtlasEntry gives a reasonable default encoding for big.Int. It is not
// included in the entries by default.
var BigIntAtlasEntry = atlas.BuildEntry(big.Int{}).Transform().
	TransformMarshal(atlas.MakeMarshalTransformFunc(
		func(i big.Int) ([]byte, error) {
			return i.Bytes(), nil
		})).
	TransformUnmarshal(atlas.MakeUnmarshalTransformFunc(
		func(x []byte) (big.Int, error) {
			return *big.NewInt(0).SetBytes(x), nil
		})).
	Complete()

// CborAtlas is the refmt.Atlas used by the CBOR IPLD decoder/encoder.
var CborAtlas atlas.Atlas
var atlasEntries = []*atlas.AtlasEntry{cidAtlasEntry}

var (
	cloner       encoding.PooledCloner
	unmarshaller PooledUnmarshaller
	marshaller   PooledMarshaller
)

func init() {
	rebuildAtlas()
}

func rebuildAtlas() {
	CborAtlas = atlas.MustBuild(atlasEntries...).
		WithMapMorphism(atlas.MapMorphism{KeySortMode: atlas.KeySortMode_RFC7049})

	marshaller = NewPooledMarshaller(CborAtlas)
	unmarshaller = NewPooledUnmarshaller(CborAtlas)
	cloner = encoding.NewPooledCloner(CborAtlas)
}

// PooledMarshaller is a thread-safe pooled CBOR marshaller.
type PooledMarshaller struct {
	pool sync.Pool
}

// NewPooledMarshaller returns a PooledMarshaller with the given atlas. Do not
// copy after use.
func NewPooledMarshaller(atl atlas.Atlas) PooledMarshaller {
	return PooledMarshaller{
		pool: sync.Pool{
			New: func() interface{} {
				return encoding.NewMarshallerAtlased(atl)
			},
		},
	}
}

// PooledUnmarshaller is a thread-safe pooled CBOR unmarshaller.
type PooledUnmarshaller struct {
	pool sync.Pool
}

// NewPooledUnmarshaller returns a PooledUnmarshaller with the given atlas. Do
// not copy after use.
func NewPooledUnmarshaller(atl atlas.Atlas) PooledUnmarshaller {
	return PooledUnmarshaller{
		pool: sync.Pool{
			New: func() interface{} {
				return encoding.NewUnmarshallerAtlased(atl)
			},
		},
	}
}

// RegisterCborType allows to register a custom cbor type
func RegisterCborType(i interface{}) {
	var entry *atlas.AtlasEntry
	if ae, ok := i.(*atlas.AtlasEntry); ok {
		entry = ae
	} else {
		entry = atlas.BuildEntry(i).StructMap().AutogenerateWithSortingScheme(atlas.KeySortMode_RFC7049).Complete()
	}
	atlasEntries = append(atlasEntries, entry)
	rebuildAtlas()
}

// From: https://github.com/ipfs/go-ipld-cbor/blob/821d2db12599a4c79963e2c7988f2d77c8e19c7e/node.go

func castBytesToCid(x []byte) (cid.Cid, error) {
	if len(x) == 0 {
		return cid.Cid{}, cbornode.ErrEmptyLink
	}

	// TODO: manually doing multibase checking here since our deps don't
	// support binary multibase yet
	if x[0] != 0 {
		return cid.Cid{}, cbornode.ErrInvalidMultibase
	}

	c, err := cid.Cast(x[1:])
	if err != nil {
		return cid.Cid{}, cbornode.ErrInvalidLink
	}

	return c, nil
}

func castCidToBytes(link cid.Cid) ([]byte, error) {
	if !link.Defined() {
		return nil, cbornode.ErrEmptyLink
	}
	return append([]byte{0}, link.Bytes()...), nil
}
