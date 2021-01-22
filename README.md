# aws-s3-peek

Preview large s3 objects before downloading them.

## Install

    go get github.com/steinfletcher/aws-s3-peek
    
## Use

    aws-s3-peek --profile prod --bucket my-bucket --key my-key --range bytes=0-1024
