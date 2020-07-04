# IPFS Mobile Helper

Send files from your mobile device to the InterPlanetary File System

## Summary

I wanted to create as a weekend project a minimally viable process to easily share files from my iPhone to IPFS.

### The Process:

- [x] Hit the share button on a file on your iPhone
- [x] Use the Shortcut "Add to IPFS"
- [x] The Shortcut uses Secure Shellfish (the app needs to be installed on the device) to SFTP the file to a web server
- [ ] Next the Shortcut calls a URL on this web server
- [ ] The URL on the server triggers this app here
- [ ] This app checks the upload directory for new files
- [x] This app uses an IPFS client to tell the IPFS daemon running on the server to add the file
- [ ] This app returns the CID to the Shortcut