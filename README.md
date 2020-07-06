# IPFS Mobile Helper

Send files from your mobile device to the InterPlanetary File System

## Summary

I wanted to create as a weekend project a minimally viable process to easily share files from my iPhone to IPFS.

### The Process:

- [x] Hit the share button on a file on your iPhone
- [x] Use the Shortcut "Add to IPFS"
- [x] The Shortcut uses Secure ShellFish (the app needs to be installed on the device) to SFTP the file to a web server
- [x] Next the Shortcut calls a URL on this web server
- [x] The URL on the server triggers this app here
- [x] This app checks the upload directory for new files
- [x] This app uses an IPFS client to tell the IPFS daemon running on the server to add the file
- [x] This app returns the CID to the Shortcut

### Even Better, Let's Share Files to an IPFS Cluster:

- [x] First, figure out how to just share to a single IPFS node
- [x] Set up an IPFS cluster
- [x] Test adding a file to the cluster manually
- [x] Add a file to the cluster programmatically via the ipfs-cluster client when this app starts up
- [ ] Move cluster adder from start-up to new http route (I can choose to share a file from my iPhone to an individual node, or to my cluster)
- [ ] Add option to encrypt files before adding them to the cluster