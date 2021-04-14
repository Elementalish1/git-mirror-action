package main

const (

	// UpstreamRemote name of the upstream remote (originalURL)
	UpstreamRemote = "upstream"
	// MirrorRemote name of the mirror remote (mirrorURL)
	MirrorRemote   = "mirror"

	// OriginalDefaultBranch default name of the original branch
	OriginalDefaultBranch = "master"
	// MirrorDefaultBranch default name of the mirror branch
	MirrorDefaultBranch   = "mirror"

	// TempDir name used as temporary directory
	TempDir = "tmp"

	// ErrURLNotHTTPS the provided link is not https
	ErrURLNotHTTPS             = "url is not https"
	// ErrURLEndNotDotGit the url doesn't end with .git
	ErrURLEndNotDotGit         = "url doesn't end with .git, adding it manually (this might be a cause of failure)"
	// ErrNoOriginalURL a url for the original repository was not provided
	ErrNoOriginalURL           = "no original repository url provided"
	// ErrNoMirrorURL a url for the mirror repository was not provided
	ErrNoMirrorURL             = "no mirror repository url provided"
	// ErrNoPAT a PAT was not provided or found
	ErrNoPAT                   = "no personal access token provided"
	// ErrFailedToBase64DecodePAT something went wrong while base64 decoding the PAT
	ErrFailedToBase64DecodePAT = "failed to decode PAT from b64"

	// InfoNoOriginalBranch no original branch was provided, will use the OriginalDefaultBranch
	InfoNoOriginalBranch = "no original branch provided, using 'master'"
	// InfoNoMirrorBranch no mirror branch was provided, will use the OriginalDefaultBranch
	InfoNoMirrorBranch   = "no mirror branch provided, using 'mirror'"
	// InfoUsingForce force input detected, will use the --force argument when pushing
	InfoUsingForce       = "git will now use --force to push"
	// InfoUsingVerbose verbose input detected, will use print out individual git commands
	InfoUsingVerbose     = "using verbose mode"
	// InfoUsingTags tag input detected, will transfer tags when pushing
	InfoUsingTags        = "transferring"


	//defaults for inputs when parsing the user supplied Github Actions config

	OriginalURLInputField    = "originalURL"
	OriginalBranchInputField = "originalBranch"

	MirrorURLInputField    = "mirrorURL"
	MirrorBranchInputField = "mirrorBranch"

	PATInputField        = "pat"
	UseForceInputField   = "force"
	UseVerboseInputField = "verbose"
	UseTagsInputField    = "tags"

	UseForceTrue   = "true"
	UseVerboseTrue = "true"
	UseTagsTrue    = "true"
)
