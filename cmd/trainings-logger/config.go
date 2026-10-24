// START_MODULE_CONTRACT
//   PURPOSE: Supply the public Dropbox app key for local and release builds.
//   SCOPE: Build-time configuration through Go linker flags.
//   DEPENDS: none
//   LINKS: M-CLI-BUILD-CONFIG, V-M-CLI-BUILD-CONFIG, M-RELEASE
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   dropboxAppKey - accept the public release app key through Go linker flags
// END_MODULE_MAP

package main

// START_BLOCK_DROPBOX_APP_KEY
var dropboxAppKey string

// END_BLOCK_DROPBOX_APP_KEY
