# IPFS Mobile Helper

Send files from your mobile device to the InterPlanetary File System

## Summary

I wanted to create as a weekend project a minimally viable process to easily share files from my iPhone to IPFS. The user can choose between adding a file to an individual IPFS node, or adding the file to an IPFS Cluster.

This is a baby step towards what I really want: an InterGenerational Archive for family data. If you want to make it easy for family members to add data to an archive, you have to make it easy for them to add it from their mobile device.

### The Process

#### Demo - https://youtu.be/_oecLoDLzbY

- [x] Hit the share button on a file on your iPhone
- [x] Use the Shortcut "Add to IPFS"
- [x] The Shortcut uses Secure ShellFish (the app needs to be installed on the device) to SFTP the file to a web server
- [x] Next the Shortcut calls a URL on this web server
- [x] The URL on the server triggers this app here
- [x] This app checks the upload directory for new files
- [x] This app uses an IPFS client to tell the IPFS daemon running on the server to add the file
- [x] This app returns the CID to the Shortcut

### Even Better, Let's Share Files to an IPFS Cluster

- [x] First, figure out how to just share to a single IPFS node
- [x] Set up an IPFS cluster
- [x] Test adding a file to the cluster manually
- [x] Add a file to the cluster programmatically via the ipfs-cluster client when this app starts up
- [x] Move cluster adder from start-up to new http route (I can choose to share a file from my iPhone to an individual node, or to my cluster)
- [ ] Add option to encrypt files before adding them to the cluster (the family archive can hold family secrets)
- [ ] Secure the 'add' endpoints
- [ ] Make it pretty

### Shortcuts Taken

My focus is on learning the IPFS ecosystem (there's a lot there!). I don't want to get side-tracked in iOS development just yet. So, to get my data off of my iPhone I used the Apple Shortcuts app. Shortcuts can be added to the system-wide share sheet, very convenient, and they can pick up capabilities from other installed apps. Secure ShellFish provides the SFTP step in Shortcuts. The shortcut transfers the files to the server running this ipfs-mobile-helper, then it triggers either the vanilla IPFS add action, or the IPFS Cluster add action.

The IPFS 'add' actions read their corresponding upload directory, and then use their respective IPFS API client to communicate with the daemons running on the same machine. There is no clean-up of the upload directories, and no security to prevent someone triggering the 'add' command on the data in the upload directories (SSH secures the upload directories themselves).
