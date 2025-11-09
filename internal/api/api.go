package api

type PackageType int

const (
	ClientRequestPkg PackageType = iota
	ServerResponsePkg
)

type ServerResponseType int

const (
	ServerSuccessResponse ServerResponseType = iota
	ServerFailedResponse
)

type ServerFailType int

const (
	ServerNoFail ServerFailType = iota
	ServerClientUnkown
	ServerClientOverdraw
	ServerBlockchainTainted
)

type ClientRequestType int

const (
	ClientCheckBalance ClientRequestType = iota
	ClientTransaction
	ClientCheckBlockchainIntegrity
)

// Defines interface for communication
type Package struct {
	PkgType PackageType
	PayloadSize int64
	Payload []byte
}

// Defines interface for client request
type ClientRequest struct {
	Type ClientRequestType
	Identifier int64
	TransactionValue int64
}

// Defines interface for server response
type ServerResponse struct {
	Type ServerResponseType
	FailType ServerFailType
	ClientBalance int64
}

func SendPackage(pkgType PackageType, payload []byte) error {
	
}

func RecvPackage(pkg *Package) error {

}