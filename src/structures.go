package main

// Errors

type ErrorResults struct {
    Error string `json:"Error"`
}








// Blockchain

// Internal Structures

type BlockchainStats struct {
	ID      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		AdjustedTime              int    `json:"adjusted_time"`
		AltBlocksCount            int    `json:"alt_blocks_count"`
		BlockSizeLimit            int    `json:"block_size_limit"`
		BlockSizeMedian           int    `json:"block_size_median"`
		BlockWeightLimit          int    `json:"block_weight_limit"`
		BlockWeightMedian         int    `json:"block_weight_median"`
		BootstrapDaemonAddress    string `json:"bootstrap_daemon_address"`
		BusySyncing               bool   `json:"busy_syncing"`
		Credits                   int    `json:"credits"`
		CumulativeDifficulty      int64  `json:"cumulative_difficulty"`
		CumulativeDifficultyTop64 int    `json:"cumulative_difficulty_top64"`
		DatabaseSize              int64  `json:"database_size"`
		Difficulty                int64  `json:"difficulty"`
		DifficultyTop64           int    `json:"difficulty_top64"`
		FreeSpace                 int64  `json:"free_space"`
		GreyPeerlistSize          int    `json:"grey_peerlist_size"`
		Height                    int    `json:"height"`
		HeightWithoutBootstrap    int    `json:"height_without_bootstrap"`
		IncomingConnectionsCount  int    `json:"incoming_connections_count"`
		Mainnet                   bool   `json:"mainnet"`
		Nettype                   string `json:"nettype"`
		Offline                   bool   `json:"offline"`
		OutgoingConnectionsCount  int    `json:"outgoing_connections_count"`
		RPCConnectionsCount       int    `json:"rpc_connections_count"`
		Stagenet                  bool   `json:"stagenet"`
		StartTime                 int    `json:"start_time"`
		Status                    string `json:"status"`
		Synchronized              bool   `json:"synchronized"`
		Target                    int    `json:"target"`
		TargetHeight              int    `json:"target_height"`
		Testnet                   bool   `json:"testnet"`
		TopBlockHash              string `json:"top_block_hash"`
		TopHash                   string `json:"top_hash"`
		TxCount                   int    `json:"tx_count"`
		TxPoolSize                int    `json:"tx_pool_size"`
		Untrusted                 bool   `json:"untrusted"`
		UpdateAvailable           bool   `json:"update_available"`
		Version                   string `json:"version"`
		WasBootstrapEverUsed      bool   `json:"was_bootstrap_ever_used"`
		WhitePeerlistSize         int    `json:"white_peerlist_size"`
		WideCumulativeDifficulty  string `json:"wide_cumulative_difficulty"`
		WideDifficulty            string `json:"wide_difficulty"`
	} `json:"result"`
}

type BlockchainBlock struct {
	ID      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Blob        string `json:"blob"`
		BlockHeader struct {
			BlockSize                 int    `json:"block_size"`
			BlockWeight               int    `json:"block_weight"`
			CumulativeDifficulty      int64  `json:"cumulative_difficulty"`
			CumulativeDifficultyTop64 int    `json:"cumulative_difficulty_top64"`
			Depth                     int    `json:"depth"`
			Difficulty                int    `json:"difficulty"`
			DifficultyTop64           int    `json:"difficulty_top64"`
			Hash                      string `json:"hash"`
			Height                    int    `json:"height"`
			LongTermWeight            int    `json:"long_term_weight"`
			MajorVersion              int    `json:"major_version"`
			MinerTxHash               string `json:"miner_tx_hash"`
			MinorVersion              int    `json:"minor_version"`
			Nonce                     int    `json:"nonce"`
			NumTxes                   int    `json:"num_txes"`
			OrphanStatus              bool   `json:"orphan_status"`
			PowHash                   string `json:"pow_hash"`
			PrevHash                  string `json:"prev_hash"`
			Reward                    int64  `json:"reward"`
			Timestamp                 int    `json:"timestamp"`
			WideCumulativeDifficulty  string `json:"wide_cumulative_difficulty"`
			WideDifficulty            string `json:"wide_difficulty"`
		} `json:"block_header"`
		Credits     int    `json:"credits"`
		JSON        string `json:"json"`
		MinerTxHash string `json:"miner_tx_hash"`
		Status      string `json:"status"`
		TopHash     string `json:"top_hash"`
		Untrusted   bool   `json:"untrusted"`
	} `json:"result"`
}


// API Structures

type v1XcashBlockchainUnauthorizedStats struct {
	Height                 int `json:"height"`
	Hash                   string `json:"hash"`
	Reward                 int64 `json:"reward"`
	Size                   int64 `json:"size"`
	Version                int `json:"version"`
	VersionBlockHeight     int `json:"versionBlockHeight"`
	NextVersionBlockHeight int `json:"nextVersionBlockHeight"`
	TotalTx                int `json:"totalTx"`
	CirculatingSupply      int64 `json:"circulatingSupply"`
	GeneratedSupply        int64 `json:"generatedSupply"`
	TotalSupply            int64 `json:"totalSupply"`
	EmissionReward         int64 `json:"emissionReward"`
        EmissionHeight         int `json:"emissionHeight"`
    	EmissionTime           int `json:"emissionTime"`
        InflationHeight        int `json:"inflationHeight"`
    	InflationTime          int `json:"inflationTime"`
}
