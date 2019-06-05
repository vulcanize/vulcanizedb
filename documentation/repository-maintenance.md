# Repository Maintenance

## Diagrams
- Diagrams were created with [draw.io](draw.io).
- To update a diagram:
  1. Go to [draw.io](draw.io).
  1. Click on *File > Open from* and choose the location of the diagram you want to update.
  1. Once open in draw.io, you may update it.
  1. Export the diagram to this repository's directory and add commit it.
  
  
## Generating the Changelog
We use [github-changelog-generator](https://github.com/github-changelog-generator/github-changelog-generator) to 
generate release Changelogs. To be consistent with previous Changelogs, the following flags should be passed to the 
command:

```
--user vulcanize
--project vulcanizedb
--token {YOUR_GITHUB_TOKEN}
--no-issues
--usernames-as-github-logins
--since-tag {PREVIOUS_RELEASE_TAG}
```

More information on why your github token is needed, and how to generate it here: 
https://github.com/github-changelog-generator/github-changelog-generator#github-token 