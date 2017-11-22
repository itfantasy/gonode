package opcode

// ----------------- lobby&room -------------------

const (
	AuthenticateOnce int = 231
	Authenticate         = 230
	JoinLobby            = 229
	LeaveLobby           = 228
	CreateGame           = 227
	JoinGame             = 226
	JoinRandomGame       = 225
	Leave                = 254
	RaiseEvent           = 253
	SetProperties        = 252
	GetProperties        = 251
	ChangeGroups         = 250
	FindFriends          = 222
	GetLobbyStats        = 221
	GetRegions           = 220
	WebRpc               = 219
	ServerSettings       = 218
	GetGameList          = 217
)
