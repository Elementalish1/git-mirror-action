package main

const (

	// name of remote urls
	UpstreamRemote = "upstream"
	MirrorRemote   = "mirror"

	// default branch names
	OriginalDefaultBranch = "master"
	MirrorDefaultBranch   = "mirror"

	// TempDir name used as temporary directory
	TempDir = "tmp"

	// error messages
	ErrURLNotHTTPS             = "url is not https"
	ErrNoOriginalURL           = "no original repository url provided"
	ErrNoMirrorURL             = "no mirror repository url provided"
	ErrNoPAT                   = "no personal access token provided"
	ErrFailedToBase64DecodePAT = "failed to decode PAT from b64"

	// info messages
	InfoNoOriginalBranch = "no original branch provided, using 'master'"
	InfoNoMirrorBranch   = "no mirror branch provided, using 'mirror'"
	InfoUsingForce       = "git will now use --force to push"
	InfoUsingVerbose     = "using verbose mode"
	InfoUsingTags        = "transferring"

	// input constants
	OriginalURLInputField    = "originalURL"
	OriginalBranchInputField = "originalBranch"

	MirrorURLInputField    = "mirrorURL"
	MirrorBranchInputField = "mirrorBranch"

	PATInputField        = "pat"
	UseForceInputField   = "force"
	UseVerboseInputField = "verbose"
	UseTagsInputField    = "tags"

	// input responses
	UseForceTrue   = "true"
	UseVerboseTrue = "true"
	UseTagsTrue    = "true"

)
