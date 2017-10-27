# Orion
A tool for running static code analysis on Github projects.

### Status: SUPER-ALPHA
This project is in very early development. The dev makes no guarantees or assurances regarding its functionality or state. Please check back periodically to see the progress Orion will make!

### App Dependencies
This tool uses and expects several other tools to be installed in order to function properly. If any of the dependencies are not installed, Orion will throw an error and exit when starting up.
- git
- gas ([Go code analysis tool](https://github.com/GoASTScanner/gas))

### Languages Currently Supported
- Go (Golang) - in development

### TODO List
GitHub
- [x] Fully integrate with Github PR webhooks
- [ ] Figure out permissions scheme for reporting via GitHub comments
- [ ] Find out if a service account is required

Analysis
- [ ] Define form & methodology/workflow for reporting analysis results
- [ ] Determine list of languages to implement (based on availability of packages/analysis tools)
- [ ] Determine what types of analysis are desirable

Database
- [ ] Establish database design for storing results (to avoid duplicate reports)
- [ ] Handle reporting when a security defect has been re-introduced to the codebase
