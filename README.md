# Github-stars-backup application 

This application clone your starred github repository with all commits, branch, tags etc. to your local disk
Based on https://github.com/Kirill/github-backup

## Dependencies

Application parameters:

    -users  <[user-or-organisation-comma-separated-list]>
    -output [local-folder-name], default: ./repos

Usage examples:

    go run . -users=kirill -output=./tmp

